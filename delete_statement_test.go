package pgsql_test

import (
	"testing"

	"github.com/jackc/pgsql"
	"github.com/stretchr/testify/assert"
)

func TestDeleteStatement(t *testing.T) {
	a := pgsql.Delete("users")

	assert.Equal(t, `delete from users`, a.String())

	a.From("people")
	assert.Equal(t, `delete from people`, a.String())

	a.Where("id=?", 42)
	assert.Equal(t, `delete from people
where id=$1`, a.String())

	b := a.Clone()
	b.Where("name=?", "foo")
	assert.Equal(t, `delete from people
where (id=$1 and name=$2)`, b.String())

	assert.Equal(t, `delete from people
where id=$1`, a.String())

	a.Where("nickname=?", "bar")
	assert.Equal(t, `delete from people
where (id=$1 and nickname=$2)`, a.String())

}
