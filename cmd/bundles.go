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
			err = lib.FormatIter(it, formatList, fieldsList)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdList.Flags().Int64VarP(&paramsBundleList.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVarP(&paramsBundleList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsBundleList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

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
	cmdFind.Flags().Int64VarP(&paramsBundleFind.Id, "id", "i", 0, "Bundle ID.")

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
	cmdCreate.Flags().Int64VarP(&paramsBundleCreate.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVarP(&paramsBundleCreate.Password, "password", "p", "", "Password for this bundle.")
	cmdCreate.Flags().Int64VarP(&paramsBundleCreate.FormFieldSetId, "form-field-set-id", "f", 0, "Id of Form Field Set to use with this bundle")
	lib.TimeVarP(cmdCreate.Flags(), &paramsBundleCreate.ExpiresAt, "expires-at", "e")
	cmdCreate.Flags().Int64VarP(&paramsBundleCreate.MaxUses, "max-uses", "a", 0, "Maximum number of times bundle can be accessed")
	cmdCreate.Flags().StringVarP(&paramsBundleCreate.Description, "description", "d", "", "Public description")
	cmdCreate.Flags().StringVarP(&paramsBundleCreate.Note, "note", "n", "", "Bundle internal note")
	cmdCreate.Flags().StringVarP(&paramsBundleCreate.Code, "code", "o", "", "Bundle code.  This code forms the end part of the Public URL.")
	cmdCreate.Flags().BoolVarP(&createPreviewOnly, "preview-only", "r", createPreviewOnly, "Restrict users to previewing files only?")
	cmdCreate.Flags().BoolVarP(&createRequireRegistration, "require-registration", "g", createRequireRegistration, "Show a registration page that captures the downloader's name and email address?")
	cmdCreate.Flags().Int64VarP(&paramsBundleCreate.ClickwrapId, "clickwrap-id", "c", 0, "ID of the clickwrap to use with this bundle.")
	cmdCreate.Flags().Int64VarP(&paramsBundleCreate.InboxId, "inbox-id", "i", 0, "ID of the associated inbox, if available.")
	cmdCreate.Flags().BoolVarP(&createRequireShareRecipient, "require-share-recipient", "", createRequireShareRecipient, "Only allow access to recipients who have explicitly received the share via an email sent through the Files.com UI?")

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
	cmdShare.Flags().Int64VarP(&paramsBundleShare.Id, "id", "i", 0, "Bundle ID.")
	cmdShare.Flags().StringVarP(&paramsBundleShare.Note, "note", "n", "", "Note to include in email.")

	cmdShare.Flags().StringVarP(&fieldsShare, "fields", "", "", "comma separated list of field names")
	cmdShare.Flags().StringVarP(&formatShare, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Bundles.AddCommand(cmdShare)
	var fieldsUpdate string
	var formatUpdate string
	updatePreviewOnly := false
	updateRequireRegistration := false
	updateRequireShareRecipient := false
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
	cmdUpdate.Flags().Int64VarP(&paramsBundleUpdate.Id, "id", "i", 0, "Bundle ID.")
	cmdUpdate.Flags().StringVarP(&paramsBundleUpdate.Password, "password", "p", "", "Password for this bundle.")
	cmdUpdate.Flags().Int64VarP(&paramsBundleUpdate.FormFieldSetId, "form-field-set-id", "f", 0, "Id of Form Field Set to use with this bundle")
	cmdUpdate.Flags().Int64VarP(&paramsBundleUpdate.ClickwrapId, "clickwrap-id", "c", 0, "ID of the clickwrap to use with this bundle.")
	cmdUpdate.Flags().StringVarP(&paramsBundleUpdate.Code, "code", "o", "", "Bundle code.  This code forms the end part of the Public URL.")
	cmdUpdate.Flags().StringVarP(&paramsBundleUpdate.Description, "description", "d", "", "Public description")
	lib.TimeVarP(cmdUpdate.Flags(), &paramsBundleUpdate.ExpiresAt, "expires-at", "e")
	cmdUpdate.Flags().Int64VarP(&paramsBundleUpdate.InboxId, "inbox-id", "n", 0, "ID of the associated inbox, if available.")
	cmdUpdate.Flags().Int64VarP(&paramsBundleUpdate.MaxUses, "max-uses", "a", 0, "Maximum number of times bundle can be accessed")
	cmdUpdate.Flags().StringVarP(&paramsBundleUpdate.Note, "note", "t", "", "Bundle internal note")
	cmdUpdate.Flags().BoolVarP(&updatePreviewOnly, "preview-only", "r", updatePreviewOnly, "Restrict users to previewing files only?")
	cmdUpdate.Flags().BoolVarP(&updateRequireRegistration, "require-registration", "g", updateRequireRegistration, "Show a registration page that captures the downloader's name and email address?")
	cmdUpdate.Flags().BoolVarP(&updateRequireShareRecipient, "require-share-recipient", "", updateRequireShareRecipient, "Only allow access to recipients who have explicitly received the share via an email sent through the Files.com UI?")

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
	cmdDelete.Flags().Int64VarP(&paramsBundleDelete.Id, "id", "i", 0, "Bundle ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Bundles.AddCommand(cmdDelete)
}
