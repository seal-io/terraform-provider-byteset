package sqlx

import (
	"context"
	stdsql "database/sql"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (stdsql.Result, error)
}

// Exec executes the given SQL with the executor.
func Exec(ctx context.Context, ex Executor, sql string, args ...any) error {
	_, err := ex.ExecContext(ctx, sql, args...)
	if err != nil {
		if !IsIgnorableError(err) {
			return err
		}
	}

	tflog.Debug(ctx, "Executed", map[string]any{"sql": sql, "args": args})

	return nil
}

var ignorableErrMessages = []string{
	"sql: no rows in result set",
	"Error 1065", // MySQL (Query was empty).
}

func IsIgnorableError(err error) bool {
	if err != nil {
		m := err.Error()
		for i := range ignorableErrMessages {
			if strings.Contains(m, ignorableErrMessages[i]) {
				return true
			}
		}
	}

	return false
}
