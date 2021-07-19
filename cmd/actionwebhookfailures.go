package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	action_webhook_failure "github.com/Files-com/files-sdk-go/actionwebhookfailure"
)

var (
	ActionWebhookFailures = &cobra.Command{}
)

func ActionWebhookFailuresInit() {
	ActionWebhookFailures = &cobra.Command{
		Use:  "action-webhook-failures [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
	var fieldsRetry string
	paramsActionWebhookFailureRetry := files_sdk.ActionWebhookFailureRetryParams{}

	cmdRetry := &cobra.Command{
		Use: "retry",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := action_webhook_failure.Client{Config: *ctx.GetConfig()}

			result, err := client.Retry(paramsActionWebhookFailureRetry)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsRetry)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdRetry.Flags().Int64VarP(&paramsActionWebhookFailureRetry.Id, "id", "i", 0, "Action Webhook Failure ID.")

	cmdRetry.Flags().StringVarP(&fieldsRetry, "fields", "", "", "comma separated list of field names")
	ActionWebhookFailures.AddCommand(cmdRetry)
}
