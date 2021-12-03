package lib

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonMarshalIter(t *testing.T) {
	a := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	p2 := Person{FirstName: "Tom", LastName: "Smith", Age: 99}
	it := MockIter{People: []Person{p1, p2}}
	buf := bytes.NewBufferString("")

	JsonMarshalIter(&it, "", nil, buf)

	a.Equal(`[{
    "age": 100,
    "first_name": "Dustin",
    "last_name": "Zeisler"
},
{
    "age": 99,
    "first_name": "Tom",
    "last_name": "Smith"
}]
`, buf.String())
}

func TestJsonMarshalIter_Filter(t *testing.T) {
	a := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	p2 := Person{FirstName: "Tom", LastName: "Smith", Age: 99}
	it := MockIter{People: []Person{p1, p2}}
	buf := bytes.NewBufferString("")
	JsonMarshalIter(&it, "", func(i interface{}) bool {
		return i.(Person).FirstName == "Dustin"
	}, buf)

	a.Equal(`[{
    "age": 100,
    "first_name": "Dustin",
    "last_name": "Zeisler"
}]
`, buf.String())
}
