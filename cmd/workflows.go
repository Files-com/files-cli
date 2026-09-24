package cmd

import (
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
	"text/tabwriter"

	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

// Guides holds the agent guides that package main embeds from this build:
// CONTEXT.md and skills/recipes/<name>/SKILL.md.
var Guides fs.FS

const (
	contextGuideName        = "context"
	contextGuideDescription = "CLI-wide invocation guidance for agents: JSON output, non-interactive use, authentication, global flags, offline discovery, bounded listing with continuation, workspaces, and errors."
)

func init() {
	RootCmd.AddCommand(Workflows())
}

type workflowGuide struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Content     string `json:"content,omitempty"`
	source      string
}

type workflowList struct {
	Workflows []workflowGuide `json:"workflows"`
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
guidance ("context") and recipes for common multi-step workflows.
Works offline, without credentials, and without reading or writing the config file.

With no subcommand, lists the guides.`,
		Example: `files-cli workflows
files-cli workflows search large folder size
files-cli workflows show recipe-searching-for-files`,
		Args:              cobra.NoArgs,
		PersistentPreRunE: offlinePreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			return writeWorkflowList(cmd, format)
		},
	}
	addDiscoveryFormatFlag(workflows, &format)

	workflows.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List workflow guides",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return writeWorkflowList(cmd, format)
		},
	})

	workflows.AddCommand(&cobra.Command{
		Use:   "show <name>",
		Short: "Show one workflow guide as Markdown",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			guides, err := loadWorkflowGuides()
			if err != nil {
				return err
			}
			var names []string
			for _, guide := range guides {
				if guide.Name == args[0] {
					return writeDiscovery(cmd, format, guide, func(out io.Writer) { fmt.Fprint(out, guide.source) })
				}
				names = append(names, guide.Name)
			}
			return clierr.Errorf(clierr.ErrorCodeUsage, "unknown workflow %q; available: %s", args[0], strings.Join(names, ", "))
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

	return workflows
}

// loadWorkflowGuides reads CONTEXT.md and the recipe skills from Guides. The
// recipe name and description come from each SKILL.md front matter, and
// Content holds the Markdown after it.
func loadWorkflowGuides() ([]workflowGuide, error) {
	if Guides == nil {
		return nil, clierr.Errorf(clierr.ErrorCodeFatal, "workflow guides are not included in this build")
	}
	contextGuide, err := fs.ReadFile(Guides, "CONTEXT.md")
	if err != nil {
		return nil, clierr.New(clierr.ErrorCodeFatal, err)
	}
	guides := []workflowGuide{{Name: contextGuideName, Description: contextGuideDescription, Content: string(contextGuide), source: string(contextGuide)}}

	recipes, err := fs.Glob(Guides, "skills/recipes/*/SKILL.md")
	if err != nil {
		return nil, clierr.New(clierr.ErrorCodeFatal, err)
	}
	for _, recipe := range recipes {
		data, err := fs.ReadFile(Guides, recipe)
		if err != nil {
			return nil, clierr.New(clierr.ErrorCodeFatal, err)
		}
		guide := parseWorkflowGuide(string(data))
		if guide.Name == "" {
			guide.Name = path.Base(path.Dir(recipe))
		}
		guides = append(guides, guide)
	}
	return guides, nil
}

// parseWorkflowGuide reads the name and description from a SKILL.md YAML
// front matter block. It supports the plain and block-scalar (|) values the
// recipes use rather than general YAML.
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
		if value == "|" || value == ">" {
			var block []string
			for i+1 < len(lines) && strings.HasPrefix(lines[i+1], " ") {
				i++
				block = append(block, strings.TrimSpace(lines[i]))
			}
			value = strings.Join(block, " ")
		}
		switch key {
		case "name":
			guide.Name = value
		case "description":
			guide.Description = value
		}
	}
	return guide
}

func summarizeWorkflows(guides []workflowGuide) []workflowGuide {
	summaries := make([]workflowGuide, len(guides))
	for i, guide := range guides {
		summaries[i] = workflowGuide{Name: guide.Name, Description: guide.Description}
	}
	return summaries
}

func writeWorkflowList(cmd *cobra.Command, format []string) error {
	guides, err := loadWorkflowGuides()
	if err != nil {
		return err
	}
	list := workflowList{Workflows: summarizeWorkflows(guides)}
	return writeDiscovery(cmd, format, list, func(out io.Writer) {
		fmt.Fprintf(out, "Workflow guides. Next: \"%s workflows show <name>\".\n\n", cmd.Root().Name())
		writeWorkflowsText(out, list.Workflows)
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
