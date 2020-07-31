package cmd
import "github.com/spf13/cobra"

var (
    FilePartUploads = &cobra.Command{
      Use: "file-part-uploads [command]",
      Args:  cobra.ExactArgs(1),
      Run: func(cmd *cobra.Command, args []string) {},
    }
)
func FilePartUploadsInit() {
}
