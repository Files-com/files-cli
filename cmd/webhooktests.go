package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	flib "github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/Files-com/files-sdk-go/v2/webhooktest"
)

var (
	WebhookTests = &cobra.Command{}
)

func WebhookTestsInit() {
	WebhookTests = &cobra.Command{
		Use:  "webhook-tests [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command webhook-tests\n\t%v", args[0])
		},
	}
	var fieldsCreate string
	var formatCreate string
	usePagerCreate := true
	createFileAsBody := false
	paramsWebhookTestCreate := files_sdk.WebhookTestCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Webhook Test`,
		Long:  `Create Webhook Test`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := webhooktest.Client{Config: *config}

			if createFileAsBody {
				paramsWebhookTestCreate.FileAsBody = flib.Bool(true)
			}

			var webhookTest interface{}
			var err error
			webhookTest, err = client.Create(ctx, paramsWebhookTestCreate)
			lib.HandleResponse(ctx, webhookTest, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdCreate.Flags().StringVar(&paramsWebhookTestCreate.Url, "url", "", "URL for testing the webhook.")
	cmdCreate.Flags().StringVar(&paramsWebhookTestCreate.Method, "method", "", "HTTP method(GET or POST).")
	cmdCreate.Flags().StringVar(&paramsWebhookTestCreate.Encoding, "encoding", "", "HTTP encoding method.  Can be JSON, XML, or RAW (form data).")
	cmdCreate.Flags().StringVar(&paramsWebhookTestCreate.RawBody, "raw-body", "", "raw body text")
	cmdCreate.Flags().BoolVar(&createFileAsBody, "file-as-body", createFileAsBody, "Send the file data as the request body?")
	cmdCreate.Flags().StringVar(&paramsWebhookTestCreate.FileFormField, "file-form-field", "", "Send the file data as a named parameter in the request POST body")
	cmdCreate.Flags().StringVar(&paramsWebhookTestCreate.Action, "action", "", "action for test body")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	WebhookTests.AddCommand(cmdCreate)
}
