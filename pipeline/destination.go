package pipeline

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"time"

	"github.com/seal-io/terraform-provider-byteset/utils/sqlx"
)

type Destination interface {
	io.Closer

	Exec(ctx context.Context, query string, args ...any) error
}

func NewDestination(ctx context.Context, addr string, opts ...Option) (Destination, error) {
	drv, db, err := sqlx.LoadDatabase(addr)
	if err != nil {
		return nil, fmt.Errorf("cannot load database from %q: %w", addr, err)
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err = sqlx.IsDatabaseConnected(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("cannot connect database on %q: %w", addr, err)
	}

	for i := range opts {
		if opts[i] == nil {
			continue
		}

		opts[i](db)
	}

	return dst{drv: drv, db: db}, nil
}

type dst struct {
	drv string
	db  *sql.DB
}

func (d dst) Close() error {
	return d.db.Close()
}

func (d dst) Exec(ctx context.Context, query string, args ...any) error {
	_, err := d.db.ExecContext(ctx, query, args...)
	return err
}
