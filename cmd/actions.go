package cmd
import "github.com/spf13/cobra"

var (
    Actions = &cobra.Command{
      Use: "actions [command]",
      Args:  cobra.ExactArgs(1),
      Run: func(cmd *cobra.Command, args []string) {},
    }
)
func ActionsInit() {
}
