package pgsql

import (
	"fmt"
	"strings"
)

const argsCountSwitchToMap = 32

type Placeholder int32

func (p Placeholder) ToSQL() string {
	return fmt.Sprintf("$%d", p)
}

type Args struct {
	values              []interface{}
	valuesToPlaceholder map[interface{}]Placeholder
}

func (a *Args) Use(v interface{}) Placeholder {
	if p, present := a.existingPlaceholder(v); present {
		return p
	}

	if len(a.values) == 0 {
		a.values = make([]interface{}, 0, 8)
	}

	// Start using a map for existingPlaceholder lookup when there are enough values to make it worth while.
	if len(a.values) == argsCountSwitchToMap {
		a.valuesToPlaceholder = make(map[interface{}]Placeholder, argsCountSwitchToMap)
		for i, v := range a.values {
			a.valuesToPlaceholder[v] = Placeholder(i + 1)
		}
	}

	a.values = append(a.values, v)
	p := Placeholder(len(a.values))

	if a.valuesToPlaceholder != nil {
		a.valuesToPlaceholder[v] = p
	}

	return p
}

func (a *Args) Values() []interface{} {
	return a.values
}

func (a *Args) SQL(s string, values ...interface{}) Raw {
	b := &strings.Builder{}

	for i := 0; ; i++ {
		pos := strings.IndexByte(s, '?')
		if pos == -1 {
			b.WriteString(s)
			break
		}

		b.WriteString(s[0:pos])
		b.WriteString(a.Use(values[i]).ToSQL())
		s = s[pos+1:]
	}

	return Raw(b.String())
}

func (a *Args) existingPlaceholder(v interface{}) (p Placeholder, present bool) {
	if a.valuesToPlaceholder != nil {
		p, present = a.valuesToPlaceholder[v]
		return p, present
	}

	for i, vv := range a.Values() {
		if v == vv {
			return Placeholder(i + 1), true
		}
	}

	return 0, false
}
