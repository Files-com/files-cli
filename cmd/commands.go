package cmd

import (
	"fmt"
	"io"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// maxFlagDescription bounds each flag description in commands describe unless
// --full or --flag asks for the complete text.
const maxFlagDescription = 200

func init() {
	RootCmd.AddCommand(Commands())
}

type commandSummary struct {
	Command         string   `json:"command"`
	Short           string   `json:"short,omitempty"`
	Aliases         []string `json:"aliases,omitempty"`
	Args            []string `json:"args,omitempty"`
	Effect          string   `json:"effect,omitempty"`
	SubcommandCount int      `json:"subcommand_count,omitempty"`
}

type commandList struct {
	Group    string           `json:"group,omitempty"`
	Commands []commandSummary `json:"commands"`
}

type commandSearch struct {
	Query     string           `json:"query"`
	Total     int              `json:"total"`
	Truncated bool             `json:"truncated"`
	Commands  []commandSummary `json:"commands"`
}

type commandDescription struct {
	Command        string            `json:"command"`
	Usage          string            `json:"usage"`
	Short          string            `json:"short,omitempty"`
	Long           string            `json:"long,omitempty"`
	Aliases        []string          `json:"aliases,omitempty"`
	Args           []string          `json:"args,omitempty"`
	Effect         string            `json:"effect,omitempty"`
	Example        string            `json:"example,omitempty"`
	Subcommands    []commandSummary  `json:"subcommands,omitempty"`
	Flags          []flagDescription `json:"flags,omitempty"`
	InheritedFlags []flagDescription `json:"inherited_flags,omitempty"`
	Response       *commandResponse  `json:"response,omitempty"`
	Workflows      []workflowLink    `json:"workflows,omitempty"`
}

type flagDescription struct {
	Name          string   `json:"name"`
	Shorthand     string   `json:"shorthand,omitempty"`
	Type          string   `json:"type"`
	Syntax        string   `json:"syntax,omitempty"`
	Default       string   `json:"default,omitempty"`
	OptionalValue string   `json:"optional_value,omitempty"`
	Required      bool     `json:"required,omitempty"`
	APIRequired   bool     `json:"api_required,omitempty"`
	Enum          []string `json:"enum,omitempty"`
	Description   string   `json:"description"`
	Truncated     bool     `json:"truncated,omitempty"`
}

// responseField is a top-level field of the records a command prints, as
// --fields and the JSON output name it.
type responseField struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

// responseEntity is the generated description of an API entity's records.
type responseEntity struct {
	Recommended []string        `json:"recommended"`
	Fields      []responseField `json:"fields"`
}

// responseDescription describes what a command prints on success: records of
// an API entity, or nothing.
type responseDescription struct {
	Type              string          `json:"type,omitempty"`
	List              bool            `json:"list,omitempty"`
	None              bool            `json:"none,omitempty"`
	RecommendedFields []string        `json:"recommended_fields,omitempty"`
	Fields            []responseField `json:"fields,omitempty"`
}

type commandResponse struct {
	responseDescription
	WithFlag *flagResponse `json:"with_flag,omitempty"`
}

// flagResponse is what a command can print instead when the flag is given.
type flagResponse struct {
	Flag string `json:"flag"`
	responseDescription
}

// workflowLink is a guide related to a command.
type workflowLink struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
}

