// Package pgsql helps build SQL queries.
package pgsql

import (
	"io"
	"strings"
)

type ToSQL interface {
	String() string
	writeToSQL(io.Writer)
}

type rawSQL string

func (r rawSQL) String() string {
	return string(r)
}

func (r rawSQL) writeToSQL(w io.Writer) {
	io.WriteString(w, string(r))
}

type SelectClause struct {
	IsDistinct         bool
	DistinctOnExprList []ToSQL
	ExprList           []ToSQL
}

func (s *SelectClause) Distinct() {
	s.IsDistinct = true
}

func (s *SelectClause) DistinctOn(sql string, args *Args, values ...interface{}) {
	s.Distinct()
	if len(values) > 0 {
		sql = args.Format(sql, values...)
	}
	s.DistinctOnExprList = append(s.DistinctOnExprList, rawSQL(sql))
}

func (s *SelectClause) Select(sql string, args *Args, values ...interface{}) {
	if len(values) > 0 {
		sql = args.Format(sql, values...)
	}
	s.ExprList = append(s.ExprList, rawSQL(sql))
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
	Value ToSQL
}

func (f *FromClause) From(sql string, args *Args, values ...interface{}) {
	f.Value = rawSQL(args.Format(sql, values...))
}

func (f *FromClause) String() string {
	sb := &strings.Builder{}
	f.writeToSQL(sb)
	return sb.String()
}

func (f *FromClause) writeToSQL(w io.Writer) {
	if f.Value == nil {
		return
	}

	io.WriteString(w, "from ")
	f.Value.writeToSQL(w)
}

type WhereClause struct {
	Cond ToSQL
}

func (wc *WhereClause) String() string {
	sb := &strings.Builder{}
	wc.writeToSQL(sb)
	return sb.String()
}

func (wc *WhereClause) writeToSQL(w io.Writer) {
	if wc.Cond == nil {
		return
	}

	io.WriteString(w, "where ")
	wc.Cond.writeToSQL(w)
}

func (wc *WhereClause) Where(sql string, args *Args, values ...interface{}) {
	if len(values) > 0 {
		sql = args.Format(sql, values...)
	}

	cond := rawSQL(sql)

	if wc.Cond == nil {
		wc.Cond = cond
	} else {
		wc.Cond = &And{Left: wc.Cond, Right: cond}
	}
}

func (wc *WhereClause) Or(sql string, args *Args, values ...interface{}) {
	if len(values) > 0 {
		sql = args.Format(sql, values...)
	}

	cond := rawSQL(sql)

	if wc.Cond == nil {
		wc.Cond = cond
	} else {
		wc.Cond = &Or{Left: wc.Cond, Right: cond}
	}
}

type And struct {
	Left  ToSQL
	Right ToSQL
}

func (a *And) String() string {
	sb := &strings.Builder{}
	a.writeToSQL(sb)
	return sb.String()
}

func (a *And) writeToSQL(w io.Writer) {
	io.WriteString(w, "(")
	a.Left.writeToSQL(w)
	io.WriteString(w, " and ")
	a.Right.writeToSQL(w)
	io.WriteString(w, ")")
}

type Or struct {
	Left  ToSQL
	Right ToSQL
}

func (o *Or) String() string {
	sb := &strings.Builder{}
	o.writeToSQL(sb)
	return sb.String()
}

func (o *Or) writeToSQL(w io.Writer) {
	io.WriteString(w, "(")
	o.Left.writeToSQL(w)
	io.WriteString(w, " or ")
	o.Right.writeToSQL(w)
	io.WriteString(w, ")")
}

type OrderByClause struct {
	ExprList []ToSQL
}

func (o *OrderByClause) OrderBy(sql string, args *Args, values ...interface{}) {
	if len(values) > 0 {
		sql = args.Format(sql, values...)
	}
	o.ExprList = append(o.ExprList, rawSQL(sql))
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
	Value ToSQL
}

func (l *LimitClause) Limit(sql string, args *Args, values ...interface{}) {
	if len(values) > 0 {
		sql = args.Format(sql, values...)
	}
	l.Value = rawSQL(sql)
}

func (l *LimitClause) String() string {
	sb := &strings.Builder{}
	l.writeToSQL(sb)
	return sb.String()
}

func (l *LimitClause) writeToSQL(w io.Writer) {
	if l.Value == nil {
		return
	}

	io.WriteString(w, "limit ")
	l.Value.writeToSQL(w)
}

type OffsetClause struct {
	Value ToSQL
}

func (o *OffsetClause) Offset(sql string, args *Args, values ...interface{}) {
	if len(values) > 0 {
		sql = args.Format(sql, values...)
	}
	o.Value = rawSQL(sql)
}

func (o *OffsetClause) String() string {
	sb := &strings.Builder{}
	o.writeToSQL(sb)
	return sb.String()
}

func (o *OffsetClause) writeToSQL(w io.Writer) {
	if o.Value == nil {
		return
	}

	io.WriteString(w, "offset ")
	o.Value.writeToSQL(w)
}

func writeExprList(w io.Writer, exprList []ToSQL) {
	for i, e := range exprList {
		if i > 0 {
			io.WriteString(w, ", ")
		}
		io.WriteString(w, e.String())
	}
}
