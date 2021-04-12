package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/webhooktest"
)

var (
	WebhookTests = &cobra.Command{}
)

func WebhookTestsInit() {
	WebhookTests = &cobra.Command{
		Use:  "webhook-tests [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
	var fieldsCreate string
	paramsWebhookTestCreate := files_sdk.WebhookTestCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := webhooktest.Client{Config: *ctx.GetConfig()}

			result, err := client.Create(paramsWebhookTestCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdCreate.Flags().StringVarP(&paramsWebhookTestCreate.Url, "url", "u", "", "URL for testing the webhook.")
	cmdCreate.Flags().StringVarP(&paramsWebhookTestCreate.Method, "method", "t", "", "HTTP method(GET or POST).")
	cmdCreate.Flags().StringVarP(&paramsWebhookTestCreate.Encoding, "encoding", "e", "", "HTTP encoding method.  Can be JSON, XML, or RAW (form data).")
	cmdCreate.Flags().StringVarP(&paramsWebhookTestCreate.Action, "action", "a", "", "action for test body")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	WebhookTests.AddCommand(cmdCreate)
}
