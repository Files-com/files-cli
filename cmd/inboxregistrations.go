package cmd

import "github.com/spf13/cobra"

var (
	InboxRegistrations = &cobra.Command{
		Use:  "inbox-registrations [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func InboxRegistrationsInit() {
}
