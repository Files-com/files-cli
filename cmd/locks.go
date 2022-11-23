package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	flib "github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/Files-com/files-sdk-go/v2/lock"
)

func init() {
	RootCmd.AddCommand(Locks())
}

func Locks() *cobra.Command {
	Locks := &cobra.Command{
		Use:  "locks [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command locks\n\t%v", args[0])
		},
	}
	var fieldsListFor []string
	var formatListFor []string
	usePagerListFor := true
	paramsLockListFor := files_sdk.LockListForParams{}
	var MaxPagesListFor int64
	listForIncludeChildren := true

	cmdListFor := &cobra.Command{
		Use:   "list-for [path]",
		Short: "List Locks by path",
		Long:  `List Locks by path`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsLockListFor
			params.MaxPages = MaxPagesListFor
			if len(args) > 0 && args[0] != "" {
				params.Path = args[0]
			}
			if cmd.Flags().Changed("include-children") {
				paramsLockListFor.IncludeChildren = flib.Bool(listForIncludeChildren)
			}

			client := lock.Client{Config: *config}
			it, err := client.ListFor(ctx, params)
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger())
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatListFor, fieldsListFor, usePagerListFor, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdListFor.Flags().StringVar(&paramsLockListFor.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdListFor.Flags().Int64Var(&paramsLockListFor.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListFor.Flags().StringVar(&paramsLockListFor.Path, "path", "", "Path to operate on.")
	cmdListFor.Flags().BoolVar(&listForIncludeChildren, "include-children", listForIncludeChildren, "Include locks from children objects?")

	cmdListFor.Flags().Int64VarP(&MaxPagesListFor, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListFor.Flags().StringSliceVar(&fieldsListFor, "fields", []string{}, "comma separated list of field names to include in response")
	cmdListFor.Flags().StringSliceVar(&formatListFor, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdListFor.Flags().BoolVar(&usePagerListFor, "use-pager", usePagerListFor, "Use $PAGER (.ie less, more, etc)")
	Locks.AddCommand(cmdListFor)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createAllowAccessByAnyUser := true
	createExclusive := true
	paramsLockCreate := files_sdk.LockCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create [path]",
		Short: `Create Lock`,
		Long:  `Create Lock`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := lock.Client{Config: *config}

			if cmd.Flags().Changed("allow-access-by-any-user") {
				paramsLockCreate.AllowAccessByAnyUser = flib.Bool(createAllowAccessByAnyUser)
			}
			if cmd.Flags().Changed("exclusive") {
				paramsLockCreate.Exclusive = flib.Bool(createExclusive)
			}

			if len(args) > 0 && args[0] != "" {
				paramsLockCreate.Path = args[0]
			}
			var lock interface{}
			var err error
			lock, err = client.Create(ctx, paramsLockCreate)
			lib.HandleResponse(ctx, Profile(cmd), lock, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdCreate.Flags().StringVar(&paramsLockCreate.Path, "path", "", "Path")
	cmdCreate.Flags().BoolVar(&createAllowAccessByAnyUser, "allow-access-by-any-user", createAllowAccessByAnyUser, "Allow lock to be updated by any user?")
	cmdCreate.Flags().BoolVar(&createExclusive, "exclusive", createExclusive, "Is lock exclusive?")
	cmdCreate.Flags().StringVar(&paramsLockCreate.Recursive, "recursive", "", "Does lock apply to subfolders?")
	cmdCreate.Flags().Int64Var(&paramsLockCreate.Timeout, "timeout", 0, "Lock timeout length")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Locks.AddCommand(cmdCreate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsLockDelete := files_sdk.LockDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete [path]",
		Short: `Delete Lock`,
		Long:  `Delete Lock`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := lock.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsLockDelete.Path = args[0]
			}
			var err error
			err = client.Delete(ctx, paramsLockDelete)
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().StringVar(&paramsLockDelete.Path, "path", "", "Path")
	cmdDelete.Flags().StringVar(&paramsLockDelete.Token, "token", "", "Lock token")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Locks.AddCommand(cmdDelete)
	return Locks
}
