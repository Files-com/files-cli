package cmd

import (
	"fmt"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	sso_strategy "github.com/Files-com/files-sdk-go/v3/ssostrategy"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(SsoStrategies())
}

func SsoStrategies() *cobra.Command {
	SsoStrategies := &cobra.Command{
		Use:  "sso-strategies [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command sso-strategies\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsSsoStrategyList := files_sdk.SsoStrategyListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Sso Strategies",
		Long:    `List Sso Strategies`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsSsoStrategyList
			params.MaxPages = MaxPagesList

			client := sso_strategy.Client{Config: config}
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
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyList) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyList, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatList), fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdList.Flags().StringVar(&paramsSsoStrategyList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsSsoStrategyList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVar(&paramsSsoStrategyList.Action, "action", "", "")
	cmdList.Flags().Int64Var(&paramsSsoStrategyList.Page, "page", 0, "")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	SsoStrategies.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsSsoStrategyFind := files_sdk.SsoStrategyFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Sso Strategy`,
		Long:  `Show Sso Strategy`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := sso_strategy.Client{Config: config}

			var ssoStrategy interface{}
			var err error
			ssoStrategy, err = client.Find(paramsSsoStrategyFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), ssoStrategy, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsSsoStrategyFind.Id, "id", 0, "Sso Strategy ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	SsoStrategies.AddCommand(cmdFind)
	var fieldsSync []string
	var formatSync []string
	usePagerSync := true
	paramsSsoStrategySync := files_sdk.SsoStrategySyncParams{}

	cmdSync := &cobra.Command{
		Use:   "sync",
		Short: `Synchronize provisioning data with the SSO remote server`,
		Long:  `Synchronize provisioning data with the SSO remote server`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := sso_strategy.Client{Config: config}

			var err error
			err = client.Sync(paramsSsoStrategySync, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdSync.Flags().Int64Var(&paramsSsoStrategySync.Id, "id", 0, "Sso Strategy ID.")

	cmdSync.Flags().StringSliceVar(&fieldsSync, "fields", []string{}, "comma separated list of field names")
	cmdSync.Flags().StringSliceVar(&formatSync, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdSync.Flags().BoolVar(&usePagerSync, "use-pager", usePagerSync, "Use $PAGER (.ie less, more, etc)")

	SsoStrategies.AddCommand(cmdSync)
	return SsoStrategies
}
