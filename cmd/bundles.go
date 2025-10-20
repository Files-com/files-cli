package cmd

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/bundle"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Bundles())
}

func Bundles() *cobra.Command {
	Bundles := &cobra.Command{
		Use:  "bundles [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command bundles\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsBundleList := files_sdk.BundleListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Share Links",
		Long:    `List Share Links`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsBundleList
			params.MaxPages = MaxPagesList

			client := bundle.Client{Config: config}
			it, err := client.List(params, files_sdk.WithContext(ctx))
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger)
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyList) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyList, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatList), fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdList.Flags().Int64Var(&paramsBundleList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsBundleList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsBundleList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Bundles.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsBundleFind := files_sdk.BundleFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Share Link`,
		Long:  `Show Share Link`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := bundle.Client{Config: config}

			var bundle interface{}
			var err error
			bundle, err = client.Find(paramsBundleFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), bundle, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsBundleFind.Id, "id", 0, "Bundle ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	Bundles.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createCreateSnapshot := true
	createDontSeparateSubmissionsByFolder := true
	createFinalizeSnapshot := true
	createRequireRegistration := true
	createRequireShareRecipient := true
	createSendEmailReceiptToUploader := true
	createSkipEmail := true
	createSkipName := true
	createSkipCompany := true
	paramsBundleCreate := files_sdk.BundleCreateParams{}
	BundleCreatePermissions := ""

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Share Link`,
		Long:  `Create Share Link`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := bundle.Client{Config: config}

			var BundleCreatePermissionsErr error
			paramsBundleCreate.Permissions, BundleCreatePermissionsErr = lib.FetchKey("permissions", paramsBundleCreate.Permissions.Enum(), BundleCreatePermissions)
			if BundleCreatePermissions != "" && BundleCreatePermissionsErr != nil {
				return BundleCreatePermissionsErr
			}

			if cmd.Flags().Changed("create-snapshot") {
				paramsBundleCreate.CreateSnapshot = flib.Bool(createCreateSnapshot)
			}
			if cmd.Flags().Changed("dont-separate-submissions-by-folder") {
				paramsBundleCreate.DontSeparateSubmissionsByFolder = flib.Bool(createDontSeparateSubmissionsByFolder)
			}
			if cmd.Flags().Changed("finalize-snapshot") {
				paramsBundleCreate.FinalizeSnapshot = flib.Bool(createFinalizeSnapshot)
			}
			if cmd.Flags().Changed("require-registration") {
				paramsBundleCreate.RequireRegistration = flib.Bool(createRequireRegistration)
			}
			if cmd.Flags().Changed("require-share-recipient") {
				paramsBundleCreate.RequireShareRecipient = flib.Bool(createRequireShareRecipient)
			}
			if cmd.Flags().Changed("send-email-receipt-to-uploader") {
				paramsBundleCreate.SendEmailReceiptToUploader = flib.Bool(createSendEmailReceiptToUploader)
			}
			if cmd.Flags().Changed("skip-email") {
				paramsBundleCreate.SkipEmail = flib.Bool(createSkipEmail)
			}
			if cmd.Flags().Changed("skip-name") {
				paramsBundleCreate.SkipName = flib.Bool(createSkipName)
			}
			if cmd.Flags().Changed("skip-company") {
				paramsBundleCreate.SkipCompany = flib.Bool(createSkipCompany)
			}

			if paramsBundleCreate.ExpiresAt.IsZero() {
				paramsBundleCreate.ExpiresAt = nil
			}
			if paramsBundleCreate.StartAccessOnDate.IsZero() {
				paramsBundleCreate.StartAccessOnDate = nil
			}

			var bundle interface{}
			var err error
			bundle, err = client.Create(paramsBundleCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), bundle, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().Int64Var(&paramsBundleCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringSliceVar(&paramsBundleCreate.Paths, "paths", []string{}, "A list of paths to include in this bundle.")
	cmdCreate.Flags().StringVar(&paramsBundleCreate.Password, "password", "", "Password for this bundle.")
	cmdCreate.Flags().Int64Var(&paramsBundleCreate.FormFieldSetId, "form-field-set-id", 0, "Id of Form Field Set to use with this bundle")
	cmdCreate.Flags().BoolVar(&createCreateSnapshot, "create-snapshot", createCreateSnapshot, "If true, create a snapshot of this bundle's contents.")
	cmdCreate.Flags().BoolVar(&createDontSeparateSubmissionsByFolder, "dont-separate-submissions-by-folder", createDontSeparateSubmissionsByFolder, "Do not create subfolders for files uploaded to this share. Note: there are subtle security pitfalls with allowing anonymous uploads from multiple users to live in the same folder. We strongly discourage use of this option unless absolutely required.")
	paramsBundleCreate.ExpiresAt = &time.Time{}
	lib.TimeVar(cmdCreate.Flags(), paramsBundleCreate.ExpiresAt, "expires-at", "Bundle expiration date/time")
	cmdCreate.Flags().BoolVar(&createFinalizeSnapshot, "finalize-snapshot", createFinalizeSnapshot, "If true, finalize the snapshot of this bundle's contents. Note that `create_snapshot` must also be true.")
	cmdCreate.Flags().Int64Var(&paramsBundleCreate.MaxUses, "max-uses", 0, "Maximum number of times bundle can be accessed")
	cmdCreate.Flags().StringVar(&paramsBundleCreate.Description, "description", "", "Public description")
	cmdCreate.Flags().StringVar(&paramsBundleCreate.Note, "note", "", "Bundle internal note")
	cmdCreate.Flags().StringVar(&paramsBundleCreate.Code, "code", "", "Bundle code.  This code forms the end part of the Public URL.")
	cmdCreate.Flags().StringVar(&paramsBundleCreate.PathTemplate, "path-template", "", "Template for creating submission subfolders. Can use the uploader's name, email address, ip, company, `strftime` directives, and any custom form data.")
	cmdCreate.Flags().StringVar(&paramsBundleCreate.PathTemplateTimeZone, "path-template-time-zone", "", "Timezone to use when rendering timestamps in path templates.")
	cmdCreate.Flags().StringVar(&BundleCreatePermissions, "permissions", "", fmt.Sprintf("Permissions that apply to Folders in this Share Link. %v", reflect.ValueOf(paramsBundleCreate.Permissions.Enum()).MapKeys()))
	cmdCreate.Flags().BoolVar(&createRequireRegistration, "require-registration", createRequireRegistration, "Show a registration page that captures the downloader's name and email address?")
	cmdCreate.Flags().Int64Var(&paramsBundleCreate.ClickwrapId, "clickwrap-id", 0, "ID of the clickwrap to use with this bundle.")
	cmdCreate.Flags().Int64Var(&paramsBundleCreate.InboxId, "inbox-id", 0, "ID of the associated inbox, if available.")
	cmdCreate.Flags().BoolVar(&createRequireShareRecipient, "require-share-recipient", createRequireShareRecipient, "Only allow access to recipients who have explicitly received the share via an email sent through the Files.com UI?")
	cmdCreate.Flags().BoolVar(&createSendEmailReceiptToUploader, "send-email-receipt-to-uploader", createSendEmailReceiptToUploader, "Send delivery receipt to the uploader. Note: For writable share only")
	cmdCreate.Flags().BoolVar(&createSkipEmail, "skip-email", createSkipEmail, "BundleRegistrations can be saved without providing email?")
	cmdCreate.Flags().BoolVar(&createSkipName, "skip-name", createSkipName, "BundleRegistrations can be saved without providing name?")
	cmdCreate.Flags().BoolVar(&createSkipCompany, "skip-company", createSkipCompany, "BundleRegistrations can be saved without providing company?")
	paramsBundleCreate.StartAccessOnDate = &time.Time{}
	lib.TimeVar(cmdCreate.Flags(), paramsBundleCreate.StartAccessOnDate, "start-access-on-date", "Date when share will start to be accessible. If `nil` access granted right after create.")
	cmdCreate.Flags().Int64Var(&paramsBundleCreate.SnapshotId, "snapshot-id", 0, "ID of the snapshot containing this bundle's contents.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Bundles.AddCommand(cmdCreate)
	var fieldsShare []string
	var formatShare []string
	usePagerShare := true
	paramsBundleShare := files_sdk.BundleShareParams{}

	cmdShare := &cobra.Command{
		Use:   "share",
		Short: `Send email(s) with a link to bundle`,
		Long:  `Send email(s) with a link to bundle`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := bundle.Client{Config: config}

			var err error
			err = client.Share(paramsBundleShare, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdShare.Flags().Int64Var(&paramsBundleShare.Id, "id", 0, "Bundle ID.")
	cmdShare.Flags().StringSliceVar(&paramsBundleShare.To, "to", []string{}, "A list of email addresses to share this bundle with. Required unless `recipients` is used.")
	cmdShare.Flags().StringVar(&paramsBundleShare.Note, "note", "", "Note to include in email.")

	cmdShare.Flags().StringSliceVar(&fieldsShare, "fields", []string{}, "comma separated list of field names")
	cmdShare.Flags().StringSliceVar(&formatShare, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdShare.Flags().BoolVar(&usePagerShare, "use-pager", usePagerShare, "Use $PAGER (.ie less, more, etc)")

	Bundles.AddCommand(cmdShare)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateCreateSnapshot := true
	updateDontSeparateSubmissionsByFolder := true
	updateFinalizeSnapshot := true
	updateRequireRegistration := true
	updateRequireShareRecipient := true
	updateSendEmailReceiptToUploader := true
	updateSkipCompany := true
	updateSkipEmail := true
	updateSkipName := true
	updateWatermarkAttachmentDelete := true
	paramsBundleUpdate := files_sdk.BundleUpdateParams{}
	BundleUpdatePermissions := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Share Link`,
		Long:  `Update Share Link`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := bundle.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.BundleUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			var BundleUpdatePermissionsErr error
			paramsBundleUpdate.Permissions, BundleUpdatePermissionsErr = lib.FetchKey("permissions", paramsBundleUpdate.Permissions.Enum(), BundleUpdatePermissions)
			if BundleUpdatePermissions != "" && BundleUpdatePermissionsErr != nil {
				return BundleUpdatePermissionsErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsBundleUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("paths") {
				lib.FlagUpdateLen(cmd, "paths", paramsBundleUpdate.Paths, mapParams)
			}
			if cmd.Flags().Changed("password") {
				lib.FlagUpdate(cmd, "password", paramsBundleUpdate.Password, mapParams)
			}
			if cmd.Flags().Changed("form-field-set-id") {
				lib.FlagUpdate(cmd, "form_field_set_id", paramsBundleUpdate.FormFieldSetId, mapParams)
			}
			if cmd.Flags().Changed("clickwrap-id") {
				lib.FlagUpdate(cmd, "clickwrap_id", paramsBundleUpdate.ClickwrapId, mapParams)
			}
			if cmd.Flags().Changed("code") {
				lib.FlagUpdate(cmd, "code", paramsBundleUpdate.Code, mapParams)
			}
			if cmd.Flags().Changed("create-snapshot") {
				mapParams["create_snapshot"] = updateCreateSnapshot
			}
			if cmd.Flags().Changed("description") {
				lib.FlagUpdate(cmd, "description", paramsBundleUpdate.Description, mapParams)
			}
			if cmd.Flags().Changed("dont-separate-submissions-by-folder") {
				mapParams["dont_separate_submissions_by_folder"] = updateDontSeparateSubmissionsByFolder
			}
			if cmd.Flags().Changed("expires-at") {
				lib.FlagUpdate(cmd, "expires_at", paramsBundleUpdate.ExpiresAt, mapParams)
			}
			if cmd.Flags().Changed("finalize-snapshot") {
				mapParams["finalize_snapshot"] = updateFinalizeSnapshot
			}
			if cmd.Flags().Changed("inbox-id") {
				lib.FlagUpdate(cmd, "inbox_id", paramsBundleUpdate.InboxId, mapParams)
			}
			if cmd.Flags().Changed("max-uses") {
				lib.FlagUpdate(cmd, "max_uses", paramsBundleUpdate.MaxUses, mapParams)
			}
			if cmd.Flags().Changed("note") {
				lib.FlagUpdate(cmd, "note", paramsBundleUpdate.Note, mapParams)
			}
			if cmd.Flags().Changed("path-template") {
				lib.FlagUpdate(cmd, "path_template", paramsBundleUpdate.PathTemplate, mapParams)
			}
			if cmd.Flags().Changed("path-template-time-zone") {
				lib.FlagUpdate(cmd, "path_template_time_zone", paramsBundleUpdate.PathTemplateTimeZone, mapParams)
			}
			if cmd.Flags().Changed("permissions") {
				lib.FlagUpdate(cmd, "permissions", paramsBundleUpdate.Permissions, mapParams)
			}
			if cmd.Flags().Changed("require-registration") {
				mapParams["require_registration"] = updateRequireRegistration
			}
			if cmd.Flags().Changed("require-share-recipient") {
				mapParams["require_share_recipient"] = updateRequireShareRecipient
			}
			if cmd.Flags().Changed("send-email-receipt-to-uploader") {
				mapParams["send_email_receipt_to_uploader"] = updateSendEmailReceiptToUploader
			}
			if cmd.Flags().Changed("skip-company") {
				mapParams["skip_company"] = updateSkipCompany
			}
			if cmd.Flags().Changed("start-access-on-date") {
				lib.FlagUpdate(cmd, "start_access_on_date", paramsBundleUpdate.StartAccessOnDate, mapParams)
			}
			if cmd.Flags().Changed("skip-email") {
				mapParams["skip_email"] = updateSkipEmail
			}
			if cmd.Flags().Changed("skip-name") {
				mapParams["skip_name"] = updateSkipName
			}
			if cmd.Flags().Changed("watermark-attachment-delete") {
				mapParams["watermark_attachment_delete"] = updateWatermarkAttachmentDelete
			}
			if cmd.Flags().Changed("watermark-attachment-file") {
			}

			if paramsBundleUpdate.ExpiresAt.IsZero() {
				paramsBundleUpdate.ExpiresAt = nil
			}
			if paramsBundleUpdate.StartAccessOnDate.IsZero() {
				paramsBundleUpdate.StartAccessOnDate = nil
			}

			var bundle interface{}
			var err error
			bundle, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), bundle, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsBundleUpdate.Id, "id", 0, "Bundle ID.")
	cmdUpdate.Flags().StringSliceVar(&paramsBundleUpdate.Paths, "paths", []string{}, "A list of paths to include in this bundle.")
	cmdUpdate.Flags().StringVar(&paramsBundleUpdate.Password, "password", "", "Password for this bundle.")
	cmdUpdate.Flags().Int64Var(&paramsBundleUpdate.FormFieldSetId, "form-field-set-id", 0, "Id of Form Field Set to use with this bundle")
	cmdUpdate.Flags().Int64Var(&paramsBundleUpdate.ClickwrapId, "clickwrap-id", 0, "ID of the clickwrap to use with this bundle.")
	cmdUpdate.Flags().StringVar(&paramsBundleUpdate.Code, "code", "", "Bundle code.  This code forms the end part of the Public URL.")
	cmdUpdate.Flags().BoolVar(&updateCreateSnapshot, "create-snapshot", updateCreateSnapshot, "If true, create a snapshot of this bundle's contents.")
	cmdUpdate.Flags().StringVar(&paramsBundleUpdate.Description, "description", "", "Public description")
	cmdUpdate.Flags().BoolVar(&updateDontSeparateSubmissionsByFolder, "dont-separate-submissions-by-folder", updateDontSeparateSubmissionsByFolder, "Do not create subfolders for files uploaded to this share. Note: there are subtle security pitfalls with allowing anonymous uploads from multiple users to live in the same folder. We strongly discourage use of this option unless absolutely required.")
	paramsBundleUpdate.ExpiresAt = &time.Time{}
	lib.TimeVar(cmdUpdate.Flags(), paramsBundleUpdate.ExpiresAt, "expires-at", "Bundle expiration date/time")
	cmdUpdate.Flags().BoolVar(&updateFinalizeSnapshot, "finalize-snapshot", updateFinalizeSnapshot, "If true, finalize the snapshot of this bundle's contents. Note that `create_snapshot` must also be true.")
	cmdUpdate.Flags().Int64Var(&paramsBundleUpdate.InboxId, "inbox-id", 0, "ID of the associated inbox, if available.")
	cmdUpdate.Flags().Int64Var(&paramsBundleUpdate.MaxUses, "max-uses", 0, "Maximum number of times bundle can be accessed")
	cmdUpdate.Flags().StringVar(&paramsBundleUpdate.Note, "note", "", "Bundle internal note")
	cmdUpdate.Flags().StringVar(&paramsBundleUpdate.PathTemplate, "path-template", "", "Template for creating submission subfolders. Can use the uploader's name, email address, ip, company, `strftime` directives, and any custom form data.")
	cmdUpdate.Flags().StringVar(&paramsBundleUpdate.PathTemplateTimeZone, "path-template-time-zone", "", "Timezone to use when rendering timestamps in path templates.")
	cmdUpdate.Flags().StringVar(&BundleUpdatePermissions, "permissions", "", fmt.Sprintf("Permissions that apply to Folders in this Share Link. %v", reflect.ValueOf(paramsBundleUpdate.Permissions.Enum()).MapKeys()))
	cmdUpdate.Flags().BoolVar(&updateRequireRegistration, "require-registration", updateRequireRegistration, "Show a registration page that captures the downloader's name and email address?")
	cmdUpdate.Flags().BoolVar(&updateRequireShareRecipient, "require-share-recipient", updateRequireShareRecipient, "Only allow access to recipients who have explicitly received the share via an email sent through the Files.com UI?")
	cmdUpdate.Flags().BoolVar(&updateSendEmailReceiptToUploader, "send-email-receipt-to-uploader", updateSendEmailReceiptToUploader, "Send delivery receipt to the uploader. Note: For writable share only")
	cmdUpdate.Flags().BoolVar(&updateSkipCompany, "skip-company", updateSkipCompany, "BundleRegistrations can be saved without providing company?")
	paramsBundleUpdate.StartAccessOnDate = &time.Time{}
	lib.TimeVar(cmdUpdate.Flags(), paramsBundleUpdate.StartAccessOnDate, "start-access-on-date", "Date when share will start to be accessible. If `nil` access granted right after create.")
	cmdUpdate.Flags().BoolVar(&updateSkipEmail, "skip-email", updateSkipEmail, "BundleRegistrations can be saved without providing email?")
	cmdUpdate.Flags().BoolVar(&updateSkipName, "skip-name", updateSkipName, "BundleRegistrations can be saved without providing name?")
	cmdUpdate.Flags().BoolVar(&updateWatermarkAttachmentDelete, "watermark-attachment-delete", updateWatermarkAttachmentDelete, "If true, will delete the file stored in watermark_attachment")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Bundles.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsBundleDelete := files_sdk.BundleDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Share Link`,
		Long:  `Delete Share Link`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := bundle.Client{Config: config}

			var err error
			err = client.Delete(paramsBundleDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsBundleDelete.Id, "id", 0, "Bundle ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Bundles.AddCommand(cmdDelete)
	return Bundles
}
