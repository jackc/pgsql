package pgsql_test

import (
	"strings"
	"testing"

	"github.com/jackc/pgsql"
	"github.com/stretchr/testify/assert"
)

func TestSelectStatement(t *testing.T) {
	s := &pgsql.SelectStatement{}

	s.Select.Add(pgsql.Raw("id"))
	assert.Equal(t, "select id", strings.TrimSpace(s.String()))

	s.From.Value = pgsql.Raw("users")
	assert.Equal(t, `select id from users`, strings.TrimSpace(s.String()))

	s.Where.And(pgsql.Raw("id=42"))
	assert.Equal(t, `select id from users where id=42`, strings.TrimSpace(s.String()))

	s.OrderBy.Add(pgsql.Raw("name"))
	assert.Equal(t, `select id from users where id=42 order by name`, strings.TrimSpace(s.String()))

	s.Limit.Expr = pgsql.Raw("10")
	assert.Equal(t, `select id from users where id=42 order by name limit 10`, strings.TrimSpace(s.String()))

	s.Offset.Expr = pgsql.Raw("5")
	assert.Equal(t, `select id from users where id=42 order by name limit 10 offset 5`, strings.TrimSpace(s.String()))
}

func TestSelectClause(t *testing.T) {
	s := &pgsql.SelectClause{}
	assert.Equal(t, "select *", s.String(), "empty")

	s.Distinct = true
	assert.Equal(t, "select distinct *", s.String())

	s.AddDistinctOn(pgsql.Raw("id"))
	assert.Equal(t, "select distinct on (id) *", s.String())

	s.AddDistinctOn(pgsql.Raw("name"))
	assert.Equal(t, "select distinct on (id, name) *", s.String())

	s.Add(pgsql.Raw("id"))
	assert.Equal(t, "select distinct on (id, name) id", s.String())

	s.Add(pgsql.Raw("name"))
	assert.Equal(t, "select distinct on (id, name) id, name", s.String())
}

func TestFromClause(t *testing.T) {
	f := &pgsql.FromClause{}
	assert.Equal(t, "", f.String(), "empty")

	f.Value = pgsql.Raw("users")
	assert.Equal(t, "from users", f.String())
}

func TestWhereClause(t *testing.T) {
	wc := &pgsql.WhereClause{}
	assert.Equal(t, "", wc.String(), "empty")

	wc.And(pgsql.Raw("true"))
	assert.Equal(t, "where true", wc.String())

	wc.And(pgsql.Raw("1=1"))
	assert.Equal(t, "where (true and 1=1)", wc.String())

	wc.Or(pgsql.Raw("1+1=2"))
	assert.Equal(t, "where ((true and 1=1) or 1+1=2)", wc.String())
}

func TestWhereAndArgs(t *testing.T) {
	args := &pgsql.Args{}

	w := &pgsql.WhereClause{}
	w.And(args.SQL("id=?", 42))
	w.Or(args.SQL("id=?", 43))

	assert.Equal(t, "where (id=$1 or id=$2)", w.String())
	assert.Equal(t, []interface{}{42, 43}, args.Values())
}

func TestOrderByClause(t *testing.T) {
	ob := &pgsql.OrderByClause{}
	assert.Equal(t, "", ob.String(), "empty")

	ob.Add(pgsql.Raw("1 asc"))
	assert.Equal(t, "order by 1 asc", ob.String())

	ob.Add(pgsql.Raw("2 desc"))
	assert.Equal(t, "order by 1 asc, 2 desc", ob.String())
}

func TestLimitClause(t *testing.T) {
	l := &pgsql.LimitClause{}
	assert.Equal(t, "", l.String(), "empty")

	l.Expr = pgsql.Raw("10")
	assert.Equal(t, "limit 10", l.String())
}

func TestOffsetClause(t *testing.T) {
	o := &pgsql.OffsetClause{}
	assert.Equal(t, "", o.String(), "empty")

	o.Expr = pgsql.Raw("10")
	assert.Equal(t, "offset 10", o.String())
}
