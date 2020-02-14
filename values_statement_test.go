package pgsql_test

import (
	"testing"

	"github.com/jackc/pgsql"
	"github.com/stretchr/testify/assert"
)

func TestValuesStatement(t *testing.T) {
	v := pgsql.Values()
	v.Row("a", "b")

	sql, args := pgsql.Build(v)
	assert.Equal(t, "values ($1,$2)", sql)
	assert.Equal(t, []interface{}{"a", "b"}, args)
}

func TestValuesStatementMultipleRows(t *testing.T) {
	v := pgsql.Values()
	v.Row("a", "b")
	v.Row("c", "d")

	sql, args := pgsql.Build(v)
	assert.Equal(t, "values ($1,$2), ($3,$4)", sql)
	assert.Equal(t, []interface{}{"a", "b", "c", "d"}, args)
}
