package cmd

import (
	"fmt"
	"time"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/snapshot"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Snapshots())
}

func Snapshots() *cobra.Command {
	Snapshots := &cobra.Command{
		Use:  "snapshots [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command snapshots\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsSnapshotList := files_sdk.SnapshotListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Snapshots",
		Long:    `List Snapshots`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsSnapshotList
			params.MaxPages = MaxPagesList

			client := snapshot.Client{Config: config}
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

	cmdList.Flags().StringVar(&paramsSnapshotList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsSnapshotList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVar(&paramsSnapshotList.Action, "action", "", "")
	cmdList.Flags().Int64Var(&paramsSnapshotList.Page, "page", 0, "")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Snapshots.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsSnapshotFind := files_sdk.SnapshotFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Snapshot`,
		Long:  `Show Snapshot`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := snapshot.Client{Config: config}

			var snapshot interface{}
			var err error
			snapshot, err = client.Find(paramsSnapshotFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), snapshot, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsSnapshotFind.Id, "id", 0, "Snapshot ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	Snapshots.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsSnapshotCreate := files_sdk.SnapshotCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Snapshot`,
		Long:  `Create Snapshot`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := snapshot.Client{Config: config}

			if paramsSnapshotCreate.ExpiresAt.IsZero() {
				paramsSnapshotCreate.ExpiresAt = nil
			}

			var snapshot interface{}
			var err error
			snapshot, err = client.Create(paramsSnapshotCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), snapshot, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	paramsSnapshotCreate.ExpiresAt = &time.Time{}
	lib.TimeVar(cmdCreate.Flags(), paramsSnapshotCreate.ExpiresAt, "expires-at", "When the snapshot expires.")
	cmdCreate.Flags().StringVar(&paramsSnapshotCreate.Name, "name", "", "A name for the snapshot.")
	cmdCreate.Flags().StringSliceVar(&paramsSnapshotCreate.Paths, "paths", []string{}, "An array of paths to add to the snapshot.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Snapshots.AddCommand(cmdCreate)
	var fieldsFinalize []string
	var formatFinalize []string
	usePagerFinalize := true
	paramsSnapshotFinalize := files_sdk.SnapshotFinalizeParams{}

	cmdFinalize := &cobra.Command{
		Use:   "finalize",
		Short: `Finalize Snapshot`,
		Long:  `Finalize Snapshot`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := snapshot.Client{Config: config}

			var err error
			err = client.Finalize(paramsSnapshotFinalize, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdFinalize.Flags().Int64Var(&paramsSnapshotFinalize.Id, "id", 0, "Snapshot ID.")

	cmdFinalize.Flags().StringSliceVar(&fieldsFinalize, "fields", []string{}, "comma separated list of field names")
	cmdFinalize.Flags().StringSliceVar(&formatFinalize, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFinalize.Flags().BoolVar(&usePagerFinalize, "use-pager", usePagerFinalize, "Use $PAGER (.ie less, more, etc)")

	Snapshots.AddCommand(cmdFinalize)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsSnapshotUpdate := files_sdk.SnapshotUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Snapshot`,
		Long:  `Update Snapshot`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := snapshot.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.SnapshotUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsSnapshotUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("expires-at") {
				lib.FlagUpdate(cmd, "expires_at", paramsSnapshotUpdate.ExpiresAt, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsSnapshotUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("paths") {
				lib.FlagUpdateLen(cmd, "paths", paramsSnapshotUpdate.Paths, mapParams)
			}

			if paramsSnapshotUpdate.ExpiresAt.IsZero() {
				paramsSnapshotUpdate.ExpiresAt = nil
			}

			var snapshot interface{}
			var err error
			snapshot, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), snapshot, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsSnapshotUpdate.Id, "id", 0, "Snapshot ID.")
	paramsSnapshotUpdate.ExpiresAt = &time.Time{}
	lib.TimeVar(cmdUpdate.Flags(), paramsSnapshotUpdate.ExpiresAt, "expires-at", "When the snapshot expires.")
	cmdUpdate.Flags().StringVar(&paramsSnapshotUpdate.Name, "name", "", "A name for the snapshot.")
	cmdUpdate.Flags().StringSliceVar(&paramsSnapshotUpdate.Paths, "paths", []string{}, "An array of paths to add to the snapshot.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Snapshots.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsSnapshotDelete := files_sdk.SnapshotDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Snapshot`,
		Long:  `Delete Snapshot`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := snapshot.Client{Config: config}

			var err error
			err = client.Delete(paramsSnapshotDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsSnapshotDelete.Id, "id", 0, "Snapshot ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Snapshots.AddCommand(cmdDelete)
	return Snapshots
}
