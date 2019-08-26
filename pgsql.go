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

type Raw string

func (r Raw) String() string {
	return string(r)
}

func (r Raw) writeToSQL(w io.Writer) {
	io.WriteString(w, string(r))
}

type SelectStatement struct {
	Select  SelectClause
	From    FromClause
	Where   WhereClause
	OrderBy OrderByClause
	Limit   LimitClause
	Offset  OffsetClause
}

func (s *SelectStatement) String() string {
	sb := &strings.Builder{}
	s.writeToSQL(sb)
	return sb.String()
}

func (s *SelectStatement) writeToSQL(w io.Writer) {
	s.Select.writeToSQL(w)
	io.WriteString(w, " ")
	s.From.writeToSQL(w)
	io.WriteString(w, " ")
	s.Where.writeToSQL(w)
	io.WriteString(w, " ")
	s.OrderBy.writeToSQL(w)
	io.WriteString(w, " ")
	s.Limit.writeToSQL(w)
	io.WriteString(w, " ")
	s.Offset.writeToSQL(w)
}

type SelectClause struct {
	Distinct           bool
	DistinctOnExprList []ToSQL
	ExprList           []ToSQL
}

func (s *SelectClause) AddDistinctOn(expr ToSQL) {
	s.Distinct = true
	s.DistinctOnExprList = append(s.DistinctOnExprList, expr)
}

func (s *SelectClause) Add(expr ToSQL) {
	s.ExprList = append(s.ExprList, expr)
}

func (s *SelectClause) String() string {
	sb := &strings.Builder{}
	s.writeToSQL(sb)
	return sb.String()
}

func (s *SelectClause) writeToSQL(w io.Writer) {
	io.WriteString(w, "select")
	if s.Distinct {
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

func (wc *WhereClause) And(cond ToSQL) {
	if wc.Cond == nil {
		wc.Cond = cond
	} else {
		wc.Cond = &And{Left: wc.Cond, Right: cond}
	}
}

func (wc *WhereClause) Or(cond ToSQL) {
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

func (o *OrderByClause) Add(expr ToSQL) {
	o.ExprList = append(o.ExprList, expr)
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
	Expr ToSQL
}

func (l *LimitClause) String() string {
	sb := &strings.Builder{}
	l.writeToSQL(sb)
	return sb.String()
}

func (l *LimitClause) writeToSQL(w io.Writer) {
	if l.Expr == nil {
		return
	}

	io.WriteString(w, "limit ")
	l.Expr.writeToSQL(w)
}

type OffsetClause struct {
	Expr ToSQL
}

func (o *OffsetClause) String() string {
	sb := &strings.Builder{}
	o.writeToSQL(sb)
	return sb.String()
}

func (o *OffsetClause) writeToSQL(w io.Writer) {
	if o.Expr == nil {
		return
	}

	io.WriteString(w, "offset ")
	o.Expr.writeToSQL(w)
}

func writeExprList(w io.Writer, exprList []ToSQL) {
	for i, e := range exprList {
		if i > 0 {
			io.WriteString(w, ", ")
		}
		io.WriteString(w, e.String())
	}
}
