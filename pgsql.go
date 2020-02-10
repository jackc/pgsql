// Package pgsql helps build SQL queries.
package pgsql

import (
	"sort"
	"strings"
)

type SQLWriter interface {
	WriteSQL(*strings.Builder, *Args)
}

type RawSQL string

func (r RawSQL) WriteSQL(sb *strings.Builder, args *Args) {
	sb.WriteString(string(r))
}

type QueryParameter struct {
	Value interface{}
}

func (qp *QueryParameter) WriteSQL(sb *strings.Builder, args *Args) {
	sb.WriteString(args.Use(qp.Value).String())
}

func Build(ab SQLWriter) (string, []interface{}) {
	sb := &strings.Builder{}
	args := &Args{}

	ab.WriteSQL(sb, args)

	return sb.String(), args.Values()
}

type BinaryExpr struct {
	Left  SQLWriter
	Op    string
	Right SQLWriter
}

func (be *BinaryExpr) WriteSQL(sb *strings.Builder, args *Args) {
	sb.WriteByte('(')
	be.Left.WriteSQL(sb, args)
	sb.WriteByte(' ')
	sb.WriteString(be.Op)
	sb.WriteByte(' ')
	be.Right.WriteSQL(sb, args)
	sb.WriteByte(')')
}

type RowMap map[string]interface{}

func (rm RowMap) InsertData() ([]string, *ValuesStatement) {
	keys := rm.sortedKeys()

	values := make([]interface{}, len(keys))
	for i, k := range keys {
		values[i] = rm[k]
	}

	vs := Values()
	vs.Row(values...)

	return keys, vs
}

func (rm RowMap) UpdateData() []*Assignment {
	keys := rm.sortedKeys()

	assignments := make([]*Assignment, len(keys))
	for i, k := range keys {
		assignments[i] = &Assignment{Left: RawSQL(k), Right: &QueryParameter{Value: rm[k]}}
	}

	return assignments
}

func (rm RowMap) sortedKeys() []string {
	keys := make([]string, 0, len(rm))
	for k, _ := range rm {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

type FormatString struct {
	s    string
	args []interface{}
}

func (fs *FormatString) WriteSQL(sb *strings.Builder, args *Args) {
	sb.WriteString(args.Format(fs.s, fs.args...))
}
