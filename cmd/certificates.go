package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/certificate"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Certificates())
}

func Certificates() *cobra.Command {
	Certificates := &cobra.Command{
		Use:   "certificates [command]",
		Short: "A Certificate is a way to store your custom SSL certificate for your Files.com site.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command certificates\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsCertificateList := files_sdk.CertificateListParams{}
	var MaxPagesList int64
	var jsonEnvelopeList bool

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Certificates",
		Long:    `List Certificates`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsCertificateList
			params.MaxPages = MaxPagesList
			var envelopeStyle string
			if jsonEnvelopeList {
				var envelopeErr error
				if envelopeStyle, envelopeErr = lib.PrepareJSONEnvelope(cmd, Profile(cmd).Current().SetResourceFormat(cmd, formatList), &params.MaxPages); envelopeErr != nil {
					return envelopeErr
				}
			}

			client := certificate.Client{Config: config}
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
			if jsonEnvelopeList {
				err = lib.JSONEnvelopeIter(it, fieldsList, listFilter, usePagerList, envelopeStyle, cmd.OutOrStdout())
			} else {
				err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatList), fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			}
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, "Client-side wildcard filtering, for example field-name=*.jpg or field-name=?ello")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-by", "field=pattern")

	cmdList.Flags().StringVar(&paramsCertificateList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsCertificateList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	cmdList.Flags().BoolVar(&jsonEnvelopeList, "json-envelope", false, lib.JSONEnvelopeHelpText)
	lib.SetCommandEffect(cmdList, lib.EffectReadOnly)
	lib.SetCommandResponse(cmdList, "Certificate", true)
	Certificates.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsCertificateFind := files_sdk.CertificateFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Certificate`,
		Long:  `Show Certificate`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := certificate.Client{Config: config}

			var certificate interface{}
			var err error
			certificate, err = client.Find(paramsCertificateFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), certificate, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsCertificateFind.Id, "id", 0, "Certificate ID.")
	lib.SetFlagAPIRequired(cmdFind.Flags(), "id")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	lib.SetCommandEffect(cmdFind, lib.EffectReadOnly)
	lib.SetCommandResponse(cmdFind, "Certificate", false)
	Certificates.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsCertificateCreate := files_sdk.CertificateCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create CSR or import existing SSL Certificate`,
		Long:  `Create CSR or import existing SSL Certificate`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := certificate.Client{Config: config}

			var certificate interface{}
			var err error
			certificate, err = client.Create(paramsCertificateCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), certificate, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsCertificateCreate.Name, "name", "", "Internal name of the SSL certificate.")
	lib.SetFlagAPIRequired(cmdCreate.Flags(), "name")
	cmdCreate.Flags().StringVar(&paramsCertificateCreate.CertificateDomain, "certificate-domain", "", "Domain for certificate.")
	cmdCreate.Flags().StringSliceVar(&paramsCertificateCreate.CertificateExtraDomains, "certificate-extra-domains", []string{}, "Additional domains for the certificate CSR.")
	cmdCreate.Flags().StringVar(&paramsCertificateCreate.CertificateCountry, "certificate-country", "", "Country.")
	cmdCreate.Flags().StringVar(&paramsCertificateCreate.CertificateStateOrProvince, "certificate-state-or-province", "", "State or province.")
	cmdCreate.Flags().StringVar(&paramsCertificateCreate.CertificateCityOrLocale, "certificate-city-or-locale", "", "City or locale.")
	cmdCreate.Flags().StringVar(&paramsCertificateCreate.CertificateCompanyName, "certificate-company-name", "", "Company name.")
	cmdCreate.Flags().StringVar(&paramsCertificateCreate.CsrOu1, "csr-ou1", "", "Department name or organization unit 1")
	cmdCreate.Flags().StringVar(&paramsCertificateCreate.CsrOu2, "csr-ou2", "", "Department name or organization unit 2")
	cmdCreate.Flags().StringVar(&paramsCertificateCreate.CsrOu3, "csr-ou3", "", "Department name or organization unit 3")
	cmdCreate.Flags().StringVar(&paramsCertificateCreate.CertificateEmailAddress, "certificate-email-address", "", "Email address for the certificate owner.")
	cmdCreate.Flags().StringVar(&paramsCertificateCreate.KeyType, "key-type", "", "Any supported key type.  Defaults to `RSA-4096`.")
	cmdCreate.Flags().StringVar(&paramsCertificateCreate.Certificate, "certificate", "", "The certificate.  PEM or PKCS7 formats accepted.")
	cmdCreate.Flags().StringVar(&paramsCertificateCreate.PrivateKey, "private-key", "", "Certificate private key.")
	cmdCreate.Flags().StringVar(&paramsCertificateCreate.Password, "password", "", "Certificate password.")
	cmdCreate.Flags().StringVar(&paramsCertificateCreate.Intermediates, "intermediates", "", "Intermediate certificates.  PEM or PKCS7 formats accepted.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	lib.SetCommandEffect(cmdCreate, lib.EffectMutating)
	lib.SetCommandResponse(cmdCreate, "Certificate", false)
	Certificates.AddCommand(cmdCreate)
	var fieldsActivate []string
	var formatActivate []string
	usePagerActivate := true
	paramsCertificateActivate := files_sdk.CertificateActivateParams{}

	cmdActivate := &cobra.Command{
		Use:   "activate",
		Short: `Activate SSL Certificate`,
		Long:  `Activate SSL Certificate`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := certificate.Client{Config: config}

			var err error
			err = client.Activate(paramsCertificateActivate, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdActivate.Flags().Int64Var(&paramsCertificateActivate.Id, "id", 0, "Certificate ID.")
	lib.SetFlagAPIRequired(cmdActivate.Flags(), "id")
	cmdActivate.Flags().StringVar(&paramsCertificateActivate.ReplaceCert, "replace-cert", "", "Leave blank or set to `any` to replace any currently active certificate, if applicable.  Set to `new_ips` to activate this certificate as an additional certificate on your Site by allocating a new set of Dedicated IPs (may require a Plan upgrade).  Set to the ID of a currently active certificate to replace that certificate on its set of dedicated IPs.  A certificate belonging to a Custom Domain activates only for that domain: leave blank, use `any`, or specify that domain's current certificate ID.")

	cmdActivate.Flags().StringSliceVar(&fieldsActivate, "fields", []string{}, "comma separated list of field names")
	cmdActivate.Flags().StringSliceVar(&formatActivate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdActivate.Flags().BoolVar(&usePagerActivate, "use-pager", usePagerActivate, "Use $PAGER (.ie less, more, etc)")

	lib.SetCommandEffect(cmdActivate, lib.EffectMutating)
	lib.SetCommandResponse(cmdActivate, "", false)
	Certificates.AddCommand(cmdActivate)
	var fieldsDeactivate []string
	var formatDeactivate []string
	usePagerDeactivate := true
	paramsCertificateDeactivate := files_sdk.CertificateDeactivateParams{}

	cmdDeactivate := &cobra.Command{
		Use:   "deactivate",
		Short: `Deactivate SSL Certificate`,
		Long:  `Deactivate SSL Certificate`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := certificate.Client{Config: config}

			var err error
			err = client.Deactivate(paramsCertificateDeactivate, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDeactivate.Flags().Int64Var(&paramsCertificateDeactivate.Id, "id", 0, "Certificate ID.")
	lib.SetFlagAPIRequired(cmdDeactivate.Flags(), "id")

	cmdDeactivate.Flags().StringSliceVar(&fieldsDeactivate, "fields", []string{}, "comma separated list of field names")
	cmdDeactivate.Flags().StringSliceVar(&formatDeactivate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDeactivate.Flags().BoolVar(&usePagerDeactivate, "use-pager", usePagerDeactivate, "Use $PAGER (.ie less, more, etc)")

	lib.SetCommandEffect(cmdDeactivate, lib.EffectMutating)
	lib.SetCommandResponse(cmdDeactivate, "", false)
	Certificates.AddCommand(cmdDeactivate)
	var fieldsReplaceWithFiles []string
	var formatReplaceWithFiles []string
	usePagerReplaceWithFiles := true
	paramsCertificateReplaceWithFiles := files_sdk.CertificateReplaceWithFilesParams{}

	cmdReplaceWithFiles := &cobra.Command{
		Use:   "replace-with-files",
		Short: `Replace this active custom SSL Certificate for the Primary Custom Domain with a Files.com SSL Certificate.  The Primary Custom Domain's CNAME record must point at the Site's subdomain.  The Files.com certificate is issued and activated in the background`,
		Long:  `Replace this active custom SSL Certificate for the Primary Custom Domain with a Files.com SSL Certificate.  The Primary Custom Domain's CNAME record must point at the Site's subdomain.  The Files.com certificate is issued and activated in the background`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := certificate.Client{Config: config}

			var err error
			err = client.ReplaceWithFiles(paramsCertificateReplaceWithFiles, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdReplaceWithFiles.Flags().Int64Var(&paramsCertificateReplaceWithFiles.Id, "id", 0, "Certificate ID.")
	lib.SetFlagAPIRequired(cmdReplaceWithFiles.Flags(), "id")

	cmdReplaceWithFiles.Flags().StringSliceVar(&fieldsReplaceWithFiles, "fields", []string{}, "comma separated list of field names")
	cmdReplaceWithFiles.Flags().StringSliceVar(&formatReplaceWithFiles, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdReplaceWithFiles.Flags().BoolVar(&usePagerReplaceWithFiles, "use-pager", usePagerReplaceWithFiles, "Use $PAGER (.ie less, more, etc)")

	lib.SetCommandEffect(cmdReplaceWithFiles, lib.EffectMutating)
	lib.SetCommandResponse(cmdReplaceWithFiles, "", false)
	Certificates.AddCommand(cmdReplaceWithFiles)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsCertificateUpdate := files_sdk.CertificateUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Certificate`,
		Long:  `Update Certificate`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := certificate.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.CertificateUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsCertificateUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsCertificateUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("intermediates") {
				lib.FlagUpdate(cmd, "intermediates", paramsCertificateUpdate.Intermediates, mapParams)
			}
			if cmd.Flags().Changed("certificate") {
				lib.FlagUpdate(cmd, "certificate", paramsCertificateUpdate.Certificate, mapParams)
			}

			var certificate interface{}
			var err error
			certificate, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), certificate, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsCertificateUpdate.Id, "id", 0, "Certificate ID.")
	lib.SetFlagAPIRequired(cmdUpdate.Flags(), "id")
	cmdUpdate.Flags().StringVar(&paramsCertificateUpdate.Name, "name", "", "Internal certificate name.")
	cmdUpdate.Flags().StringVar(&paramsCertificateUpdate.Intermediates, "intermediates", "", "Any intermediate certificates.  PEM or PKCS7 formats accepted.")
	cmdUpdate.Flags().StringVar(&paramsCertificateUpdate.Certificate, "certificate", "", "The certificate.  PEM or PKCS7 formats accepted.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	lib.SetCommandEffect(cmdUpdate, lib.EffectDestructive)
	lib.SetCommandResponse(cmdUpdate, "Certificate", false)
	Certificates.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsCertificateDelete := files_sdk.CertificateDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Certificate`,
		Long:  `Delete Certificate`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := certificate.Client{Config: config}

			var err error
			err = client.Delete(paramsCertificateDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsCertificateDelete.Id, "id", 0, "Certificate ID.")
	lib.SetFlagAPIRequired(cmdDelete.Flags(), "id")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	lib.SetCommandEffect(cmdDelete, lib.EffectDestructive)
	lib.SetCommandResponse(cmdDelete, "", false)
	Certificates.AddCommand(cmdDelete)
	return Certificates
}
