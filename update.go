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
	setf        *FormatString
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
	us.setf = nil
	return us
}

func (us *UpdateStatement) Setf(s string, args ...interface{}) *UpdateStatement {
	us.setf = &FormatString{s: s, args: args}
	us.assignments = nil
	return us
}

func (us *UpdateStatement) Where(cond SQLWriter) *UpdateStatement {
	us.where = whereAnd(us.where, cond)
	return us
}

func (us *UpdateStatement) Wheref(s string, args ...interface{}) *UpdateStatement {
	return us.Where(&FormatString{s: s, args: args})
}

func (us *UpdateStatement) WriteSQL(sb *strings.Builder, args *Args) {
	sb.WriteString("update ")
	sb.WriteString(us.tableName)
	sb.WriteString("\nset ")

	if us.setf != nil {
		us.setf.WriteSQL(sb, args)
	} else {
		for i, a := range us.assignments {
			if i > 0 {
				sb.WriteString(",\n")
			}
			a.Left.WriteSQL(sb, args)
			sb.WriteString(" = ")
			a.Right.WriteSQL(sb, args)
		}
	}

	if us.where != nil {
		sb.WriteString("\nwhere ")
		us.where.WriteSQL(sb, args)
	}
}
