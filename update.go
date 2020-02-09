package pgsql

import (
	"strings"
)

type Assignment struct {
	Left  AppendBuilder
	Right AppendBuilder
}

type UpdateStatement struct {
	tableName   string
	assignments []*Assignment
}

func Update(tableName string) *UpdateStatement {
	return &UpdateStatement{tableName: tableName}
}

type Updateable interface {
	UpdateData() []*Assignment
}

func (us *UpdateStatement) Set(data Updateable) *UpdateStatement {
	us.assignments = data.UpdateData()
	return us
}

func (us *UpdateStatement) AppendBuild(sb *strings.Builder, args *Args) {
	sb.WriteString("update ")
	sb.WriteString(us.tableName)
	sb.WriteString("\nset ")

	for i, a := range us.assignments {
		if i > 0 {
			sb.WriteString(",\n")
		}
		a.Left.AppendBuild(sb, args)
		sb.WriteString(" = ")
		a.Right.AppendBuild(sb, args)
	}

	sb.WriteByte('\n')
}
