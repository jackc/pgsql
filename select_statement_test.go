package pgsql_test

import (
	"testing"

	"github.com/jackc/pgsql"
	"github.com/stretchr/testify/assert"
)

func TestSelect(t *testing.T) {
	a := pgsql.Select("a, b, c")
	sql, args := pgsql.Build(a)
	assert.Equal(t, "select a, b, c", sql)
	assert.Empty(t, args)
}

func TestFrom(t *testing.T) {
	a := pgsql.From("people")
	sql, args := pgsql.Build(a)
	assert.Equal(t, "select * from people", sql)
	assert.Empty(t, args)
}

func TestSelectStatementDistinct(t *testing.T) {
	a := pgsql.Select("a, b, c").Distinct(true)
	sql, args := pgsql.Build(a)
	assert.Equal(t, "select distinct a, b, c", sql)
	assert.Empty(t, args)
}

func TestSelectStatementDistinctOn(t *testing.T) {
	a := pgsql.Select("a, b, c").DistinctOn("a")
	sql, args := pgsql.Build(a)
	assert.Equal(t, "select distinct on (a) a, b, c", sql)
	assert.Empty(t, args)

	a.DistinctOn("b")
	sql, args = pgsql.Build(a)
	assert.Equal(t, "select distinct on (a, b) a, b, c", sql)
	assert.Empty(t, args)
}

func TestSelectStatementMultipleSelect(t *testing.T) {
	a := pgsql.Select("a").Select("b")
	sql, args := pgsql.Build(a)
	assert.Equal(t, "select a, b", sql)
	assert.Empty(t, args)
}

func TestSelectStatementWhere(t *testing.T) {
	a := pgsql.Select("a, b, c").From("t").Where("foo=?", 42)
	sql, args := pgsql.Build(a)
	assert.Equal(t, "select a, b, c from t where foo=$1", sql)
	assert.Equal(t, []interface{}{42}, args)
}

func TestSelectStatementOrder(t *testing.T) {
	a := pgsql.Select("a, b, c").From("t").Order("c desc")
	sql, args := pgsql.Build(a)
	assert.Equal(t, "select a, b, c from t order by c desc", sql)
	assert.Empty(t, args)
}

func TestSelectStatementLimitAndOffset(t *testing.T) {
	a := pgsql.Select("a, b, c").From("t").Order("c desc")
	a.Limit(5)
	sql, args := pgsql.Build(a)
	assert.Equal(t, "select a, b, c from t order by c desc limit 5", sql)
	assert.Empty(t, args)

	a.Offset(10)
	sql, args = pgsql.Build(a)
	assert.Equal(t, "select a, b, c from t order by c desc limit 5 offset 10", sql)
	assert.Empty(t, args)
}
