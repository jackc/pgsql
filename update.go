package pgsql

import (
	"strings"
)

type Assignment struct {
	Left  SQLWriter
	Right SQLWriter
}

type UpdateStatement struct {
	tableName   string
	assignments []*Assignment
	where       SQLWriter
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

func (us *UpdateStatement) Where(cond SQLWriter) *UpdateStatement {
	if us.where == nil {
		us.where = cond
	} else {
		us.where = &BinaryExpr{Left: us.where, Op: "and", Right: cond}
	}
	return us
}

func (us *UpdateStatement) WriteSQL(sb *strings.Builder, args *Args) {
	sb.WriteString("update ")
	sb.WriteString(us.tableName)
	sb.WriteString("\nset ")

	for i, a := range us.assignments {
		if i > 0 {
			sb.WriteString(",\n")
		}
		a.Left.WriteSQL(sb, args)
		sb.WriteString(" = ")
		a.Right.WriteSQL(sb, args)
	}

	if us.where != nil {
		sb.WriteString("\nwhere ")
		us.where.WriteSQL(sb, args)
	}
}
