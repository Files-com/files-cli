package cmd

import (
	"context"
	"testing"

	cliLib "github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadTempProfile(t *testing.T, dir string) *cliLib.Profiles {
	t.Helper()
	profile := &cliLib.Profiles{ConfigDir: dir}
	require.NoError(t, profile.Load(&files_sdk.Config{}, ""))
	return profile
}

func runConfigCmd(t *testing.T, profile *cliLib.Profiles, args ...string) {
	t.Helper()
	cmd := Config()
	cmd.SetArgs(args)
	ctx := context.WithValue(context.Background(), "profile", profile)
	require.NoError(t, cmd.ExecuteContext(ctx))
}

func TestConfigSetDirectTransfers(t *testing.T) {
	dir := t.TempDir()
	profile := loadTempProfile(t, dir)
	profile.Current().DisableDirectTransfers = true
	require.NoError(t, profile.Save())

	runConfigCmd(t, profile, "set", "--direct-transfers")
	assert.False(t, profile.Current().DisableDirectTransfers)

	reloaded := loadTempProfile(t, dir)
	assert.False(t, reloaded.Current().DisableDirectTransfers)
}

func TestConfigSetDirectTransfersExplicitFalse(t *testing.T) {
	dir := t.TempDir()
	profile := loadTempProfile(t, dir)

	runConfigCmd(t, profile, "set", "--direct-transfers=false")
	assert.True(t, profile.Current().DisableDirectTransfers)

	reloaded := loadTempProfile(t, dir)
	assert.True(t, reloaded.Current().DisableDirectTransfers)
}

func TestConfigSetWithoutDirectTransfersLeavesValue(t *testing.T) {
	dir := t.TempDir()
	profile := loadTempProfile(t, dir)
	profile.Current().DisableDirectTransfers = true
	require.NoError(t, profile.Save())

	runConfigCmd(t, profile, "set", "--subdomain", "testdomain")
	assert.True(t, profile.Current().DisableDirectTransfers)
	assert.Equal(t, "testdomain", profile.Current().Subdomain)

	reloaded := loadTempProfile(t, dir)
	assert.True(t, reloaded.Current().DisableDirectTransfers)
}

func TestConfigResetDirectTransfers(t *testing.T) {
	dir := t.TempDir()
	profile := loadTempProfile(t, dir)
	profile.Current().DisableDirectTransfers = true
	profile.Current().ConcurrentConnectionLimit = 5
	require.NoError(t, profile.Save())

	runConfigCmd(t, profile, "reset", "--direct-transfers")
	assert.False(t, profile.Current().DisableDirectTransfers)
	assert.Equal(t, 5, profile.Current().ConcurrentConnectionLimit, "reset only touches the requested option")

	reloaded := loadTempProfile(t, dir)
	assert.False(t, reloaded.Current().DisableDirectTransfers)
	assert.Equal(t, 5, reloaded.Current().ConcurrentConnectionLimit)
}
