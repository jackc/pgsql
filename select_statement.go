// Package pgsql helps build SQL queries.
package pgsql

import (
	"strconv"
	"strings"
)

type SelectStatement struct {
	selectClause  fooselectClause
	fromClause    foofromClause
	where         SQLWriter
	orderByClause fooorderClause
	limit         int64
	offset        int64
}

func Select(s string, args ...interface{}) *SelectStatement {
	return (&SelectStatement{}).Select(s, args...)
}

func From(s string, args ...interface{}) *SelectStatement {
	return (&SelectStatement{}).From(s, args...)
}

func (ss *SelectStatement) Select(s string, args ...interface{}) *SelectStatement {
	ss.selectClause.exprList = append(ss.selectClause.exprList, &FormatString{s: s, args: args})
	return ss
}

func (ss *SelectStatement) Distinct(b bool) *SelectStatement {
	ss.selectClause.isDistinct = b
	if !b {
		ss.selectClause.distinctOnExprList = nil
	}

	return ss
}

func (ss *SelectStatement) DistinctOn(s string, args ...interface{}) *SelectStatement {
	ss.selectClause.isDistinct = true
	ss.selectClause.distinctOnExprList = append(ss.selectClause.distinctOnExprList, &FormatString{s: s, args: args})
	return ss
}

func (ss *SelectStatement) From(s string, args ...interface{}) *SelectStatement {
	ss.fromClause = foofromClause{
		from: &FormatString{s: s, args: args},
	}

	return ss
}

func (ss *SelectStatement) Where(s string, args ...interface{}) *SelectStatement {
	ss.where = whereAnd(ss.where, &FormatString{s: s, args: args})
	return ss
}

func (ss *SelectStatement) Order(s string, args ...interface{}) *SelectStatement {
	ss.orderByClause.exprList = append(ss.orderByClause.exprList, &FormatString{s: s, args: args})
	return ss
}

func (ss *SelectStatement) Limit(n int64) *SelectStatement {
	ss.limit = n
	return ss
}

func (ss *SelectStatement) Offset(n int64) *SelectStatement {
	ss.offset = n
	return ss
}

func (ss *SelectStatement) WriteSQL(sb *strings.Builder, args *Args) {
	ss.selectClause.WriteSQL(sb, args)
	ss.fromClause.WriteSQL(sb, args)
	if ss.where != nil {
		sb.WriteString(" where ")
		ss.where.WriteSQL(sb, args)
	}
	ss.orderByClause.WriteSQL(sb, args)
	if ss.limit != 0 {
		sb.WriteString(" limit ")
		sb.WriteString(strconv.FormatInt(ss.limit, 10))
	}
	if ss.offset != 0 {
		sb.WriteString(" offset ")
		sb.WriteString(strconv.FormatInt(ss.offset, 10))
	}
}

type fooselectClause struct {
	isDistinct         bool
	distinctOnExprList []SQLWriter
	exprList           []SQLWriter
}

func (s fooselectClause) WriteSQL(sb *strings.Builder, args *Args) {
	sb.WriteString("select")
	if s.isDistinct {
		sb.WriteString(" distinct")
	}

	if len(s.distinctOnExprList) > 0 {
		sb.WriteString(" on (")
		for i, e := range s.distinctOnExprList {
			if i > 0 {
				sb.WriteString(", ")
			}
			e.WriteSQL(sb, args)
		}
		sb.WriteString(")")
	}

	sb.WriteString(" ")
	if len(s.exprList) > 0 {
		for i, e := range s.exprList {
			if i > 0 {
				sb.WriteString(", ")
			}
			e.WriteSQL(sb, args)
		}
	} else {
		sb.WriteString("*")
	}
}

type foofromClause struct {
	from SQLWriter
}

func (f foofromClause) WriteSQL(sb *strings.Builder, args *Args) {
	if f.from == nil {
		return
	}

	sb.WriteString(" from ")
	f.from.WriteSQL(sb, args)
}

type fooorderClause struct {
	exprList []SQLWriter
}

func (o fooorderClause) WriteSQL(sb *strings.Builder, args *Args) {
	if len(o.exprList) == 0 {
		return
	}

	sb.WriteString(" order by ")
	for i, e := range o.exprList {
		if i > 0 {
			sb.WriteString(", ")
		}
		e.WriteSQL(sb, args)
	}
}
