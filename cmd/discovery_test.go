package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
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
// flag values never carry over. Guides and response fields come from the data
// embedded in the package, as in the binary.
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
	rootPreRun := RootCmd.PersistentPreRunE
	RootCmd.PersistentPreRunE = func(*cobra.Command, []string) error { return errRootPreRun }
	t.Cleanup(func() {
		RootCmd.PersistentPreRunE = rootPreRun
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
		{"commands", "describe", "permissions", "create", "--response-field=path"},
		{"workflows"},
		{"workflows", "list", "--domains"},
		{"workflows", "search", "folder", "size"},
		{"workflows", "show", "context"},
		{"workflows", "show", "filescom-permissions"},
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
		assert.NotEqual(t, guideKindDomain, guide.Kind, "domain guides are listed with --domains")
	}
	assert.Contains(t, names, "context")
	assert.Contains(t, names, "recipe-searching-for-files")

	domains := discoverJSON[workflowList](t, "workflows", "list", "--domains")
	assert.Equal(t, list.DomainGuideCount, len(domains.Workflows))
	assert.Zero(t, domains.DomainGuideCount, "the domain listing omits nothing to count")
	domainsText, err := runDiscovery(t, "workflows", "list", "--domains")
	require.NoError(t, err)
	assert.NotContains(t, domainsText, "workflows list --domains")
	domainGroups := make(map[string]string)
	for _, guide := range domains.Workflows {
		assert.Equal(t, guideKindDomain, guide.Kind, guide.Name)
		assert.NotEmpty(t, guide.Description, guide.Name)
		assert.Empty(t, guide.Content, "the list carries no guide content")
		domainGroups[guide.Name] = guide.Group
	}
	assert.Equal(t, "permissions", domainGroups["filescom-permissions"])
	assert.Equal(t, "history-exports", domainGroups["filescom-history-exports"])

	t.Run("domain guides carry the authored guidance and the group's commands, not flag tables", func(t *testing.T) {
		for name, authored := range map[string]string{
			"filescom-permissions": "To revoke access, delete the permission by its ID.",
			"filescom-users":       "Disable vs delete.",
			"filescom-workspaces":  "Resources cannot be moved between Workspaces.",
			"filescom-bundles":     "Sending the link is a separate step from creating it.",
		} {
			guide := discoverJSON[workflowGuide](t, "workflows", "show", name)

			assert.Contains(t, guide.Content, "## Limitations and considerations", name)
			assert.Contains(t, guide.Content, authored, name)
			assert.Contains(t, guide.Content, "## Commands", name)
			assert.NotContains(t, guide.Content, "| Flag |", name)
		}

		permissions := discoverJSON[workflowGuide](t, "workflows", "show", "filescom-permissions")
		assert.Contains(t, permissions.Content, "## Common patterns")
		assert.Contains(t, permissions.Content, "- `permissions delete`: Delete Permission (destructive)")

		text, err := runDiscovery(t, "workflows", "show", "filescom-permissions")
		require.NoError(t, err)
		assert.True(t, strings.HasPrefix(text, "---\nname: filescom-permissions\n"), "text is the Markdown source with its front matter")
		assert.Contains(t, text, "run `files-cli commands describe permissions <command>`")
	})

	guide := discoverJSON[workflowGuide](t, "workflows", "show", "recipe-searching-for-files")
	assert.Equal(t, "recipe-searching-for-files", guide.Name)
	assert.NotEmpty(t, guide.Description)
	assert.True(t, strings.HasPrefix(guide.Content, "# recipe-searching-for-files"), guide.Content)
	assert.NotContains(t, guide.Content, "# recipe-share-and-notify")

	search := discoverJSON[workflowSearch](t, "workflows", "search", "folder size")
	require.NotEmpty(t, search.Workflows)
	assert.Equal(t, "recipe-folder-size-and-counts", search.Workflows[0].Name)

	_, err = runDiscovery(t, "workflows", "show", "missing")
	assert.Equal(t, clierr.ErrorCodeUsage, clierr.From(err).Code)
}

