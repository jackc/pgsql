package pgsql

import (
	"strings"
)

type ValuesStatement struct {
	rows [][]AppendBuilder
}

func Values() *ValuesStatement {
	return &ValuesStatement{}
}

func (vs *ValuesStatement) Row(values ...interface{}) *ValuesStatement {
	row := make([]AppendBuilder, len(values))
	for i := range values {
		in := values[i]
		out, ok := in.(AppendBuilder)
		if !ok {
			out = &QueryParameter{Value: in}
		}
		row[i] = out
	}

	vs.rows = append(vs.rows, row)

	return vs
}

func (vs *ValuesStatement) AppendBuild(sb *strings.Builder, args *Args) {
	sb.WriteString("values ")
	for i, row := range vs.rows {
		if i > 0 {
			sb.WriteString(",\n")
		}
		sb.WriteByte('(')
		for j, v := range row {
			if j > 0 {
				sb.WriteByte(',')
			}
			v.AppendBuild(sb, args)
		}
		sb.WriteByte(')')
	}
}
