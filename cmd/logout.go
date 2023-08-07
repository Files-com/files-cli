package cmd

import (
	"fmt"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/session"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(LogOut())
}

func LogOut() *cobra.Command {
	return &cobra.Command{
		Use:  "logout",
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			client := session.Client{Config: *ctx.Value("config").(*files_sdk.Config)}
			deleteErr := client.Delete(files_sdk.WithContext(cmd.Context()))
			if deleteErr != nil {
				fmt.Println(deleteErr)
			}
			Profile(cmd).Current().SessionId = ""
			err := Profile(cmd).Save()
			if err != nil {
				fmt.Println(err)
			}

			if deleteErr != nil || err != nil {
				os.Exit(1)
			}
		},
	}
}
