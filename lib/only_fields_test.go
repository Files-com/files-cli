package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	A string `json:"a,omitempty"`
	B string `json:"b,omitempty"`
	C string `json:"c,omitempty"`
}

func TestOnlyFields(t *testing.T) {
	assert := assert.New(t)
	a := testStruct{A: "hello", B: "I'm B", C: "I'm C"}
	result, _ := OnlyFields("", a)

	assert.Equal("hello", result["a"])
	assert.Equal("I'm B", result["b"])
	assert.Equal("I'm C", result["c"])

	result2, _ := OnlyFields("b,c", a)

	assert.Equal(nil, result2["a"])
	assert.Equal("I'm B", result2["b"])
	assert.Equal("I'm C", result2["c"])

	_, err := OnlyFields("d", a)

	assert.EqualError(err, "field: `d` is not valid.")

	b := testStruct{A: "hello", B: "I'm B"}

	result3, err := OnlyFields("c", b)

	assert.NoError(err)

	assert.Equal(nil, result3["c"])
	assert.Equal(nil, result3["a"])
	assert.Equal(nil, result3["b"])

	keys := make([]string, 0)
	for key := range result3 {
		keys = append(keys, key)
	}

	assert.Equal([]string{}, keys, "returns no keys")

	result4, err := OnlyFields("a,c", b)

	assert.NoError(err)

	keys = make([]string, 0)
	for key := range result4 {
		keys = append(keys, key)
	}

	assert.Equal([]string{"a"}, keys, "returns only valid keys")
	assert.Equal("hello", result4["a"])
}
