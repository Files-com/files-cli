package cmd

import "github.com/spf13/cobra"

var (
	InvoiceLineItems = &cobra.Command{}
)

func InvoiceLineItemsInit() {
	InvoiceLineItems = &cobra.Command{
		Use:  "invoice-line-items [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
}
