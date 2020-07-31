package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/lock"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = lock.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
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
			it := lock.ListFor(params)

			lib.JsonMarshalIter(it, fieldsListFor)
		},
	}
	cmdListFor.Flags().IntVarP(&paramsLockListFor.Page, "page", "p", 0, "List Locks by path")
	cmdListFor.Flags().IntVarP(&paramsLockListFor.PerPage, "per-page", "e", 0, "List Locks by path")
	cmdListFor.Flags().StringVarP(&paramsLockListFor.Action, "action", "a", "", "List Locks by path")
	cmdListFor.Flags().StringVarP(&paramsLockListFor.Path, "path", "t", "", "List Locks by path")
	cmdListFor.Flags().IntVarP(&MaxPagesListFor, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdListFor.Flags().StringVarP(&fieldsListFor, "fields", "f", "", "comma separated list of field names to include in response")
	Locks.AddCommand(cmdListFor)
	var fieldsCreate string
	paramsLockCreate := files_sdk.LockCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := lock.Create(paramsLockCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().StringVarP(&paramsLockCreate.Path, "path", "p", "", "Create Lock")
	cmdCreate.Flags().IntVarP(&paramsLockCreate.Timeout, "timeout", "t", 0, "Create Lock")
	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
	Locks.AddCommand(cmdCreate)
	var fieldsDelete string
	paramsLockDelete := files_sdk.LockDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := lock.Delete(paramsLockDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}
	cmdDelete.Flags().StringVarP(&paramsLockDelete.Path, "path", "p", "", "Delete Lock")
	cmdDelete.Flags().StringVarP(&paramsLockDelete.Token, "token", "t", "", "Delete Lock")
	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
	Locks.AddCommand(cmdDelete)
}
