// Package pgsql helps build SQL queries.
package pgsql

import (
	"fmt"
	"io"
	"strings"
)

type SelectStatement struct {
	Args          *Args
	SelectClause  *SelectClause
	FromClause    *FromClause
	WhereClause   *WhereClause
	OrderByClause *OrderByClause
	LimitClause   *LimitClause
	OffsetClause  *OffsetClause
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

	if s.SelectClause != nil {
		f(s.SelectClause)
	}

	if s.FromClause != nil {
		f(s.FromClause)
	}

	if s.WhereClause != nil {
		f(s.WhereClause)
	}

	if s.OrderByClause != nil {
		f(s.OrderByClause)
	}

	if s.LimitClause != nil {
		f(s.LimitClause)
	}

	if s.OffsetClause != nil {
		f(s.OffsetClause)
	}
}

func (s *SelectStatement) ensureArgs() {
	if s.Args == nil {
		s.Args = &Args{}
	}
}

func (s *SelectStatement) ensureSelectClause() {
	if s.SelectClause == nil {
		s.SelectClause = &SelectClause{}
	}
}

func (s *SelectStatement) Distinct() *SelectStatement {
	s.ensureSelectClause()
	s.SelectClause.Distinct()
	return s
}

func (s *SelectStatement) DistinctOn(sql string, values ...interface{}) *SelectStatement {
	s.ensureSelectClause()
	s.ensureArgs()
	s.SelectClause.DistinctOn(sql, s.Args, values...)
	return s
}

func (s *SelectStatement) Select(sql string, values ...interface{}) *SelectStatement {
	s.ensureSelectClause()
	s.ensureArgs()
	s.SelectClause.Select(sql, s.Args, values...)
	return s
}

func (s *SelectStatement) ensureFromClause() {
	if s.FromClause == nil {
		s.FromClause = &FromClause{}
	}
}

func (s *SelectStatement) From(sql string, values ...interface{}) *SelectStatement {
	s.ensureFromClause()
	s.ensureArgs()
	s.FromClause.From(sql, s.Args, values...)
	return s
}

func (s *SelectStatement) ensureWhereClause() {
	if s.WhereClause == nil {
		s.WhereClause = &WhereClause{}
	}
}

func (s *SelectStatement) Where(sql string, values ...interface{}) *SelectStatement {
	s.ensureWhereClause()
	s.ensureArgs()
	s.WhereClause.Where(sql, s.Args, values...)
	return s
}

func (s *SelectStatement) WhereOr(sql string, values ...interface{}) *SelectStatement {
	s.ensureWhereClause()
	s.ensureArgs()
	s.WhereClause.Or(sql, s.Args, values...)
	return s
}

func (s *SelectStatement) ensureOrderByClause() {
	if s.OrderByClause == nil {
		s.OrderByClause = &OrderByClause{}
	}
}

func (s *SelectStatement) OrderBy(sql string, values ...interface{}) *SelectStatement {
	s.ensureOrderByClause()
	s.ensureArgs()
	s.OrderByClause.OrderBy(sql, s.Args, values...)
	return s
}

func (s *SelectStatement) ensureLimitClause() {
	if s.LimitClause == nil {
		s.LimitClause = &LimitClause{}
	}
}

func (s *SelectStatement) Limit(sql string, values ...interface{}) *SelectStatement {
	s.ensureLimitClause()
	s.ensureArgs()
	s.LimitClause.Limit(sql, s.Args, values...)
	return s
}

func (s *SelectStatement) ensureOffsetClause() {
	if s.OffsetClause == nil {
		s.OffsetClause = &OffsetClause{}
	}
}

func (s *SelectStatement) Offset(sql string, values ...interface{}) *SelectStatement {
	s.ensureOffsetClause()
	s.ensureArgs()
	s.OffsetClause.Offset(sql, s.Args, values...)
	return s
}
