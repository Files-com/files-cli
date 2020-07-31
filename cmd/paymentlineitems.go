package cmd

import "github.com/spf13/cobra"

var (
	PaymentLineItems = &cobra.Command{
		Use:  "payment-line-items [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func PaymentLineItemsInit() {
}
