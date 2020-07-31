package cmd

import "github.com/spf13/cobra"

var (
	PublicIpAddresses = &cobra.Command{
		Use:  "public-ip-addresses [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func PublicIpAddressesInit() {
}
