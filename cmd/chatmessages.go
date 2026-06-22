package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(ChatMessages())
}

func ChatMessages() *cobra.Command {
	ChatMessages := &cobra.Command{
		Use:  "chat-messages [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command chat-messages\n\t%v", args[0])
		},
	}
	return ChatMessages
}
