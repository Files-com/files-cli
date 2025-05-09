package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/lock"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Locks())
}

func Locks() *cobra.Command {
	Locks := &cobra.Command{
		Use:  "locks [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command locks\n\t%v", args[0])
		},
	}
	var fieldsListFor []string
	var formatListFor []string
	usePagerListFor := true
	filterbyListFor := make(map[string]string)
	paramsLockListFor := files_sdk.LockListForParams{}
	var MaxPagesListFor int64
	listForIncludeChildren := true

	cmdListFor := &cobra.Command{
		Use:     "list-for [path]",
		Short:   "List Locks by Path",
		Long:    `List Locks by Path`,
		Args:    cobra.RangeArgs(0, 1),
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsLockListFor
			params.MaxPages = MaxPagesListFor
			if len(args) > 0 && args[0] != "" {
				params.Path = args[0]
			}

			if cmd.Flags().Changed("include-children") {
				params.IncludeChildren = flib.Bool(listForIncludeChildren)
			}

			client := lock.Client{Config: config}
			it, err := client.ListFor(params, files_sdk.WithContext(ctx))
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
			if len(filterbyListFor) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyListFor, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatListFor), fieldsListFor, usePagerListFor, listFilter, cmd.OutOrStdout())
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdListFor.Flags().StringToStringVar(&filterbyListFor, "filter-by", filterbyListFor, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdListFor.Flags().StringVar(&paramsLockListFor.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdListFor.Flags().Int64Var(&paramsLockListFor.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListFor.Flags().StringVar(&paramsLockListFor.Path, "path", "", "Path to operate on.")
	cmdListFor.Flags().BoolVar(&listForIncludeChildren, "include-children", listForIncludeChildren, "Include locks from children objects?")

	cmdListFor.Flags().Int64VarP(&MaxPagesListFor, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListFor.Flags().StringSliceVar(&fieldsListFor, "fields", []string{}, "comma separated list of field names to include in response")
	cmdListFor.Flags().StringSliceVar(&formatListFor, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdListFor.Flags().BoolVar(&usePagerListFor, "use-pager", usePagerListFor, "Use $PAGER (.ie less, more, etc)")
	Locks.AddCommand(cmdListFor)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createAllowAccessByAnyUser := true
	createExclusive := true
	createRecursive := true
	paramsLockCreate := files_sdk.LockCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create [path]",
		Short: `Create Lock`,
		Long:  `Create Lock`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := lock.Client{Config: config}

			if cmd.Flags().Changed("allow-access-by-any-user") {
				paramsLockCreate.AllowAccessByAnyUser = flib.Bool(createAllowAccessByAnyUser)
			}
			if cmd.Flags().Changed("exclusive") {
				paramsLockCreate.Exclusive = flib.Bool(createExclusive)
			}
			if cmd.Flags().Changed("recursive") {
				paramsLockCreate.Recursive = flib.Bool(createRecursive)
			}

			if len(args) > 0 && args[0] != "" {
				paramsLockCreate.Path = args[0]
			}
			var lock interface{}
			var err error
			lock, err = client.Create(paramsLockCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), lock, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsLockCreate.Path, "path", "", "Path")
	cmdCreate.Flags().BoolVar(&createAllowAccessByAnyUser, "allow-access-by-any-user", createAllowAccessByAnyUser, "Can lock be modified by users other than its creator?")
	cmdCreate.Flags().BoolVar(&createExclusive, "exclusive", createExclusive, "Is lock exclusive?")
	cmdCreate.Flags().BoolVar(&createRecursive, "recursive", createRecursive, "Does lock apply to subfolders?")
	cmdCreate.Flags().Int64Var(&paramsLockCreate.Timeout, "timeout", 0, "Lock timeout in seconds")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
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
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := lock.Client{Config: config}

			if len(args) > 0 && args[0] != "" {
				paramsLockDelete.Path = args[0]
			}
			var err error
			err = client.Delete(paramsLockDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().StringVar(&paramsLockDelete.Path, "path", "", "Path")
	cmdDelete.Flags().StringVar(&paramsLockDelete.Token, "token", "", "Lock token")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Locks.AddCommand(cmdDelete)
	return Locks
}
