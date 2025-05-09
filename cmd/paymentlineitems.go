package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(PaymentLineItems())
}

func PaymentLineItems() *cobra.Command {
	PaymentLineItems := &cobra.Command{
		Use:  "payment-line-items [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command payment-line-items\n\t%v", args[0])
		},
	}
	return PaymentLineItems
}
