package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/session"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = session.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
)

var (
	Sessions = &cobra.Command{
		Use:  "sessions [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func SessionsInit() {
	var fieldsCreate string
	paramsSessionCreate := files_sdk.SessionCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := session.Create(paramsSessionCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().StringVarP(&paramsSessionCreate.Username, "username", "u", "", "Create user session (log in)")
	cmdCreate.Flags().StringVarP(&paramsSessionCreate.Password, "password", "a", "", "Create user session (log in)")
	cmdCreate.Flags().StringVarP(&paramsSessionCreate.Otp, "otp", "o", "", "Create user session (log in)")
	cmdCreate.Flags().StringVarP(&paramsSessionCreate.PartialSessionId, "partial-session-id", "p", "", "Create user session (log in)")
	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
	Sessions.AddCommand(cmdCreate)
	var fieldsDelete string
	paramsSessionDelete := files_sdk.SessionDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := session.Delete(paramsSessionDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}
	cmdDelete.Flags().StringVarP(&paramsSessionDelete.Format, "format", "o", "", "Delete user session (log out)")
	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
	Sessions.AddCommand(cmdDelete)
}
