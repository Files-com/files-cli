package cmd

import (
	"fmt"
	"time"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/restore"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Restores())
}

func Restores() *cobra.Command {
	Restores := &cobra.Command{
		Use:  "restores [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command restores\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsRestoreList := files_sdk.RestoreListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Restores",
		Long:    `List Restores`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsRestoreList
			params.MaxPages = MaxPagesList

			client := restore.Client{Config: config}
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

	cmdList.Flags().StringVar(&paramsRestoreList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsRestoreList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Restores.AddCommand(cmdList)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createRestoreDeletedPermissions := true
	createRestoreInPlace := true
	paramsRestoreCreate := files_sdk.RestoreCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Restore`,
		Long:  `Create Restore`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := restore.Client{Config: config}

			if cmd.Flags().Changed("restore-deleted-permissions") {
				paramsRestoreCreate.RestoreDeletedPermissions = flib.Bool(createRestoreDeletedPermissions)
			}
			if cmd.Flags().Changed("restore-in-place") {
				paramsRestoreCreate.RestoreInPlace = flib.Bool(createRestoreInPlace)
			}

			if paramsRestoreCreate.EarliestDate.IsZero() {
				paramsRestoreCreate.EarliestDate = nil
			}

			var restore interface{}
			var err error
			restore, err = client.Create(paramsRestoreCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), restore, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	paramsRestoreCreate.EarliestDate = &time.Time{}
	lib.TimeVar(cmdCreate.Flags(), paramsRestoreCreate.EarliestDate, "earliest-date", "Restore all files deleted after this date/time. Don't set this earlier than you need. Can not be greater than 365")
	cmdCreate.Flags().BoolVar(&createRestoreDeletedPermissions, "restore-deleted-permissions", createRestoreDeletedPermissions, "If true, we will also restore any Permissions that match the same path prefix from the same dates.")
	cmdCreate.Flags().BoolVar(&createRestoreInPlace, "restore-in-place", createRestoreInPlace, "If true, we will restore the files in place (into their original paths). If false, we will create a new restoration folder in the root and restore files there.")
	cmdCreate.Flags().StringVar(&paramsRestoreCreate.Prefix, "prefix", "", "Prefix of the files/folders to restore. To restore a folder, add a trailing slash to the folder name. Do not use a leading slash.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Restores.AddCommand(cmdCreate)
	return Restores
}
