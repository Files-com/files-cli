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
	t.Run("it deals with a struct", func(t *testing.T) {
		a := testStruct{A: "hello", B: "I'm B", C: "I'm C"}
		result, _, _ := OnlyFields([]string{}, a)

		assert.Equal("hello", result["a"])
		assert.Equal("I'm B", result["b"])
		assert.Equal("I'm C", result["c"])

		result2, _, _ := OnlyFields([]string{"b", "C"}, a)

		assert.Equal(nil, result2["a"])
		assert.Equal("I'm B", result2["b"])
		assert.Equal("I'm C", result2["c"])

		_, _, err1 := OnlyFields([]string{"d"}, a)

		assert.EqualError(err1, "field: `d` is not valid. - status (1)")

		b := testStruct{A: "hello", B: "I'm B"}

		result3, _, err3 := OnlyFields([]string{"c"}, b)

		assert.NoError(err3)

		assert.Equal(nil, result3["c"])
		assert.Equal(nil, result3["a"])
		assert.Equal(nil, result3["b"])

		keys := make([]string, 0)
		for key := range result3 {
			keys = append(keys, key)
		}

		assert.Equal([]string{}, keys, "returns no keys")
		result4, _, err4 := OnlyFields([]string{"a", "c"}, b)

		assert.NoError(err4)

		keys = make([]string, 0)
		for key := range result4 {
			keys = append(keys, key)
		}

		assert.Equal([]string{"a"}, keys, "returns only valid keys")
		assert.Equal("hello", result4["a"])
	})

	t.Run("it deals with a Map", func(t *testing.T) {
		m := make(map[string]interface{})
		m["key"] = "value"
		results, orderedKeys, err := OnlyFields([]string{}, m)
		assert.NoError(err)
		assert.Equal(m, results)
		assert.Equal([]string{"key"}, orderedKeys)
	})

	t.Run("it subtracts a field", func(t *testing.T) {
		a := testStruct{A: "hello", B: "I'm B", C: "I'm C"}
		results, orderedKeys, err := OnlyFields([]string{"-a"}, a)
		assert.NoError(err)
		assert.Equal(map[string]interface{}{"b": "I'm B", "c": "I'm C"}, results)
		assert.Equal([]string{"b", "c"}, orderedKeys)
	})
}
