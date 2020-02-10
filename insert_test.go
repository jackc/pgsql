package pgsql_test

import (
	"testing"

	"github.com/jackc/pgsql"
	"github.com/stretchr/testify/assert"
)

func TestInsertStatement(t *testing.T) {
	a := pgsql.Insert("people")
	sql, args := pgsql.Build(a)
	assert.Equal(t, "insert into people ", sql)
	assert.Empty(t, args)

	a.Columns("name", "age")
	sql, args = pgsql.Build(a)
	assert.Equal(t, "insert into people (name, age)", sql)
	assert.Empty(t, args)

	vs := pgsql.Values()
	vs.Row("Alice", 30)
	vs.Row("Bob", 32)
	a.Values(vs)

	sql, args = pgsql.Build(a)
	assert.Equal(t, "insert into people (name, age) values ($1,$2), ($3,$4)", sql)
	assert.Equal(t, []interface{}{"Alice", 30, "Bob", 32}, args)
}

func TestInsertStatementData(t *testing.T) {
	a := pgsql.Insert("people")
	sql, args := pgsql.Build(a)
	assert.Equal(t, "insert into people ", sql)
	assert.Empty(t, args)

	a.Data(pgsql.RowMap{"name": "Alice", "age": 30})
	sql, args = pgsql.Build(a)
	assert.Equal(t, "insert into people (age, name) values ($1,$2)", sql)
	assert.Equal(t, []interface{}{30, "Alice"}, args)
}