func TestCommandsDescribeEffectHints(t *testing.T) {
	for command, want := range map[string]string{
		"users list":                  "read_only",   // GET
		"permissions create":          "mutating",    // POST, not known to be destructive
		"permissions delete":          "destructive", // DELETE
		"api-keys update-current":     "destructive", // PATCH not named update
		"files copy":                  "destructive", // POST that can overwrite the destination
		"files copy-to-remote-server": "destructive", // reviewed handwritten copy
		"commands describe":           "read_only",   // reviewed handwritten discovery
		"upload":                      "",            // handwritten, none declared
	} {
		t.Run(command, func(t *testing.T) {
			assert.Equal(t, want, discoverJSON[commandDescription](t, "commands", "describe", command).Effect)
		})
	}

	t.Run("summaries carry the same hints", func(t *testing.T) {
		effects := make(map[string]string)
		for _, command := range discoverJSON[commandList](t, "commands", "list", "permissions").Commands {
			effects[command.Command] = command.Effect
		}
		assert.Equal(t, "read_only", effects["permissions list"])
		assert.Equal(t, "mutating", effects["permissions create"])
		assert.Equal(t, "destructive", effects["permissions delete"])
	})

	t.Run("text says when no effect is declared", func(t *testing.T) {
		text, err := runDiscovery(t, "commands", "describe", "upload")

		require.NoError(t, err)
		assert.Contains(t, text, "Effect: not declared for this command.")
	})
}

func responseFieldNames(fields []responseField) []string {
	names := make([]string, len(fields))
	for i, field := range fields {
		names[i] = field.Name
	}
	return names
}

// assertSelectable checks that --fields accepts every name for records of
// the Go SDK type the command prints.
func assertSelectable(t *testing.T, names []string, record interface{}) {
	t.Helper()
	require.NotEmpty(t, names)
	_, _, err := lib.OnlyFields(names, record)
	assert.NoError(t, err)
}

