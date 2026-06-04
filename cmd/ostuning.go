package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Files-com/files-sdk-go/v3/lib/ostuning"
	"github.com/spf13/cobra"
)

const commandNameOSTuning = "os-tuning"

func init() {
	osTuning := OSTuning()
	RootCmd.AddCommand(osTuning)
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, osTuning.Use)
}

func OSTuning() *cobra.Command {
	cmd := &cobra.Command{
		Use:   commandNameOSTuning,
		Short: "Print OS tuning recommendations for high-throughput transfers",
	}

	cmd.AddCommand(osTuningHighThroughput())
	return cmd
}

type osTuningHighThroughputOptions struct {
	targetOS      string
	interfaceName string
	commandsOnly  bool
	apply         bool
	networkTest   bool
	binaryName    string
}

func osTuningHighThroughput() *cobra.Command {
	options := &osTuningHighThroughputOptions{}

	cmd := &cobra.Command{
		Use:     "high-throughput",
		Short:   "Inspect, plan, or repair high-throughput upload OS tuning",
		Example: osTuningHighThroughputExamples(runtime.GOOS),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runOSTuningPlan(cmd, options)
		},
	}
	addOSTuningPlanFlags(cmd, options)

	plan := &cobra.Command{
		Use:   "plan",
		Short: "Show the dry-run tuning plan",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runOSTuningPlan(cmd, options)
		},
	}
	addOSTuningPlanFlags(plan, options)

	verify := &cobra.Command{
		Use:   "verify",
		Short: "Show or run non-mutating verification checks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			plan, err := ostuning.HighThroughputUploadPlan(ostuning.Options{OS: options.targetOS, InterfaceName: options.interfaceName, IncludeNetworkTest: options.networkTest})
			if err != nil {
				return err
			}
			steps := plan.VerificationSteps()
			if options.commandsOnly {
				renderOSTuningCommandMode(cmd.OutOrStdout(), plan, "Verify", steps, true)
				return nil
			}
			return runOSTuningStepCommands(cmd, plan, steps, options, "verify")
		},
	}
	addOSTuningPlanFlags(verify, options)

	repair := &cobra.Command{
		Use:   "repair",
		Short: "Show or apply privileged tuning changes",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			plan, err := ostuning.HighThroughputUploadPlan(ostuning.Options{OS: options.targetOS, InterfaceName: options.interfaceName, IncludeNetworkTest: options.networkTest})
			if err != nil {
				return err
			}
			steps := plan.RepairSteps()
			if options.apply {
				return runOSTuningRepairCommands(cmd, plan, options)
			}
			if options.commandsOnly {
				renderOSTuningCommandMode(cmd.OutOrStdout(), plan, "Repair workflow", steps, true)
				fmt.Fprintf(cmd.OutOrStdout(), "\nNo changes were applied. Re-run with --apply to execute the steps allowed by the current process privileges.\n")
				return nil
			}
			renderOSTuningCommandMode(cmd.OutOrStdout(), plan, "Repair dry run", steps, options.commandsOnly)
			fmt.Fprintf(cmd.OutOrStdout(), "\nNo changes were applied. Re-run with --apply to execute the steps allowed by the current process privileges.\n")
			return nil
		},
	}
	addOSTuningPlanFlags(repair, options)
	repair.Flags().BoolVar(&options.apply, "apply", false, "Apply privileged repair commands. Only supported when --os matches this host.")

	restore := &cobra.Command{
		Use:   "restore",
		Short: "Show or apply restore from snapshot or system defaults",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			plan, err := ostuning.HighThroughputUploadPlan(ostuning.Options{OS: options.targetOS, InterfaceName: options.interfaceName, IncludeNetworkTest: options.networkTest})
			if err != nil {
				return err
			}
			steps := plan.RestorePlanSteps()
			if options.apply {
				return runOSTuningStepCommands(cmd, plan, steps, options, "restore")
			}
			renderOSTuningCommandMode(cmd.OutOrStdout(), plan, "Restore dry run", steps, options.commandsOnly)
			fmt.Fprintf(cmd.OutOrStdout(), "\nNo changes were applied. Re-run with --apply to execute the restore steps allowed by the current process privileges.\n")
			return nil
		},
	}
	addOSTuningPlanFlags(restore, options)
	restore.Flags().BoolVar(&options.apply, "apply", false, "Apply restore commands. Only supported when --os matches this host.")

	cmd.AddCommand(plan, verify, repair, restore)
	return cmd
}

