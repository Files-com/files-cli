package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/usercipheruse"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = user_cipher_use.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
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
			it := user_cipher_use.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsUserCipherUseList.UserId, "user-id", "u", 0, "List User Cipher Uses")
	cmdList.Flags().IntVarP(&paramsUserCipherUseList.Page, "page", "p", 0, "List User Cipher Uses")
	cmdList.Flags().IntVarP(&paramsUserCipherUseList.PerPage, "per-page", "e", 0, "List User Cipher Uses")
	cmdList.Flags().StringVarP(&paramsUserCipherUseList.Action, "action", "a", "", "List User Cipher Uses")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
	UserCipherUses.AddCommand(cmdList)
}
