package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	webhook_test "github.com/Files-com/files-sdk-go/webhooktest"
)

var (
	WebhookTests = &cobra.Command{
		Use:  "webhook-tests [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func WebhookTestsInit() {
	var fieldsCreate string
	paramsWebhookTestCreate := files_sdk.WebhookTestCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			client := webhook_test.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Create(paramsWebhookTestCreate)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(err)
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
