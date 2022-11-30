package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	action_webhook_failure "github.com/Files-com/files-sdk-go/v2/actionwebhookfailure"
)

func init() {
	RootCmd.AddCommand(ActionWebhookFailures())
}

func ActionWebhookFailures() *cobra.Command {
	ActionWebhookFailures := &cobra.Command{
		Use:  "action-webhook-failures [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command action-webhook-failures\n\t%v", args[0])
		},
	}
	var fieldsRetry []string
	var formatRetry []string
	usePagerRetry := true
	paramsActionWebhookFailureRetry := files_sdk.ActionWebhookFailureRetryParams{}

	cmdRetry := &cobra.Command{
		Use:   "retry",
		Short: `retry Action Webhook Failure`,
		Long:  `retry Action Webhook Failure`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := action_webhook_failure.Client{Config: *config}

			var err error
			err = client.Retry(ctx, paramsActionWebhookFailureRetry)
			if err != nil {
				return lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdRetry.Flags().Int64Var(&paramsActionWebhookFailureRetry.Id, "id", 0, "Action Webhook Failure ID.")

	cmdRetry.Flags().StringSliceVar(&fieldsRetry, "fields", []string{}, "comma separated list of field names")
	cmdRetry.Flags().StringSliceVar(&formatRetry, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdRetry.Flags().BoolVar(&usePagerRetry, "use-pager", usePagerRetry, "Use $PAGER (.ie less, more, etc)")

	ActionWebhookFailures.AddCommand(cmdRetry)
	return ActionWebhookFailures
}
