package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// runDiscovery executes args through the real root command, as the binary
// does, and returns what it wrote to stdout.
func runDiscovery(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var stdout bytes.Buffer
	err := runDiscoveryTo(t, &stdout, args...)
	return stdout.String(), err
}

// errRootPreRun stands in for the root pre-run, which loads the user's
// profile and may check the version and authenticate, while discovery tests
// run.
var errRootPreRun = errors.New("root pre-run called")

// runDiscoveryTo executes args through the real root command with stdout set
// to out, an HTTP client that fails the test on any request, and no
// credentials. The root pre-run is replaced by one that fails with
// errRootPreRun, so a command that does not bypass it fails without touching
// the user's profile. The discovery subcommands are rebuilt for each run so
// flag values never carry over, and the guides are read from the generated
// CLI.
func runDiscoveryTo(t *testing.T, out io.Writer, args ...string) error {
	t.Helper()
	for _, fresh := range []*cobra.Command{Commands(), Workflows()} {
		for _, existing := range RootCmd.Commands() {
			if existing.Name() == fresh.Name() {
				RootCmd.RemoveCommand(existing)
			}
		}
		RootCmd.AddCommand(fresh)
	}
	previousGuides := Guides
	Guides = os.DirFS("..")
	rootPreRun := RootCmd.PersistentPreRunE
	RootCmd.PersistentPreRunE = func(*cobra.Command, []string) error { return errRootPreRun }
	t.Cleanup(func() {
		RootCmd.PersistentPreRunE = rootPreRun
		Guides = previousGuides
		resetRootCommandState()
		APIKey = ""
		RootCmd.PersistentFlags().Lookup(flagNameApiKey).Changed = false
	})

	config := newTestConfig(func(req *http.Request) (*http.Response, error) {
		t.Errorf("unexpected request: %s %s", req.Method, req.URL)
		return nil, errors.New("requests are not allowed in this test")
	})
	config.APIKey = ""
	RootCmd.SetOut(out)
	RootCmd.SetErr(io.Discard)
	RootCmd.SetArgs(args)
	return execute(config)
}

// discoverJSON runs a discovery command with --format=json and decodes it.
func discoverJSON[T any](t *testing.T, args ...string) T {
	t.Helper()
	stdout, err := runDiscovery(t, append(args, "--format=json")...)
	require.NoError(t, err)
	var decoded T
	require.NoError(t, json.Unmarshal([]byte(stdout), &decoded), stdout)
	return decoded
}

func flagNamed(t *testing.T, flags []flagDescription, name string) flagDescription {
	t.Helper()
	for _, flag := range flags {
		if flag.Name == name {
			return flag
		}
	}
	require.Failf(t, "missing flag", "--%s", name)
	return flagDescription{}
}

func TestDiscoveryBypassesRootPreRun(t *testing.T) {
	require.ErrorIs(t, runDiscoveryTo(t, io.Discard, "version"), errRootPreRun, "other commands still run the root pre-run")

	for _, args := range [][]string{
		{"commands"},
		{"commands", "list", "users"},
		{"commands", "search", "share", "link"},
		{"commands", "describe", "users", "list", "--format=json"},
		{"workflows"},
		{"workflows", "search", "folder", "size"},
		{"workflows", "show", "context"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			stdout, err := runDiscovery(t, append([]string{"--non-interactive"}, args...)...)

			require.NoError(t, err)
			assert.NotEmpty(t, stdout)
		})
	}
}

func TestCommandsIndexAndSearch(t *testing.T) {
	index := discoverJSON[commandList](t, "commands")
	listed := make(map[string]commandSummary)
	for _, command := range index.Commands {
		listed[command.Command] = command
	}
	for _, name := range []string{"upload", "download", "sync", "users", "workflows"} {
		assert.Contains(t, listed, name)
	}
	for _, name := range []string{"agent", "raw-download-url", "generate-fig-spec", "help", "account-line-items"} {
		assert.NotContains(t, listed, name, "hidden commands and groups without commands are not listed")
	}
	assert.Equal(t, 2, listed["sync"].SubcommandCount)
	assert.Equal(t, []string{"[remote-path]", "[local-path]"}, listed["download"].Args)

	group := discoverJSON[commandList](t, "commands", "list", "sync")
	assert.Equal(t, "sync", group.Group)
	assert.Len(t, group.Commands, 2)

	search := discoverJSON[commandSearch](t, "commands", "search", "share link", "--limit=3")
	assert.Len(t, search.Commands, 3)
	assert.Greater(t, search.Total, 3)
	assert.True(t, search.Truncated)
	for _, command := range search.Commands {
		assert.Contains(t, strings.ToLower(command.Short), "share link")
	}

	_, err := runDiscovery(t, "commands", "describe", "agent")
	assert.Equal(t, clierr.ErrorCodeUsage, clierr.From(err).Code, "hidden commands are not described")
}

