package pipeline

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"time"

	"github.com/golang-sql/sqlexp"
	"github.com/sourcegraph/conc/pool"

	"github.com/seal-io/terraform-provider-byteset/utils/sqlx"
)

type Destination interface {
	io.Closer

	// Exec executes the given query,
	// detects the given query and pick the proper execution(sync/async) to finish the job.
	Exec(ctx context.Context, query string) error
}

func NewDestination(ctx context.Context, addr string, opts ...Option) (Destination, error) {
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

	dbStats := db.Stats()

	return &dst{drv: drv, db: db, dbStats: dbStats}, nil
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

func (in *dst) Exec(ctx context.Context, query string) (err error) {
	if in.dbStats.MaxOpenConnections == 1 {
		err = in.execSync(ctx, query)
	} else {
		err = in.execAsync(ctx, query)
	}

	if err != nil {
		return fmt.Errorf("cannot execute %q: %w", query, err)
	}

	return
}

func (in *dst) execSync(ctx context.Context, query string) error {
	return exec(ctx, in.db, query)
}

func (in *dst) execAsync(ctx context.Context, query string) error {
	if sqlx.IsDCL(query) {
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

		// Execute DCL in any one connection.
		return exec(ctx, in.db, query)
	}

	if sqlx.IsDML(query) {
		// Create a go pool if not existed.
		if in.gp == nil {
			in.gp = pool.New().
				WithMaxGoroutines(in.dbStats.MaxOpenConnections).
				WithContext(ctx).
				WithFirstError()
		}

		// Execute DML in async.
		in.gp.Go(func(ctx context.Context) error {
			return exec(ctx, in.db, query)
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

	if sqlx.IsDDL(query) {
		// Execute DDL in any one connection.
		return exec(ctx, in.db, query)
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

		err = exec(ctx, c, query)
		if err != nil {
			return err
		}
	}

	return nil
}

func exec(ctx context.Context, db sqlexp.Querier, query string) error {
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		if sqlx.IsEmptyError(err) {
			err = nil
		}
	}

	return err
}