func TestCommandsDescribeResponseFields(t *testing.T) {
	describe := func(t *testing.T, args ...string) commandDescription {
		t.Helper()
		description := discoverJSON[commandDescription](t, append([]string{"commands", "describe"}, args...)...)
		require.NotNil(t, description.Response, args)
		return description
	}

	t.Run("list command names the fields --fields selects, without descriptions by default", func(t *testing.T) {
		list := describe(t, "users", "list")

		assert.Equal(t, "User", list.Response.Type)
		assert.True(t, list.Response.List)
		assert.Contains(t, list.Response.Fields, responseField{Name: "username", Type: "string"})
		assertSelectable(t, responseFieldNames(list.Response.Fields), files_sdk.User{})

		full := describe(t, "users", "list", "--full")
		require.Len(t, full.Response.Fields, len(list.Response.Fields))
		assert.NotEmpty(t, full.Response.Fields[slices.Index(responseFieldNames(full.Response.Fields), "username")].Description)
	})

	t.Run("recommended fields for history and remote server records", func(t *testing.T) {
		history := describe(t, "histories", "list")

		assert.Equal(t, "Action", history.Response.Type)
		assert.Equal(t, []string{"id", "path", "when", "destination", "display", "ip", "source", "targets", "user_id", "username", "action", "failure_type", "interface"}, history.Response.RecommendedFields)
		assertSelectable(t, history.Response.RecommendedFields, files_sdk.Action{})

		remote := describe(t, "remote-servers", "find")
		assert.Equal(t, []string{"id", "name", "server_type"}, remote.Response.RecommendedFields)
		assertSelectable(t, remote.Response.RecommendedFields, files_sdk.RemoteServer{})

		agentNodes := describe(t, "remote-servers", "agent-nodes")
		assert.Equal(t, "AgentNode", agentNodes.Response.Type)
		assert.Empty(t, agentNodes.Response.RecommendedFields, "recommendations belong to the record type, not the command group")
	})

	t.Run("fields of the operation's entity, not its group's", func(t *testing.T) {
		usage := describe(t, "sites", "get-usage")

		assert.Equal(t, "UsageSnapshot", usage.Response.Type)
		assert.False(t, usage.Response.List)
		assertSelectable(t, responseFieldNames(usage.Response.Fields), files_sdk.UsageSnapshot{})
		assert.NotContains(t, responseFieldNames(usage.Response.Fields), "subdomain", "a Site field")
	})

	t.Run("--block can print a FileMigration instead", func(t *testing.T) {
		for _, command := range [][]string{{"files", "copy"}, {"files", "move-to-remote-server"}} {
			copied := describe(t, command...)

			assert.Equal(t, "FileAction", copied.Response.Type, command)
			assertSelectable(t, responseFieldNames(copied.Response.Fields), files_sdk.FileAction{})
			require.NotNil(t, copied.Response.WithFlag, command)
			assert.Equal(t, "block", copied.Response.WithFlag.Flag)
			assert.Equal(t, "FileMigration", copied.Response.WithFlag.Type)
			assertSelectable(t, responseFieldNames(copied.Response.WithFlag.Fields), files_sdk.FileMigration{})
			flagNamed(t, copied.Flags, "block")
		}
	})

	t.Run("commands that print nothing, and handwritten commands with no response declared", func(t *testing.T) {
		deleted := describe(t, "permissions", "delete")
		assert.True(t, deleted.Response.None)
		assert.Empty(t, deleted.Response.Fields)

		assert.Nil(t, discoverJSON[commandDescription](t, "commands", "describe", "upload").Response)
		_, err := runDiscovery(t, "commands", "describe", "upload", "--response-field=id")
		assert.Equal(t, clierr.ErrorCodeUsage, clierr.From(err).Code)
	})

	t.Run("--response-field shows only the named fields, with descriptions", func(t *testing.T) {
		remote := describe(t, "remote-servers", "find", "--response-field=server-type")

		require.Len(t, remote.Response.Fields, 1)
		assert.Equal(t, "server_type", remote.Response.Fields[0].Name)
		assert.NotEmpty(t, remote.Response.Fields[0].Description)
		assert.Empty(t, remote.Flags)
		assert.Empty(t, remote.InheritedFlags)

		_, err := runDiscovery(t, "commands", "describe", "remote-servers", "find", "--response-field=missing")
		assert.Equal(t, clierr.ErrorCodeUsage, clierr.From(err).Code)
	})

	t.Run("--flag alone lists no response fields", func(t *testing.T) {
		list := describe(t, "users", "list", "--flag=cursor")

		assert.Equal(t, "User", list.Response.Type)
		assert.Empty(t, list.Response.Fields)
	})
}

func TestEveryDeclaredResponseEntityIsDescribed(t *testing.T) {
	entities, err := loadResponseEntities()
	require.NoError(t, err)

	declared := 0
	var visit func(parent *cobra.Command)
	visit = func(parent *cobra.Command) {
		for _, command := range parent.Commands() {
			if entity, _, ok := lib.GetCommandResponse(command); ok && entity != "" {
				declared++
				assert.Contains(t, entities, entity, commandName(command))
			}
			if _, entity, _, ok := lib.GetCommandResponseWithFlag(command); ok {
				assert.Contains(t, entities, entity, commandName(command))
			}
			visit(command)
		}
	}
	visit(RootCmd)
	assert.Greater(t, declared, 300, "generated commands declare their response")

	for name, entity := range entities {
		assert.NotEmpty(t, entity.Fields, name)
		assert.Subset(t, responseFieldNames(entity.Fields), entity.Recommended, "recommended fields of %s", name)
	}
}

// corruptResponseFieldsEnv makes TestCorruptResponseFieldsFailOnlyResponseDescriptions
// run its checks in the child test process it starts.
const corruptResponseFieldsEnv = "FILES_CLI_TEST_CORRUPT_RESPONSE_FIELDS"

