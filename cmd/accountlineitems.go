package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	AccountLineItems = &cobra.Command{}
)

func AccountLineItemsInit() {
	AccountLineItems = &cobra.Command{
		Use:  "account-line-items [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command account-line-items\n\t%v", args[0])
		},
	}
}
