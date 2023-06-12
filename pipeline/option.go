package pipeline

import (
	"database/sql"
	"time"
)

type Option func(db *sql.DB)

func WithConnMaxOpen(o int) Option {
	return func(db *sql.DB) {
		if o <= 0 {
			o = 5
		}

		db.SetMaxOpenConns(o)
	}
}

func WithConnMaxIdle(i int) Option {
	return func(db *sql.DB) {
		if i <= 0 {
			i = 5
		}

		db.SetMaxIdleConns(i)
	}
}

func WithConnMaxLife(d time.Duration) Option {
	return func(db *sql.DB) {
		db.SetConnMaxLifetime(d)
	}
}
