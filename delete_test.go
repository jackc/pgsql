package pgsql_test

import (
	"testing"

	"github.com/jackc/pgsql"
	"github.com/stretchr/testify/assert"
)

func TestDeleteStatement(t *testing.T) {
	a := pgsql.Delete("people")
	sql, args := pgsql.Build(a)
	assert.Equal(t, `delete from people`, sql)
	assert.Empty(t, args)
}

func TestDeleteStatementWhere(t *testing.T) {
	a := pgsql.Delete("people")
	a.Where("foo=?", 43)
	sql, args := pgsql.Build(a)
	assert.Equal(t, `delete from people where (foo=$1)`, sql)
	assert.Equal(t, []interface{}{43}, args)
}

func TestDeleteStatementReturning(t *testing.T) {
	a := pgsql.Delete("people")
	a.Returning("id")
	sql, args := pgsql.Build(a)
	assert.Equal(t, `delete from people returning id`, sql)
	assert.Empty(t, args)
}

func TestDeleteStatementApply(t *testing.T) {
	a := pgsql.Delete("people")
	a.Where("foo=?", 43)
	sql, args := pgsql.Build(a)
	assert.Equal(t, `delete from people where (foo=$1)`, sql)
	assert.Equal(t, []interface{}{43}, args)

	a.Apply(pgsql.Where("bar=?", 7))
	sql, args = pgsql.Build(a)
	assert.Equal(t, `delete from people where (foo=$1) and (bar=$2)`, sql)
	assert.Equal(t, []interface{}{43, 7}, args)
}
