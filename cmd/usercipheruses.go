package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"
	"os"

	user_cipher_use "github.com/Files-com/files-sdk-go/usercipheruse"
)

var (
	UserCipherUses = &cobra.Command{
		Use:  "user-cipher-uses [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func UserCipherUsesInit() {
	var fieldsList string
	paramsUserCipherUseList := files_sdk.UserCipherUseListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsUserCipherUseList
			params.MaxPages = MaxPagesList
			client := user_cipher_use.Client{Config: files_sdk.GlobalConfig}
			it, err := client.List(params)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdList.Flags().Int64VarP(&paramsUserCipherUseList.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVarP(&paramsUserCipherUseList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&paramsUserCipherUseList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	UserCipherUses.AddCommand(cmdList)
}
