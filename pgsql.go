// Package pgsql helps build SQL queries.
package pgsql

import (
	"fmt"
	"strings"
)

type SelectClause struct {
	IsDistinct         bool
	DistinctOnExprList string
	ExprList           string
}

func (s *SelectClause) Distinct() {
	s.IsDistinct = true
}

func (s *SelectClause) DistinctOn(sql string, args *Args, values ...interface{}) {
	s.Distinct()
	if len(values) > 0 {
		sql = args.Format(sql, values...)
	}

	if len(s.DistinctOnExprList) > 0 {
		s.DistinctOnExprList += ", " + sql
	} else {
		s.DistinctOnExprList = sql
	}
}

func (s *SelectClause) Select(sql string, args *Args, values ...interface{}) {
	if len(values) > 0 {
		sql = args.Format(sql, values...)
	}

	if len(s.ExprList) > 0 {
		s.ExprList += ", " + sql
	} else {
		s.ExprList = sql
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

func (f *FromClause) From(sql string, args *Args, values ...interface{}) {
	*f = FromClause(args.Format(sql, values...))
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

func (wc *WhereClause) Where(sql string, args *Args, values ...interface{}) {
	if len(values) > 0 {
		sql = args.Format(sql, values...)
	}

	if len(*wc) == 0 {
		*wc = WhereClause(sql)
	} else {
		*wc = WhereClause(fmt.Sprintf("(%s and %s)", string(*wc), sql))
	}
}

func (wc *WhereClause) Or(sql string, args *Args, values ...interface{}) {
	if len(values) > 0 {
		sql = args.Format(sql, values...)
	}

	if len(*wc) == 0 {
		*wc = WhereClause(sql)
	} else {
		*wc = WhereClause(fmt.Sprintf("(%s or %s)", string(*wc), sql))
	}
}

type OrderByClause string

func (o *OrderByClause) OrderBy(sql string, args *Args, values ...interface{}) {
	if len(values) > 0 {
		sql = args.Format(sql, values...)
	}

	if len(*o) > 0 {
		*o = OrderByClause(string(*o) + ", " + sql)
	} else {
		*o = OrderByClause(sql)
	}
}

func (o OrderByClause) String() string {
	if len(o) == 0 {
		return ""
	}

	return "order by " + string(o)
}

type LimitClause string

func (l *LimitClause) Limit(sql string, args *Args, values ...interface{}) {
	if len(values) > 0 {
		sql = args.Format(sql, values...)
	}
	*l = LimitClause(sql)
}

func (l LimitClause) String() string {
	if len(l) == 0 {
		return ""
	}

	return "limit " + string(l)
}

type OffsetClause string

func (o *OffsetClause) Offset(sql string, args *Args, values ...interface{}) {
	if len(values) > 0 {
		sql = args.Format(sql, values...)
	}
	*o = OffsetClause(sql)
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
