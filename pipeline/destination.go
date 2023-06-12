package pipeline

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"time"

	"github.com/sourcegraph/conc/pool"

	"github.com/seal-io/terraform-provider-byteset/utils/sqlx"
)

type Destination interface {
	io.Closer

	// Exec executes the given statement with arguments,
	// detects the given statement and pick the proper execution(sync/async) to finish the job.
	Exec(ctx context.Context, statement string, args ...any) error
}

func NewDestination(ctx context.Context, addr string, opts ...Option) (Destination, error) {
	drv, db, err := sqlx.LoadDatabase(addr)
	if err != nil {
		return nil, fmt.Errorf("cannot load database from %q: %w", addr, err)
	}

	// Configure.
	opts = append(opts, WithConnMaxLife(0))
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

	return &dst{drv: drv, db: db, dbStats: db.Stats()}, nil
}

type dst struct {
	drv     string
	db      *sql.DB
	dbStats sql.DBStats

	gp *pool.ContextPool
}

func (in *dst) Close() error {
	return in.db.Close()
}

func (in *dst) Exec(ctx context.Context, statement string, args ...any) (err error) {
	if in.dbStats.MaxOpenConnections == 1 {
		err = in.exec(ctx, statement, args...)
	} else {
		err = in.execAsync(ctx, statement, args...)
	}

	if err != nil {
		return fmt.Errorf("cannot execute %q: %w", statement, err)
	}

	return
}

func (in *dst) exec(ctx context.Context, statement string, args ...any) (err error) {
	_, err = in.db.ExecContext(ctx, statement, args...)
	if err != nil {
		if sqlx.IsEmptyError(err) {
			err = nil
		}
	}

	return
}

func (in *dst) execAsync(ctx context.Context, statement string, args ...any) error {
	if sqlx.IsDCL(statement) {
		if in.gp != nil {
			// Wait for all previous DDL finishing.
			if err := in.gp.Wait(); err != nil {
				return err
			}
			in.gp = nil
		}

		// Switch to sync mode.
		in.db.SetMaxIdleConns(1)
		in.db.SetMaxOpenConns(1)
		in.dbStats = in.db.Stats()

		return in.exec(ctx, statement, args...)
	}

	if sqlx.IsDML(statement) {
		// Execute DML in async mode.
		if in.gp == nil {
			in.gp = pool.New().
				WithMaxGoroutines(in.dbStats.MaxOpenConnections).
				WithContext(ctx).
				WithFirstError()
		}

		in.gp.Go(func(ctx context.Context) error {
			return in.exec(ctx, statement, args...)
		})

		return nil
	}

	if in.gp != nil {
		// Wait for all previous DDL finishing.
		if err := in.gp.Wait(); err != nil {
			return err
		}
		in.gp = nil
	}

	if sqlx.IsDDL(statement) {
		// Execute DDL in one connection.
		return in.exec(ctx, statement, args...)
	}

	// Execute DCL for all connections.
	cs := make([]*sql.Conn, 0, in.dbStats.MaxOpenConnections)
	defer func() {
		for i := range cs {
			_ = cs[i].Close()
		}
	}()

	for i := 0; i < in.dbStats.MaxOpenConnections; i++ {
		c, err := in.db.Conn(ctx)
		if err != nil {
			return err
		}

		cs = append(cs, c)

		_, err = c.ExecContext(ctx, statement, args...)
		if err != nil {
			if sqlx.IsEmptyError(err) {
				continue
			}

			return err
		}
	}

	return nil
}
