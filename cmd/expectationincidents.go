package cmd

import (
	"time"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	expectation_incident "github.com/Files-com/files-sdk-go/v3/expectationincident"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(ExpectationIncidents())
}

func ExpectationIncidents() *cobra.Command {
	ExpectationIncidents := &cobra.Command{
		Use:  "expectation-incidents [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command expectation-incidents\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsExpectationIncidentList := files_sdk.ExpectationIncidentListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Expectation Incidents",
		Long:    `List Expectation Incidents`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsExpectationIncidentList
			params.MaxPages = MaxPagesList

			client := expectation_incident.Client{Config: config}
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

	cmdList.Flags().StringVar(&paramsExpectationIncidentList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsExpectationIncidentList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	ExpectationIncidents.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsExpectationIncidentFind := files_sdk.ExpectationIncidentFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Expectation Incident`,
		Long:  `Show Expectation Incident`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := expectation_incident.Client{Config: config}

			var expectationIncident interface{}
			var err error
			expectationIncident, err = client.Find(paramsExpectationIncidentFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), expectationIncident, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsExpectationIncidentFind.Id, "id", 0, "Expectation Incident ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	ExpectationIncidents.AddCommand(cmdFind)
	var fieldsResolve []string
	var formatResolve []string
	usePagerResolve := true
	paramsExpectationIncidentResolve := files_sdk.ExpectationIncidentResolveParams{}

	cmdResolve := &cobra.Command{
		Use:   "resolve",
		Short: `Resolve an expectation incident`,
		Long:  `Resolve an expectation incident`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := expectation_incident.Client{Config: config}

			var expectationIncident interface{}
			var err error
			expectationIncident, err = client.Resolve(paramsExpectationIncidentResolve, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), expectationIncident, err, Profile(cmd).Current().SetResourceFormat(cmd, formatResolve), fieldsResolve, usePagerResolve, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdResolve.Flags().Int64Var(&paramsExpectationIncidentResolve.Id, "id", 0, "Expectation Incident ID.")

	cmdResolve.Flags().StringSliceVar(&fieldsResolve, "fields", []string{}, "comma separated list of field names")
	cmdResolve.Flags().StringSliceVar(&formatResolve, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdResolve.Flags().BoolVar(&usePagerResolve, "use-pager", usePagerResolve, "Use $PAGER (.ie less, more, etc)")

	ExpectationIncidents.AddCommand(cmdResolve)
	var fieldsSnooze []string
	var formatSnooze []string
	usePagerSnooze := true
	paramsExpectationIncidentSnooze := files_sdk.ExpectationIncidentSnoozeParams{}

	cmdSnooze := &cobra.Command{
		Use:   "snooze",
		Short: `Snooze an expectation incident until a specified time`,
		Long:  `Snooze an expectation incident until a specified time`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := expectation_incident.Client{Config: config}

			if paramsExpectationIncidentSnooze.SnoozedUntil.IsZero() {
				paramsExpectationIncidentSnooze.SnoozedUntil = nil
			}

			var expectationIncident interface{}
			var err error
			expectationIncident, err = client.Snooze(paramsExpectationIncidentSnooze, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), expectationIncident, err, Profile(cmd).Current().SetResourceFormat(cmd, formatSnooze), fieldsSnooze, usePagerSnooze, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdSnooze.Flags().Int64Var(&paramsExpectationIncidentSnooze.Id, "id", 0, "Expectation Incident ID.")
	paramsExpectationIncidentSnooze.SnoozedUntil = &time.Time{}
	lib.TimeVar(cmdSnooze.Flags(), paramsExpectationIncidentSnooze.SnoozedUntil, "snoozed-until", "Time until which the incident should remain snoozed.")

	cmdSnooze.Flags().StringSliceVar(&fieldsSnooze, "fields", []string{}, "comma separated list of field names")
	cmdSnooze.Flags().StringSliceVar(&formatSnooze, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdSnooze.Flags().BoolVar(&usePagerSnooze, "use-pager", usePagerSnooze, "Use $PAGER (.ie less, more, etc)")

	ExpectationIncidents.AddCommand(cmdSnooze)
	var fieldsAcknowledge []string
	var formatAcknowledge []string
	usePagerAcknowledge := true
	paramsExpectationIncidentAcknowledge := files_sdk.ExpectationIncidentAcknowledgeParams{}

	cmdAcknowledge := &cobra.Command{
		Use:   "acknowledge",
		Short: `Acknowledge an expectation incident`,
		Long:  `Acknowledge an expectation incident`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := expectation_incident.Client{Config: config}

			var expectationIncident interface{}
			var err error
			expectationIncident, err = client.Acknowledge(paramsExpectationIncidentAcknowledge, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), expectationIncident, err, Profile(cmd).Current().SetResourceFormat(cmd, formatAcknowledge), fieldsAcknowledge, usePagerAcknowledge, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdAcknowledge.Flags().Int64Var(&paramsExpectationIncidentAcknowledge.Id, "id", 0, "Expectation Incident ID.")

	cmdAcknowledge.Flags().StringSliceVar(&fieldsAcknowledge, "fields", []string{}, "comma separated list of field names")
	cmdAcknowledge.Flags().StringSliceVar(&formatAcknowledge, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdAcknowledge.Flags().BoolVar(&usePagerAcknowledge, "use-pager", usePagerAcknowledge, "Use $PAGER (.ie less, more, etc)")

	ExpectationIncidents.AddCommand(cmdAcknowledge)
	return ExpectationIncidents
}
