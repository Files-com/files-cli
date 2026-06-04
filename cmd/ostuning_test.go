package cmd

import (
	"bytes"
	"runtime"
	"strings"
	"testing"

	"github.com/Files-com/files-sdk-go/v3/lib/ostuning"
	"github.com/stretchr/testify/require"
)

func TestOSTuningHighThroughputLinuxPlan(t *testing.T) {
	root := OSTuning()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"high-throughput", "--os", "linux", "--interface", "ens4"})

	require.NoError(t, root.Execute())

	output := out.String()
	require.Contains(t, output, "High-throughput upload OS tuning")
	require.Contains(t, output, "Target OS: linux")
	require.Contains(t, output, "Files Agent already has Linux UDP buffer tuning")
	require.Contains(t, output, "Snapshot before repair")
	require.Contains(t, output, "/var/lib/files.com/os-tuning/high-throughput-upload.snapshot")
	require.Contains(t, output, "Restore")
	require.Contains(t, output, "Requires: root or Administrator privileges")
	require.Contains(t, output, "sysctl --system")
	require.NotContains(t, output, "sudo sysctl --system")
	require.Contains(t, output, "ens4")
}

func TestOSTuningHighThroughputHelpIncludesCommonUsage(t *testing.T) {
	root := OSTuning()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"high-throughput", "--help"})

	require.NoError(t, root.Execute())

	output := out.String()
	require.Contains(t, output, "Examples:")
	require.Contains(t, output, "files-cli os-tuning high-throughput verify")
	require.Contains(t, output, "files-cli os-tuning high-throughput restore --apply")
	switch runtime.GOOS {
	case "darwin":
		require.Contains(t, output, "files-cli os-tuning high-throughput repair --include-network-test --apply")
		require.NotContains(t, output, "--interface ens4")
	case "linux":
		require.Contains(t, output, "files-cli os-tuning high-throughput repair --interface ens4 --commands-only")
		require.NotContains(t, output, "networkQuality")
	case "windows":
		require.Contains(t, output, `files-cli os-tuning high-throughput repair --interface "Ethernet 2" --commands-only`)
		require.NotContains(t, output, "networkQuality")
	}
}

func TestOSTuningHighThroughputExamplesAreHostSpecific(t *testing.T) {
	darwinExamples := osTuningHighThroughputExamples("darwin")
	require.Contains(t, darwinExamples, "networkQuality")
	require.NotContains(t, darwinExamples, "--interface ens4")
	require.NotContains(t, darwinExamples, `"Ethernet 2"`)

	linuxExamples := osTuningHighThroughputExamples("linux")
	require.Contains(t, linuxExamples, "--interface ens4")
	require.NotContains(t, linuxExamples, "networkQuality")
	require.NotContains(t, linuxExamples, `"Ethernet 2"`)

	windowsExamples := osTuningHighThroughputExamples("windows")
	require.Contains(t, windowsExamples, `"Ethernet 2"`)
	require.NotContains(t, windowsExamples, "--interface ens4")
	require.NotContains(t, windowsExamples, "networkQuality")
}

func TestOSTuningHighThroughputCommandsOnly(t *testing.T) {
	root := OSTuning()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"high-throughput", "repair", "--os", "windows", "--interface", "Ethernet 2", "--commands-only"})

	require.NoError(t, root.Execute())

	output := out.String()
	require.Contains(t, output, "# Repair workflow")
	require.Contains(t, output, "Snapshot current high-throughput tuning values")
	require.Contains(t, output, "# Requires root or Administrator privileges.")
	require.Contains(t, output, "netsh interface tcp set global autotuninglevel=normal rss=enabled")
	require.Contains(t, output, `Enable-NetAdapterRss -Name 'Ethernet 2'`)
	require.Contains(t, output, `Enable-NetAdapterRsc -Name 'Ethernet 2' -IPv4`)
	require.Contains(t, output, "No changes were applied")
}

func TestOSTuningHighThroughputVerifyMode(t *testing.T) {
	root := OSTuning()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"high-throughput", "verify", "--os", "linux", "--commands-only"})

	require.NoError(t, root.Execute())

	output := out.String()
	require.Contains(t, output, "Verify")
	require.Contains(t, output, "Inspect current TCP settings")
	require.Contains(t, output, "Verify persist high-throughput TCP defaults")
}

func TestRenderOSTuningRunningStepIncludesHumanReadableContext(t *testing.T) {
	var out bytes.Buffer
	renderOSTuningRunningStep(&out, ostuning.Step{
		Title:           "Inspect current TCP buffer settings",
		Description:     "Read macOS socket buffer and TCP auto-buffer ceilings.",
		ExpectedOutcome: "Supported macOS releases report larger TCP buffer ceilings.",
	})

	output := out.String()
	require.Contains(t, output, "Inspect current TCP buffer settings")
	require.Contains(t, output, "Read macOS socket buffer and TCP auto-buffer ceilings.")
	require.Contains(t, output, "Expected: Supported macOS releases report larger TCP buffer ceilings.")
}

func TestOSTuningHighThroughputVerifyRejectsUnsupportedOS(t *testing.T) {
	root := OSTuning()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"high-throughput", "verify", "--os", "plan9"})

	err := root.Execute()
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported OS")
}

func TestOSTuningHighThroughputVerifyRejectsHostMismatchWhenRunning(t *testing.T) {
	targetOS := "linux"
	if runtime.GOOS == "linux" {
		targetOS = "darwin"
	}

	root := OSTuning()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"high-throughput", "verify", "--os", targetOS})

	err := root.Execute()
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot run "+targetOS+" commands on "+runtime.GOOS)
	require.Contains(t, err.Error(), "--commands-only")
}

func TestOSTuningHighThroughputDarwinNetworkQualityIsOptIn(t *testing.T) {
	root := OSTuning()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"high-throughput", "verify", "--os", "darwin", "--commands-only"})

	require.NoError(t, root.Execute())
	require.NotContains(t, out.String(), "networkQuality")

	root = OSTuning()
	out.Reset()
	root.SetOut(&out)
	root.SetArgs([]string{"high-throughput", "verify", "--os", "darwin", "--include-network-test", "--commands-only"})

	require.NoError(t, root.Execute())
	require.Contains(t, out.String(), "networkQuality -v")

	root = OSTuning()
	out.Reset()
	root.SetOut(&out)
	root.SetArgs([]string{"high-throughput", "repair", "--os", "darwin", "--include-network-test", "--commands-only"})

	require.NoError(t, root.Execute())
	output := out.String()
	require.Contains(t, output, "Measure Apple network quality before repair")
	require.Contains(t, output, "Measure Apple network quality after repair")
	require.Equal(t, 2, strings.Count(output, "networkQuality -v"))
}
