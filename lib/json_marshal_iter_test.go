package lib

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonMarshalIter(t *testing.T) {
	a := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	p2 := Person{FirstName: "Tom", LastName: "Smith", Age: 99}
	it := MockIter{SliceIter: SliceIter{Items: []interface{}{p1, p2}}}
	buf := bytes.NewBufferString("")

	JsonMarshalIter(context.Background(), &it, []string{}, nil, false, "", buf)

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

func TestJsonMarshalIter_Raw(t *testing.T) {
	a := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	p2 := Person{FirstName: "Tom&Jerry", LastName: "Smith", Age: 99}
	it := MockIter{SliceIter: SliceIter{Items: []interface{}{p1, p2}}}
	buf := bytes.NewBufferString("")

	JsonMarshalIter(context.Background(), &it, []string{}, nil, false, "raw", buf)

	a.Equal(`[{"age":100,"first_name":"Dustin","last_name":"Zeisler"},{"age":99,"first_name":"Tom&Jerry","last_name":"Smith"}]
`, buf.String())
}

type StringBufferWithClose struct {
	*bytes.Buffer
	closed bool
}

func (s *StringBufferWithClose) Init() *StringBufferWithClose {
	s.Buffer = bytes.NewBufferString("")
	return s
}

func (s *StringBufferWithClose) Close() error {
	s.closed = true
	return nil
}

func (s *StringBufferWithClose) Write(p []byte) (n int, err error) {
	if !s.closed {
		return s.Buffer.Write(p)
	}
	return 0, nil
}

func TestJsonMarshalIter_UsePager(t *testing.T) {
	a := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	p2 := Person{FirstName: "Tom", LastName: "Smith", Age: 99}
	it := MockIter{SliceIter: SliceIter{Items: []interface{}{p1, p2}}}
	buf := (&StringBufferWithClose{}).Init()

	JsonMarshalIter(context.Background(), &it, []string{}, nil, true, "", buf)

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
	it := MockIter{SliceIter: SliceIter{Items: []interface{}{p1, p2}}}
	buf := bytes.NewBufferString("")
	JsonMarshalIter(context.Background(), &it, []string{}, func(i interface{}) (interface{}, bool, error) {
		return i, i.(Person).FirstName == "Dustin", nil
	}, false, "", buf)

	a.Equal(`[{
    "age": 100,
    "first_name": "Dustin",
    "last_name": "Zeisler"
}]
`, buf.String())
}