func Commands() *cobra.Command {
	var format []string
	commands := &cobra.Command{
		Use:   "commands",
		Short: "Discover commands, arguments, and flags offline",
		Long: `Discover commands, arguments, and flags from this build of the CLI.
Works offline, without credentials, and without reading or writing the config file.

With no subcommand, lists the top-level commands and command groups.`,
		Example: `files-cli commands
files-cli commands list folders
files-cli commands search share link
files-cli commands describe folders list-for --format json
files-cli commands describe users list --flag cursor
files-cli commands describe remote-servers list --response-field server-type`,
		Args:              cobra.NoArgs,
		PersistentPreRunE: offlinePreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			return writeCommandList(cmd, format, cmd.Root())
		},
	}
	addDiscoveryFormatFlag(commands, &format)

	commands.AddCommand(&cobra.Command{
		Use:   "list [group]",
		Short: "List the top-level commands, or the commands in a group",
		RunE: func(cmd *cobra.Command, args []string) error {
			group, err := lookupDiscoverable(cmd.Root(), args)
			if err != nil {
				return err
			}
			if !group.HasAvailableSubCommands() {
				return clierr.Errorf(clierr.ErrorCodeUsage, "%q has no subcommands; run \"%s commands describe %s\"", commandName(group), cmd.Root().Name(), commandName(group))
			}
			return writeCommandList(cmd, format, group)
		},
	})

	var limit int
	search := &cobra.Command{
		Use:   "search <words...>",
		Short: "Find commands by keywords in their names, descriptions, aliases, and flag names",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if limit < 0 {
				return clierr.Errorf(clierr.ErrorCodeUsage, "--limit must be 0 or greater")
			}
			return writeCommandSearch(cmd, format, args, limit)
		},
	}
	search.Flags().IntVar(&limit, "limit", 20, "Maximum number of results; 0 returns all matches")
	commands.AddCommand(search)

	var full bool
	var onlyFlags []string
	var onlyFields []string
	describe := &cobra.Command{
		Use:   "describe <command...>",
		Short: "Describe a command's usage, arguments, flags, effect, response fields, and related guides",
		Long: `Describe a command's usage, arguments, subcommands, and local and inherited flags,
the effect of its API request, the fields of the records it prints, and related
workflow guides.

Flag types and defaults are the static definitions, never values from the current
invocation or config. A default is the CLI flag's, not necessarily the API's:
generated flags for optional API booleans, such as files copy --overwrite and
--structure, are sent only when given, so leaving one out leaves the choice to the
API. Pass true or false explicitly for the behavior you intend. required marks
flags the CLI rejects the command without.
api_required marks parameters the Files.com API requires; the CLI sends the request
without checking them. Descriptions longer than 200 characters are truncated and
marked; use --full, or --flag for specific flags, to see complete descriptions.

effect is a hint about the command's Files.com API request, from the API operation
or a reviewed handwritten command; it is not enforced and does not cover local
effects such as saved sessions or config. read_only only reads data. destructive
can delete, overwrite, or replace existing data. mutating changes data or starts
an action and is not known to be destructive, which does not make it additive or
reversible. Commands without an effect have none declared.

response lists the records a command prints on success and their top-level
fields, which --fields selects, with any recommended_fields; none means it prints
nothing. Add --full, or --response-field for specific fields, to see their
descriptions. Commands without a response have none declared.

--flag and --response-field show only the named flags and response fields.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := lookupDiscoverable(cmd.Root(), args)
			if err != nil {
				return err
			}
			description, err := describeCommand(target, full, onlyFlags, onlyFields)
			if err != nil {
				return err
			}
			if description.Workflows, err = relatedWorkflows(target); err != nil {
				return err
			}
			root := cmd.Root().Name()
			return writeDiscovery(cmd, format, description, func(out io.Writer) { writeDescriptionText(out, root, description) })
		},
	}
	describe.Flags().BoolVar(&full, "full", false, "Show complete flag and response field descriptions")
	describe.Flags().StringSliceVar(&onlyFlags, "flag", []string{}, "Show only these flags, with complete descriptions")
	describe.Flags().StringSliceVar(&onlyFields, "response-field", []string{}, "Show only these response fields, with descriptions")
	commands.AddCommand(describe)

	markReadOnly(commands)
	return commands
}

// discoverable reports whether command discovery shows cmd: available (not
// hidden, deprecated, or the help command) and not a resource group without
// any subcommands.
func discoverable(cmd *cobra.Command) bool {
	if !cmd.IsAvailableCommand() {
		return false
	}
	return cmd.HasAvailableSubCommands() || !strings.HasSuffix(cmd.Use, " [command]")
}

func discoverableSubcommands(cmd *cobra.Command) []*cobra.Command {
	var subcommands []*cobra.Command
	for _, subcommand := range cmd.Commands() {
		if discoverable(subcommand) {
			subcommands = append(subcommands, subcommand)
		}
	}
	return subcommands
}

// lookupDiscoverable resolves command words, given as separate arguments or
// one quoted string, by name or alias. A leading root command name is ignored
// so a copied usage line resolves too.
func lookupDiscoverable(root *cobra.Command, args []string) (*cobra.Command, error) {
	words := strings.Fields(strings.Join(args, " "))
	if len(words) > 0 && words[0] == root.Name() {
		words = words[1:]
	}
	current := root
	for _, word := range words {
		var next *cobra.Command
		for _, subcommand := range discoverableSubcommands(current) {
			if subcommand.Name() == word || subcommand.HasAlias(word) {
				next = subcommand
				break
			}
		}
		if next == nil {
			return nil, clierr.Errorf(clierr.ErrorCodeUsage, "unknown command %q; run \"%s commands search %s\"", strings.Join(words, " "), root.Name(), word)
		}
		current = next
	}
	return current, nil
}

// commandName is the command path below the root, as accepted by describe.
func commandName(cmd *cobra.Command) string {
	return strings.TrimPrefix(strings.TrimPrefix(cmd.CommandPath(), cmd.Root().Name()), " ")
}

// commandArgs returns the positional argument placeholders from a leaf
// command's Use line.
func commandArgs(cmd *cobra.Command) []string {
	if cmd.HasAvailableSubCommands() {
		return nil
	}
	return strings.Fields(cmd.Use)[1:]
}

func commandAliases(cmd *cobra.Command) []string {
	var aliases []string
	for _, alias := range cmd.Aliases {
		if alias != cmd.Name() {
			aliases = append(aliases, alias)
		}
	}
	return aliases
}

func summarizeCommand(cmd *cobra.Command) commandSummary {
	return commandSummary{
		Command:         commandName(cmd),
		Short:           cmd.Short,
		Aliases:         commandAliases(cmd),
		Args:            commandArgs(cmd),
		Effect:          commandEffect(cmd),
		SubcommandCount: len(discoverableSubcommands(cmd)),
	}
}

// commandEffect returns the effect hint recorded for cmd, or "" when it has
// none.
func commandEffect(cmd *cobra.Command) string {
	effect, _ := lib.GetCommandEffect(cmd)
	return string(effect)
}

// markReadOnly records that cmd and its subcommands only read data. Used for
// the discovery commands, which read the command tree and embedded guides.
func markReadOnly(cmd *cobra.Command) {
	lib.SetCommandEffect(cmd, lib.EffectReadOnly)
	for _, subcommand := range cmd.Commands() {
		markReadOnly(subcommand)
	}
}

func summarizeCommands(commands []*cobra.Command) []commandSummary {
	summaries := make([]commandSummary, len(commands))
	for i, command := range commands {
		summaries[i] = summarizeCommand(command)
	}
	return summaries
}

func writeCommandList(cmd *cobra.Command, format []string, group *cobra.Command) error {
	list := commandList{Group: commandName(group), Commands: summarizeCommands(discoverableSubcommands(group))}
	root := cmd.Root().Name()
	return writeDiscovery(cmd, format, list, func(out io.Writer) {
		if list.Group == "" {
			fmt.Fprintf(out, "Commands and command groups. Next: \"%s commands list <group>\", \"%s commands describe <command>\", \"%s commands search <words>\", \"%s workflows list\".\n\n", root, root, root, root)
		} else {
			fmt.Fprintf(out, "Commands in %s. Next: \"%s commands describe %s <command>\".\n\n", list.Group, root, list.Group)
		}
		writeSummariesText(out, list.Commands)
	})
}

func writeCommandSearch(cmd *cobra.Command, format []string, args []string, limit int) error {
	terms := searchTerms(args)
	var results []searchResult[commandSummary]
	var visit func(parent *cobra.Command)
	visit = func(parent *cobra.Command) {
		for _, command := range discoverableSubcommands(parent) {
			name := commandName(command)
			var flagNames []string
			command.LocalFlags().VisitAll(func(flag *pflag.Flag) {
				if describableFlag(flag) {
					flagNames = append(flagNames, "--"+flag.Name)
				}
			})
			matched, score := searchScore(terms,
				searchField{name, 3},
				searchField{strings.Join(command.Aliases, " "), 3},
				searchField{command.Short, 2},
				searchField{command.Long, 1},
				searchField{strings.Join(flagNames, " "), 1},
			)
			results = append(results, searchResult[commandSummary]{item: summarizeCommand(command), name: name, matched: matched, score: score})
			visit(command)
		}
	}
	visit(cmd.Root())

	matches, total := rankSearch(results, limit)
	search := commandSearch{Query: strings.Join(terms, " "), Total: total, Truncated: total > len(matches), Commands: matches}
	return writeDiscovery(cmd, format, search, func(out io.Writer) {
		fmt.Fprintf(out, "%d of %d matches for %q", len(matches), total, search.Query)
		if search.Truncated {
			fmt.Fprint(out, "; use --limit for more")
		}
		fmt.Fprint(out, ".\n\n")
		writeSummariesText(out, matches)
	})
}

func writeSummariesText(out io.Writer, summaries []commandSummary) {
	table := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	for _, summary := range summaries {
		name := strings.Join(append([]string{summary.Command}, summary.Args...), " ")
		switch summary.SubcommandCount {
		case 0:
		case 1:
			name += " (1 command)"
		default:
			name += fmt.Sprintf(" (%d commands)", summary.SubcommandCount)
		}
		fmt.Fprintf(table, "  %s\t%s\n", name, summary.Short)
	}
	table.Flush()
}

// describableFlag excludes flags help would not show and the help flag,
// which Cobra adds only to the command being executed.
func describableFlag(flag *pflag.Flag) bool {
	return !flag.Hidden && flag.Deprecated == "" && flag.Name != "help"
}

func describeCommand(cmd *cobra.Command, full bool, onlyFlags []string, onlyFields []string) (commandDescription, error) {
	description := commandDescription{
		Command:     commandName(cmd),
		Usage:       cmd.UseLine(),
		Short:       cmd.Short,
		Aliases:     commandAliases(cmd),
		Args:        commandArgs(cmd),
		Effect:      commandEffect(cmd),
		Example:     cmd.Example,
		Subcommands: summarizeCommands(discoverableSubcommands(cmd)),
	}
	if strings.TrimSpace(cmd.Long) != strings.TrimSpace(cmd.Short) {
		description.Long = cmd.Long
	}

	selected := make(map[string]bool)
	for _, name := range onlyFlags {
		selected[strings.TrimLeft(name, "-")] = false
	}
	describeFlags := func(flags *pflag.FlagSet) []flagDescription {
		var descriptions []flagDescription
		flags.VisitAll(func(flag *pflag.Flag) {
			if !describableFlag(flag) {
				return
			}
			if len(selected) > 0 {
				if _, ok := selected[flag.Name]; !ok {
					return
				}
				selected[flag.Name] = true
			}
			descriptions = append(descriptions, describeFlag(flag, full || len(selected) > 0))
		})
		return descriptions
	}
	// --response-field alone selects no flags.
	if len(onlyFields) == 0 || len(onlyFlags) > 0 {
		description.Flags = describeFlags(cmd.LocalFlags())
		description.InheritedFlags = describeFlags(cmd.InheritedFlags())
	}

	for name, found := range selected {
		if !found {
			return description, clierr.Errorf(clierr.ErrorCodeUsage, "unknown flag --%s for %q", name, description.Command)
		}
	}

	var err error
	description.Response, err = describeResponse(cmd, full, onlyFlags, onlyFields)
	return description, err
}

// describeResponse describes what cmd prints, from the response recorded on
// the command and the generated entity fields. Field descriptions are included
// with full or when fields are selected; --flag alone selects none.
func describeResponse(cmd *cobra.Command, full bool, onlyFlags []string, onlyFields []string) (*commandResponse, error) {
	entity, list, ok := lib.GetCommandResponse(cmd)
	if !ok {
		if len(onlyFields) > 0 {
			return nil, clierr.Errorf(clierr.ErrorCodeUsage, "%q has no described response fields", commandName(cmd))
		}
		return nil, nil
	}
	entities, err := loadResponseEntities()
	if err != nil {
		return nil, err
	}

	selected := make(map[string]bool)
	for _, name := range onlyFields {
		selected[strings.ToLower(strings.ReplaceAll(name, "-", "_"))] = false
	}
	showFields := len(onlyFields) > 0 || len(onlyFlags) == 0
	describe := func(entity string, list bool) responseDescription {
		if entity == "" {
			return responseDescription{None: true}
		}
		description := responseDescription{Type: entity, List: list, RecommendedFields: entities[entity].Recommended}
		if !showFields {
			return description
		}
		for _, field := range entities[entity].Fields {
			if len(selected) > 0 {
				if _, ok := selected[field.Name]; !ok {
					continue
				}
				selected[field.Name] = true
			} else if !full {
				field.Description = ""
			}
			description.Fields = append(description.Fields, field)
		}
		return description
	}

	response := &commandResponse{responseDescription: describe(entity, list)}
	if flag, entity, list, ok := lib.GetCommandResponseWithFlag(cmd); ok {
		response.WithFlag = &flagResponse{Flag: flag, responseDescription: describe(entity, list)}
	}
	for name, found := range selected {
		if !found {
			return nil, clierr.Errorf(clierr.ErrorCodeUsage, "unknown response field %q for %q", name, commandName(cmd))
		}
	}
	return response, nil
}

// relatedWorkflows returns the guides for cmd: the domain guide of its command
// group, then the recipes that list the command.
func relatedWorkflows(cmd *cobra.Command) ([]workflowLink, error) {
	guides, err := loadWorkflowGuides()
	if err != nil {
		return nil, err
	}
	name := commandName(cmd)
	group, _, _ := strings.Cut(name, " ")
	var domains, recipes []workflowLink
	for _, guide := range guides {
		link := workflowLink{Name: guide.Name, Kind: guide.Kind, Description: guide.Description}
		if guide.Group != "" && guide.Group == group {
			domains = append(domains, link)
		} else if slices.Contains(guide.Commands, name) {
			recipes = append(recipes, link)
		}
	}
	return append(domains, recipes...), nil
}

// describeFlag reports the flag's static definition. DefValue is the default
// given when the flag was defined; the parsed Value, which may hold an API key
// or other input from this invocation, is never read.
func describeFlag(flag *pflag.Flag, full bool) flagDescription {
	valueType, syntax := lib.FlagTypes(flag)
	description := flagDescription{
		Name:        flag.Name,
		Shorthand:   flag.Shorthand,
		Type:        valueType,
		Syntax:      syntax,
		Required:    lib.FlagRequired(flag),
		APIRequired: lib.FlagAPIRequired(flag),
		Enum:        lib.FlagEnum(flag),
		Description: flag.Usage,
	}
	if flag.DefValue != "" && flag.DefValue != "[]" {
		description.Default = flag.DefValue
	}
	if valueType != "bool" {
		description.OptionalValue = flag.NoOptDefVal
	}
	if !full {
		description.Description, description.Truncated = truncateText(flag.Usage, maxFlagDescription)
	}
	return description
}

func writeDescriptionText(out io.Writer, root string, description commandDescription) {
	fmt.Fprintf(out, "%s\n", description.Usage)
	if description.Short != "" {
		fmt.Fprintf(out, "  %s\n", description.Short)
	}
	if description.Long != "" {
		fmt.Fprintf(out, "\n%s\n", description.Long)
	}
	if effect := effectText(description.Effect); effect != "" {
		fmt.Fprintf(out, "\nEffect (a hint, not enforced): %s\n", effect)
	} else if len(description.Subcommands) == 0 {
		fmt.Fprint(out, "\nEffect: not declared for this command.\n")
	}
	if len(description.Aliases) > 0 {
		fmt.Fprintf(out, "\nAliases: %s\n", strings.Join(description.Aliases, ", "))
	}
	if description.Example != "" {
		fmt.Fprintf(out, "\nExamples:\n%s\n", strings.TrimRight(description.Example, "\n"))
	}
	if len(description.Subcommands) > 0 {
		fmt.Fprint(out, "\nCommands:\n")
		writeSummariesText(out, description.Subcommands)
	}

	truncated := false
	writeFlags := func(title string, flags []flagDescription) {
		if len(flags) == 0 {
			return
		}
		fmt.Fprintf(out, "\n%s:\n", title)
		for _, flag := range flags {
			name := "--" + flag.Name
			if flag.Shorthand != "" {
				name = "-" + flag.Shorthand + ", " + name
			}
			var attributes []string
			if flag.Syntax != "" {
				attributes = append(attributes, "syntax: "+flag.Syntax)
			}
			if flag.Required {
				attributes = append(attributes, "required")
			} else if flag.APIRequired {
				attributes = append(attributes, "required by the API, not checked by the CLI")
			}
			if flag.Default != "" {
				attributes = append(attributes, "default: "+flag.Default)
			}
			if flag.OptionalValue != "" {
				attributes = append(attributes, "value optional, default when omitted: "+flag.OptionalValue)
			}
			if len(flag.Enum) > 0 {
				attributes = append(attributes, "one of: "+strings.Join(flag.Enum, ", "))
			}
			fmt.Fprintf(out, "  %s %s", name, flag.Type)
			if len(attributes) > 0 {
				fmt.Fprintf(out, " (%s)", strings.Join(attributes, "; "))
			}
			fmt.Fprintf(out, "\n      %s\n", oneLine(flag.Description))
			truncated = truncated || flag.Truncated
		}
	}
	writeFlags("Flags", description.Flags)
	writeFlags("Inherited flags", description.InheritedFlags)
	if truncated {
		fmt.Fprint(out, "\nDescriptions ending in … are truncated; use --full or --flag NAME for the complete text.\n")
	}

	if description.Response != nil {
		writeResponseText(out, description.Response)
	}
	if len(description.Workflows) > 0 {
		fmt.Fprintf(out, "\nWorkflow guides; read one with \"%s workflows show <name>\":\n", root)
		table := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
		for _, link := range description.Workflows {
			summary, _ := truncateText(oneLine(link.Description), maxFlagDescription)
			fmt.Fprintf(table, "  %s\t%s\n", link.Name, summary)
		}
		table.Flush()
	}
}

func effectText(effect string) string {
	switch lib.CommandEffect(effect) {
	case lib.EffectReadOnly:
		return "read_only, reads data without changing it."
	case lib.EffectMutating:
		return "mutating, changes data or starts an action. It is not known to be destructive, which does not make it additive or reversible."
	case lib.EffectDestructive:
		return "destructive, can delete, overwrite, or replace existing data."
	}
	return ""
}

func writeResponseText(out io.Writer, response *commandResponse) {
	fmt.Fprintf(out, "\nResponse: %s.\n", responseSummary(response.responseDescription))
	fields := response.Fields
	writeResponseFields(out, response.responseDescription)
	if response.WithFlag != nil {
		fmt.Fprintf(out, "With --%s, it can print %s instead.\n", response.WithFlag.Flag, responseSummary(response.WithFlag.responseDescription))
		writeResponseFields(out, response.WithFlag.responseDescription)
		fields = append(fields, response.WithFlag.Fields...)
	}
	if len(fields) > 0 && !slices.ContainsFunc(fields, func(field responseField) bool { return field.Description != "" }) {
		fmt.Fprint(out, "Add --full, or --response-field NAME, for field descriptions.\n")
	}
}

func responseSummary(response responseDescription) string {
	switch {
	case response.None:
		return "prints nothing on success"
	case response.List:
		return fmt.Sprintf("a list of %s records", response.Type)
	}
	return fmt.Sprintf("one %s record", response.Type)
}

// writeResponseFields lists the fields one per line when they carry
// descriptions, otherwise as a wrapped "name type" list.
func writeResponseFields(out io.Writer, response responseDescription) {
	if len(response.RecommendedFields) > 0 {
		fmt.Fprintf(out, "  Recommended: --fields=%s\n", strings.Join(response.RecommendedFields, ","))
	}
	if len(response.Fields) == 0 {
		return
	}
	fmt.Fprint(out, "  Fields, selected with --fields=NAME,...:\n")
	if slices.ContainsFunc(response.Fields, func(field responseField) bool { return field.Description != "" }) {
		for _, field := range response.Fields {
			fmt.Fprintf(out, "    %s %s\n", field.Name, field.Type)
			if field.Description != "" {
				fmt.Fprintf(out, "        %s\n", oneLine(field.Description))
			}
		}
		return
	}
	items := make([]string, len(response.Fields))
	for i, field := range response.Fields {
		items[i] = field.Name + " " + field.Type
	}
	writeWrapped(out, "    ", items, 100)
}

// writeWrapped writes items separated by commas, wrapping lines before width.
func writeWrapped(out io.Writer, indent string, items []string, width int) {
	line := indent
	for i, item := range items {
		if i < len(items)-1 {
			item += ","
		}
		if len(line) > len(indent) {
			if len(line)+1+len(item) > width {
				fmt.Fprintln(out, line)
				line = indent
			} else {
				line += " "
			}
		}
		line += item
	}
	fmt.Fprintln(out, line)
}