func osTuningHighThroughputExamples(targetOS string) string {
	common := `  # Inspect the current host without changing anything.
  files-cli os-tuning high-throughput verify

  # Show the full dry-run tuning plan.
  files-cli os-tuning high-throughput plan

  # Apply what the current process can change and print an elevated follow-up command for the rest.
  files-cli os-tuning high-throughput repair --apply

  # Restore from the saved snapshot or the safest OS defaults.
  files-cli os-tuning high-throughput restore --apply`

	switch targetOS {
	case "darwin":
		return common + `

  # On macOS, include before/after networkQuality measurements around repair.
  files-cli os-tuning high-throughput repair --include-network-test --apply`
	case "linux":
		return common + `

  # On Linux, include the active NIC when printing repair commands for review or automation.
  files-cli os-tuning high-throughput repair --interface ens4 --commands-only`
	case "windows":
		return common + `

  # On Windows, include the adapter name when printing repair commands for review or automation.
  files-cli os-tuning high-throughput repair --interface "Ethernet 2" --commands-only`
	default:
		return common
	}
}

func addOSTuningPlanFlags(cmd *cobra.Command, options *osTuningHighThroughputOptions) {
	cmd.Flags().StringVar(&options.targetOS, "os", "", "Target OS to render: linux, darwin, windows. Defaults to this host.")
	cmd.Flags().StringVar(&options.interfaceName, "interface", "", "Network interface or adapter name for OS-specific NIC tuning commands.")
	cmd.Flags().BoolVar(&options.commandsOnly, "commands-only", false, "Only print the commands, grouped by required privilege.")
	cmd.Flags().BoolVar(&options.networkTest, "include-network-test", false, "Include OS-supported active network measurement commands where available.")
}

func runOSTuningPlan(cmd *cobra.Command, options *osTuningHighThroughputOptions) error {
	plan, err := ostuning.HighThroughputUploadPlan(ostuning.Options{
		OS:                 options.targetOS,
		InterfaceName:      options.interfaceName,
		IncludeNetworkTest: options.networkTest,
	})
	if err != nil {
		return err
	}
	if options.commandsOnly {
		renderOSTuningCommands(cmd.OutOrStdout(), plan)
		return nil
	}
	renderOSTuningPlan(cmd.OutOrStdout(), plan)
	return nil
}

func runOSTuningStepCommands(cmd *cobra.Command, plan ostuning.Plan, steps []ostuning.Step, options *osTuningHighThroughputOptions, action string) error {
	if plan.OS != runtime.GOOS {
		return fmt.Errorf("cannot run %s commands on %s; use --commands-only to print commands for another OS", plan.OS, runtime.GOOS)
	}

	runOptions := *options
	runOptions.binaryName = rootCommandName(cmd)

	ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Minute)
	defer cancel()

	results := ostuning.RunSteps(ctx, steps, ostuning.RunOptions{
		Stdout:      cmd.OutOrStdout(),
		Stderr:      cmd.ErrOrStderr(),
		StopOnError: true,
		BeforeStep: func(step ostuning.Step) {
			renderOSTuningRunningStep(cmd.OutOrStdout(), step)
		},
		BeforeCommand: func(command ostuning.Command) {
			fmt.Fprintf(cmd.OutOrStdout(), "\n$ %s\n", command.CommandLine)
		},
	})

	var skipped []ostuning.Step
	for _, result := range results {
		if result.SkippedForPrivilege {
			skipped = append(skipped, result.Step)
			continue
		}
		if result.Err == nil {
			continue
		}
		if result.SoftFailed {
			fmt.Fprintf(cmd.ErrOrStderr(), "Warning: %s failed: %v\n", result.Step.Title, result.Err)
			continue
		}
		renderOSTuningSkippedSteps(cmd.OutOrStdout(), plan, skipped, &runOptions, action)
		return result.Err
	}
	renderOSTuningSkippedSteps(cmd.OutOrStdout(), plan, skipped, &runOptions, action)
	return nil
}

