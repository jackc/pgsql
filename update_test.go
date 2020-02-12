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
	assert.Equal(t, `update people set age = $1, name = $2`, sql)
	assert.Equal(t, []interface{}{30, "Alice"}, args)
}

func TestUpdateStatementSetf(t *testing.T) {
	a := pgsql.Update("people")
	a.Setf("name=?, age=?", "Alice", 30)
	sql, args := pgsql.Build(a)
	assert.Equal(t, `update people set name=$1, age=$2`, sql)
	assert.Equal(t, []interface{}{"Alice", 30}, args)
}

func TestUpdateStatementWhere(t *testing.T) {
	a := pgsql.Update("people")
	a.Set(pgsql.RowMap{"name": "Alice", "age": 30})
	a.Where("id=?", 42)
	sql, args := pgsql.Build(a)
	assert.Equal(t, `update people set age = $1, name = $2 where (id=$3)`, sql)
	assert.Equal(t, []interface{}{30, "Alice", 42}, args)

	a.Where("foo=?", 43)
	sql, args = pgsql.Build(a)
	assert.Equal(t, `update people set age = $1, name = $2 where (id=$3) and (foo=$4)`, sql)
	assert.Equal(t, []interface{}{30, "Alice", 42, 43}, args)
}

func TestUpdateStatementReturning(t *testing.T) {
	a := pgsql.Update("people")
	a.Set(pgsql.RowMap{"name": "Alice", "age": 30})
	a.Where("id=?", 42)
	a.Returning("id")
	sql, args := pgsql.Build(a)
	assert.Equal(t, `update people set age = $1, name = $2 where (id=$3) returning id`, sql)
	assert.Equal(t, []interface{}{30, "Alice", 42}, args)
}

func TestUpdateStatementMerge(t *testing.T) {
	a := pgsql.Update("people")
	a.Set(pgsql.RowMap{"name": "Alice", "age": 30})
	a.Where("id=?", 42)
	sql, args := pgsql.Build(a)
	assert.Equal(t, `update people set age = $1, name = $2 where (id=$3)`, sql)
	assert.Equal(t, []interface{}{30, "Alice", 42}, args)

	a.Apply(pgsql.Where("foo=?", 43))
	sql, args = pgsql.Build(a)
	assert.Equal(t, `update people set age = $1, name = $2 where (id=$3) and (foo=$4)`, sql)
	assert.Equal(t, []interface{}{30, "Alice", 42, 43}, args)
}
