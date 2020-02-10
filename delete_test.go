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

func TestDeleteStatementWheref(t *testing.T) {
	a := pgsql.Delete("people")
	a.Wheref("foo=?", 43)
	sql, args := pgsql.Build(a)
	assert.Equal(t, `delete from people where foo=$1`, sql)
	assert.Equal(t, []interface{}{43}, args)
}
