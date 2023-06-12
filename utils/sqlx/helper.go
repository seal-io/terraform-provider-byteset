package sqlx

import (
	"strings"
)

var dmlPrefixes = []string{
	"INSERT ",
	"DELETE ",
	"UPDATE ",
}

func IsDML(s string) bool {
	if s != "" {
		for i := range dmlPrefixes {
			if strings.HasPrefix(s, dmlPrefixes[i]) {
				return true
			}
		}
	}

	return false
}

var ddlPrefixes = []string{
	"DROP ",
	"ALTER ",
	"CREATE ",
}

func IsDDL(s string) bool {
	if s != "" {
		for i := range ddlPrefixes {
			if strings.HasPrefix(s, ddlPrefixes[i]) {
				return true
			}
		}
	}

	return false
}

var dclPrefixes = []string{
	"LOCK ", "UNLOCK ", // MySQL.
	"COPY ", // Postgres.
}

func IsDCL(s string) bool {
	if s != "" {
		for i := range dclPrefixes {
			if strings.HasPrefix(s, dclPrefixes[i]) {
				return true
			}
		}
	}

	return false
}

var emptyErrMessages = []string{
	"sql: no rows in result set",
	"Error 1065", // MySQL (Query was empty).
}

func IsEmptyError(err error) bool {
	if err != nil {
		m := err.Error()
		for i := range emptyErrMessages {
			if strings.Contains(m, emptyErrMessages[i]) {
				return true
			}
		}
	}

	return false
}
