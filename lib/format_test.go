package lib

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	a := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	p2 := Person{FirstName: "Tom", LastName: "Smith", Age: 99}
	buf := bytes.NewBufferString("")
	result := []interface{}{p1, p2}
	Format(context.Background(), result, "json", "", false, buf)

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