func TestCommandsDescribeReportsTheCommandTree(t *testing.T) {
	describe := func(t *testing.T, args ...string) commandDescription {
		t.Helper()
		return discoverJSON[commandDescription](t, append([]string{"commands", "describe"}, args...)...)
	}

	t.Run("handwritten command", func(t *testing.T) {
		upload := describe(t, "upload")

		assert.Equal(t, "files-cli upload [*source-path] [remote-path] [flags]", upload.Usage)
		assert.Equal(t, []string{"[*source-path]", "[remote-path]"}, upload.Args)
		assert.Equal(t, "bool", flagNamed(t, upload.Flags, "dry-run").Type)
		assert.Equal(t, "o", flagNamed(t, upload.InheritedFlags, "output").Shorthand)
	})

	t.Run("generated list command by alias", func(t *testing.T) {
		list := describe(t, "users", "ls")

		assert.Equal(t, "users list", list.Command)
		assert.Equal(t, []string{"ls"}, list.Aliases)
		filter := flagNamed(t, list.Flags, "filter")
		assert.Equal(t, "stringArray", filter.Type)
		assert.Equal(t, "field=value", filter.Syntax)
		maxPages := flagNamed(t, list.Flags, "max-pages")
		assert.Equal(t, "m", maxPages.Shorthand)
		assert.Equal(t, "0", maxPages.Default)
		assert.True(t, flagNamed(t, list.Flags, "cursor").Truncated)
		assert.Equal(t, "bool", flagNamed(t, list.Flags, "json-envelope").Type)
	})

	t.Run("flags the CLI requires", func(t *testing.T) {
		upload := describe(t, "upload-to-child-site")

		siteID := flagNamed(t, upload.Flags, "site-id")
		assert.True(t, siteID.Required)
		assert.False(t, siteID.APIRequired)
	})

	t.Run("API-required and enum flags", func(t *testing.T) {
		create := describe(t, "users", "create")

		username := flagNamed(t, create.Flags, "username")
		assert.True(t, username.APIRequired)
		assert.False(t, username.Required, "the CLI does not enforce API requirements")
		assert.False(t, flagNamed(t, create.Flags, "name").APIRequired)
		assert.Contains(t, flagNamed(t, create.Flags, "authentication-method").Enum, "password")
	})

	t.Run("path given positionally is not an API-required flag", func(t *testing.T) {
		listFor := describe(t, "folders", "list-for")

		assert.Equal(t, []string{"[path]"}, listFor.Args)
		assert.False(t, flagNamed(t, listFor.Flags, "path").APIRequired)
	})

	t.Run("hidden and deprecated flags are omitted", func(t *testing.T) {
		push := describe(t, "sync", "push")

		flagNamed(t, push.InheritedFlags, "delete-source-files")
		for _, flag := range append(push.Flags, push.InheritedFlags...) {
			assert.NotContains(t, []string{"delete-source", "environment", "help"}, flag.Name)
		}
	})

	t.Run("--flag shows only the named flags in full", func(t *testing.T) {
		list := describe(t, "users", "list", "--flag=cursor")

		require.Len(t, list.Flags, 1)
		assert.Empty(t, list.InheritedFlags)
		assert.False(t, list.Flags[0].Truncated)
		assert.Greater(t, len(list.Flags[0].Description), maxFlagDescription)
	})

	t.Run("values given to this invocation are never described", func(t *testing.T) {
		stdout, err := runDiscovery(t, "--api-key=CANARY_KEY_VALUE", "commands", "describe", "users", "list", "--format=json")

		require.NoError(t, err)
		assert.NotContains(t, stdout, "CANARY_KEY_VALUE")
	})
}

func TestWorkflowsReturnOnlyTheRequestedGuide(t *testing.T) {
	list := discoverJSON[workflowList](t, "workflows")
	var names []string
	for _, guide := range list.Workflows {
		names = append(names, guide.Name)
		assert.NotEmpty(t, guide.Description, guide.Name)
		assert.Empty(t, guide.Content, "the list carries no guide content")
	}
	assert.Contains(t, names, "context")
	assert.Contains(t, names, "recipe-searching-for-files")

	guide := discoverJSON[workflowGuide](t, "workflows", "show", "recipe-searching-for-files")
	assert.Equal(t, "recipe-searching-for-files", guide.Name)
	assert.NotEmpty(t, guide.Description)
	assert.True(t, strings.HasPrefix(guide.Content, "# recipe-searching-for-files"), guide.Content)
	assert.NotContains(t, guide.Content, "# recipe-share-and-notify")

	search := discoverJSON[workflowSearch](t, "workflows", "search", "folder size")
	require.NotEmpty(t, search.Workflows)
	assert.Equal(t, "recipe-folder-size-and-counts", search.Workflows[0].Name)

	_, err := runDiscovery(t, "workflows", "show", "missing")
	assert.Equal(t, clierr.ErrorCodeUsage, clierr.From(err).Code)
}

type failingWriter struct{}

var errOutputFailed = errors.New("synthetic output failure")

func (failingWriter) Write([]byte) (int, error) {
	return 0, errOutputFailed
}

func TestDiscoveryReportsOutputFailure(t *testing.T) {
	for _, format := range []string{"text", "json"} {
		t.Run(format, func(t *testing.T) {
			err := runDiscoveryTo(t, failingWriter{}, "commands", "--format="+format)

			require.ErrorIs(t, err, errOutputFailed)
			assert.Equal(t, clierr.ErrorCodeFatal, clierr.From(err).Code)
		})
	}
}
