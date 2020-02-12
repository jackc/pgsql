package pgsql

import (
	"strings"
)

type DeleteStatement struct {
	tableName string
	whereList whereList
}

func Delete(tableName string) *DeleteStatement {
	return &DeleteStatement{tableName: tableName}
}

func (ds *DeleteStatement) Where(s string, args ...interface{}) *DeleteStatement {
	ds.whereList = append(ds.whereList, &FormatString{s: s, args: args})
	return ds
}

func (ds *DeleteStatement) WriteSQL(sb *strings.Builder, args *Args) {
	sb.WriteString("delete from ")
	sb.WriteString(ds.tableName)
	ds.whereList.WriteSQL(sb, args)
}

func (ds *DeleteStatement) Apply(others ...*SelectStatement) *DeleteStatement {
	for _, other := range others {
		ds.whereList = append(ds.whereList, other.whereList...)
	}

	return ds
}