func runOSTuningRepairCommands(cmd *cobra.Command, plan ostuning.Plan, options *osTuningHighThroughputOptions) error {
	if plan.OS != runtime.GOOS {
		return fmt.Errorf("cannot run %s commands on %s; use --commands-only to print commands for another OS", plan.OS, runtime.GOOS)
	}

	runOptions := *options
	runOptions.binaryName = rootCommandName(cmd)

	ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Minute)
	defer cancel()

	runner := func(steps []ostuning.Step) []ostuning.StepResult {
		return ostuning.RunSteps(ctx, steps, ostuning.RunOptions{
			Stdout:      cmd.OutOrStdout(),
			Stderr:      cmd.ErrOrStderr(),
			StopOnError: true,
			BeforeStep: func(step ostuning.Step) {
				renderOSTuningRunningStep(cmd.OutOrStdout(), step)
			},
			BeforeCommand: func(command ostuning.Command) {
				fmt.Fprintf(cmd.OutOrStdout(), "\n$ %s\n", command.CommandLine)
			},
		})
	}

	var skipped []ostuning.Step
	adminStepIDs := ostuningStepIDSet(plan.AdminSteps)

	preflightSummary, err := processOSTuningResults(cmd, runner(plan.RepairPreflightSteps()), nil)
	skipped = append(skipped, preflightSummary.skipped...)
	if err != nil {
		renderOSTuningSkippedSteps(cmd.OutOrStdout(), plan, skipped, &runOptions, "repair")
		return err
	}

	changeSummary, err := processOSTuningResults(cmd, runner(plan.RepairChangeSteps()), adminStepIDs)
	skipped = append(skipped, changeSummary.skipped...)
	if err != nil {
		renderOSTuningSkippedSteps(cmd.OutOrStdout(), plan, skipped, &runOptions, "repair")
		return err
	}

	postRepairSteps := plan.RepairPostChangeSteps()
	if len(postRepairSteps) > 0 {
		if changeSummary.trackedStepRan {
			postSummary, err := processOSTuningResults(cmd, runner(postRepairSteps), nil)
			skipped = append(skipped, postSummary.skipped...)
			if err != nil {
				renderOSTuningSkippedSteps(cmd.OutOrStdout(), plan, skipped, &runOptions, "repair")
				return err
			}
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "\nSkipping post-repair network measurement because no privileged tuning commands were applied.\n")
			fmt.Fprintf(cmd.OutOrStdout(), "Re-run the elevated follow-up command to capture a true before/after comparison.\n")
		}
	}

	renderOSTuningSkippedSteps(cmd.OutOrStdout(), plan, skipped, &runOptions, "repair")
	return nil
}

type ostuningRunSummary struct {
	skipped        []ostuning.Step
	trackedStepRan bool
}

func processOSTuningResults(cmd *cobra.Command, results []ostuning.StepResult, trackedStepIDs map[string]struct{}) (ostuningRunSummary, error) {
	var summary ostuningRunSummary
	for _, result := range results {
		if result.SkippedForPrivilege {
			summary.skipped = append(summary.skipped, result.Step)
			continue
		}
		if _, tracked := trackedStepIDs[result.Step.ID]; tracked && len(result.CommandResults) > 0 {
			summary.trackedStepRan = true
		}
		if result.Err == nil {
			continue
		}
		if result.SoftFailed {
			fmt.Fprintf(cmd.ErrOrStderr(), "Warning: %s failed: %v\n", result.Step.Title, result.Err)
			continue
		}
		return summary, result.Err
	}
	return summary, nil
}

func ostuningStepIDSet(steps []ostuning.Step) map[string]struct{} {
	ids := make(map[string]struct{}, len(steps))
	for _, step := range steps {
		ids[step.ID] = struct{}{}
	}
	return ids
}

