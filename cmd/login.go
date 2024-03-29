package cmd

import (
	"fmt"
	"os"

	"github.com/Files-com/files-cli/lib"
	files "github.com/Files-com/files-sdk-go/v3"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Login())
}

func Login() *cobra.Command {
	return &cobra.Command{
		Use:  "login",
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			Profile(cmd).Overrides = lib.Overrides{Out: cmd.OutOrStdout(), In: cmd.InOrStdin()}.Init()
			err := lib.CreateSession(cmd.Context(), files.SessionCreateParams{}, Profile(cmd))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
}
