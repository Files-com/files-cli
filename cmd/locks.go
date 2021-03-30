package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/lock"
)

var (
	Locks = &cobra.Command{
		Use:  "locks [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func LocksInit() {
	var fieldsListFor string
	paramsLockListFor := files_sdk.LockListForParams{}
	var MaxPagesListFor int
	cmdListFor := &cobra.Command{
		Use:   "list-for [path]",
		Short: "list-for",
		Long:  `list-for`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsLockListFor
			params.MaxPages = MaxPagesListFor
			if len(args) > 0 && args[0] != "" {
				params.Path = args[0]
			}
			client := lock.Client{Config: files_sdk.GlobalConfig}
			it, err := client.ListFor(params)
			if err != nil {
				lib.ClientError(err)
			}
			err = lib.JsonMarshalIter(it, fieldsListFor)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdListFor.Flags().StringVarP(&paramsLockListFor.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdListFor.Flags().IntVarP(&paramsLockListFor.PerPage, "per-page", "e", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListFor.Flags().StringVarP(&paramsLockListFor.Path, "path", "p", "", "Path to operate on.")
	cmdListFor.Flags().IntVarP(&MaxPagesListFor, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListFor.Flags().StringVarP(&fieldsListFor, "fields", "", "", "comma separated list of field names to include in response")
	Locks.AddCommand(cmdListFor)
	var fieldsCreate string
	paramsLockCreate := files_sdk.LockCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			client := lock.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Create(paramsLockCreate)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdCreate.Flags().StringVarP(&paramsLockCreate.Path, "path", "p", "", "Path")
	cmdCreate.Flags().IntVarP(&paramsLockCreate.Timeout, "timeout", "t", 0, "Lock timeout length")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Locks.AddCommand(cmdCreate)
	var fieldsDelete string
	paramsLockDelete := files_sdk.LockDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			client := lock.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Delete(paramsLockDelete)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdDelete.Flags().StringVarP(&paramsLockDelete.Path, "path", "p", "", "Path")
	cmdDelete.Flags().StringVarP(&paramsLockDelete.Token, "token", "t", "", "Lock token")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Locks.AddCommand(cmdDelete)
}
