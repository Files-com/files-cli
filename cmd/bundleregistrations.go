package cmd

import "github.com/spf13/cobra"

var (
	BundleRegistrations = &cobra.Command{
		Use:  "bundle-registrations [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func BundleRegistrationsInit() {
}
