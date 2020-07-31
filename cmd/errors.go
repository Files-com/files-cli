package cmd
import "github.com/spf13/cobra"

var (
    Errors = &cobra.Command{
      Use: "errors [command]",
      Args:  cobra.ExactArgs(1),
      Run: func(cmd *cobra.Command, args []string) {},
    }
)
func ErrorsInit() {
}
