package lib

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableMarshalIter_q(t *testing.T) {
	assert := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	p2 := Person{FirstName: "Tom", LastName: "Smith", Age: 99}
	it := &MockIter{eofPage: true, People: []Person{p1, p2}}
	out := strings.Builder{}
	in := strings.NewReader("q")
	TableMarshalIter("", it, "", &out, in, nil)

	assert.Equal(`
┌────────────┬───────────┬─────┐
│ FIRST_NAME │ LAST_NAME │ AGE │
├────────────┼───────────┼─────┤
│ Dustin     │ Zeisler   │ 100 │
└────────────┴───────────┴─────┘
:
`, "\n"+out.String()+"\n")
}

func TestTableMarshalIter_newline(t *testing.T) {
	assert := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	p2 := Person{FirstName: "Tom", LastName: "Smith", Age: 99}
	it := &MockIter{eofPage: true, People: []Person{p1, p2}}
	out := strings.Builder{}
	in := strings.NewReader(" \n")
	TableMarshalIter("", it, "", &out, in, nil)

	assert.Equal(`
┌────────────┬───────────┬─────┐
│ FIRST_NAME │ LAST_NAME │ AGE │
├────────────┼───────────┼─────┤
│ Dustin     │ Zeisler   │ 100 │
└────────────┴───────────┴─────┘
:┌────────────┬───────────┬─────┐
│ FIRST_NAME │ LAST_NAME │ AGE │
├────────────┼───────────┼─────┤
│ Tom        │ Smith     │ 99  │
└────────────┴───────────┴─────┘
:`, "\n"+out.String())
}

func TestTableMarshalIter_FilterIter(t *testing.T) {
	assert := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	p2 := Person{FirstName: "Tom", LastName: "Smith", Age: 99}
	it := &MockIter{eofPage: true, People: []Person{p1, p2}}
	out := strings.Builder{}
	in := strings.NewReader(" \n")
	TableMarshalIter("", it, "", &out, in, func(i interface{}) bool {
		if i.(Person).FirstName == "Dustin" {
			return true
		}
		return false
	})

	assert.Equal(`
┌────────────┬───────────┬─────┐
│ FIRST_NAME │ LAST_NAME │ AGE │
├────────────┼───────────┼─────┤
│ Dustin     │ Zeisler   │ 100 │
└────────────┴───────────┴─────┘
:`, "\n"+out.String())
}
