package cmd

import (
	"github.com/Files-com/files-cli/lib"
	files "github.com/Files-com/files-sdk-go/v2"
	"github.com/spf13/cobra"

	"fmt"
	"os"
)

func init() {
	RootCmd.AddCommand(Login())
}

func Login() *cobra.Command {
	return &cobra.Command{
		Use:  "login",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			Profile(cmd).Overrides = lib.Overrides{Out: cmd.OutOrStdout(), In: cmd.InOrStdin()}.Init()
			err := lib.CreateSession(files.SessionCreateParams{}, Profile(cmd))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
}