// The guides and the response fields are decoded separately, so unreadable
// response fields break only the output that describes responses. The checks
// run in a fresh test process, where nothing has decoded either asset yet.
func TestCorruptResponseFieldsFailOnlyResponseDescriptions(t *testing.T) {
	if os.Getenv(corruptResponseFieldsEnv) != "1" {
		child := exec.Command(os.Args[0], "-test.run=^TestCorruptResponseFieldsFailOnlyResponseDescriptions$", "-test.count=1", "-test.v")
		child.Env = append(os.Environ(), corruptResponseFieldsEnv+"=1")
		output, err := child.CombinedOutput()

		require.NoError(t, err, string(output))
		require.Contains(t, string(output), "--- PASS: TestCorruptResponseFieldsFailOnlyResponseDescriptions", string(output))
		return
	}

	responseData = []byte("not gzip")

	for _, args := range [][]string{
		{"commands"},
		{"commands", "search", "share", "link"},
		{"commands", "describe", "upload"},
		{"workflows", "show", "filescom-permissions"},
	} {
		stdout, err := runDiscovery(t, args...)
		require.NoError(t, err, args)
		assert.NotEmpty(t, stdout, args)
	}

	_, err := runDiscovery(t, "commands", "describe", "users", "list")
	assert.Equal(t, clierr.ErrorCodeFatal, clierr.From(err).Code)
	assert.ErrorContains(t, err, "the response fields built into this binary are unreadable")
}

func TestCommandsDescribeLinksRelatedGuides(t *testing.T) {
	links := func(t *testing.T, command string) []string {
		t.Helper()
		var names []string
		for _, link := range discoverJSON[commandDescription](t, "commands", "describe", command).Workflows {
			assert.NotEmpty(t, link.Description, link.Name)
			names = append(names, link.Name)
		}
		return names
	}

	assert.Equal(t, []string{"filescom-permissions"}, links(t, "permissions create"))
	assert.Equal(t, []string{"filescom-users"}, links(t, "users list"), "no recipe lists users list")
	assert.Equal(t, []string{"filescom-folders", "recipe-folder-size-and-counts", "recipe-searching-for-files"}, links(t, "folders ls"), "an alias links the command's guides")
	for _, command := range []string{"history-exports create", "history-exports find"} {
		assert.Contains(t, links(t, command), "recipe-generating-reports", command)
	}
	assert.Equal(t, []string{"filescom-files"}, links(t, "files copy-to-remote-server"), "a handwritten command links its group's guide")
	assert.Empty(t, links(t, "upload"))

	text, err := runDiscovery(t, "commands", "describe", "permissions", "create")
	require.NoError(t, err)
	assert.Contains(t, text, `read one with "files-cli workflows show <name>"`)
	assert.Contains(t, text, "\n  filescom-permissions  ")
}

// guideReference matches a guide name quoted in guide text.
var guideReference = regexp.MustCompile("`((?:filescom|recipe)-[a-z0-9-]+)`")

// guideReferenceProblems reports the declared guide links and guide names that
// do not resolve: a domain guide group that is not a command group, a recipe
// command that is not a command's canonical name, and a `filescom-…` or
// `recipe-…` name in a guide's description or text that is not a guide.
// Example invocations in the prose are not checked.
func guideReferenceProblems(root *cobra.Command, guides []workflowGuide) []string {
	var problems []string
	names := make(map[string]bool)
	for _, guide := range guides {
		if guide.Name == "" || names[guide.Name] {
			problems = append(problems, fmt.Sprintf("guide name %q is empty or repeated", guide.Name))
		}
		names[guide.Name] = true
	}
	for _, guide := range guides {
		if guide.Kind == guideKindDomain && guide.Group == "" {
			problems = append(problems, fmt.Sprintf("%s: domain guide has no group", guide.Name))
		}
		if guide.Group != "" {
			group, err := lookupDiscoverable(root, []string{guide.Group})
			if err != nil || commandName(group) != guide.Group || !group.HasAvailableSubCommands() {
				problems = append(problems, fmt.Sprintf("%s: group %q is not a command group", guide.Name, guide.Group))
			}
		}
		for _, name := range guide.Commands {
			command, err := lookupDiscoverable(root, []string{name})
			if err != nil || commandName(command) != name || command.HasAvailableSubCommands() {
				problems = append(problems, fmt.Sprintf("%s: %q is not a command's name", guide.Name, name))
			}
		}
		for _, reference := range guideReference.FindAllStringSubmatch(guide.Description+"\n"+guide.Content, -1) {
			if !names[reference[1]] {
				problems = append(problems, fmt.Sprintf("%s: unknown guide %q", guide.Name, reference[1]))
			}
		}
	}
	return problems
}

