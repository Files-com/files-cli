package cmd
import "github.com/spf13/cobra"

var (
    AccountLineItems = &cobra.Command{
      Use: "account-line-items [command]",
      Args:  cobra.ExactArgs(1),
      Run: func(cmd *cobra.Command, args []string) {},
    }
)
func AccountLineItemsInit() {
}
