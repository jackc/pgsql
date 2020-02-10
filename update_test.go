package pgsql_test

import (
	"testing"

	"github.com/jackc/pgsql"
	"github.com/stretchr/testify/assert"
)

func TestUpdateStatement(t *testing.T) {
	a := pgsql.Update("people")
	a.Set(pgsql.RowMap{"name": "Alice", "age": 30})
	sql, args := pgsql.Build(a)
	assert.Equal(t, `update people
set age = $1,
name = $2`, sql)
	assert.Equal(t, []interface{}{30, "Alice"}, args)
}

func TestUpdateStatementWhere(t *testing.T) {
	a := pgsql.Update("people")
	a.Set(pgsql.RowMap{"name": "Alice", "age": 30})
	a.Where(&pgsql.BinaryExpr{Left: pgsql.RawSQL("id"), Op: "=", Right: &pgsql.QueryParameter{Value: 42}})
	sql, args := pgsql.Build(a)
	assert.Equal(t, `update people
set age = $1,
name = $2
where (id = $3)`, sql)
	assert.Equal(t, []interface{}{30, "Alice", 42}, args)

	a.Wheref("foo=?", 43)
	sql, args = pgsql.Build(a)
	assert.Equal(t, `update people
set age = $1,
name = $2
where ((id = $3) and foo=$4)`, sql)
	assert.Equal(t, []interface{}{30, "Alice", 42, 43}, args)
}