func renderOSTuningCommandMode(writer io.Writer, plan ostuning.Plan, title string, steps []ostuning.Step, commandsOnly bool) {
	if commandsOnly {
		fmt.Fprintf(writer, "# Files.com high-throughput upload OS tuning\n")
		fmt.Fprintf(writer, "# Target OS: %s\n", plan.OS)
		fmt.Fprintf(writer, "# Mode: %s\n", title)
		renderOSTuningCommandGroup(writer, title, steps)
		return
	}

	fmt.Fprintf(writer, "%s\n\n", title)
	fmt.Fprintf(writer, "Target OS: %s\n", plan.OS)
	fmt.Fprintf(writer, "Profile: %s\n", plan.Profile)
	renderOSTuningStepGroup(writer, title, steps)
}

func renderOSTuningPlan(writer io.Writer, plan ostuning.Plan) {
	fmt.Fprintf(writer, "High-throughput upload OS tuning\n\n")
	fmt.Fprintf(writer, "Target OS: %s\n", plan.OS)
	fmt.Fprintf(writer, "Profile: %s\n", plan.Profile)
	if plan.InterfaceName != "" {
		fmt.Fprintf(writer, "Interface: %s\n", plan.InterfaceName)
	}
	if plan.SnapshotPath != "" {
		fmt.Fprintf(writer, "Snapshot: %s\n", plan.SnapshotPath)
	}
	fmt.Fprintf(writer, "\n%s\n", plan.Summary)

	renderOSTuningStepGroup(writer, "User-level checks", plan.UserSteps)
	renderOSTuningStepGroup(writer, "Network tests", plan.NetworkTests)
	renderOSTuningStepGroup(writer, "Snapshot before repair", plan.SnapshotSteps)
	renderOSTuningStepGroup(writer, "Privileged changes", plan.AdminSteps)
	renderOSTuningStepGroup(writer, "Restore", plan.RestoreSteps)

	if len(plan.Warnings) > 0 {
		fmt.Fprintf(writer, "\nWarnings\n")
		for _, warning := range plan.Warnings {
			fmt.Fprintf(writer, "- %s\n", warning)
		}
	}
	if len(plan.Notes) > 0 {
		fmt.Fprintf(writer, "\nNotes\n")
		for _, note := range plan.Notes {
			fmt.Fprintf(writer, "- %s\n", note)
		}
	}
	if len(plan.References) > 0 {
		fmt.Fprintf(writer, "\nReferences\n")
		for _, reference := range plan.References {
			fmt.Fprintf(writer, "- %s: %s\n", reference.Title, reference.URL)
		}
	}
}

func renderOSTuningStepGroup(writer io.Writer, title string, steps []ostuning.Step) {
	if len(steps) == 0 {
		return
	}

	fmt.Fprintf(writer, "\n%s\n", title)
	for index, step := range steps {
		fmt.Fprintf(writer, "%d. %s\n", index+1, step.Title)
		fmt.Fprintf(writer, "   %s\n", step.Description)
		if step.Privilege == ostuning.PrivilegeAdministrator {
			fmt.Fprintf(writer, "   Requires: root or Administrator privileges\n")
		}
		if step.RuntimeOnly {
			fmt.Fprintf(writer, "   Runtime-only: yes\n")
		}
		if step.RequiresReboot {
			fmt.Fprintf(writer, "   Requires reboot: yes\n")
		}
		if step.CanFailSoftly {
			fmt.Fprintf(writer, "   If unsupported: continue after recording the error\n")
		}
		if step.ExpectedOutcome != "" {
			fmt.Fprintf(writer, "   Expected: %s\n", step.ExpectedOutcome)
		}
		for _, command := range step.Commands {
			fmt.Fprintf(writer, "   %s> %s\n", command.Shell, indentContinuation(command.CommandLine, "      "))
		}
		if len(step.Verification) > 0 {
			fmt.Fprintf(writer, "   Verify:\n")
			for _, command := range step.Verification {
				fmt.Fprintf(writer, "   %s> %s\n", command.Shell, indentContinuation(command.CommandLine, "      "))
			}
		}
	}
}

