package pgsql_test

import (
	"testing"

	"github.com/jackc/pgsql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelectStatement(t *testing.T) {
	a := pgsql.NewStatement()

	err := a.Apply(pgsql.Select("id"))
	require.NoError(t, err)
	assert.Equal(t, "select id", a.String())

	err = a.Apply(pgsql.From("users"))
	require.NoError(t, err)
	assert.Equal(t, `select id
from users`, a.String())

	err = a.Apply(pgsql.Where("id=?", 42))
	require.NoError(t, err)
	assert.Equal(t, `select id
from users
where id=$1`, a.String())

	err = a.Apply(pgsql.OrderBy("name"))
	require.NoError(t, err)
	assert.Equal(t, `select id
from users
where id=$1
order by name`, a.String())

	err = a.Apply(pgsql.Limit("10"))
	require.NoError(t, err)
	assert.Equal(t, `select id
from users
where id=$1
order by name
limit 10`, a.String())

	err = a.Apply(pgsql.Offset("5"))
	require.NoError(t, err)
	assert.Equal(t, `select id
from users
where id=$1
order by name
limit 10
offset 5`, a.String())

	b := a.Clone()
	err = b.Apply(pgsql.Where("name=?", "foo"))
	require.NoError(t, err)
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

	err = a.Apply(pgsql.Where("nickname=?", "bar"))
	require.NoError(t, err)
	assert.Equal(t, `select id
from users
where (id=$1 and nickname=$2)
order by name
limit 10
offset 5`, a.String())

}

func TestSelectClause(t *testing.T) {
	s := pgsql.NewStatement()

	assert.Equal(t, "select *", s.String(), "empty")

	err := s.Apply(pgsql.Distinct())
	require.NoError(t, err)
	assert.Equal(t, "select distinct *", s.String())

	err = s.Apply(pgsql.DistinctOn("id"))
	require.NoError(t, err)
	assert.Equal(t, "select distinct on (id) *", s.String())

	err = s.Apply(pgsql.DistinctOn("name"))
	require.NoError(t, err)
	assert.Equal(t, "select distinct on (id, name) *", s.String())

	err = s.Apply(pgsql.Select("id"))
	require.NoError(t, err)
	assert.Equal(t, "select distinct on (id, name) id", s.String())

	err = s.Apply(pgsql.Select("name"))
	require.NoError(t, err)
	assert.Equal(t, "select distinct on (id, name) id, name", s.String())
}

func TestFromClause(t *testing.T) {
	s := pgsql.NewStatement()

	err := s.Apply(pgsql.From("users"))
	require.NoError(t, err)
	assert.Equal(t, "select *\nfrom users", s.String())
}

func TestWhereClause(t *testing.T) {
	s := pgsql.NewStatement()

	err := s.Apply(pgsql.Where("true"))
	require.NoError(t, err)
	assert.Equal(t, "select *\nwhere true", s.String())

	err = s.Apply(pgsql.Where("1=1"))
	require.NoError(t, err)
	assert.Equal(t, "select *\nwhere (true and 1=1)", s.String())

	err = s.Apply(pgsql.Or("1+1=2"))
	require.NoError(t, err)
	assert.Equal(t, "select *\nwhere ((true and 1=1) or 1+1=2)", s.String())
}

func TestOrderByClause(t *testing.T) {
	s := pgsql.NewStatement()

	err := s.Apply(pgsql.OrderBy("1 asc"))
	require.NoError(t, err)
	assert.Equal(t, "select *\norder by 1 asc", s.String())

	err = s.Apply(pgsql.OrderBy("2 desc"))
	require.NoError(t, err)
	assert.Equal(t, "select *\norder by 1 asc, 2 desc", s.String())
}

func TestLimitClause(t *testing.T) {
	s := pgsql.NewStatement()

	err := s.Apply(pgsql.Limit("10"))
	require.NoError(t, err)
	assert.Equal(t, "select *\nlimit 10", s.String())
}

func TestOffsetClause(t *testing.T) {
	s := pgsql.NewStatement()

	err := s.Apply(pgsql.Offset("10"))
	require.NoError(t, err)
	assert.Equal(t, "select *\noffset 10", s.String())
}
