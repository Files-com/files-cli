package cmd
import "github.com/spf13/cobra"

var (
    Images = &cobra.Command{
      Use: "images [command]",
      Args:  cobra.ExactArgs(1),
      Run: func(cmd *cobra.Command, args []string) {},
    }
)
func ImagesInit() {
}
