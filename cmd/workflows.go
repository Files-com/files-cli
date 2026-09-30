package cmd

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

const (
	contextGuideName        = "context"
	contextGuideDescription = "CLI-wide invocation guidance for agents: JSON output, non-interactive use, authentication, global flags, offline discovery, bounded listing with continuation, workspaces, and errors."
)

// Guide kinds: CONTEXT.md, the recipes for multi-step workflows, and one
// generated domain guide per command group.
const (
	guideKindContext = "context"
	guideKindRecipe  = "recipe"
	guideKindDomain  = "domain"
)

// Front matter metadata that links guides to commands. A domain guide names the
// command group it covers; a recipe lists, comma-separated, the commands it
// explains.
const (
	guideMetadataGroup    = "files-cli-group"
	guideMetadataCommands = "files-cli-commands"
)

func init() {
	RootCmd.AddCommand(Workflows())
}

type workflowGuide struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Kind        string   `json:"kind"`
	Group       string   `json:"group,omitempty"`
	Commands    []string `json:"commands,omitempty"`
	Content     string   `json:"content,omitempty"`
	source      string
}

type workflowList struct {
	Workflows        []workflowGuide `json:"workflows"`
	DomainGuideCount int             `json:"domain_guide_count,omitempty"`
}

type workflowSearch struct {
	Query     string          `json:"query"`
	Total     int             `json:"total"`
	Truncated bool            `json:"truncated"`
	Workflows []workflowGuide `json:"workflows"`
}

func Workflows() *cobra.Command {
	var format []string
	workflows := &cobra.Command{
		Use:   "workflows",
		Short: "Read task guides for common workflows offline",
		Long: `Read the task guides shipped with this build of the CLI: the CLI-wide agent
guidance ("context"), recipes for common multi-step workflows, and a domain
guide for each command group (filescom-<group>).
Works offline, without credentials, and without reading or writing the config file.

With no subcommand, lists the context guide and the recipes.`,
		Example: `files-cli workflows
files-cli workflows list --domains
files-cli workflows search large folder size
files-cli workflows show recipe-searching-for-files
files-cli workflows show filescom-permissions`,
		Args:              cobra.NoArgs,
		PersistentPreRunE: offlinePreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			return writeWorkflowList(cmd, format, false)
		},
	}
	addDiscoveryFormatFlag(workflows, &format)

	var domains bool
	list := &cobra.Command{
		Use:   "list",
		Short: "List workflow guides",
		Long: `List the context guide and the recipes, or with --domains the domain guides,
one per command group.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return writeWorkflowList(cmd, format, domains)
		},
	}
	list.Flags().BoolVar(&domains, "domains", false, "List the domain guides, one per command group")
	workflows.AddCommand(list)

	workflows.AddCommand(&cobra.Command{
		Use:   "show <name>",
		Short: "Show one workflow guide as Markdown",
		Long: `Show one workflow guide as Markdown. A domain guide ends with the commands of
its group in this build.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			guides, err := loadWorkflowGuides()
			if err != nil {
				return err
			}
			var names []string
			for _, guide := range guides {
				if guide.Name == args[0] {
					if guide.Kind == guideKindDomain {
						guide = withGroupCommands(cmd.Root(), guide)
					}
					return writeDiscovery(cmd, format, guide, func(out io.Writer) { fmt.Fprint(out, guide.source) })
				}
				if guide.Kind != guideKindDomain {
					names = append(names, guide.Name)
				}
			}
			return clierr.Errorf(clierr.ErrorCodeUsage, "unknown workflow %q; available: %s, and the domain guides filescom-<group> (\"%s workflows list --domains\")", args[0], strings.Join(names, ", "), cmd.Root().Name())
		},
	})

	var limit int
	search := &cobra.Command{
		Use:   "search <words...>",
		Short: "Find workflow guides by keywords",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if limit < 0 {
				return clierr.Errorf(clierr.ErrorCodeUsage, "--limit must be 0 or greater")
			}
			return writeWorkflowSearch(cmd, format, args, limit)
		},
	}
	search.Flags().IntVar(&limit, "limit", 5, "Maximum number of results; 0 returns all matches")
	workflows.AddCommand(search)

	markReadOnly(workflows)
	return workflows
}

// loadWorkflowGuides reads the embedded guides. A recipe or domain guide takes
// its name, description, and command links from its front matter, and Content
// holds the Markdown after it.
func loadWorkflowGuides() ([]workflowGuide, error) {
	bundledGuides, err := loadGuides()
	if err != nil {
		return nil, err
	}
	guides := make([]workflowGuide, 0, len(bundledGuides))
	for _, bundled := range bundledGuides {
		guide := parseWorkflowGuide(bundled.Source)
		guide.Kind = bundled.Kind
		if guide.Kind == guideKindContext {
			guide.Name, guide.Description = contextGuideName, contextGuideDescription
		}
		guides = append(guides, guide)
	}
	return guides, nil
}

