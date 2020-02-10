// Package pgsql helps build SQL queries.
package pgsql

import (
	"strconv"
	"strings"
)

type SelectStatement struct {
	// select clause
	isDistinct     bool
	distinctOnList []SQLWriter
	selectList     []SQLWriter

	from        SQLWriter
	whereList   whereList
	orderByList []SQLWriter
	limit       int64
	offset      int64
}

func Select(s string, args ...interface{}) *SelectStatement {
	return (&SelectStatement{}).Select(s, args...)
}

func From(s string, args ...interface{}) *SelectStatement {
	return (&SelectStatement{}).From(s, args...)
}

func Where(s string, args ...interface{}) *SelectStatement {
	return (&SelectStatement{}).Where(s, args...)
}

func (ss *SelectStatement) Select(s string, args ...interface{}) *SelectStatement {
	ss.selectList = append(ss.selectList, &FormatString{s: s, args: args})
	return ss
}

func (ss *SelectStatement) Distinct(b bool) *SelectStatement {
	ss.isDistinct = b
	if !b {
		ss.distinctOnList = nil
	}

	return ss
}

func (ss *SelectStatement) DistinctOn(s string, args ...interface{}) *SelectStatement {
	ss.isDistinct = true
	ss.distinctOnList = append(ss.distinctOnList, &FormatString{s: s, args: args})
	return ss
}

func (ss *SelectStatement) From(s string, args ...interface{}) *SelectStatement {
	ss.from = &FormatString{s: s, args: args}
	return ss
}

func (ss *SelectStatement) Where(s string, args ...interface{}) *SelectStatement {
	ss.whereList = append(ss.whereList, &FormatString{s: s, args: args})
	return ss
}

func (ss *SelectStatement) Order(s string, args ...interface{}) *SelectStatement {
	ss.orderByList = append(ss.orderByList, &FormatString{s: s, args: args})
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

// Merge merges other's from, where, order, limit and offset if they are set.
func (ss *SelectStatement) Merge(other *SelectStatement) *SelectStatement {
	if other.from != nil {
		ss.from = other.from
	}

	ss.whereList = append(ss.whereList, other.whereList...)
	ss.orderByList = append(ss.orderByList, other.orderByList...)

	if other.limit != 0 {
		ss.limit = other.limit
	}

	if other.offset != 0 {
		ss.offset = other.offset
	}

	return ss
}

func (ss *SelectStatement) WriteSQL(sb *strings.Builder, args *Args) {
	sb.WriteString("select")
	if ss.isDistinct {
		sb.WriteString(" distinct")
	}

	if len(ss.distinctOnList) > 0 {
		sb.WriteString(" on (")
		for i, e := range ss.distinctOnList {
			if i > 0 {
				sb.WriteString(", ")
			}
			e.WriteSQL(sb, args)
		}
		sb.WriteString(")")
	}

	sb.WriteString(" ")
	if len(ss.selectList) > 0 {
		for i, e := range ss.selectList {
			if i > 0 {
				sb.WriteString(", ")
			}
			e.WriteSQL(sb, args)
		}
	} else {
		sb.WriteString("*")
	}

	if ss.from != nil {
		sb.WriteString(" from ")
		ss.from.WriteSQL(sb, args)
	}

	ss.whereList.WriteSQL(sb, args)

	if len(ss.orderByList) > 0 {
		sb.WriteString(" order by ")
		for i, e := range ss.orderByList {
			if i > 0 {
				sb.WriteString(", ")
			}
			e.WriteSQL(sb, args)
		}
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
