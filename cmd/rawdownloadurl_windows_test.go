//go:build windows

package cmd

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows"
)

// daclSummary reports whether path has a protected DACL and whether any allow
// entry names Everyone.
func daclSummary(t *testing.T, path string) (protected bool, everyone bool) {
	t.Helper()
	sd, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	require.NoError(t, err)
	control, _, err := sd.Control()
	require.NoError(t, err)
	dacl, _, err := sd.DACL()
	require.NoError(t, err)
	if dacl == nil {
		return control&windows.SE_DACL_PROTECTED != 0, true
	}
	world, err := windows.CreateWellKnownSid(windows.WinWorldSid)
	require.NoError(t, err)
	for i := uint32(0); i < uint32(dacl.AceCount); i++ {
		var ace *windows.ACCESS_ALLOWED_ACE
		require.NoError(t, windows.GetAce(dacl, i, &ace))
		if ace.Header.AceType == windows.ACCESS_ALLOWED_ACE_TYPE && (*windows.SID)(unsafe.Pointer(&ace.SidStart)).Equals(world) {
			everyone = true
		}
	}
	return control&windows.SE_DACL_PROTECTED != 0, everyone
}

func setDACL(t *testing.T, path string, sddl string, flags windows.SECURITY_INFORMATION) {
	t.Helper()
	sd, err := windows.SecurityDescriptorFromString(sddl)
	require.NoError(t, err)
	acl, _, err := sd.DACL()
	require.NoError(t, err)
	require.NoError(t, windows.SetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION|flags, nil, nil, acl, nil))
}

// An existing destination with its own restricted list keeps it after the
// diagnostic download replaces its content, even in a directory that hands
// Everyone read access down to new files; a fresh destination inherits from
// the directory.
func TestRawDownloadURL_publishedAccessOnWindows(t *testing.T) {
	const content = "diagnostic download bytes"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, "file.bin", time.Unix(1, 0), strings.NewReader(content))
	}))
	defer server.Close()
	dir := t.TempDir()
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	require.NoError(t, err)
	// The directory hands Everyone read access down to what is created inside.
	setDACL(t, dir, "D:(A;OICI;FA;;;SY)(A;OICI;FA;;;BA)(A;OICI;FA;;;"+user.User.Sid.String()+")(A;OICI;GR;;;WD)", windows.PROTECTED_DACL_SECURITY_INFORMATION)

	restricted := filepath.Join(dir, "restricted.bin")
	require.NoError(t, os.WriteFile(restricted, []byte("old private content"), 0o600))
	setDACL(t, restricted, "D:P(A;;FA;;;SY)(A;;FA;;;BA)(A;;FA;;;"+user.User.Sid.String()+")", windows.PROTECTED_DACL_SECURITY_INFORMATION)
	options := rawDownloadURLOptions{partSizeMiB: 1, initialConcurrency: 1, maxConcurrency: 1}
	require.NoError(t, runRawDownloadURL(context.Background(), io.Discard, server.URL, restricted, options))
	protected, everyone := daclSummary(t, restricted)
	assert.True(t, protected, "the replaced restricted file keeps its own list")
	assert.False(t, everyone, "the replaced restricted file must not gain Everyone access")
	data, err := os.ReadFile(restricted)
	require.NoError(t, err)
	assert.Equal(t, content, string(data))

	fresh := filepath.Join(dir, "fresh.bin")
	require.NoError(t, runRawDownloadURL(context.Background(), io.Discard, server.URL, fresh, options))
	protected, everyone = daclSummary(t, fresh)
	assert.False(t, protected, "a fresh destination inherits from its directory")
	assert.True(t, everyone, "a fresh destination has what the directory hands down")
}
