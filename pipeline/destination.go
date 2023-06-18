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

	"github.com/seal-io/terraform-provider-byteset/utils/sqlx"
)

type Destination interface {
	io.Closer

	// Flush executes all caching sql.
	Flush(ctx context.Context) error

	// Exec executes the given sql.
	Exec(ctx context.Context, sql string) error
}

func NewDestination(ctx context.Context, addr string, addrConnMax, bufItemMax int) (Destination, error) {
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
		drv:        drv,
		db:         db,
		dbConnMax:  db.Stats().MaxOpenConnections,
		buf:        map[string][]string{},
		bufItemMax: bufItemMax,
	}, nil
}

type dst struct {
	drv       string
	db        *stdsql.DB
	dbConnMax int

	buf        map[string][]string
	bufItemMax int

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

	var ex sqlx.Executor = in.db
	if in.st != nil {
		ex = in.st
	}

	for prefix := range in.buf {
		// Construct DML(insert).
		var sb strings.Builder

		sb.WriteString(prefix)
		sb.WriteString("VALUES ")

		for i := range in.buf[prefix] {
			if i != 0 {
				sb.WriteString(", ")
			}

			sb.WriteString(in.buf[prefix][i])
		}

		// Execute DML(insert) in one session.
		if err := sqlx.Exec(ctx, ex, sb.String()); err != nil {
			return err
		}
	}

	in.buf = map[string][]string{}

	return nil
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
			// Cache inserts.
			in.buf[inst.Prefix] = append(in.buf[inst.Prefix], inst.Values...)

			// Flush if reaches limitation.
			if len(in.buf[inst.Prefix]) >= in.bufItemMax || len(in.buf) >= in.dbConnMax {
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

			if err = sqlx.Exec(ctx, c, sql); err != nil {
				return err
			}
		}

		return nil
	}

	return errors.New("nothing to do")
}
