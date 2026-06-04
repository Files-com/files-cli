//go:build windows

package cmd

import (
	"bytes"
	"testing"

	"github.com/Files-com/files-sdk-go/v3/lib/ostuning"
	"github.com/stretchr/testify/require"
)

func TestOSTuningWindowsFollowUpUsesAdministratorPowerShell(t *testing.T) {
	plan, err := ostuning.HighThroughputUploadPlan(ostuning.Options{OS: "windows", InterfaceName: "Ethernet 2"})
	require.NoError(t, err)

	var out bytes.Buffer
	renderOSTuningSkippedSteps(&out, plan, plan.AdminSteps[:1], &osTuningHighThroughputOptions{interfaceName: "Ethernet 2"}, "repair")

	output := out.String()
	require.Contains(t, output, "open PowerShell as Administrator")
	require.Contains(t, output, "files-cli os-tuning high-throughput repair --apply --os 'windows' --interface 'Ethernet 2'")
	require.NotContains(t, output, "sudo")
}
