//go:build linux || darwin

package cmd

import (
	"bytes"
	"testing"

	"github.com/Files-com/files-sdk-go/v3/lib/ostuning"
	"github.com/stretchr/testify/require"
)

func TestOSTuningSkippedPrivilegedStepsPrintsPOSIXFollowUp(t *testing.T) {
	plan, err := ostuning.HighThroughputUploadPlan(ostuning.Options{OS: "linux", InterfaceName: "ens4"})
	require.NoError(t, err)

	var out bytes.Buffer
	renderOSTuningSkippedSteps(&out, plan, plan.AdminSteps, &osTuningHighThroughputOptions{interfaceName: "ens4"}, "repair")

	output := out.String()
	require.Contains(t, output, "Some OS tuning steps require elevated privileges and were not applied")
	require.Contains(t, output, "Persist high-throughput TCP defaults")
	require.Contains(t, output, "sudo files-cli os-tuning high-throughput repair --apply --os 'linux' --interface 'ens4'")
}

func TestOSTuningSkippedPrivilegedStepsUsesConfiguredBinaryName(t *testing.T) {
	plan, err := ostuning.HighThroughputUploadPlan(ostuning.Options{OS: "linux", InterfaceName: "ens4"})
	require.NoError(t, err)

	var out bytes.Buffer
	renderOSTuningSkippedSteps(&out, plan, plan.AdminSteps, &osTuningHighThroughputOptions{interfaceName: "ens4", binaryName: "files-cli-dev"}, "repair")

	require.Contains(t, out.String(), "sudo files-cli-dev os-tuning high-throughput repair --apply --os 'linux' --interface 'ens4'")
}

func TestOSTuningSkippedPrivilegedStepsPreservesNetworkTestFollowUp(t *testing.T) {
	plan, err := ostuning.HighThroughputUploadPlan(ostuning.Options{OS: "darwin", IncludeNetworkTest: true})
	require.NoError(t, err)

	var out bytes.Buffer
	renderOSTuningSkippedSteps(&out, plan, plan.AdminSteps, &osTuningHighThroughputOptions{networkTest: true}, "repair")

	require.Contains(t, out.String(), "sudo files-cli os-tuning high-throughput repair --apply --os 'darwin' --include-network-test")
}

func TestOSTuningSkippedPrivilegedStepsQuotePOSIXFollowUp(t *testing.T) {
	interfaceName := "eth$Primary's`name`"
	plan, err := ostuning.HighThroughputUploadPlan(ostuning.Options{OS: "linux", InterfaceName: interfaceName})
	require.NoError(t, err)

	var out bytes.Buffer
	renderOSTuningSkippedSteps(&out, plan, plan.AdminSteps, &osTuningHighThroughputOptions{interfaceName: interfaceName}, "repair")

	output := out.String()
	require.Contains(t, output, "sudo files-cli os-tuning high-throughput repair --apply --os 'linux' --interface 'eth$Primary'\\''s`name`'")
	require.NotContains(t, output, `"eth$Primary's`)
}

func TestOSTuningRestoreSkippedPrivilegedStepsPrintsPOSIXRestoreFollowUp(t *testing.T) {
	plan, err := ostuning.HighThroughputUploadPlan(ostuning.Options{OS: "linux", InterfaceName: "ens4"})
	require.NoError(t, err)

	var out bytes.Buffer
	renderOSTuningSkippedSteps(&out, plan, plan.RestoreSteps, &osTuningHighThroughputOptions{interfaceName: "ens4"}, "restore")

	output := out.String()
	require.Contains(t, output, "Restore high-throughput tuning from snapshot or defaults")
	require.Contains(t, output, "sudo files-cli os-tuning high-throughput restore --apply --os 'linux' --interface 'ens4'")
}
