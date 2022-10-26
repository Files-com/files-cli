package cmd

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/session"
	"github.com/spf13/cobra"

	"fmt"
	"os"
)

func init() {
	RootCmd.AddCommand(LogOut())
}

func LogOut() *cobra.Command {
	return &cobra.Command{
		Use:  "logout",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			client := session.Client{Config: *ctx.Value("config").(*files_sdk.Config)}
			err := client.Delete(cmd.Context())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			Profile(cmd).Current().SessionId = ""
			Profile(cmd).Save()
		},
	}
}
