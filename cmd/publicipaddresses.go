package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	PublicIpAddresses = &cobra.Command{}
)

func PublicIpAddressesInit() {
	PublicIpAddresses = &cobra.Command{
		Use:  "public-ip-addresses [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command public-ip-addresses\n\t%v", args[0])
		},
	}
}
