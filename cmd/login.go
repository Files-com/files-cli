package cmd

import (
	"github.com/Files-com/files-cli/lib"
	files "github.com/Files-com/files-sdk-go"
	"github.com/spf13/cobra"

	"fmt"
	"os"
)

var Login = &cobra.Command{
	Use:  "login",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		config := &lib.Config{}
		err := config.Load()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = lib.CreateSession(files.SessionCreateParams{}, *config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
