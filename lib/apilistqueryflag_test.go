package lib

import (
	"testing"

	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAPIListQueryFlag(t *testing.T) {
	t.Run("returns nil map when unset", func(t *testing.T) {
		parsed, err := ParseAPIListQueryFlag("filter-gt", nil)
		require.NoError(t, err)
		assert.Nil(t, parsed)
	})

	t.Run("parses repeated values", func(t *testing.T) {
		parsed, err := ParseAPIListQueryFlag("filter-lteq", []string{"created_at=2024-01-01", "id=42"})
		require.NoError(t, err)
		assert.Equal(t, map[string]interface{}{
			"created_at": "2024-01-01",
			"id":         "42",
		}, parsed)
	})

	t.Run("parses exact filters", func(t *testing.T) {
		parsed, err := ParseAPIListQueryFlag("filter", []string{"not_site_admin=true"})
		require.NoError(t, err)
		assert.Equal(t, map[string]interface{}{
			"not_site_admin": "true",
		}, parsed)
	})

	t.Run("rejects malformed values", func(t *testing.T) {
		_, err := ParseAPIListQueryFlag("filter-prefix", []string{"username"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected key=value")
		assert.Equal(t, clierr.ErrorCodeUsage, clierr.From(err).Code)
	})

	t.Run("rejects empty keys", func(t *testing.T) {
		_, err := ParseAPIListQueryFlag("filter-gt", []string{"=2024-01-01"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "key cannot be empty")
		assert.Equal(t, clierr.ErrorCodeUsage, clierr.From(err).Code)
	})
}

func TestParseAPIListSortFlag(t *testing.T) {
	t.Run("returns nil map when unset", func(t *testing.T) {
		parsed, err := ParseAPIListSortFlag("sort-by", "")
		require.NoError(t, err)
		assert.Nil(t, parsed)
	})

	t.Run("parses a single sort value", func(t *testing.T) {
		parsed, err := ParseAPIListSortFlag("sort-by", "username=asc")
		require.NoError(t, err)
		assert.Equal(t, map[string]interface{}{"username": "asc"}, parsed)
	})

	t.Run("rejects malformed sort values", func(t *testing.T) {
		_, err := ParseAPIListSortFlag("sort-by", "username")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected key=value")
		assert.Equal(t, clierr.ErrorCodeUsage, clierr.From(err).Code)
	})

	t.Run("rejects invalid sort directions", func(t *testing.T) {
		_, err := ParseAPIListSortFlag("sort-by", "username=descending")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "sort direction must be asc or desc")
		assert.Equal(t, clierr.ErrorCodeUsage, clierr.From(err).Code)
	})
}
