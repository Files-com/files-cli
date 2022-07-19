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
	usePagerList := true
	paramsSsoStrategyList := files_sdk.SsoStrategyListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List Sso Strategies",
		Long:  `List Sso Strategies`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsSsoStrategyList
			params.MaxPages = MaxPagesList

			client := sso_strategy.Client{Config: *config}
			it, err := client.List(ctx, params)
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, formatList, config.Logger())
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}

	cmdList.Flags().StringVar(&paramsSsoStrategyList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsSsoStrategyList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	SsoStrategies.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	usePagerFind := true
	paramsSsoStrategyFind := files_sdk.SsoStrategyFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Sso Strategy`,
		Long:  `Show Sso Strategy`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := sso_strategy.Client{Config: *config}

			var ssoStrategy interface{}
			var err error
			ssoStrategy, err = client.Find(ctx, paramsSsoStrategyFind)
			lib.HandleResponse(ctx, ssoStrategy, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdFind.Flags().Int64Var(&paramsSsoStrategyFind.Id, "id", 0, "Sso Strategy ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	SsoStrategies.AddCommand(cmdFind)
	var fieldsSync string
	var formatSync string
	usePagerSync := true
	paramsSsoStrategySync := files_sdk.SsoStrategySyncParams{}

	cmdSync := &cobra.Command{
		Use:   "sync",
		Short: `Synchronize provisioning data with the SSO remote server`,
		Long:  `Synchronize provisioning data with the SSO remote server`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := sso_strategy.Client{Config: *config}

			var err error
			err = client.Sync(ctx, paramsSsoStrategySync)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}
	cmdSync.Flags().Int64Var(&paramsSsoStrategySync.Id, "id", 0, "Sso Strategy ID.")

	cmdSync.Flags().StringVarP(&fieldsSync, "fields", "", "", "comma separated list of field names")
	cmdSync.Flags().StringVarP(&formatSync, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdSync.Flags().BoolVar(&usePagerSync, "use-pager", usePagerSync, "Use $PAGER (.ie less, more, etc)")

	SsoStrategies.AddCommand(cmdSync)
}
