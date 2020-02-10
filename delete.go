package pgsql

import (
	"strings"
)

type DeleteStatement struct {
	tableName string
	where     SQLWriter
}

func Delete(tableName string) *DeleteStatement {
	return &DeleteStatement{tableName: tableName}
}

func (ds *DeleteStatement) Where(cond SQLWriter) *DeleteStatement {
	if ds.where == nil {
		ds.where = cond
	} else {
		ds.where = &binaryExpr{left: ds.where, op: "and", right: cond}
	}
	return ds
}

func (ds *DeleteStatement) Wheref(s string, args ...interface{}) *DeleteStatement {
	return ds.Where(&FormatString{s: s, args: args})
}

func (ds *DeleteStatement) WriteSQL(sb *strings.Builder, args *Args) {
	sb.WriteString("delete from ")
	sb.WriteString(ds.tableName)

	if ds.where != nil {
		sb.WriteString(" where ")
		ds.where.WriteSQL(sb, args)
	}
}
