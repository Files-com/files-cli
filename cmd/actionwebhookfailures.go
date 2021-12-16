package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	action_webhook_failure "github.com/Files-com/files-sdk-go/v2/actionwebhookfailure"
)

var (
	ActionWebhookFailures = &cobra.Command{}
)

func ActionWebhookFailuresInit() {
	ActionWebhookFailures = &cobra.Command{
		Use:  "action-webhook-failures [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command action-webhook-failures\n\t%v", args[0])
		},
	}
	var fieldsRetry string
	var formatRetry string
	paramsActionWebhookFailureRetry := files_sdk.ActionWebhookFailureRetryParams{}

	cmdRetry := &cobra.Command{
		Use: "retry",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := action_webhook_failure.Client{Config: *config}

			var result interface{}
			var err error
			result, err = client.Retry(ctx, paramsActionWebhookFailureRetry)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatRetry, fieldsRetry, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdRetry.Flags().Int64Var(&paramsActionWebhookFailureRetry.Id, "id", 0, "Action Webhook Failure ID.")

	cmdRetry.Flags().StringVarP(&fieldsRetry, "fields", "", "", "comma separated list of field names")
	cmdRetry.Flags().StringVarP(&formatRetry, "format", "", "table", "json, csv, table, table-dark, table-bright")
	ActionWebhookFailures.AddCommand(cmdRetry)
}
