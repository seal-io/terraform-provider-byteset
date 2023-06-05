package pipeline

import (
	"database/sql"
	"time"
)

type Option func(db *sql.DB)

func WithConnMaxOpen(o int) Option {
	return func(db *sql.DB) {
		if o <= 0 {
			return
		}

		db.SetMaxOpenConns(o)
	}
}

func WithConnMaxIdle(i int) Option {
	return func(db *sql.DB) {
		if i <= 0 {
			return
		}

		db.SetMaxIdleConns(i)
	}
}

func WithConnMaxLife(d time.Duration) Option {
	return func(db *sql.DB) {
		if d <= 0 {
			return
		}

		db.SetConnMaxLifetime(d)
	}
}
