package cmd

import (
	"github.com/Files-com/files-cli/lib"
	files "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/session"
	"github.com/spf13/cobra"

	"fmt"
	"os"
)

var LogOut = &cobra.Command{
	Use:  "logout",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := session.Delete(files.SessionDeleteParams{})
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
