package lib

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableMarshalIter_q(t *testing.T) {
	assert := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	p2 := Person{FirstName: "Tom", LastName: "Smith", Age: 99}
	it := &MockIter{SliceIter: SliceIter{Items: []interface{}{p1, p2}}, eofPage: func(iter *MockIter) bool {
		return true
	}}
	out := strings.Builder{}
	in := strings.NewReader("q")
	TableMarshalIter("", it, "", &out, in, nil)
	assert.Equal(strings.TrimSpace(`
┌────────────┬───────────┬─────┐
│ FIRST_NAME │ LAST_NAME │ AGE │
├────────────┼───────────┼─────┤
│ Dustin     │ Zeisler   │ 100 │
└────────────┴───────────┴─────┘
`), strings.TrimSpace(sanitizeOutput(out.String())))
}

func TestTableMarshalIter_newline(t *testing.T) {
	assert := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	p2 := Person{FirstName: "Tom", LastName: "Smith", Age: 99}
	it := &MockIter{SliceIter: SliceIter{Items: []interface{}{p1, p2}}, eofPage: func(iter *MockIter) bool {
		return true
	}}
	out := strings.Builder{}
	in := strings.NewReader(" \n\n\n\n")
	TableMarshalIter("", it, "", &out, in, nil)

	assert.Equal(strings.TrimSpace(`
┌────────────┬───────────┬─────┐
│ FIRST_NAME │ LAST_NAME │ AGE │
├────────────┼───────────┼─────┤
│ Dustin     │ Zeisler   │ 100 │
└────────────┴───────────┴─────┘
┌────────────┬───────────┬─────┐
│ FIRST_NAME │ LAST_NAME │ AGE │
├────────────┼───────────┼─────┤
│ Tom        │ Smith     │ 99  │
└────────────┴───────────┴─────┘
`), strings.TrimSpace(sanitizeOutput(out.String())))
}

func TestTableMarshalIter_FilterIter(t *testing.T) {
	assert := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	p2 := Person{FirstName: "Tom", LastName: "Smith", Age: 99}
	it := &MockIter{SliceIter: SliceIter{Items: []interface{}{p1, p2}}}
	out := strings.Builder{}
	in := strings.NewReader(" \n")
	TableMarshalIter("", it, "", &out, in, func(i interface{}) bool {
		if i.(Person).FirstName == "Dustin" {
			return true
		}
		return false
	})

	assert.Equal(strings.TrimSpace(`
┌────────────┬───────────┬─────┐
│ FIRST_NAME │ LAST_NAME │ AGE │
├────────────┼───────────┼─────┤
│ Dustin     │ Zeisler   │ 100 │
└────────────┴───────────┴─────┘
`), strings.TrimSpace(sanitizeOutput(out.String())))
}

func sanitizeOutput(str string) string {
	r, _ := regexp.Compile(`(┌[^┘]*┘)[^┌]*(┌[^┘]*┘)?`) // https://regoio.herokuapp.com
	matches := r.FindSubmatch([]byte(str))
	var newStr string
	for _, m := range matches[1:] {
		newStr += "\n" + string(m)
	}
	return newStr
}
