// Package pgsql helps build SQL queries.
package pgsql

import (
	"fmt"
	"strings"
)

type SelectClause struct {
	IsDistinct         bool
	DistinctOnExprList []string
	ExprList           []string
}

func (s *SelectClause) Distinct() {
	s.IsDistinct = true
}

func (s *SelectClause) DistinctOn(sql string, args *Args, values ...interface{}) {
	s.Distinct()
	if len(values) > 0 {
		sql = args.Format(sql, values...)
	}
	s.DistinctOnExprList = append(s.DistinctOnExprList, sql)
}

func (s *SelectClause) Select(sql string, args *Args, values ...interface{}) {
	if len(values) > 0 {
		sql = args.Format(sql, values...)
	}
	s.ExprList = append(s.ExprList, sql)
}

func (s *SelectClause) String() string {
	sb := &strings.Builder{}
	sb.WriteString("select")
	if s.IsDistinct {
		sb.WriteString(" distinct")
	}

	if len(s.DistinctOnExprList) > 0 {
		sb.WriteString(" on (")
		writeExprList(sb, s.DistinctOnExprList)
		sb.WriteString(")")
	}

	sb.WriteString(" ")
	if len(s.ExprList) > 0 {
		writeExprList(sb, s.ExprList)
	} else {
		sb.WriteString("*")
	}

	return sb.String()
}

type FromClause struct {
	Value string
}

func (f *FromClause) From(sql string, args *Args, values ...interface{}) {
	f.Value = (args.Format(sql, values...))
}

func (f *FromClause) String() string {
	if f.Value == "" {
		return ""
	}

	return "from " + f.Value
}

type WhereClause struct {
	Cond string
}

func (wc *WhereClause) String() string {
	if wc.Cond == "" {
		return ""
	}

	return "where " + wc.Cond
}

func (wc *WhereClause) Where(sql string, args *Args, values ...interface{}) {
	if len(values) > 0 {
		sql = args.Format(sql, values...)
	}

	if wc.Cond == "" {
		wc.Cond = sql
	} else {
		wc.Cond = fmt.Sprintf("(%s and %s)", wc.Cond, sql)
	}
}

func (wc *WhereClause) Or(sql string, args *Args, values ...interface{}) {
	if len(values) > 0 {
		sql = args.Format(sql, values...)
	}

	if wc.Cond == "" {
		wc.Cond = sql
	} else {
		wc.Cond = fmt.Sprintf("(%s or %s)", wc.Cond, sql)
	}
}

type OrderByClause struct {
	ExprList []string
}

func (o *OrderByClause) OrderBy(sql string, args *Args, values ...interface{}) {
	if len(values) > 0 {
		sql = args.Format(sql, values...)
	}
	o.ExprList = append(o.ExprList, sql)
}

func (o *OrderByClause) String() string {
	if len(o.ExprList) == 0 {
		return ""
	}

	sb := &strings.Builder{}
	sb.WriteString("order by ")
	writeExprList(sb, o.ExprList)
	return sb.String()
}

type LimitClause struct {
	Value string
}

func (l *LimitClause) Limit(sql string, args *Args, values ...interface{}) {
	if len(values) > 0 {
		sql = args.Format(sql, values...)
	}
	l.Value = sql
}

func (l *LimitClause) String() string {
	if l.Value == "" {
		return ""
	}

	return "limit " + l.Value
}

type OffsetClause struct {
	Value string
}

func (o *OffsetClause) Offset(sql string, args *Args, values ...interface{}) {
	if len(values) > 0 {
		sql = args.Format(sql, values...)
	}
	o.Value = (sql)
}

func (o *OffsetClause) String() string {
	if o.Value == "" {
		return ""
	}

	return "offset " + o.Value
}

func writeExprList(sb *strings.Builder, exprList []string) {
	for i, e := range exprList {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(e)
	}
}
