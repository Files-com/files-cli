package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	shared_workflow "github.com/Files-com/files-sdk-go/v3/sharedworkflow"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(SharedWorkflows())
}

func SharedWorkflows() *cobra.Command {
	SharedWorkflows := &cobra.Command{
		Use:   "shared-workflows [command]",
		Short: "",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command shared-workflows\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsSharedWorkflowList := files_sdk.SharedWorkflowListParams{}
	var MaxPagesList int64
	var jsonEnvelopeList bool
	var listSortByArgs string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Shared Workflows",
		Long:    `List Shared Workflows`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsSharedWorkflowList
			params.MaxPages = MaxPagesList
			var envelopeStyle string
			if jsonEnvelopeList {
				var envelopeErr error
				if envelopeStyle, envelopeErr = lib.PrepareJSONEnvelope(cmd, Profile(cmd).Current().SetResourceFormat(cmd, formatList), &params.MaxPages); envelopeErr != nil {
					return envelopeErr
				}
			}

			parsedListSortBy, parseListSortByErr := lib.ParseAPIListSortFlag("sort-by", listSortByArgs)
			if parseListSortByErr != nil {
				return parseListSortByErr
			}
			if parsedListSortBy != nil {
				params.SortBy = parsedListSortBy
			}

			client := shared_workflow.Client{Config: config}
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
			if jsonEnvelopeList {
				err = lib.JSONEnvelopeIter(it, fieldsList, listFilter, usePagerList, envelopeStyle, cmd.OutOrStdout())
			} else {
				err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatList), fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			}
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, "Client-side wildcard filtering, for example field-name=*.jpg or field-name=?ello")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-by", "field=pattern")
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort shared workflows by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")

	cmdList.Flags().StringVar(&paramsSharedWorkflowList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsSharedWorkflowList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	cmdList.Flags().BoolVar(&jsonEnvelopeList, "json-envelope", false, lib.JSONEnvelopeHelpText)
	lib.SetCommandEffect(cmdList, lib.EffectReadOnly)
	lib.SetCommandResponse(cmdList, "SharedWorkflow", true)
	SharedWorkflows.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsSharedWorkflowFind := files_sdk.SharedWorkflowFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Shared Workflow`,
		Long:  `Show Shared Workflow`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := shared_workflow.Client{Config: config}

			var sharedWorkflow interface{}
			var err error
			sharedWorkflow, err = client.Find(paramsSharedWorkflowFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), sharedWorkflow, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsSharedWorkflowFind.Id, "id", 0, "Shared Workflow ID.")
	lib.SetFlagAPIRequired(cmdFind.Flags(), "id")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	lib.SetCommandEffect(cmdFind, lib.EffectReadOnly)
	lib.SetCommandResponse(cmdFind, "SharedWorkflow", false)
	SharedWorkflows.AddCommand(cmdFind)
	var fieldsFindRun []string
	var formatFindRun []string
	usePagerFindRun := true
	paramsSharedWorkflowFindRun := files_sdk.SharedWorkflowFindRunParams{}

	cmdFindRun := &cobra.Command{
		Use:   "find-run",
		Short: `Show a Shared Workflow run started by the recipient`,
		Long:  `Show a Shared Workflow run started by the recipient`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := shared_workflow.Client{Config: config}

			var automationRun interface{}
			var err error
			automationRun, err = client.FindRun(paramsSharedWorkflowFindRun, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), automationRun, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFindRun), fieldsFindRun, usePagerFindRun, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFindRun.Flags().Int64Var(&paramsSharedWorkflowFindRun.Id, "id", 0, "Shared Workflow ID.")
	lib.SetFlagAPIRequired(cmdFindRun.Flags(), "id")
	cmdFindRun.Flags().Int64Var(&paramsSharedWorkflowFindRun.RunId, "run-id", 0, "ID of the run returned by the Shared Workflow submission.")
	lib.SetFlagAPIRequired(cmdFindRun.Flags(), "run-id")

	cmdFindRun.Flags().StringSliceVar(&fieldsFindRun, "fields", []string{}, "comma separated list of field names")
	cmdFindRun.Flags().StringSliceVar(&formatFindRun, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFindRun.Flags().BoolVar(&usePagerFindRun, "use-pager", usePagerFindRun, "Use $PAGER (.ie less, more, etc)")

	lib.SetCommandEffect(cmdFindRun, lib.EffectReadOnly)
	lib.SetCommandResponse(cmdFindRun, "AutomationRun", false)
	SharedWorkflows.AddCommand(cmdFindRun)
	var fieldsRunWorkflow []string
	var formatRunWorkflow []string
	usePagerRunWorkflow := true
	paramsSharedWorkflowRunWorkflow := files_sdk.SharedWorkflowRunWorkflowParams{}

	cmdRunWorkflow := &cobra.Command{
		Use:   "run-workflow",
		Short: `Run a Shared Workflow`,
		Long:  `Run a Shared Workflow`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := shared_workflow.Client{Config: config}

			var automationRun interface{}
			var err error
			automationRun, err = client.RunWorkflow(paramsSharedWorkflowRunWorkflow, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), automationRun, err, Profile(cmd).Current().SetResourceFormat(cmd, formatRunWorkflow), fieldsRunWorkflow, usePagerRunWorkflow, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdRunWorkflow.Flags().Int64Var(&paramsSharedWorkflowRunWorkflow.Id, "id", 0, "Shared Workflow ID.")
	lib.SetFlagAPIRequired(cmdRunWorkflow.Flags(), "id")
	cmdRunWorkflow.Flags().StringVar(&paramsSharedWorkflowRunWorkflow.SubmissionRevision, "submission-revision", "", "Revision returned when opening the Shared Workflow.")
	lib.SetFlagAPIRequired(cmdRunWorkflow.Flags(), "submission-revision")
	cmdRunWorkflow.Flags().StringVar(&paramsSharedWorkflowRunWorkflow.SubmissionId, "submission-id", "", "Client-generated UUID used to retry the same submission without starting another run.")
	lib.SetFlagAPIRequired(cmdRunWorkflow.Flags(), "submission-id")
	cmdRunWorkflow.Flags().StringVar(&paramsSharedWorkflowRunWorkflow.FilePath, "file-path", "", "Path of a file the recipient can read, including a file uploaded to a Snapshot created for this workflow.")
	lib.SetFlagAPIRequired(cmdRunWorkflow.Flags(), "file-path")

	cmdRunWorkflow.Flags().StringSliceVar(&fieldsRunWorkflow, "fields", []string{}, "comma separated list of field names")
	cmdRunWorkflow.Flags().StringSliceVar(&formatRunWorkflow, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdRunWorkflow.Flags().BoolVar(&usePagerRunWorkflow, "use-pager", usePagerRunWorkflow, "Use $PAGER (.ie less, more, etc)")

	lib.SetCommandEffect(cmdRunWorkflow, lib.EffectMutating)
	lib.SetCommandResponse(cmdRunWorkflow, "AutomationRun", false)
	SharedWorkflows.AddCommand(cmdRunWorkflow)
	return SharedWorkflows
}
