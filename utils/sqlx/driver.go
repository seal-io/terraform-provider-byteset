package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/seal-io/terraform-provider-byteset/utils/wait"

	_ "github.com/denisenkom/go-mssqldb" // Db = mssql.
	_ "github.com/go-sql-driver/mysql"   // Db = mysql/mariadb.
	_ "github.com/lib/pq"                // Db = postgres.
	_ "github.com/sijms/go-ora/v2"       // Db = oracle.
)

const (
	MySQLDialect     = "mysql"
	MariaDBDialect   = MySQLDialect
	PostgresDialect  = "postgres"
	OracleDialect    = "oracle"
	SQLServerDialect = "mssql"
)

func ParseAddress(addr string) (drv, dsn string, err error) {
	if addr == "" {
		err = errors.New("blank data source address")
		return
	}

	switch {
	case strings.HasPrefix(addr, "mysql://"):
		drv = MySQLDialect
		dsn = strings.TrimPrefix(addr, "mysql://")
	case strings.HasPrefix(addr, "maria://"):
		drv = MariaDBDialect
		dsn = strings.TrimPrefix(addr, "maria://")
	case strings.HasPrefix(addr, "mariadb://"):
		drv = MariaDBDialect
		dsn = strings.TrimPrefix(addr, "mariadb://")
	case strings.HasPrefix(addr, "postgres://"):
		drv = PostgresDialect
		dsn = addr
	case strings.HasPrefix(addr, "postgresql://"):
		drv = PostgresDialect
		dsn = "postgres://" + strings.TrimPrefix(addr, "postgresql://")
	case strings.HasPrefix(addr, "oracle://"):
		drv = OracleDialect
		dsn = addr
	case strings.HasPrefix(addr, "sqlserver://"):
		drv = SQLServerDialect
		dsn = addr
	case strings.HasPrefix(addr, "mssql://"):
		drv = SQLServerDialect
		dsn = "sqlserver://" + strings.TrimPrefix(addr, "mssql://")
	}

	if drv == "" {
		err = errors.New("cannot recognize driver from database address")
	}

	return
}

func LoadDatabase(addr string) (drv string, db *sql.DB, err error) {
	drv, dsn, err := ParseAddress(addr)
	if err != nil {
		return
	}

	db, err = sql.Open(drv, dsn)
	if err != nil {
		return
	}

	return
}

func IsDatabaseConnected(ctx context.Context, db *sql.DB) (perr error) {
	err := wait.PollImmediateUntil(2*time.Second,
		func() (bool, error) {
			perr = db.PingContext(ctx)
			if perr != nil {
				tflog.Error(ctx, "Cannot ping database", map[string]any{"error": perr})
			}

			return perr == nil, ctx.Err()
		},
		ctx.Done(),
	)
	if err != nil {
		if perr == nil {
			perr = err
		}
	}

	return
}
