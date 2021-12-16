package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	automation_run "github.com/Files-com/files-sdk-go/v2/automationrun"
)

var (
	AutomationRuns = &cobra.Command{}
)

func AutomationRunsInit() {
	AutomationRuns = &cobra.Command{
		Use:  "automation-runs [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command automation-runs\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsAutomationRunList := files_sdk.AutomationRunListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsAutomationRunList
			params.MaxPages = MaxPagesList

			client := automation_run.Client{Config: *config}
			it, err := client.List(ctx, params)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(it, formatList, fieldsList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}

	cmdList.Flags().Int64Var(&paramsAutomationRunList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsAutomationRunList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsAutomationRunList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64Var(&paramsAutomationRunList.AutomationId, "automation-id", 0, "ID of the associated Automation.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	AutomationRuns.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	paramsAutomationRunFind := files_sdk.AutomationRunFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := automation_run.Client{Config: *config}

			var result interface{}
			var err error
			result, err = client.Find(ctx, paramsAutomationRunFind)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatFind, fieldsFind, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdFind.Flags().Int64Var(&paramsAutomationRunFind.Id, "id", 0, "Automation Run ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright")
	AutomationRuns.AddCommand(cmdFind)
}
