package pipeline

import (
	"database/sql"
	"time"
)

type Option func(db *sql.DB)

func WithConnMax(i int) Option {
	return func(db *sql.DB) {
		if db == nil {
			return
		}

		if i <= 0 {
			i = 25
		}

		db.SetConnMaxIdleTime(0)
		db.SetConnMaxLifetime(0)
		db.SetMaxOpenConns(i)
		db.SetMaxIdleConns(i)
	}
}

func WithConnMaxOpen(i int) Option {
	return func(db *sql.DB) {
		if db == nil {
			return
		}

		if i <= 0 {
			i = 25
		}

		db.SetMaxOpenConns(i)
	}
}

func WithConnMaxIdle(i int) Option {
	return func(db *sql.DB) {
		if db == nil {
			return
		}

		if i <= 0 {
			i = 25
		}

		db.SetMaxIdleConns(i)
	}
}

func WithConnMaxLife(d time.Duration) Option {
	return func(db *sql.DB) {
		if db == nil {
			return
		}

		db.SetConnMaxIdleTime(0)
		db.SetConnMaxLifetime(d)
	}
}
