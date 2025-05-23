package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/webhooktest"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(WebhookTests())
}

func WebhookTests() *cobra.Command {
	WebhookTests := &cobra.Command{
		Use:  "webhook-tests [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command webhook-tests\n\t%v", args[0])
		},
	}
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createFileAsBody := true
	createUseDedicatedIps := true
	paramsWebhookTestCreate := files_sdk.WebhookTestCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Webhook Test`,
		Long:  `Create Webhook Test`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := webhooktest.Client{Config: config}

			if cmd.Flags().Changed("file-as-body") {
				paramsWebhookTestCreate.FileAsBody = flib.Bool(createFileAsBody)
			}
			if cmd.Flags().Changed("use-dedicated-ips") {
				paramsWebhookTestCreate.UseDedicatedIps = flib.Bool(createUseDedicatedIps)
			}

			var webhookTest interface{}
			var err error
			webhookTest, err = client.Create(paramsWebhookTestCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), webhookTest, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsWebhookTestCreate.Url, "url", "", "URL for testing the webhook.")
	cmdCreate.Flags().StringVar(&paramsWebhookTestCreate.Method, "method", "", "HTTP method(GET or POST).")
	cmdCreate.Flags().StringVar(&paramsWebhookTestCreate.Encoding, "encoding", "", "HTTP encoding method.  Can be JSON, XML, or RAW (form data).")
	cmdCreate.Flags().StringVar(&paramsWebhookTestCreate.RawBody, "raw-body", "", "raw body text")
	cmdCreate.Flags().BoolVar(&createFileAsBody, "file-as-body", createFileAsBody, "Send the file data as the request body?")
	cmdCreate.Flags().StringVar(&paramsWebhookTestCreate.FileFormField, "file-form-field", "", "Send the file data as a named parameter in the request POST body")
	cmdCreate.Flags().StringVar(&paramsWebhookTestCreate.Action, "action", "", "action for test body")
	cmdCreate.Flags().BoolVar(&createUseDedicatedIps, "use-dedicated-ips", createUseDedicatedIps, "Use dedicated IPs for sending the webhook?")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	WebhookTests.AddCommand(cmdCreate)
	return WebhookTests
}
