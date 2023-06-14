package pipeline

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/seal-io/terraform-provider-byteset/utils/sqlx"
	"github.com/seal-io/terraform-provider-byteset/utils/strx"
)

type Source interface {
	io.Closer

	Pipe(ctx context.Context, destination Destination) error
}

func NewSource(ctx context.Context, addr string, opts ...Option) (Source, error) {
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

	drv, db, err := sqlx.LoadDatabase(addr)
	if err != nil {
		return nil, fmt.Errorf("cannot load database from %q: %w", addr, err)
	}

	// Configure.
	for i := range opts {
		if opts[i] == nil {
			continue
		}

		opts[i](db)
	}

	// Detect connectivity.
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err = sqlx.IsDatabaseConnected(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("cannot connect database on %q: %w", addr, err)
	}

	return &srcDatabase{drv: drv, db: db}, nil
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
		s := ss.Text()

		err := dst.Exec(ctx, s)
		if err != nil {
			return err
		}

		tflog.Debug(ctx, "Executed", map[string]any{"query": s})
	}

	return nil
}

type srcDatabase struct {
	drv string
	db  *sql.DB
}

func (in *srcDatabase) Close() error {
	return in.db.Close()
}

func (in *srcDatabase) Pipe(ctx context.Context, dst Destination) error {
	return nil
}
