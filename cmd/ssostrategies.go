package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	sso_strategy "github.com/Files-com/files-sdk-go/v2/ssostrategy"
)

var (
	SsoStrategies = &cobra.Command{}
)

func SsoStrategiesInit() {
	SsoStrategies = &cobra.Command{
		Use:  "sso-strategies [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command sso-strategies\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsSsoStrategyList := files_sdk.SsoStrategyListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsSsoStrategyList
			params.MaxPages = MaxPagesList

			client := sso_strategy.Client{Config: *config}
			it, err := client.List(ctx, params)
			if err != nil {
				lib.ClientError(ctx, err)
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(it, formatList, fieldsList, listFilter)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}

	cmdList.Flags().StringVar(&paramsSsoStrategyList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64Var(&paramsSsoStrategyList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	SsoStrategies.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	paramsSsoStrategyFind := files_sdk.SsoStrategyFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := sso_strategy.Client{Config: *config}

			result, err := client.Find(ctx, paramsSsoStrategyFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatFind, fieldsFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdFind.Flags().Int64Var(&paramsSsoStrategyFind.Id, "id", 0, "Sso Strategy ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright")
	SsoStrategies.AddCommand(cmdFind)
}
