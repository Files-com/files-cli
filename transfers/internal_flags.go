//go:build internal

package transfers

import "github.com/spf13/cobra"

func (t *Transfers) addInternalFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&t.TestProgressBarOut, "test-progress-bar-out", "", "redirect progress bar to file for testing.")
	cmd.Flags().StringVar(&t.CPUProfilePath, "cpu-profile", "", "Write a Go CPU profile for benchmark or PGO analysis.")
	cmd.Flags().MarkHidden("test-progress-bar-out")
	cmd.Flags().MarkHidden("cpu-profile")
}
