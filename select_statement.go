// Package pgsql helps build SQL queries.
package pgsql

import (
	"fmt"
	"io"
	"strings"
)

type SelectStatement struct {
	Args          Args
	SelectClause  SelectClause
	FromClause    FromClause
	WhereClause   WhereClause
	OrderByClause OrderByClause
	LimitClause   LimitClause
	OffsetClause  OffsetClause
}

func (s *SelectStatement) String() string {
	sb := &strings.Builder{}
	s.writeToSQL(sb)
	return sb.String()
}

func (s *SelectStatement) writeToSQL(w io.Writer) {
	writeCount := 0

	f := func(clause fmt.Stringer) {
		if writeCount > 0 {
			io.WriteString(w, "\n")
		}
		io.WriteString(w, clause.String())
		writeCount += 1
	}

	f(s.SelectClause)

	if len(s.FromClause) != 0 {
		f(s.FromClause)
	}

	if len(s.WhereClause) != 0 {
		f(s.WhereClause)
	}

	if len(s.OrderByClause) != 0 {
		f(s.OrderByClause)
	}

	if len(s.LimitClause) != 0 {
		f(s.LimitClause)
	}

	if len(s.OffsetClause) != 0 {
		f(s.OffsetClause)
	}
}

func (s *SelectStatement) Distinct() *SelectStatement {
	s.SelectClause.Distinct()
	return s
}

func (s *SelectStatement) DistinctOn(sql string, values ...interface{}) *SelectStatement {
	s.SelectClause.DistinctOn(sql, &s.Args, values...)
	return s
}

func (s *SelectStatement) Select(sql string, values ...interface{}) *SelectStatement {
	s.SelectClause.Select(sql, &s.Args, values...)
	return s
}

func (s *SelectStatement) From(sql string, values ...interface{}) *SelectStatement {
	s.FromClause.From(sql, &s.Args, values...)
	return s
}

func (s *SelectStatement) Where(sql string, values ...interface{}) *SelectStatement {
	s.WhereClause.Where(sql, &s.Args, values...)
	return s
}

func (s *SelectStatement) WhereOr(sql string, values ...interface{}) *SelectStatement {
	s.WhereClause.Or(sql, &s.Args, values...)
	return s
}

func (s *SelectStatement) OrderBy(sql string, values ...interface{}) *SelectStatement {
	s.OrderByClause.OrderBy(sql, &s.Args, values...)
	return s
}

func (s *SelectStatement) Limit(sql string, values ...interface{}) *SelectStatement {
	s.LimitClause.Limit(sql, &s.Args, values...)
	return s
}

func (s *SelectStatement) Offset(sql string, values ...interface{}) *SelectStatement {
	s.OffsetClause.Offset(sql, &s.Args, values...)
	return s
}
