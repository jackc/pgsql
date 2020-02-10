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
	orderByClause *fooorderClause
	limit         int64
	offset        int64
}

func Selectf(s string, args ...interface{}) *SelectStatement {
	return (&SelectStatement{}).Selectf(s, args...)
}

func Fromf(s string, args ...interface{}) *SelectStatement {
	return (&SelectStatement{}).Fromf(s, args...)
}

func (ss *SelectStatement) Selectf(s string, args ...interface{}) *SelectStatement {
	ss.selectClause = fooselectClause{
		exprList: []SQLWriter{&FormatString{s: s, args: args}},
	}

	return ss
}

func (ss *SelectStatement) Distinct(b bool) *SelectStatement {
	ss.selectClause.isDistinct = b
	if !b {
		ss.selectClause.distinctOnExprList = nil
	}

	return ss
}

func (ss *SelectStatement) DistinctOnf(s string, args ...interface{}) *SelectStatement {
	ss.selectClause.isDistinct = true
	ss.selectClause.distinctOnExprList = []SQLWriter{&FormatString{s: s, args: args}}

	return ss
}

func (ss *SelectStatement) Fromf(s string, args ...interface{}) *SelectStatement {
	ss.fromClause = foofromClause{
		from: &FormatString{s: s, args: args},
	}

	return ss
}

func (ss *SelectStatement) Where(cond SQLWriter) *SelectStatement {
	if ss.where == nil {
		ss.where = cond
	} else {
		ss.where = &binaryExpr{left: ss.where, op: "and", right: cond}
	}
	return ss
}

func (ss *SelectStatement) Wheref(s string, args ...interface{}) *SelectStatement {
	return ss.Where(&FormatString{s: s, args: args})
}

func (ss *SelectStatement) Order(order SQLWriter) *SelectStatement {
	ss.orderByClause = &fooorderClause{exprList: []SQLWriter{order}}
	return ss
}

func (ss *SelectStatement) Orderf(s string, args ...interface{}) *SelectStatement {
	return ss.Order(&FormatString{s: s, args: args})
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
	if ss.orderByClause != nil {
		ss.orderByClause.WriteSQL(sb, args)
	}
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
	sb.WriteString(" order by ")
	for i, e := range o.exprList {
		if i > 0 {
			sb.WriteString(", ")
		}
		e.WriteSQL(sb, args)
	}
}
