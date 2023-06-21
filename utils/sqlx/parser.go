package sqlx

import (
	vp "vitess.io/vitess/go/vt/sqlparser"

	cp "github.com/cockroachdb/cockroach/pkg/sql/parser"
	cpt "github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

type TCLLevel = uint

const (
	StartTCL TCLLevel = iota + 1
	EndTCL
)

type DMLLevel = uint

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

func Parse(drv, sql string) Parsed {
	return &parsed{
		drv:      drv,
		raw:      sql,
		stmtType: Preview(sql),
	}
}

type parsed struct {
	drv      string
	raw      string
	stmtType StatementType
}

func (p parsed) Origin() string {
	return p.raw
}

func (p parsed) Unknown() bool {
	return p.stmtType == StatementTypeUnknown
}

func (p parsed) TCL() (TCLLevel, bool) {
	b := p.stmtType & StatementTypeTCL
	if b == 0 {
		return 0, false
	}
	return p.stmtType - StatementTypeTCL + 1, true
}

func (p parsed) DCL() bool {
	return p.stmtType == StatementTypeDCL
}

func (p parsed) DDL() bool {
	return p.stmtType == StatementTypeDDL
}

func (p parsed) DML() (DMLLevel, bool) {
	b := p.stmtType & StatementTypeDML
	if b == 0 {
		return 0, false
	}
	return p.stmtType - StatementTypeDML + 1, true
}

func (p parsed) AsDMLInsert() (DMLInsert, bool) {
	if p.stmtType != StatementTypeDMLSingle {
		return DMLInsert{}, false
	}

	if p.drv == PostgresDialect {
		return parsePostgres(p.raw)
	}

	return parse(p.raw)
}

func parsePostgres(raw string) (DMLInsert, bool) {
	stmt, err := cp.ParseOne(raw)
	if err != nil {
		return DMLInsert{}, false
	}

	switch in := stmt.AST.(type) {
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

func parse(raw string) (DMLInsert, bool) {
	stmt, err := vp.Parse(raw)
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
