package cmd

import (
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/session"
	"github.com/spf13/cobra"

	"fmt"
	"os"
)

var LogOut = &cobra.Command{
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

		config := lib.Config{}
		err = config.Load()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		config.SessionId = ""
		config.Save()
	},
}
