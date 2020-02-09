package pgsql

import (
	"strings"
)

type InsertStatement struct {
	tableName string
	columns   []string
	values    AppendBuilder
}

func Insert(tableName string) *InsertStatement {
	return &InsertStatement{tableName: tableName}
}

type Insertable interface {
	InsertData() ([]string, *ValuesStatement)
}

func (is *InsertStatement) Data(data Insertable) *InsertStatement {
	columns, values := data.InsertData()
	is.Columns(columns...)
	is.Values(values)
	return is
}

func (is *InsertStatement) Columns(columns ...string) *InsertStatement {
	is.columns = columns
	return is
}

func (is *InsertStatement) Values(vs *ValuesStatement) *InsertStatement {
	is.values = vs
	return is
}

func (is *InsertStatement) AppendBuild(sb *strings.Builder, args *Args) {
	sb.WriteString("insert into ")
	sb.WriteString(is.tableName)
	sb.WriteByte(' ')

	if len(is.columns) > 0 {
		sb.WriteByte('(')
		for i, c := range is.columns {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(c)
		}
		sb.WriteByte(')')
	}

	if is.values != nil {
		sb.WriteByte('\n')
		is.values.AppendBuild(sb, args)
	}
}
