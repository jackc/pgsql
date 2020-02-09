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
name = $2
`, sql)
	assert.Equal(t, []interface{}{30, "Alice"}, args)
}
