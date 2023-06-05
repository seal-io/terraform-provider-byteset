package pipeline

import (
	"bufio"
	"bytes"
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
)

type Source interface {
	io.Closer

	Pipe(ctx context.Context, destination Destination) error
}

func NewSource(
	ctx context.Context,
	addr string,
	opts ...Option,
) (Source, error) {
	switch {
	case strings.HasPrefix(addr, "file://"):
		addr = strings.TrimPrefix(addr, "file://")

		local, err := os.Open(addr)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot open local file from %s: %w",
				addr,
				err,
			)
		}

		return srcFile{f: local}, nil

	case strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://"):
		remote, err := http.Get(addr)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot open remote file from %s: %w",
				addr,
				err,
			)
		}

		return srcFile{f: remote.Body}, nil

	default:
	}

	drv, db, err := sqlx.LoadDatabase(addr)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot load database from %s: %w",
			addr,
			err,
		)
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err = sqlx.IsDatabaseConnected(ctx, db)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot connect database on %s: %w",
			addr,
			err,
		)
	}

	for i := range opts {
		if opts[i] == nil {
			continue
		}

		opts[i](db)
	}

	return srcDatabase{drv: drv, db: db}, nil
}

type srcFile struct {
	f io.ReadCloser
}

func (s srcFile) Close() error {
	return s.f.Close()
}

func (s srcFile) Pipe(ctx context.Context, dst Destination) error {
	line := func(data []byte, eof bool) (int, []byte, error) {
		if eof && len(data) == 0 {
			return 0, nil, nil
		}

		var (
			i int
			d = data
		)

		for {
			if j := bytes.IndexByte(d, '\n'); j >= 0 {
				if (j == 1 && d[j-1] == ';') ||
					(j > 1 && (d[j-1] == ';' || d[j-2] == ';')) {
					return i + j + 1, bytes.TrimLeft(
						bytes.TrimRight(data[0:i+j], ";\r"),
						"\n",
					), nil
				}

				if j+1 >= len(d) {
					break
				}
				d = d[j+1:]
				i += j + 1
			}
		}

		if eof {
			return len(
					data,
				), bytes.TrimLeft(
					bytes.TrimRight(data, ";\r"),
					"\n",
				), nil
		}

		return 0, nil, nil
	}

	qs := bufio.NewScanner(s.f)
	qs.Split(line)

	for qs.Scan() {
		q := qs.Text()

		err := dst.Exec(ctx, q)
		if err != nil {
			return fmt.Errorf("cannot execute %q: %w", q, err)
		}

		tflog.Debug(ctx, "Executed", map[string]any{"query": q})
	}

	return nil
}

type srcDatabase struct {
	drv string
	db  *sql.DB
}

func (s srcDatabase) Close() error {
	return s.db.Close()
}

func (s srcDatabase) Pipe(
	ctx context.Context,
	dst Destination,
) error {
	return nil
}
