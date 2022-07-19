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
	usePagerList := true
	paramsBundleList := files_sdk.BundleListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List Bundles",
		Long:  `List Bundles`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsBundleList
			params.MaxPages = MaxPagesList

			client := bundle.Client{Config: *config}
			it, err := client.List(ctx, params)
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger())
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}

	cmdList.Flags().Int64Var(&paramsBundleList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsBundleList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsBundleList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Bundles.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	usePagerFind := true
	paramsBundleFind := files_sdk.BundleFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Bundle`,
		Long:  `Show Bundle`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := bundle.Client{Config: *config}

			var bundle interface{}
			var err error
			bundle, err = client.Find(ctx, paramsBundleFind)
			lib.HandleResponse(ctx, bundle, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdFind.Flags().Int64Var(&paramsBundleFind.Id, "id", 0, "Bundle ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	Bundles.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	usePagerCreate := true
	createPreviewOnly := false
	createRequireRegistration := false
	createRequireShareRecipient := false
	createSkipEmail := false
	createSkipName := false
	createSkipCompany := false
	paramsBundleCreate := files_sdk.BundleCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Bundle`,
		Long:  `Create Bundle`,
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
			if createSkipEmail {
				paramsBundleCreate.SkipEmail = flib.Bool(true)
			}
			if createSkipName {
				paramsBundleCreate.SkipName = flib.Bool(true)
			}
			if createSkipCompany {
				paramsBundleCreate.SkipCompany = flib.Bool(true)
			}

			var bundle interface{}
			var err error
			bundle, err = client.Create(ctx, paramsBundleCreate)
			lib.HandleResponse(ctx, bundle, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdCreate.Flags().Int64Var(&paramsBundleCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringSliceVar(&paramsBundleCreate.Paths, "paths", []string{}, "A list of paths to include in this bundle.")
	cmdCreate.Flags().StringVar(&paramsBundleCreate.Password, "password", "", "Password for this bundle.")
	cmdCreate.Flags().Int64Var(&paramsBundleCreate.FormFieldSetId, "form-field-set-id", 0, "Id of Form Field Set to use with this bundle")
	lib.TimeVar(cmdCreate.Flags(), paramsBundleCreate.ExpiresAt, "expires-at")
	cmdCreate.Flags().Int64Var(&paramsBundleCreate.MaxUses, "max-uses", 0, "Maximum number of times bundle can be accessed")
	cmdCreate.Flags().StringVar(&paramsBundleCreate.Description, "description", "", "Public description")
	cmdCreate.Flags().StringVar(&paramsBundleCreate.Note, "note", "", "Bundle internal note")
	cmdCreate.Flags().StringVar(&paramsBundleCreate.Code, "code", "", "Bundle code.  This code forms the end part of the Public URL.")
	cmdCreate.Flags().BoolVar(&createPreviewOnly, "preview-only", createPreviewOnly, "Restrict users to previewing files only?")
	cmdCreate.Flags().BoolVar(&createRequireRegistration, "require-registration", createRequireRegistration, "Show a registration page that captures the downloader's name and email address?")
	cmdCreate.Flags().Int64Var(&paramsBundleCreate.ClickwrapId, "clickwrap-id", 0, "ID of the clickwrap to use with this bundle.")
	cmdCreate.Flags().Int64Var(&paramsBundleCreate.InboxId, "inbox-id", 0, "ID of the associated inbox, if available.")
	cmdCreate.Flags().BoolVar(&createRequireShareRecipient, "require-share-recipient", createRequireShareRecipient, "Only allow access to recipients who have explicitly received the share via an email sent through the Files.com UI?")
	cmdCreate.Flags().BoolVar(&createSkipEmail, "skip-email", createSkipEmail, "BundleRegistrations can be saved without providing email?")
	cmdCreate.Flags().BoolVar(&createSkipName, "skip-name", createSkipName, "BundleRegistrations can be saved without providing name?")
	cmdCreate.Flags().BoolVar(&createSkipCompany, "skip-company", createSkipCompany, "BundleRegistrations can be saved without providing company?")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Bundles.AddCommand(cmdCreate)
	var fieldsShare string
	var formatShare string
	usePagerShare := true
	paramsBundleShare := files_sdk.BundleShareParams{}

	cmdShare := &cobra.Command{
		Use:   "share",
		Short: `Send email(s) with a link to bundle`,
		Long:  `Send email(s) with a link to bundle`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := bundle.Client{Config: *config}

			var err error
			err = client.Share(ctx, paramsBundleShare)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}
	cmdShare.Flags().Int64Var(&paramsBundleShare.Id, "id", 0, "Bundle ID.")
	cmdShare.Flags().StringSliceVar(&paramsBundleShare.To, "to", []string{}, "A list of email addresses to share this bundle with. Required unless `recipients` is used.")
	cmdShare.Flags().StringVar(&paramsBundleShare.Note, "note", "", "Note to include in email.")

	cmdShare.Flags().StringVarP(&fieldsShare, "fields", "", "", "comma separated list of field names")
	cmdShare.Flags().StringVarP(&formatShare, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdShare.Flags().BoolVar(&usePagerShare, "use-pager", usePagerShare, "Use $PAGER (.ie less, more, etc)")

	Bundles.AddCommand(cmdShare)
	var fieldsUpdate string
	var formatUpdate string
	usePagerUpdate := true
	updatePreviewOnly := false
	updateRequireRegistration := false
	updateRequireShareRecipient := false
	updateSkipEmail := false
	updateSkipName := false
	updateSkipCompany := false
	updateWatermarkAttachmentDelete := false
	paramsBundleUpdate := files_sdk.BundleUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Bundle`,
		Long:  `Update Bundle`,
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
			if updateSkipEmail {
				paramsBundleUpdate.SkipEmail = flib.Bool(true)
			}
			if updateSkipName {
				paramsBundleUpdate.SkipName = flib.Bool(true)
			}
			if updateSkipCompany {
				paramsBundleUpdate.SkipCompany = flib.Bool(true)
			}
			if updateWatermarkAttachmentDelete {
				paramsBundleUpdate.WatermarkAttachmentDelete = flib.Bool(true)
			}

			var bundle interface{}
			var err error
			bundle, err = client.Update(ctx, paramsBundleUpdate)
			lib.HandleResponse(ctx, bundle, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsBundleUpdate.Id, "id", 0, "Bundle ID.")
	cmdUpdate.Flags().StringSliceVar(&paramsBundleUpdate.Paths, "paths", []string{}, "A list of paths to include in this bundle.")
	cmdUpdate.Flags().StringVar(&paramsBundleUpdate.Password, "password", "", "Password for this bundle.")
	cmdUpdate.Flags().Int64Var(&paramsBundleUpdate.FormFieldSetId, "form-field-set-id", 0, "Id of Form Field Set to use with this bundle")
	cmdUpdate.Flags().Int64Var(&paramsBundleUpdate.ClickwrapId, "clickwrap-id", 0, "ID of the clickwrap to use with this bundle.")
	cmdUpdate.Flags().StringVar(&paramsBundleUpdate.Code, "code", "", "Bundle code.  This code forms the end part of the Public URL.")
	cmdUpdate.Flags().StringVar(&paramsBundleUpdate.Description, "description", "", "Public description")
	lib.TimeVar(cmdUpdate.Flags(), paramsBundleUpdate.ExpiresAt, "expires-at")
	cmdUpdate.Flags().Int64Var(&paramsBundleUpdate.InboxId, "inbox-id", 0, "ID of the associated inbox, if available.")
	cmdUpdate.Flags().Int64Var(&paramsBundleUpdate.MaxUses, "max-uses", 0, "Maximum number of times bundle can be accessed")
	cmdUpdate.Flags().StringVar(&paramsBundleUpdate.Note, "note", "", "Bundle internal note")
	cmdUpdate.Flags().BoolVar(&updatePreviewOnly, "preview-only", updatePreviewOnly, "Restrict users to previewing files only?")
	cmdUpdate.Flags().BoolVar(&updateRequireRegistration, "require-registration", updateRequireRegistration, "Show a registration page that captures the downloader's name and email address?")
	cmdUpdate.Flags().BoolVar(&updateRequireShareRecipient, "require-share-recipient", updateRequireShareRecipient, "Only allow access to recipients who have explicitly received the share via an email sent through the Files.com UI?")
	cmdUpdate.Flags().BoolVar(&updateSkipEmail, "skip-email", updateSkipEmail, "BundleRegistrations can be saved without providing email?")
	cmdUpdate.Flags().BoolVar(&updateSkipName, "skip-name", updateSkipName, "BundleRegistrations can be saved without providing name?")
	cmdUpdate.Flags().BoolVar(&updateSkipCompany, "skip-company", updateSkipCompany, "BundleRegistrations can be saved without providing company?")
	cmdUpdate.Flags().BoolVar(&updateWatermarkAttachmentDelete, "watermark-attachment-delete", updateWatermarkAttachmentDelete, "If true, will delete the file stored in watermark_attachment")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Bundles.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	usePagerDelete := true
	paramsBundleDelete := files_sdk.BundleDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Bundle`,
		Long:  `Delete Bundle`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := bundle.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsBundleDelete)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}
	cmdDelete.Flags().Int64Var(&paramsBundleDelete.Id, "id", 0, "Bundle ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Bundles.AddCommand(cmdDelete)
}
