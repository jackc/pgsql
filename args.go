package pgsql

import (
	"strconv"
	"strings"
)

const argsCountSwitchToMap = 32

type Placeholder int32

func (p Placeholder) String() string {
	return "$" + strconv.FormatInt(int64(p), 10)
}

type Args struct {
	values []interface{}
}

func (a *Args) Use(v interface{}) Placeholder {
	if len(a.values) == 0 {
		a.values = make([]interface{}, 0, 8)
	}

	a.values = append(a.values, v)
	p := Placeholder(len(a.values))

	return p
}

func (a *Args) Values() []interface{} {
	return a.values
}

func (a *Args) Format(s string, values ...interface{}) string {
	b := &strings.Builder{}

	for i := 0; ; i++ {
		pos := strings.IndexByte(s, '?')
		if pos == -1 {
			b.WriteString(s)
			break
		}

		b.WriteString(s[0:pos])
		b.WriteString(a.Use(values[i]).String())
		s = s[pos+1:]
	}

	return b.String()
}

func (a *Args) Clone() *Args {
	b := &Args{}

	b.values = make([]interface{}, len(a.values))
	copy(b.values, a.values)

	return b
}
