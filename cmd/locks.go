package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	flib "github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/Files-com/files-sdk-go/v2/lock"
)

var (
	Locks = &cobra.Command{}
)

func LocksInit() {
	Locks = &cobra.Command{
		Use:  "locks [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command locks\n\t%v", args[0])
		},
	}
	var fieldsListFor string
	var formatListFor string
	paramsLockListFor := files_sdk.LockListForParams{}
	var MaxPagesListFor int64
	listForIncludeChildren := false

	cmdListFor := &cobra.Command{
		Use:   "list-for [path]",
		Short: "List Locks by path",
		Long:  `List Locks by path`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsLockListFor
			params.MaxPages = MaxPagesListFor
			if len(args) > 0 && args[0] != "" {
				params.Path = args[0]
			}
			if listForIncludeChildren {
				paramsLockListFor.IncludeChildren = flib.Bool(true)
			}

			client := lock.Client{Config: *config}
			it, err := client.ListFor(ctx, params)
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, formatListFor, config.Logger())
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
			err = lib.FormatIter(it, formatListFor, fieldsListFor, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}

	cmdListFor.Flags().StringVar(&paramsLockListFor.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdListFor.Flags().Int64Var(&paramsLockListFor.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListFor.Flags().StringVar(&paramsLockListFor.Path, "path", "", "Path to operate on.")
	cmdListFor.Flags().BoolVar(&listForIncludeChildren, "include-children", listForIncludeChildren, "Include locks from children objects?")

	cmdListFor.Flags().Int64VarP(&MaxPagesListFor, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListFor.Flags().StringVarP(&fieldsListFor, "fields", "", "", "comma separated list of field names to include in response")
	cmdListFor.Flags().StringVarP(&formatListFor, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Locks.AddCommand(cmdListFor)
	var fieldsCreate string
	var formatCreate string
	createAllowAccessByAnyUser := false
	createExclusive := false
	paramsLockCreate := files_sdk.LockCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create [path]",
		Short: `Create Lock`,
		Long:  `Create Lock`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := lock.Client{Config: *config}

			if createAllowAccessByAnyUser {
				paramsLockCreate.AllowAccessByAnyUser = flib.Bool(true)
			}
			if createExclusive {
				paramsLockCreate.Exclusive = flib.Bool(true)
			}

			if len(args) > 0 && args[0] != "" {
				paramsLockCreate.Path = args[0]
			}
			var lock interface{}
			var err error
			lock, err = client.Create(ctx, paramsLockCreate)
			lib.HandleResponse(ctx, lock, err, formatCreate, fieldsCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdCreate.Flags().StringVar(&paramsLockCreate.Path, "path", "", "Path")
	cmdCreate.Flags().BoolVar(&createAllowAccessByAnyUser, "allow-access-by-any-user", createAllowAccessByAnyUser, "Allow lock to be updated by any user?")
	cmdCreate.Flags().BoolVar(&createExclusive, "exclusive", createExclusive, "Is lock exclusive?")
	cmdCreate.Flags().StringVar(&paramsLockCreate.Recursive, "recursive", "", "Does lock apply to subfolders?")
	cmdCreate.Flags().Int64Var(&paramsLockCreate.Timeout, "timeout", 0, "Lock timeout length")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Locks.AddCommand(cmdCreate)
	var fieldsDelete string
	var formatDelete string
	paramsLockDelete := files_sdk.LockDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete [path]",
		Short: `Delete Lock`,
		Long:  `Delete Lock`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := lock.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsLockDelete.Path = args[0]
			}
			var err error
			err = client.Delete(ctx, paramsLockDelete)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}
	cmdDelete.Flags().StringVar(&paramsLockDelete.Path, "path", "", "Path")
	cmdDelete.Flags().StringVar(&paramsLockDelete.Token, "token", "", "Lock token")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Locks.AddCommand(cmdDelete)
}
