// Package pgsql helps build SQL queries.
package pgsql

import (
	"fmt"
	"io"
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
	s.writeToSQL(sb)
	return sb.String()
}

func (s *SelectClause) writeToSQL(w io.Writer) {
	io.WriteString(w, "select")
	if s.IsDistinct {
		io.WriteString(w, " distinct")
	}

	if len(s.DistinctOnExprList) > 0 {
		io.WriteString(w, " on (")
		writeExprList(w, s.DistinctOnExprList)
		io.WriteString(w, ")")
	}

	io.WriteString(w, " ")
	if len(s.ExprList) > 0 {
		writeExprList(w, s.ExprList)
	} else {
		io.WriteString(w, "*")
	}
}

type FromClause struct {
	Value string
}

func (f *FromClause) From(sql string, args *Args, values ...interface{}) {
	f.Value = (args.Format(sql, values...))
}

func (f *FromClause) String() string {
	sb := &strings.Builder{}
	f.writeToSQL(sb)
	return sb.String()
}

func (f *FromClause) writeToSQL(w io.Writer) {
	if f.Value == "" {
		return
	}

	io.WriteString(w, "from ")
	io.WriteString(w, f.Value)
}

type WhereClause struct {
	Cond string
}

func (wc *WhereClause) String() string {
	sb := &strings.Builder{}
	wc.writeToSQL(sb)
	return sb.String()
}

func (wc *WhereClause) writeToSQL(w io.Writer) {
	if wc.Cond == "" {
		return
	}

	io.WriteString(w, "where ")
	io.WriteString(w, wc.Cond)
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
	sb := &strings.Builder{}
	o.writeToSQL(sb)
	return sb.String()
}

func (o *OrderByClause) writeToSQL(w io.Writer) {
	if len(o.ExprList) == 0 {
		return
	}

	io.WriteString(w, "order by ")
	writeExprList(w, o.ExprList)
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
	sb := &strings.Builder{}
	l.writeToSQL(sb)
	return sb.String()
}

func (l *LimitClause) writeToSQL(w io.Writer) {
	if l.Value == "" {
		return
	}

	io.WriteString(w, "limit ")
	io.WriteString(w, l.Value)
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
	sb := &strings.Builder{}
	o.writeToSQL(sb)
	return sb.String()
}

func (o *OffsetClause) writeToSQL(w io.Writer) {
	if o.Value == "" {
		return
	}

	io.WriteString(w, "offset ")
	io.WriteString(w, o.Value)
}

func writeExprList(w io.Writer, exprList []string) {
	for i, e := range exprList {
		if i > 0 {
			io.WriteString(w, ", ")
		}
		io.WriteString(w, e)
	}
}
