package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	ai_assistant_personality "github.com/Files-com/files-sdk-go/v3/aiassistantpersonality"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(AiAssistantPersonalities())
}

func AiAssistantPersonalities() *cobra.Command {
	AiAssistantPersonalities := &cobra.Command{
		Use:  "ai-assistant-personalities [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command ai-assistant-personalities\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsAiAssistantPersonalityList := files_sdk.AiAssistantPersonalityListParams{}
	var MaxPagesList int64
	var listSortByArgs string
	var listFilterArgs []string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Ai Assistant Personalities",
		Long:    `List Ai Assistant Personalities`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsAiAssistantPersonalityList
			params.MaxPages = MaxPagesList

			parsedListSortBy, parseListSortByErr := lib.ParseAPIListSortFlag("sort-by", listSortByArgs)
			if parseListSortByErr != nil {
				return parseListSortByErr
			}
			if parsedListSortBy != nil {
				params.SortBy = parsedListSortBy
			}
			parsedListFilter, parseListFilterErr := lib.ParseAPIListQueryFlag("filter", listFilterArgs)
			if parseListFilterErr != nil {
				return parseListFilterErr
			}
			if parsedListFilter != nil {
				params.Filter = parsedListFilter
			}

			client := ai_assistant_personality.Client{Config: config}
			it, err := client.List(params, files_sdk.WithContext(ctx))
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger)
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyList) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyList, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatList), fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, "Client-side wildcard filtering, for example field-name=*.jpg or field-name=?ello")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-by", "field=pattern")
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort ai assistant personalities by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")
	cmdList.Flags().StringArrayVar(&listFilterArgs, "filter", []string{}, "Find ai assistant personalities where field exactly matches value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter", "field=value")

	cmdList.Flags().StringVar(&paramsAiAssistantPersonalityList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsAiAssistantPersonalityList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	AiAssistantPersonalities.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsAiAssistantPersonalityFind := files_sdk.AiAssistantPersonalityFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Ai Assistant Personality`,
		Long:  `Show Ai Assistant Personality`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := ai_assistant_personality.Client{Config: config}

			var aiAssistantPersonality interface{}
			var err error
			aiAssistantPersonality, err = client.Find(paramsAiAssistantPersonalityFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), aiAssistantPersonality, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsAiAssistantPersonalityFind.Id, "id", 0, "Ai Assistant Personality ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	AiAssistantPersonalities.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createApplyToAllWorkspaces := true
	createUseByDefault := true
	paramsAiAssistantPersonalityCreate := files_sdk.AiAssistantPersonalityCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Ai Assistant Personality`,
		Long:  `Create Ai Assistant Personality`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := ai_assistant_personality.Client{Config: config}

			if cmd.Flags().Changed("apply-to-all-workspaces") {
				paramsAiAssistantPersonalityCreate.ApplyToAllWorkspaces = flib.Bool(createApplyToAllWorkspaces)
			}
			if cmd.Flags().Changed("use-by-default") {
				paramsAiAssistantPersonalityCreate.UseByDefault = flib.Bool(createUseByDefault)
			}

			var aiAssistantPersonality interface{}
			var err error
			aiAssistantPersonality, err = client.Create(paramsAiAssistantPersonalityCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), aiAssistantPersonality, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().BoolVar(&createApplyToAllWorkspaces, "apply-to-all-workspaces", createApplyToAllWorkspaces, "If true, this default-workspace personality can apply to users in all workspaces.")
	cmdCreate.Flags().StringVar(&paramsAiAssistantPersonalityCreate.SystemPrompt, "system-prompt", "", "System prompt injected into the in-app AI Assistant.")
	cmdCreate.Flags().BoolVar(&createUseByDefault, "use-by-default", createUseByDefault, "Whether this personality is the default personality for the Workspace.")
	cmdCreate.Flags().Int64Var(&paramsAiAssistantPersonalityCreate.WorkspaceId, "workspace-id", 0, "Workspace ID. `0` means the default workspace.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	AiAssistantPersonalities.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateApplyToAllWorkspaces := true
	updateUseByDefault := true
	paramsAiAssistantPersonalityUpdate := files_sdk.AiAssistantPersonalityUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Ai Assistant Personality`,
		Long:  `Update Ai Assistant Personality`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := ai_assistant_personality.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.AiAssistantPersonalityUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsAiAssistantPersonalityUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("apply-to-all-workspaces") {
				mapParams["apply_to_all_workspaces"] = updateApplyToAllWorkspaces
			}
			if cmd.Flags().Changed("system-prompt") {
				lib.FlagUpdate(cmd, "system_prompt", paramsAiAssistantPersonalityUpdate.SystemPrompt, mapParams)
			}
			if cmd.Flags().Changed("use-by-default") {
				mapParams["use_by_default"] = updateUseByDefault
			}
			if cmd.Flags().Changed("workspace-id") {
				lib.FlagUpdate(cmd, "workspace_id", paramsAiAssistantPersonalityUpdate.WorkspaceId, mapParams)
			}

			var aiAssistantPersonality interface{}
			var err error
			aiAssistantPersonality, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), aiAssistantPersonality, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsAiAssistantPersonalityUpdate.Id, "id", 0, "Ai Assistant Personality ID.")
	cmdUpdate.Flags().BoolVar(&updateApplyToAllWorkspaces, "apply-to-all-workspaces", updateApplyToAllWorkspaces, "If true, this default-workspace personality can apply to users in all workspaces.")
	cmdUpdate.Flags().StringVar(&paramsAiAssistantPersonalityUpdate.SystemPrompt, "system-prompt", "", "System prompt injected into the in-app AI Assistant.")
	cmdUpdate.Flags().BoolVar(&updateUseByDefault, "use-by-default", updateUseByDefault, "Whether this personality is the default personality for the Workspace.")
	cmdUpdate.Flags().Int64Var(&paramsAiAssistantPersonalityUpdate.WorkspaceId, "workspace-id", 0, "Workspace ID. `0` means the default workspace.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	AiAssistantPersonalities.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsAiAssistantPersonalityDelete := files_sdk.AiAssistantPersonalityDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Ai Assistant Personality`,
		Long:  `Delete Ai Assistant Personality`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := ai_assistant_personality.Client{Config: config}

			var err error
			err = client.Delete(paramsAiAssistantPersonalityDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsAiAssistantPersonalityDelete.Id, "id", 0, "Ai Assistant Personality ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	AiAssistantPersonalities.AddCommand(cmdDelete)
	return AiAssistantPersonalities
}
