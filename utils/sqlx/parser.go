package sqlx

import (
	"strings"

	vp "vitess.io/vitess/go/vt/sqlparser"

	cp "github.com/cockroachdb/cockroach/pkg/sql/parser"
	cps "github.com/cockroachdb/cockroach/pkg/sql/parser/statements"
	cpt "github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

type TCLLevel uint

const (
	StartTCL TCLLevel = iota + 1
	EndTCL
)

type DMLLevel uint

const (
	SingleSessionDML DMLLevel = iota + 1
	MultiSessionsDML
)

type DMLInsert struct {
	Prefix string
	Values []string
}

type Parsed interface {
	// Origin returns the original SQL.
	Origin() string
	// Unknown returns true if not matched any languages.
	Unknown() bool
	// TCL returns the TCLLevel and true if the Origin is transaction control language.
	TCL() (TCLLevel, bool)
	// DCL returns true if the Origin is data control language.
	DCL() bool
	// DDL returns true if the Origin is data definition language.
	DDL() bool
	// DML returns the DMLLevel and true if the Origin is data manipulation language.
	DML() (DMLLevel, bool)
	// AsDMLInsert returns the structuring insert statement if possible.
	AsDMLInsert() (DMLInsert, bool)
}

func Parse(drv, sql string) (Parsed, error) {
	if drv == PostgresDialect {
		// TODO parse COPY FROM statement.
		stmt, err := cp.ParseOne(sql)
		if err != nil {
			if err.Error() != "expected 1 statement, but found 0" {
				return nil, err
			}
		}

		var stmtType cpt.StatementType = -1
		if stmt.AST != nil {
			stmtType = stmt.AST.StatementType()
		}

		return cockroachParsed{
			original: sql,
			stmt:     stmt,
			stmtType: stmtType,
		}, nil
	}

	stmtType := vp.Preview(sql)

	return vitessParsed{
		original: sql,
		stmtType: stmtType,
	}, nil
}

type vitessParsed struct {
	original string
	stmtType vp.StatementType
}

func (p vitessParsed) Origin() string {
	return p.original
}

func (p vitessParsed) Unknown() bool {
	switch p.stmtType {
	case vp.StmtCommentOnly, vp.StmtUnknown:
		return true
	}

	return false
}

func (p vitessParsed) TCL() (TCLLevel, bool) {
	switch p.stmtType {
	case vp.StmtLockTables,
		vp.StmtBegin,
		vp.StmtSavepoint:
		return StartTCL, true
	case vp.StmtUnlockTables,
		vp.StmtCommit, vp.StmtRollback,
		vp.StmtSRollback, vp.StmtRelease:
		return EndTCL, true
	}

	return 0, false
}

func (p vitessParsed) DCL() bool {
	switch p.stmtType {
	case vp.StmtRevert, vp.StmtFlush,
		vp.StmtStream, vp.StmtVStream, vp.StmtCallProc:
		return true
	}

	return false
}

func (p vitessParsed) DDL() bool {
	switch p.stmtType {
	case vp.StmtDDL, vp.StmtPriv:
		return true
	}

	return false
}

func (p vitessParsed) DML() (DMLLevel, bool) {
	switch p.stmtType {
	case vp.StmtSelect, vp.StmtExplain, vp.StmtShow, vp.StmtOther,
		vp.StmtInsert, vp.StmtReplace, vp.StmtUpdate, vp.StmtDelete:
		return SingleSessionDML, true
	case vp.StmtSet, vp.StmtUse, vp.StmtComment:
		return MultiSessionsDML, true
	}

	return 0, false
}

func (p vitessParsed) AsDMLInsert() (DMLInsert, bool) {
	if typ, _ := p.DML(); typ != SingleSessionDML {
		return DMLInsert{}, false
	}

	stmt, err := vp.Parse(p.original)
	if err != nil {
		return DMLInsert{}, false
	}

	in, ok := stmt.(*vp.Insert)
	if !ok ||
		in.Action != vp.InsertAct ||
		len(in.OnDup) != 0 ||
		in.Rows == nil {
		return DMLInsert{}, false
	}

	var is DMLInsert

	switch t := in.Rows.(type) {
	default:
		return DMLInsert{}, false
	case vp.Values:
		for _, v := range t {
			vb := vp.NewTrackedBuffer(nil)
			v.Format(vb)
			is.Values = append(is.Values, vb.String())
		}
	}

	pb := vp.NewTrackedBuffer(nil)
	pb.SetEscapeAllIdentifiers(true)

	pb.WriteString("INSERT ")

	if in.Comments != nil {
		in.Comments.Format(pb)
	}

	if in.Ignore {
		pb.WriteString("IGNORE ")
	}

	pb.WriteString("INTO ")
	in.Table.Format(pb)
	pb.WriteString(" ")
	in.Partitions.Format(pb)
	in.Columns.Format(pb)
	pb.WriteString(" ")
	is.Prefix = pb.String()

	return is, true
}

type cockroachParsed struct {
	original string
	stmt     cps.Statement[cpt.Statement]
	stmtType cpt.StatementType
}

func (p cockroachParsed) Origin() string {
	return p.original
}

func (p cockroachParsed) Unknown() bool {
	return p.stmtType == -1
}

func (p cockroachParsed) TCL() (TCLLevel, bool) {
	if p.stmtType == cpt.TypeTCL {
		switch p.stmt.AST.(type) {
		case *cpt.BeginTransaction:
			return StartTCL, true
		case *cpt.CommitTransaction, *cpt.RollbackTransaction:
			return EndTCL, true
		}
	}

	return 0, false
}

func (p cockroachParsed) DCL() bool {
	return p.stmtType == cpt.TypeDCL
}

func (p cockroachParsed) DDL() bool {
	return p.stmtType == cpt.TypeDDL
}

func (p cockroachParsed) DML() (DMLLevel, bool) {
	if p.stmtType != cpt.TypeDML {
		return 0, false
	}

	st := p.stmt.AST.StatementTag()
	if st == "SET" || strings.HasPrefix(st, "SET ") {
		return MultiSessionsDML, true
	}

	return SingleSessionDML, true
}

func (p cockroachParsed) AsDMLInsert() (DMLInsert, bool) {
	if typ, _ := p.DML(); typ != SingleSessionDML {
		return DMLInsert{}, false
	}

	switch in := p.stmt.AST.(type) {
	case *cpt.Insert:
		if in.OnConflict.IsUpsertAlias() ||
			cpt.HasReturningClause(in.Returning) ||
			in.Rows == nil || in.Rows.Select == nil {
			return DMLInsert{}, false
		}

		var is DMLInsert

		switch t := in.Rows.Select.(type) {
		default:
			return DMLInsert{}, false
		case *cpt.ValuesClause:
			for _, v := range t.Rows {
				vb := cpt.NewFmtCtx(cpt.FmtSimple)
				vb.WriteString("(")
				v.Format(vb)
				vb.WriteString(")")
				is.Values = append(is.Values, vb.String())
			}
		}

		pb := cpt.NewFmtCtx(cpt.FmtSimple)

		pb.WriteString("INSERT INTO ")

		pb.FormatNode(in.Table)
		pb.WriteString(" ")

		if in.Columns != nil {
			pb.WriteString("(")
			pb.FormatNode(&in.Columns)
			pb.WriteString(") ")
		}

		if in.DefaultValues() {
			pb.WriteString("DEFAULT")
		}
		is.Prefix = pb.String()

		return is, true

	case *cpt.CopyFrom:
		// TODO.
	}

	return DMLInsert{}, false
}
