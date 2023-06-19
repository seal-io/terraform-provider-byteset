package pipeline

import (
	"context"
	stdsql "database/sql"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sourcegraph/conc/pool"

	"github.com/seal-io/terraform-provider-byteset/utils/sqlx"
)

type Destination interface {
	io.Closer

	// Flush executes all caching sql.
	Flush(ctx context.Context) error

	// Exec executes the given sql.
	Exec(ctx context.Context, sql string) error
}

func NewDestination(ctx context.Context, addr string, addrConnMax, bufSegCap int) (Destination, error) {
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

	return &dst{
		drv:       drv,
		db:        db,
		dbConnMax: db.Stats().MaxOpenConnections,
		buf:       map[string][][]string{},
		bufSegCap: bufSegCap,
	}, nil
}

type dst struct {
	drv       string
	db        *stdsql.DB
	dbConnMax int

	buf       map[string][][]string
	bufSegCap int

	st  *stdsql.Conn
	stC int
}

func (in *dst) Close() error {
	if in.st != nil {
		_ = in.st.Close()
	}

	return in.db.Close()
}

func (in *dst) Flush(ctx context.Context) error {
	if len(in.buf) == 0 {
		return nil
	}

	// Construct DML(insert).
	sqls := make([]string, 0,
		func() (s int) {
			for i := range in.buf {
				s += len(in.buf[i])
			}

			return
		}())

	for p := range in.buf {
		for i := 0; i < len(in.buf[p]); i++ {
			sqls = append(sqls,
				p+"VALUES "+strings.Join(in.buf[p][i], ", "))
		}
	}

	in.buf = map[string][][]string{}

	// Execute DML(insert).
	var ex sqlx.Executor = in.db
	if in.st != nil {
		ex = in.st
	}

	gp := pool.New().
		WithMaxGoroutines(in.dbConnMax).
		WithContext(ctx).
		WithFirstError()

	for i := 0; i < len(sqls); i++ {
		sql := sqls[i]

		gp.Go(func(ctx context.Context) error {
			return sqlx.Exec(ctx, ex, sql)
		})
	}

	return gp.Wait()
}

func (in *dst) Exec(ctx context.Context, sql string) error {
	sqlp, err := sqlx.Parse(in.drv, sql)
	if err != nil {
		return fmt.Errorf("failed to parse sql %q: %w", sql, err)
	}

	if sqlp.Unknown() {
		tflog.Trace(ctx, "Ignored", map[string]any{"sql": sql})
		return nil
	}

	err = in.exec(ctx, sqlp, sql)
	if err != nil {
		return fmt.Errorf("failed to execute sql %q: %w", sql, err)
	}

	return nil
}

func (in *dst) exec(ctx context.Context, sqlp sqlx.Parsed, sql string) error {
	if typ, ok := sqlp.TCL(); ok {
		// Prepare sentry.
		if in.st == nil {
			conn, err := in.db.Conn(ctx)
			if err != nil {
				return err
			}
			in.st = conn
		}

		// Flush.
		if err := in.Flush(ctx); err != nil {
			return err
		}

		// Execute TCL in sentry session.
		if err := sqlx.Exec(ctx, in.st, sql); err != nil {
			return err
		}

		if typ == sqlx.StartTCL {
			in.stC += 1
		} else {
			in.stC -= 1
		}

		// Release sentry.
		if in.stC <= 0 {
			err := in.st.Close()
			if err != nil {
				return err
			}
			in.st = nil
		}

		return nil
	}

	if sqlp.DCL() || sqlp.DDL() {
		// Flush.
		if err := in.Flush(ctx); err != nil {
			return err
		}

		var ex sqlx.Executor = in.db
		if in.st != nil {
			ex = in.st
		}

		// Execute DDL/DCL in one session.
		return sqlx.Exec(ctx, ex, sql)
	}

	if typ, ok := sqlp.DML(); ok {
		inst, ok := sqlp.AsDMLInsert()
		if ok {
			// Prepare first buffer segment.
			if len(in.buf[inst.Prefix]) == 0 {
				in.buf[inst.Prefix] = append(in.buf[inst.Prefix], nil)
			}

			lsi := len(in.buf[inst.Prefix]) - 1

			// Append latest buffer segment.
			in.buf[inst.Prefix][lsi] = append(in.buf[inst.Prefix][lsi], inst.Values...)

			// Increase segment of buffer.
			if len(in.buf[inst.Prefix][lsi]) >= in.bufSegCap {
				if len(in.buf[inst.Prefix]) < in.dbConnMax {
					in.buf[inst.Prefix] = append(in.buf[inst.Prefix], nil)
					lsi += 1
				}
			}

			// Flush if reaches limitations.
			isBufFull := in.dbConnMax != 1 && len(in.buf) >= in.dbConnMax
			isBufSegFull := lsi >= in.dbConnMax && len(in.buf[inst.Prefix][lsi]) >= in.bufSegCap

			if isBufFull || isBufSegFull {
				return in.Flush(ctx)
			}

			return nil
		}

		// Flush.
		if err := in.Flush(ctx); err != nil {
			return err
		}

		// Execute DML in sentry session if found.
		if in.st != nil {
			return sqlx.Exec(ctx, in.st, sql)
		}

		// Execute DML in single session.
		if typ == sqlx.SingleSessionDML {
			return sqlx.Exec(ctx, in.db, sql)
		}

		// Execute DML in multiple sessions.
		cs := make([]*stdsql.Conn, 0, in.dbConnMax)
		defer func() {
			for i := range cs {
				_ = cs[i].Close()
			}
		}()

		for i := 0; i < in.dbConnMax; i++ {
			c, err := in.db.Conn(ctx)
			if err != nil {
				return err
			}

			cs = append(cs, c)
		}

		gp := pool.New().
			WithMaxGoroutines(in.dbConnMax).
			WithContext(ctx).
			WithFirstError()

		for i := 0; i < in.dbConnMax; i++ {
			c := cs[i]

			gp.Go(func(ctx context.Context) error {
				return sqlx.Exec(ctx, c, sql)
			})
		}

		return gp.Wait()
	}

	return errors.New("nothing to do")
}