func renderOSTuningCommands(writer io.Writer, plan ostuning.Plan) {
	fmt.Fprintf(writer, "# Files.com high-throughput upload OS tuning\n")
	fmt.Fprintf(writer, "# Target OS: %s\n", plan.OS)
	fmt.Fprintf(writer, "# Profile: %s\n", plan.Profile)
	if plan.SnapshotPath != "" {
		fmt.Fprintf(writer, "# Snapshot: %s\n", plan.SnapshotPath)
	}
	renderOSTuningCommandGroup(writer, "User-level checks", plan.UserSteps)
	renderOSTuningCommandGroup(writer, "Network tests", plan.NetworkTests)
	renderOSTuningCommandGroup(writer, "Snapshot before repair", plan.SnapshotSteps)
	renderOSTuningCommandGroup(writer, "Privileged changes", plan.AdminSteps)
	renderOSTuningCommandGroup(writer, "Restore", plan.RestoreSteps)
}

func renderOSTuningCommandGroup(writer io.Writer, title string, steps []ostuning.Step) {
	if len(steps) == 0 {
		return
	}
	fmt.Fprintf(writer, "\n# %s\n", title)
	for _, step := range steps {
		fmt.Fprintf(writer, "\n# %s\n", step.Title)
		if step.Privilege == ostuning.PrivilegeAdministrator {
			fmt.Fprintf(writer, "# Requires root or Administrator privileges.\n")
		}
		for _, command := range step.Commands {
			fmt.Fprintf(writer, "%s\n", command.CommandLine)
		}
		if len(step.Verification) > 0 {
			fmt.Fprintf(writer, "# Verify %s\n", lowerFirst(step.Title))
			for _, command := range step.Verification {
				fmt.Fprintf(writer, "%s\n", command.CommandLine)
			}
		}
	}
}

func renderOSTuningSkippedSteps(writer io.Writer, plan ostuning.Plan, skipped []ostuning.Step, options *osTuningHighThroughputOptions, action string) {
	if len(skipped) == 0 {
		return
	}

	fmt.Fprintf(writer, "\nSome OS tuning steps require elevated privileges and were not applied:\n")
	for _, step := range skipped {
		fmt.Fprintf(writer, "- %s\n", step.Title)
	}

	command := elevatedOSTuningCommand(plan, options, action)
	renderOSTuningElevationHint(writer, command)
}

func renderOSTuningRunningStep(writer io.Writer, step ostuning.Step) {
	fmt.Fprintf(writer, "\n%s\n", step.Title)
	if step.Description != "" {
		fmt.Fprintf(writer, "%s\n", step.Description)
	}
	if step.ExpectedOutcome != "" {
		fmt.Fprintf(writer, "Expected: %s\n", step.ExpectedOutcome)
	}
}

func elevatedOSTuningCommand(plan ostuning.Plan, options *osTuningHighThroughputOptions, action string) string {
	if action == "" || action == "verify" {
		action = "repair"
	}
	parts := []string{ostuningBinaryName(options), "os-tuning", "high-throughput", action, "--apply", "--os", quoteCLIArg(plan.OS)}
	if options != nil && options.interfaceName != "" {
		parts = append(parts, "--interface", quoteCLIArg(options.interfaceName))
	}
	if options != nil && options.networkTest {
		parts = append(parts, "--include-network-test")
	}
	return strings.Join(parts, " ")
}

func ostuningBinaryName(options *osTuningHighThroughputOptions) string {
	if options != nil && options.binaryName != "" {
		return options.binaryName
	}
	if RootCmd.Name() != "" {
		return RootCmd.Name()
	}
	return ""
}

func rootCommandName(cmd *cobra.Command) string {
	if executable := filepath.Base(os.Args[0]); executable != "." {
		return executable
	}
	if cmd != nil && cmd.Root() != nil {
		return cmd.Root().Name()
	}
	return ""
}

func lowerFirst(value string) string {
	if value == "" {
		return value
	}
	return strings.ToLower(value[:1]) + value[1:]
}

func indentContinuation(value string, prefix string) string {
	return strings.ReplaceAll(value, "\n", "\n"+prefix)
}
