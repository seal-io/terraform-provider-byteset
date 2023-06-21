package sqlx

import (
	"strings"
	"unicode"

	vp "vitess.io/vitess/go/vt/sqlparser"
)

type StatementType = uint

const (
	_ StatementType = 1 << iota
	StatementTypeUnknown
	StatementTypeTCL
	StatementTypeDCL
	StatementTypeDDL
	StatementTypeDML
)

const (
	StatementTypeTCLBegin = iota + StatementTypeTCL
	StatementTypeTCLEnd
)

const (
	StatementTypeDMLSingle = iota + StatementTypeDML
	StatementTypeDMLMultiple
)

// Preview analyzes the beginning of the query using a simpler and faster
// textual comparison to identify the statement type,
// borrows from the vitess.io/vitess/go/vt/sqlparser.
func Preview(sql string) StatementType {
	trimmed := vp.StripLeadingComments(sql)

	if trimmed == "" {
		return StatementTypeUnknown
	}

	if strings.HasPrefix(trimmed, "/*!") {
		// MySQL command.
		return StatementTypeDMLMultiple
	}

	isNotLetter := func(r rune) bool { return !unicode.IsLetter(r) }
	firstWord := strings.TrimLeftFunc(trimmed, isNotLetter)

	if end := strings.IndexFunc(firstWord, unicode.IsSpace); end != -1 {
		firstWord = firstWord[:end]
	}

	// Comparison is done in order of priority.
	loweredFirstWord := strings.ToLower(firstWord)
	switch loweredFirstWord {
	case "select", "insert", "replace", "update", "delete",
		"copy":
		return StatementTypeDMLSingle
	case "stream", "vstream", "revert":
		return StatementTypeDCL
	case "savepoint", "lock":
		return StatementTypeTCLBegin
	case "unlock":
		return StatementTypeTCLEnd
	}

	// For the following statements it is not sufficient to rely
	// on loweredFirstWord. This is because they are not statements
	// in the grammar and we are relying on Preview to parse them.
	// For instance, we don't want: "BEGIN JUNK" to be parsed
	// as StmtBegin.
	trimmedNoComments, _ := vp.SplitMarginComments(trimmed)
	switch strings.ToLower(trimmedNoComments) {
	case "begin", "start transaction":
		return StatementTypeTCLBegin
	case "commit":
		return StatementTypeTCLEnd
	case "rollback":
		return StatementTypeTCLEnd
	}
	switch loweredFirstWord {
	case "create", "alter", "rename", "drop", "truncate",
		"grant", "revoke":
		return StatementTypeDDL
	case "flush":
		return StatementTypeDCL
	case "set", "use":
		return StatementTypeDMLMultiple
	case "show",
		"describe", "desc", "explain",
		"analyze", "repair", "optimize":
		return StatementTypeDMLSingle
	case "release", "rollback":
		return StatementTypeTCLEnd
	}

	return StatementTypeUnknown
}
