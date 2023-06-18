package pipeline

import (
	"bufio"
	"context"
	stdsql "database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/seal-io/terraform-provider-byteset/utils/sqlx"
	"github.com/seal-io/terraform-provider-byteset/utils/strx"
)

type Source interface {
	io.Closer

	Pipe(ctx context.Context, destination Destination) error
}

func NewSource(ctx context.Context, addr string, addrConnMax int) (Source, error) {
	switch {
	case strings.HasPrefix(addr, "file://"):
		addr = strings.TrimPrefix(addr, "file://")

		local, err := os.Open(addr)
		if err != nil {
			return nil, fmt.Errorf("cannot open local file from %q: %w", addr, err)
		}

		return &srcFile{f: local}, nil

	case strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://"):
		remote, err := http.Get(addr)
		if err != nil {
			return nil, fmt.Errorf("cannot open remote file from %q: %w", addr, err)
		}

		return &srcFile{f: remote.Body}, nil

	case strings.HasPrefix(addr, "raw://"):
		raw := addr[len("raw://"):]
		return &srcFile{f: io.NopCloser(strings.NewReader(raw))}, nil

	case strings.HasPrefix(addr, "raw+base64://"):
		raw, err := strx.DecodeBase64(addr[len("raw+base64://"):])
		if err != nil {
			return nil, fmt.Errorf("cannot decode raw base64 content: %w", err)
		}

		return &srcFile{f: io.NopCloser(strings.NewReader(raw))}, nil

	default:
	}

	// Load database.
	drv, db, err := sqlx.LoadDatabase(addr, addrConnMax)
	if err != nil {
		return nil, fmt.Errorf("cannot load database from %q: %w", addr, err)
	}

	// Detect connectivity.
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	err = sqlx.IsDatabaseConnected(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("cannot connect database on %q: %w", addr, err)
	}

	return &srcDatabase{
		drv: drv,
		db:  db,
	}, nil
}

type srcFile struct {
	f io.ReadCloser
}

func (in *srcFile) Close() error {
	return in.f.Close()
}

func (in *srcFile) Pipe(ctx context.Context, dst Destination) error {
	ss := bufio.NewScanner(in.f)
	ss.Split(split)

	for ss.Scan() {
		err := dst.Exec(ctx, ss.Text())
		if err != nil {
			return err
		}
	}

	return dst.Flush(ctx)
}

type srcDatabase struct {
	drv string
	db  *stdsql.DB
}

func (in *srcDatabase) Close() error {
	return in.db.Close()
}

func (in *srcDatabase) Pipe(ctx context.Context, dst Destination) error {
	return nil
}
