// Package pgsql helps build SQL queries.
package pgsql

import (
	"fmt"
	"strings"
)

type statementType int

const (
	TypeSelect statementType = iota
	TypeInsert
	TypeUpdate
	TypeDelete
)

func (st statementType) String() string {
	return [...]string{"select", "insert", "update", "delete"}[st]
}

type Statement struct {
	Type          statementType
	SelectClause  SelectClause
	FromClause    FromClause
	WhereClause   WhereClause
	OrderByClause OrderByClause
	LimitClause   LimitClause
	OffsetClause  OffsetClause
	Args          *Args
}

func NewStatement() *Statement {
	return &Statement{Args: &Args{}}
}

type StatementOption func(*Statement) error

func (s *Statement) Apply(funcs ...StatementOption) error {
	for _, f := range funcs {
		err := f(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Statement) Clone() *Statement {
	clone := *s
	clone.Args = s.Args.Clone()
	return &clone
}

func (s *Statement) String() string {
	sb := &strings.Builder{}

	writeCount := 0
	f := func(s string) {
		if s == "" {
			return
		}

		if writeCount > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(s)
		writeCount += 1
	}

	switch s.Type {
	case TypeSelect:
		f(s.SelectClause.String())
		f(s.FromClause.String())
		f(s.WhereClause.String())
		f(s.OrderByClause.String())
		f(s.LimitClause.String())
		f(s.OffsetClause.String())
	default:
		sb.WriteString("unknown type")
	}

	return sb.String()
}

type SelectClause struct {
	IsDistinct         bool
	DistinctOnExprList string
	ExprList           string
}

func Distinct() StatementOption {
	return func(s *Statement) error {
		s.SelectClause.IsDistinct = true
		return nil
	}
}

func DistinctOn(sql string, values ...interface{}) StatementOption {
	return func(s *Statement) error {
		s.SelectClause.IsDistinct = true

		if len(values) > 0 {
			sql = s.Args.Format(sql, values...)
		}

		if len(s.SelectClause.DistinctOnExprList) > 0 {
			s.SelectClause.DistinctOnExprList += ", " + sql
		} else {
			s.SelectClause.DistinctOnExprList = sql
		}

		return nil
	}
}

func Select(sql string, values ...interface{}) StatementOption {
	return func(s *Statement) error {
		s.Type = TypeSelect
		if len(values) > 0 {
			sql = s.Args.Format(sql, values...)
		}

		if len(s.SelectClause.ExprList) > 0 {
			s.SelectClause.ExprList += ", " + sql
		} else {
			s.SelectClause.ExprList = sql
		}

		return nil
	}
}

func (s SelectClause) String() string {
	sb := &strings.Builder{}
	sb.WriteString("select")
	if s.IsDistinct {
		sb.WriteString(" distinct")
	}

	if len(s.DistinctOnExprList) > 0 {
		sb.WriteString(" on (")
		sb.WriteString(s.DistinctOnExprList)
		sb.WriteString(")")
	}

	sb.WriteString(" ")
	if len(s.ExprList) > 0 {
		sb.WriteString(s.ExprList)
	} else {
		sb.WriteString("*")
	}

	return sb.String()
}

type FromClause string

func From(sql string, values ...interface{}) StatementOption {
	return func(s *Statement) error {
		s.FromClause = FromClause(s.Args.Format(sql, values...))
		return nil
	}
}

func (f FromClause) String() string {
	if len(f) == 0 {
		return ""
	}

	return "from " + string(f)
}

type WhereClause string

func (wc WhereClause) String() string {
	if len(wc) == 0 {
		return ""
	}

	return "where " + string(wc)
}

func Where(sql string, values ...interface{}) StatementOption {
	return func(s *Statement) error {
		if len(values) > 0 {
			sql = s.Args.Format(sql, values...)
		}

		if len(s.WhereClause) == 0 {
			s.WhereClause = WhereClause(sql)
		} else {
			s.WhereClause = WhereClause(fmt.Sprintf("(%s and %s)", string(s.WhereClause), sql))
		}

		return nil
	}
}

func Or(sql string, values ...interface{}) StatementOption {
	return func(s *Statement) error {
		if len(values) > 0 {
			sql = s.Args.Format(sql, values...)
		}

		if len(s.WhereClause) == 0 {
			s.WhereClause = WhereClause(sql)
		} else {
			s.WhereClause = WhereClause(fmt.Sprintf("(%s or %s)", string(s.WhereClause), sql))
		}

		return nil
	}
}

type OrderByClause string

func OrderBy(sql string, values ...interface{}) StatementOption {
	return func(s *Statement) error {
		if len(values) > 0 {
			sql = s.Args.Format(sql, values...)
		}

		if len(s.OrderByClause) > 0 {
			s.OrderByClause = OrderByClause(string(s.OrderByClause) + ", " + sql)
		} else {
			s.OrderByClause = OrderByClause(sql)
		}

		return nil
	}
}

func (o OrderByClause) String() string {
	if len(o) == 0 {
		return ""
	}

	return "order by " + string(o)
}

type LimitClause string

func Limit(sql string, values ...interface{}) StatementOption {
	return func(s *Statement) error {
		if len(values) > 0 {
			sql = s.Args.Format(sql, values...)
		}
		s.LimitClause = LimitClause(sql)

		return nil
	}
}

func (l LimitClause) String() string {
	if len(l) == 0 {
		return ""
	}

	return "limit " + string(l)
}

type OffsetClause string

func Offset(sql string, values ...interface{}) StatementOption {
	return func(s *Statement) error {
		if len(values) > 0 {
			sql = s.Args.Format(sql, values...)
		}
		s.OffsetClause = OffsetClause(sql)

		return nil
	}
}

func (o OffsetClause) String() string {
	if len(o) == 0 {
		return ""
	}

	return "offset " + string(o)
}

func writeExprList(sb *strings.Builder, exprList []string) {
	for i, e := range exprList {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(e)
	}
}
