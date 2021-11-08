package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	"github.com/Files-com/files-sdk-go/v2/bundle"
	flib "github.com/Files-com/files-sdk-go/v2/lib"
)

var (
	Bundles = &cobra.Command{}
)

func BundlesInit() {
	Bundles = &cobra.Command{
		Use:  "bundles [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command bundles\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsBundleList := files_sdk.BundleListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsBundleList
			params.MaxPages = MaxPagesList

			client := bundle.Client{Config: *config}
			it, err := client.List(ctx, params)
			if err != nil {
				lib.ClientError(ctx, err)
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(it, formatList, fieldsList, listFilter)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}

	cmdList.Flags().Int64Var(&paramsBundleList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsBundleList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64Var(&paramsBundleList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Bundles.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	paramsBundleFind := files_sdk.BundleFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := bundle.Client{Config: *config}

			result, err := client.Find(ctx, paramsBundleFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatFind, fieldsFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdFind.Flags().Int64Var(&paramsBundleFind.Id, "id", 0, "Bundle ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Bundles.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	createPreviewOnly := false
	createRequireRegistration := false
	createRequireShareRecipient := false
	paramsBundleCreate := files_sdk.BundleCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := bundle.Client{Config: *config}

			if createPreviewOnly {
				paramsBundleCreate.PreviewOnly = flib.Bool(true)
			}
			if createRequireRegistration {
				paramsBundleCreate.RequireRegistration = flib.Bool(true)
			}
			if createRequireShareRecipient {
				paramsBundleCreate.RequireShareRecipient = flib.Bool(true)
			}

			result, err := client.Create(ctx, paramsBundleCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatCreate, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().Int64Var(&paramsBundleCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringSliceVar(&paramsBundleCreate.Paths, "paths", []string{}, "A list of paths to include in this bundle.")
	cmdCreate.Flags().StringVar(&paramsBundleCreate.Password, "password", "", "Password for this bundle.")
	cmdCreate.Flags().Int64Var(&paramsBundleCreate.FormFieldSetId, "form-field-set-id", 0, "Id of Form Field Set to use with this bundle")
	lib.TimeVar(cmdCreate.Flags(), &paramsBundleCreate.ExpiresAt, "expires-at")
	cmdCreate.Flags().Int64Var(&paramsBundleCreate.MaxUses, "max-uses", 0, "Maximum number of times bundle can be accessed")
	cmdCreate.Flags().StringVar(&paramsBundleCreate.Description, "description", "", "Public description")
	cmdCreate.Flags().StringVar(&paramsBundleCreate.Note, "note", "", "Bundle internal note")
	cmdCreate.Flags().StringVar(&paramsBundleCreate.Code, "code", "", "Bundle code.  This code forms the end part of the Public URL.")
	cmdCreate.Flags().BoolVar(&createPreviewOnly, "preview-only", createPreviewOnly, "Restrict users to previewing files only?")
	cmdCreate.Flags().BoolVar(&createRequireRegistration, "require-registration", createRequireRegistration, "Show a registration page that captures the downloader's name and email address?")
	cmdCreate.Flags().Int64Var(&paramsBundleCreate.ClickwrapId, "clickwrap-id", 0, "ID of the clickwrap to use with this bundle.")
	cmdCreate.Flags().Int64Var(&paramsBundleCreate.InboxId, "inbox-id", 0, "ID of the associated inbox, if available.")
	cmdCreate.Flags().BoolVar(&createRequireShareRecipient, "require-share-recipient", createRequireShareRecipient, "Only allow access to recipients who have explicitly received the share via an email sent through the Files.com UI?")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Bundles.AddCommand(cmdCreate)
	var fieldsShare string
	var formatShare string
	paramsBundleShare := files_sdk.BundleShareParams{}

	cmdShare := &cobra.Command{
		Use: "share",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := bundle.Client{Config: *config}

			result, err := client.Share(ctx, paramsBundleShare)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatShare, fieldsShare)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdShare.Flags().Int64Var(&paramsBundleShare.Id, "id", 0, "Bundle ID.")
	cmdShare.Flags().StringSliceVar(&paramsBundleShare.To, "to", []string{}, "A list of email addresses to share this bundle with. Required unless `recipients` is used.")
	cmdShare.Flags().StringVar(&paramsBundleShare.Note, "note", "", "Note to include in email.")
	cmdShare.Flags().StringSliceVar(&paramsBundleShare.Recipients, "recipients", []string{}, "A list of recipients to share this bundle with. Required unless `to` is used.")

	cmdShare.Flags().StringVarP(&fieldsShare, "fields", "", "", "comma separated list of field names")
	cmdShare.Flags().StringVarP(&formatShare, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Bundles.AddCommand(cmdShare)
	var fieldsUpdate string
	var formatUpdate string
	updatePreviewOnly := false
	updateRequireRegistration := false
	updateRequireShareRecipient := false
	updateWatermarkAttachmentDelete := false
	paramsBundleUpdate := files_sdk.BundleUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := bundle.Client{Config: *config}

			if updatePreviewOnly {
				paramsBundleUpdate.PreviewOnly = flib.Bool(true)
			}
			if updateRequireRegistration {
				paramsBundleUpdate.RequireRegistration = flib.Bool(true)
			}
			if updateRequireShareRecipient {
				paramsBundleUpdate.RequireShareRecipient = flib.Bool(true)
			}
			if updateWatermarkAttachmentDelete {
				paramsBundleUpdate.WatermarkAttachmentDelete = flib.Bool(true)
			}

			result, err := client.Update(ctx, paramsBundleUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatUpdate, fieldsUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsBundleUpdate.Id, "id", 0, "Bundle ID.")
	cmdUpdate.Flags().StringSliceVar(&paramsBundleUpdate.Paths, "paths", []string{}, "A list of paths to include in this bundle.")
	cmdUpdate.Flags().StringVar(&paramsBundleUpdate.Password, "password", "", "Password for this bundle.")
	cmdUpdate.Flags().Int64Var(&paramsBundleUpdate.FormFieldSetId, "form-field-set-id", 0, "Id of Form Field Set to use with this bundle")
	cmdUpdate.Flags().Int64Var(&paramsBundleUpdate.ClickwrapId, "clickwrap-id", 0, "ID of the clickwrap to use with this bundle.")
	cmdUpdate.Flags().StringVar(&paramsBundleUpdate.Code, "code", "", "Bundle code.  This code forms the end part of the Public URL.")
	cmdUpdate.Flags().StringVar(&paramsBundleUpdate.Description, "description", "", "Public description")
	lib.TimeVar(cmdUpdate.Flags(), &paramsBundleUpdate.ExpiresAt, "expires-at")
	cmdUpdate.Flags().Int64Var(&paramsBundleUpdate.InboxId, "inbox-id", 0, "ID of the associated inbox, if available.")
	cmdUpdate.Flags().Int64Var(&paramsBundleUpdate.MaxUses, "max-uses", 0, "Maximum number of times bundle can be accessed")
	cmdUpdate.Flags().StringVar(&paramsBundleUpdate.Note, "note", "", "Bundle internal note")
	cmdUpdate.Flags().BoolVar(&updatePreviewOnly, "preview-only", updatePreviewOnly, "Restrict users to previewing files only?")
	cmdUpdate.Flags().BoolVar(&updateRequireRegistration, "require-registration", updateRequireRegistration, "Show a registration page that captures the downloader's name and email address?")
	cmdUpdate.Flags().BoolVar(&updateRequireShareRecipient, "require-share-recipient", updateRequireShareRecipient, "Only allow access to recipients who have explicitly received the share via an email sent through the Files.com UI?")
	cmdUpdate.Flags().BoolVar(&updateWatermarkAttachmentDelete, "watermark-attachment-delete", updateWatermarkAttachmentDelete, "If true, will delete the file stored in watermark_attachment")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Bundles.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	paramsBundleDelete := files_sdk.BundleDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := bundle.Client{Config: *config}

			result, err := client.Delete(ctx, paramsBundleDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatDelete, fieldsDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdDelete.Flags().Int64Var(&paramsBundleDelete.Id, "id", 0, "Bundle ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Bundles.AddCommand(cmdDelete)
}
