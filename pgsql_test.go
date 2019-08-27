package pgsql_test

import (
	"testing"

	"github.com/jackc/pgsql"
	"github.com/stretchr/testify/assert"
)

func TestSelectStatement(t *testing.T) {
	a := &pgsql.SelectStatement{}

	a.Select("id")
	assert.Equal(t, "select id", a.String())

	a.From("users")
	assert.Equal(t, `select id
from users`, a.String())

	a.Where("id=?", 42)
	assert.Equal(t, `select id
from users
where id=$1`, a.String())

	a.OrderBy("name")
	assert.Equal(t, `select id
from users
where id=$1
order by name`, a.String())

	a.Limit("10")
	assert.Equal(t, `select id
from users
where id=$1
order by name
limit 10`, a.String())

	a.Offset("5")
	assert.Equal(t, `select id
from users
where id=$1
order by name
limit 10
offset 5`, a.String())

	b := a.Clone()
	b.Where("name=?", "foo")
	assert.Equal(t, `select id
from users
where (id=$1 and name=$2)
order by name
limit 10
offset 5`, b.String())

	assert.Equal(t, `select id
from users
where id=$1
order by name
limit 10
offset 5`, a.String())

	a.Where("nickname=?", "bar")
	assert.Equal(t, `select id
from users
where (id=$1 and nickname=$2)
order by name
limit 10
offset 5`, a.String())

}

func TestSelectClause(t *testing.T) {
	s := &pgsql.SelectClause{}
	assert.Equal(t, "select *", s.String(), "empty")

	s.Distinct()
	assert.Equal(t, "select distinct *", s.String())

	s.DistinctOn("id", nil)
	assert.Equal(t, "select distinct on (id) *", s.String())

	s.DistinctOn("name", nil)
	assert.Equal(t, "select distinct on (id, name) *", s.String())

	s.Select("id", nil)
	assert.Equal(t, "select distinct on (id, name) id", s.String())

	s.Select("name", nil)
	assert.Equal(t, "select distinct on (id, name) id, name", s.String())
}

func TestFromClause(t *testing.T) {
	var f pgsql.FromClause
	assert.Equal(t, "", f.String(), "empty")

	f.From("users", nil)
	assert.Equal(t, "from users", f.String())
}

func TestWhereClause(t *testing.T) {
	var wc pgsql.WhereClause
	assert.Equal(t, "", wc.String(), "empty")

	wc.Where("true", nil)
	assert.Equal(t, "where true", wc.String())

	wc.Where("1=1", nil)
	assert.Equal(t, "where (true and 1=1)", wc.String())

	wc.Or("1+1=2", nil)
	assert.Equal(t, "where ((true and 1=1) or 1+1=2)", wc.String())
}

func TestWhereAndArgs(t *testing.T) {
	args := &pgsql.Args{}

	var w pgsql.WhereClause
	w.Where("id=?", args, 42)
	w.Or("id=?", args, 43)

	assert.Equal(t, "where (id=$1 or id=$2)", w.String())
	assert.Equal(t, []interface{}{42, 43}, args.Values())
}

func TestOrderByClause(t *testing.T) {
	var ob pgsql.OrderByClause
	assert.Equal(t, "", ob.String(), "empty")

	ob.OrderBy("1 asc", nil)
	assert.Equal(t, "order by 1 asc", ob.String())

	ob.OrderBy("2 desc", nil)
	assert.Equal(t, "order by 1 asc, 2 desc", ob.String())
}

func TestLimitClause(t *testing.T) {
	var l pgsql.LimitClause
	assert.Equal(t, "", l.String(), "empty")

	l.Limit("10", nil)
	assert.Equal(t, "limit 10", l.String())
}

func TestOffsetClause(t *testing.T) {
	var o pgsql.OffsetClause
	assert.Equal(t, "", o.String(), "empty")

	o.Offset("10", nil)
	assert.Equal(t, "offset 10", o.String())
}
