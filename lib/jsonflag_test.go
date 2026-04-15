package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseJSONObjectFlag(t *testing.T) {
	t.Run("parses object", func(t *testing.T) {
		parsed, err := ParseJSONObjectFlag("custom-metadata", `{"label":"sample","count":2}`)
		require.NoError(t, err)
		require.Equal(t, map[string]interface{}{
			"label": "sample",
			"count": float64(2),
		}, parsed)
	})

	t.Run("allows null", func(t *testing.T) {
		parsed, err := ParseJSONObjectFlag("custom-metadata", "null")
		require.NoError(t, err)
		assert.Nil(t, parsed)
	})

	t.Run("rejects invalid json", func(t *testing.T) {
		_, err := ParseJSONObjectFlag("custom-metadata", `{"label":`)
		require.Error(t, err)
		assert.ErrorContains(t, err, "expected valid JSON")
	})

	t.Run("rejects non object json", func(t *testing.T) {
		_, err := ParseJSONObjectFlag("custom-metadata", `["sample"]`)
		require.Error(t, err)
		assert.ErrorContains(t, err, "expected JSON object")
	})
}

func TestParseJSONArrayObjectFlag(t *testing.T) {
	t.Run("parses array of objects", func(t *testing.T) {
		parsed, err := ParseJSONArrayObjectFlag("members", `[{"user_id":1},{"group_id":2}]`)
		require.NoError(t, err)
		require.Equal(t, []map[string]interface{}{
			{"user_id": float64(1)},
			{"group_id": float64(2)},
		}, parsed)
	})

	t.Run("allows null", func(t *testing.T) {
		parsed, err := ParseJSONArrayObjectFlag("members", "null")
		require.NoError(t, err)
		assert.Nil(t, parsed)
	})

	t.Run("rejects invalid json", func(t *testing.T) {
		_, err := ParseJSONArrayObjectFlag("members", `[{"user_id":1}`)
		require.Error(t, err)
		assert.ErrorContains(t, err, "expected valid JSON")
	})

	t.Run("rejects non array json", func(t *testing.T) {
		_, err := ParseJSONArrayObjectFlag("members", `{"user_id":1}`)
		require.Error(t, err)
		assert.ErrorContains(t, err, "expected JSON array of objects")
	})

	t.Run("rejects non object members", func(t *testing.T) {
		_, err := ParseJSONArrayObjectFlag("members", `[{"user_id":1},42]`)
		require.Error(t, err)
		assert.ErrorContains(t, err, "item 2 must be a JSON object")
	})
}
