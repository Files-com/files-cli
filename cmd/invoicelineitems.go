package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	InvoiceLineItems = &cobra.Command{}
)

func InvoiceLineItemsInit() {
	InvoiceLineItems = &cobra.Command{
		Use:  "invoice-line-items [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command invoice-line-items\n\t%v", args[0])
		},
	}
}
