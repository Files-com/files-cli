package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func findSubcommand(t *testing.T, parent *cobra.Command, name string) *cobra.Command {
	t.Helper()

	for _, command := range parent.Commands() {
		if command.Name() == name {
			return command
		}
	}

	require.Failf(t, "missing subcommand", "expected subcommand %q on %q", name, parent.Name())
	return nil
}

func TestUsersListFlags(t *testing.T) {
	list := findSubcommand(t, Users(), "list")

	for _, flagName := range []string{
		"filter-by",
		"sort-by",
		"filter",
		"filter-gt",
		"filter-gteq",
		"filter-prefix",
		"filter-lt",
		"filter-lteq",
	} {
		assert.NotNil(t, list.Flags().Lookup(flagName), flagName)
	}
}

func TestFoldersListForFlags(t *testing.T) {
	listFor := findSubcommand(t, Folders(), "list-for")

	assert.NotNil(t, listFor.Flags().Lookup("filter-by"))
	assert.NotNil(t, listFor.Flags().Lookup("sort-by"))
	assert.Nil(t, listFor.Flags().Lookup("filter"))
	assert.Nil(t, listFor.Flags().Lookup("filter-gt"))
	assert.Nil(t, listFor.Flags().Lookup("filter-gteq"))
	assert.Nil(t, listFor.Flags().Lookup("filter-prefix"))
	assert.Nil(t, listFor.Flags().Lookup("filter-lt"))
	assert.Nil(t, listFor.Flags().Lookup("filter-lteq"))
}