// parseWorkflowGuide reads the name, description, and metadata from a guide's
// YAML front matter block. It supports the plain values, block scalars (|),
// and metadata map the guides use rather than general YAML.
func parseWorkflowGuide(source string) workflowGuide {
	guide := workflowGuide{Content: source, source: source}
	rest, ok := strings.CutPrefix(source, "---\n")
	if !ok {
		return guide
	}
	frontMatter, body, ok := strings.Cut(rest, "\n---\n")
	if !ok {
		return guide
	}
	guide.Content = strings.TrimLeft(body, "\n")

	lines := strings.Split(frontMatter, "\n")
	for i := 0; i < len(lines); i++ {
		key, value, _ := strings.Cut(lines[i], ":")
		value = strings.TrimSpace(value)
		var nested []string
		for i+1 < len(lines) && strings.HasPrefix(lines[i+1], " ") {
			i++
			nested = append(nested, strings.TrimSpace(lines[i]))
		}
		if value == "|" || value == ">" {
			value = strings.Join(nested, " ")
		}
		switch key {
		case "name":
			guide.Name = value
		case "description":
			guide.Description = value
		case "metadata":
			for _, line := range nested {
				metadataKey, metadataValue, _ := strings.Cut(line, ":")
				switch strings.TrimSpace(metadataKey) {
				case guideMetadataGroup:
					guide.Group = strings.TrimSpace(metadataValue)
				case guideMetadataCommands:
					for _, command := range strings.Split(metadataValue, ",") {
						if command = strings.Join(strings.Fields(command), " "); command != "" {
							guide.Commands = append(guide.Commands, command)
						}
					}
				}
			}
		}
	}
	return guide
}

// withGroupCommands appends the commands of the guide's group, read from the
// command tree, so the guide lists what this build runs.
func withGroupCommands(root *cobra.Command, guide workflowGuide) workflowGuide {
	group, err := lookupDiscoverable(root, []string{guide.Group})
	if err != nil {
		return guide
	}
	var section strings.Builder
	fmt.Fprintf(&section, "\n## Commands\n\nThe %s commands in this build. For arguments, flags, effect, and response fields, run `%s commands describe %s <command>`.\n\n", guide.Group, root.Name(), guide.Group)
	for _, command := range summarizeCommands(discoverableSubcommands(group)) {
		fmt.Fprintf(&section, "- `%s`", strings.Join(append([]string{command.Command}, command.Args...), " "))
		if command.Short != "" {
			fmt.Fprintf(&section, ": %s", oneLine(command.Short))
		}
		if command.Effect != "" {
			fmt.Fprintf(&section, " (%s)", command.Effect)
		}
		section.WriteString("\n")
	}
	guide.Content += section.String()
	guide.source += section.String()
	return guide
}

func summarizeWorkflows(guides []workflowGuide) []workflowGuide {
	summaries := make([]workflowGuide, len(guides))
	for i, guide := range guides {
		summaries[i] = workflowGuide{Name: guide.Name, Description: guide.Description, Kind: guide.Kind, Group: guide.Group, Commands: guide.Commands}
	}
	return summaries
}

// writeWorkflowList lists the context guide and the recipes with a count of
// the domain guides, or with domains only the domain guides.
func writeWorkflowList(cmd *cobra.Command, format []string, domains bool) error {
	guides, err := loadWorkflowGuides()
	if err != nil {
		return err
	}
	var list workflowList
	for _, guide := range guides {
		switch {
		case (guide.Kind == guideKindDomain) == domains:
			list.Workflows = append(list.Workflows, guide)
		case guide.Kind == guideKindDomain:
			list.DomainGuideCount++
		}
	}
	list.Workflows = summarizeWorkflows(list.Workflows)
	root := cmd.Root().Name()
	return writeDiscovery(cmd, format, list, func(out io.Writer) {
		if domains {
			fmt.Fprintf(out, "Domain guides, one per command group. Next: \"%s workflows show <name>\".\n\n", root)
		} else {
			fmt.Fprintf(out, "Workflow guides. Next: \"%s workflows show <name>\".\n\n", root)
		}
		writeWorkflowsText(out, list.Workflows)
		if list.DomainGuideCount > 0 {
			fmt.Fprintf(out, "\n%d domain guides, one per command group and named filescom-<group>, are listed by \"%s workflows list --domains\". \"%s commands describe <command>\" links the guides for a command.\n", list.DomainGuideCount, root, root)
		}
	})
}

func writeWorkflowSearch(cmd *cobra.Command, format []string, args []string, limit int) error {
	guides, err := loadWorkflowGuides()
	if err != nil {
		return err
	}
	terms := searchTerms(args)
	results := make([]searchResult[workflowGuide], len(guides))
	for i, guide := range guides {
		matched, score := searchScore(terms, searchField{guide.Name, 3}, searchField{guide.Description, 2}, searchField{guide.Content, 1})
		results[i] = searchResult[workflowGuide]{item: guide, name: guide.Name, matched: matched, score: score}
	}
	matches, total := rankSearch(results, limit)
	search := workflowSearch{Query: strings.Join(terms, " "), Total: total, Truncated: total > len(matches), Workflows: summarizeWorkflows(matches)}
	return writeDiscovery(cmd, format, search, func(out io.Writer) {
		fmt.Fprintf(out, "%d of %d matches for %q", len(matches), total, search.Query)
		if search.Truncated {
			fmt.Fprint(out, "; use --limit for more")
		}
		fmt.Fprintf(out, ". Next: \"%s workflows show <name>\".\n\n", cmd.Root().Name())
		writeWorkflowsText(out, search.Workflows)
	})
}

func writeWorkflowsText(out io.Writer, guides []workflowGuide) {
	table := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	for _, guide := range guides {
		fmt.Fprintf(table, "  %s\t%s\n", guide.Name, guide.Description)
	}
	table.Flush()
}