func TestWorkflowGuideReferencesResolve(t *testing.T) {
	guides, err := loadWorkflowGuides()
	require.NoError(t, err)

	assert.Empty(t, guideReferenceProblems(RootCmd, guides))

	t.Run("stale references are reported", func(t *testing.T) {
		problems := guideReferenceProblems(RootCmd, []workflowGuide{
			{Name: "filescom-renamed", Kind: guideKindDomain, Group: "no-such-group"},
			{Name: "recipe-stale", Kind: guideKindRecipe, Commands: []string{"folders ls", "users nope", "users"}, Content: "See `filescom-missing` and `recipe-stale`."},
		})

		assert.ElementsMatch(t, []string{
			`filescom-renamed: group "no-such-group" is not a command group`,
			`recipe-stale: "folders ls" is not a command's name`,
			`recipe-stale: "users nope" is not a command's name`,
			`recipe-stale: "users" is not a command's name`,
			`recipe-stale: unknown guide "filescom-missing"`,
		}, problems)
	})
}

// The binary serves the same context and recipes that the repository
// publishes, and one domain guide for each published domain skill.
func TestEmbeddedGuidesMatchThePublishedFiles(t *testing.T) {
	guides, err := loadWorkflowGuides()
	require.NoError(t, err)

	var recipes, domains []string
	for _, guide := range guides {
		var published string
		switch guide.Kind {
		case guideKindContext:
			published = "../CONTEXT.md"
		case guideKindRecipe:
			recipes = append(recipes, guide.Name)
			published = "../skills/recipes/" + guide.Name + "/SKILL.md"
		case guideKindDomain:
			domains = append(domains, guide.Name)
			continue
		}
		data, err := os.ReadFile(published)
		require.NoError(t, err, guide.Name)
		assert.Equal(t, string(data), guide.source, guide.Name)
	}

	publishedNames := func(pattern string) []string {
		paths, err := filepath.Glob(pattern)
		require.NoError(t, err)
		names := make([]string, len(paths))
		for i, path := range paths {
			names[i] = filepath.Base(filepath.Dir(path))
		}
		return names
	}
	assert.ElementsMatch(t, publishedNames("../skills/recipes/*/SKILL.md"), recipes)
	assert.ElementsMatch(t, publishedNames("../skills/filescom-*/SKILL.md"), domains)
}

type failingWriter struct{}

var errOutputFailed = errors.New("synthetic output failure")

func (failingWriter) Write([]byte) (int, error) {
	return 0, errOutputFailed
}

func TestDiscoveryReportsOutputFailure(t *testing.T) {
	for _, args := range [][]string{{"commands"}, {"workflows", "show", "filescom-permissions"}} {
		for _, format := range []string{"text", "json"} {
			t.Run(strings.Join(args, " ")+" "+format, func(t *testing.T) {
				err := runDiscoveryTo(t, failingWriter{}, append(args, "--format="+format)...)

				require.ErrorIs(t, err, errOutputFailed)
				assert.Equal(t, clierr.ErrorCodeFatal, clierr.From(err).Code)
			})
		}
	}
}
