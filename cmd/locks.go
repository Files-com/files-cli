package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"

	flib "github.com/Files-com/files-sdk-go/lib"
	"github.com/Files-com/files-sdk-go/lock"
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
	paramsLockListFor := files_sdk.LockListForParams{}
	var MaxPagesListFor int64
	listForIncludeChildren := false

	cmdListFor := &cobra.Command{
		Use:   "list-for [path]",
		Short: "list-for",
		Long:  `list-for`,
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
			if err != nil {
				lib.ClientError(ctx, err)
			}
			err = lib.JsonMarshalIter(it, fieldsListFor)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdListFor.Flags().StringVarP(&paramsLockListFor.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdListFor.Flags().Int64VarP(&paramsLockListFor.PerPage, "per-page", "e", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListFor.Flags().StringVarP(&paramsLockListFor.Path, "path", "p", "", "Path to operate on.")
	cmdListFor.Flags().BoolVarP(&listForIncludeChildren, "include-children", "i", listForIncludeChildren, "Include locks from children objects?")

	cmdListFor.Flags().Int64VarP(&MaxPagesListFor, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListFor.Flags().StringVarP(&fieldsListFor, "fields", "", "", "comma separated list of field names to include in response")
	Locks.AddCommand(cmdListFor)
	var fieldsCreate string
	createAllowAccessByAnyUser := false
	createExclusive := false
	paramsLockCreate := files_sdk.LockCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create [path]",
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

			result, err := client.Create(ctx, paramsLockCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().StringVarP(&paramsLockCreate.Path, "path", "p", "", "Path")
	cmdCreate.Flags().BoolVarP(&createAllowAccessByAnyUser, "allow-access-by-any-user", "a", createAllowAccessByAnyUser, "Allow lock to be updated by any user?")
	cmdCreate.Flags().BoolVarP(&createExclusive, "exclusive", "e", createExclusive, "Is lock exclusive?")
	cmdCreate.Flags().StringVarP(&paramsLockCreate.Recursive, "recursive", "r", "", "Does lock apply to subfolders?")
	cmdCreate.Flags().Int64VarP(&paramsLockCreate.Timeout, "timeout", "t", 0, "Lock timeout length")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Locks.AddCommand(cmdCreate)
	var fieldsDelete string
	paramsLockDelete := files_sdk.LockDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := lock.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsLockDelete.Path = args[0]
			}

			result, err := client.Delete(ctx, paramsLockDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdDelete.Flags().StringVarP(&paramsLockDelete.Path, "path", "p", "", "Path")
	cmdDelete.Flags().StringVarP(&paramsLockDelete.Token, "token", "t", "", "Lock token")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Locks.AddCommand(cmdDelete)
}
