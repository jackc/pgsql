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

func (ds *DeleteStatement) Where(s string, args ...interface{}) *DeleteStatement {
	ds.where = whereAnd(ds.where, &FormatString{s: s, args: args})
	return ds
}

func (ds *DeleteStatement) WriteSQL(sb *strings.Builder, args *Args) {
	sb.WriteString("delete from ")
	sb.WriteString(ds.tableName)

	if ds.where != nil {
		sb.WriteString(" where ")
		ds.where.WriteSQL(sb, args)
	}
}
