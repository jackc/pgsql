package pgsql

import (
	"strings"
)

func Delete(sql string, values ...interface{}) *DeleteStatement {
	s := &DeleteStatement{Args: &Args{}}
	return s.From(sql, values...)
}

type DeleteStatement struct {
	FromClause  FromClause
	WhereClause WhereClause
	Args        *Args
}

func (s *DeleteStatement) String() string {
	sb := &strings.Builder{}
	sb.WriteString("delete ")

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

	f(s.FromClause.String())
	f(s.WhereClause.String())

	return sb.String()
}

func (s *DeleteStatement) Clone() *DeleteStatement {
	clone := *s
	clone.Args = s.Args.Clone()
	return &clone
}

func (s *DeleteStatement) From(sql string, values ...interface{}) *DeleteStatement {
	s.FromClause.From(sql, s.Args, values...)
	return s
}

func (s *DeleteStatement) Where(sql string, values ...interface{}) *DeleteStatement {
	s.WhereClause.Where(sql, s.Args, values...)
	return s
}

func (s *DeleteStatement) WhereOr(sql string, values ...interface{}) *DeleteStatement {
	s.WhereClause.Or(sql, s.Args, values...)
	return s
}
