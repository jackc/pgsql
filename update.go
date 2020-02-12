package pgsql

import (
	"strings"
)

type Assignment struct {
	Left  SQLWriter
	Right SQLWriter
}

type Assignments []*Assignment

func (a Assignments) UpdateData() []*Assignment {
	return []*Assignment(a)
}

type UpdateStatement struct {
	tableName     string
	setf          *FormatString
	assignments   []*Assignment
	whereList     whereList
	returningList returningList
}

func Update(tableName string) *UpdateStatement {
	return &UpdateStatement{tableName: tableName}
}

func (us *UpdateStatement) UpdateStatement() (*UpdateStatement, error) {
	return us, nil
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

func (us *UpdateStatement) Where(s string, args ...interface{}) *UpdateStatement {
	us.whereList = append(us.whereList, &FormatString{s: s, args: args})
	return us
}

func (us *UpdateStatement) Returning(s string, args ...interface{}) *UpdateStatement {
	us.returningList = append(us.returningList, &FormatString{s: s, args: args})
	return us
}

func (us *UpdateStatement) WriteSQL(sb *strings.Builder, args *Args) {
	sb.WriteString("update ")
	sb.WriteString(us.tableName)
	sb.WriteString(" set ")

	if us.setf != nil {
		us.setf.WriteSQL(sb, args)
	} else {
		for i, a := range us.assignments {
			if i > 0 {
				sb.WriteString(", ")
			}
			a.Left.WriteSQL(sb, args)
			sb.WriteString(" = ")
			a.Right.WriteSQL(sb, args)
		}
	}

	us.whereList.WriteSQL(sb, args)
	us.returningList.WriteSQL(sb, args)
}

func (us *UpdateStatement) Apply(others ...*SelectStatement) *UpdateStatement {
	for _, other := range others {
		us.whereList = append(us.whereList, other.whereList...)
	}

	return us
}
