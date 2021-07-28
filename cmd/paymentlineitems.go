package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	PaymentLineItems = &cobra.Command{}
)

func PaymentLineItemsInit() {
	PaymentLineItems = &cobra.Command{
		Use:  "payment-line-items [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command payment-line-items\n\t%v", args[0])
		},
	}
}
