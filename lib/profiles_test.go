package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/samber/lo"

	"github.com/stretchr/testify/require"

	"github.com/Files-com/files-cli/lib/version"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/dnaeon/go-vcr/cassette"
	recorder "github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
)

func pipeInput(input string, f func()) {
	inputBytes := []byte(input)

	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	_, err = w.Write(inputBytes)
	if err != nil {
		panic(err)
	}
	err = w.Close()
	if err != nil {
		panic(err)
	}
	stdin := os.Stdin
	defer func() { os.Stdin = stdin }()
	os.Stdin = r
	f()
}

func createTempConfig(sdkConfig *files_sdk.Config) (string, *Profiles) {
	dir, err := os.MkdirTemp(os.TempDir(), "file-cli-config-test")
	if err != nil {
		log.Fatal(err)
	}
	config := &Profiles{ConfigDir: dir}
	config.Load(sdkConfig, "")
	return filepath.Join(dir, "files-cli"), config
}

func createRecorder(fixture string) (sdkConfig *files_sdk.Config, r *recorder.Recorder) {
	var err error
	sdkConfig = &files_sdk.Config{}
	if os.Getenv("GITLAB") != "" {
		fmt.Println("using ModeReplaying")
		r, err = recorder.NewAsMode(filepath.Join("fixtures", fixture), recorder.ModeReplaying, nil)
	} else {
		r, err = recorder.New(filepath.Join("fixtures", fixture))
	}
	if err != nil {
		panic(err)
	}
	httpClient := &http.Client{
		Transport: r,
	}
	sdkConfig.SetHttpClient(httpClient)

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "X-Filesapi-Auth")
		return nil
	})
	return sdkConfig, r
}

type StubInput struct {
	index  int
	inputs []string
}

func (s *StubInput) Read(b []byte) (int, error) {
	if s.index+1 > len(s.inputs) {
		return 0, nil
	}
	defer func() { s.index = s.index + 1 }()
	value := s.inputs[s.index]
	if value == "" {
		return 0, io.EOF
	}
	return bytes.NewBufferString(s.inputs[s.index]).Read(b)
}

func TestProfiles_Load(t *testing.T) {
	dir, err := os.MkdirTemp(os.TempDir(), "file-cli-config-test")
	if err != nil {
		log.Fatal(err)
	}

	t.Log("creates default profile")
	{
		config := &files_sdk.Config{}
		config.APIKey = "123456789"
		profile := &Profiles{ConfigDir: dir}
		profile.Load(config, "")
	}

	t.Log("creates custom profile")
	{
		config := &files_sdk.Config{}
		config.APIKey = "xxxxxxxxxx"
		profile := &Profiles{ConfigDir: dir}
		profile.Load(config, "custom")
	}

	t.Log("creates custom profile with environment")
	{
		config := &files_sdk.Config{}
		config.APIKey = "zzzzzzzzzzz"
		config.Environment = files_sdk.Staging
		profile := &Profiles{ConfigDir: dir}
		profile.Load(config, "staging")
	}

	config := &files_sdk.Config{}
	profile := &Profiles{ConfigDir: dir}
	profile.Load(config, "")
	profileNames := lo.Keys[string, *Profile](profile.Profiles)
	sort.Strings(profileNames)
	require.Equal(t, []string{"custom", "default", "staging"}, profileNames)

	assert.Equal(t, "123456789", profile.Profiles["default"].APIKey)
	assert.Equal(t, "xxxxxxxxxx", profile.Profiles["custom"].APIKey)
	assert.Equal(t, "zzzzzzzzzzz", profile.Profiles["staging"].APIKey)

	t.Log("loads default profile")
	{
		config := &files_sdk.Config{}
		profile := &Profiles{ConfigDir: dir}
		profile.Load(config, "")
		assert.Equal(t, "123456789", config.APIKey)
		assert.Equal(t, files_sdk.Production, config.Environment)
	}

	t.Log("loads custom profile")
	{
		config := &files_sdk.Config{}
		profile := &Profiles{ConfigDir: dir}
		profile.Load(config, "custom")
		assert.Equal(t, "xxxxxxxxxx", config.APIKey)
		assert.Equal(t, files_sdk.Production, config.Environment)
	}

	t.Log("loads custom profile with environment")
	{
		config := &files_sdk.Config{}
		profile := &Profiles{ConfigDir: dir}
		profile.Load(config, "staging")
		assert.Equal(t, "zzzzzzzzzzz", config.APIKey)
		assert.Equal(t, files_sdk.Staging, config.Environment)
	}

	t.Log("updates custom profile with environment")
	{
		config := &files_sdk.Config{APIKey: "yyyyyyyy"}
		profile := &Profiles{ConfigDir: dir}
		profile.Load(config, "staging")
		assert.Equal(t, "yyyyyyyy", config.APIKey)
		assert.Equal(t, files_sdk.Staging, config.Environment)
	}

	t.Log("upgrades v1 to v2")
	{
		v1Profile := Profile{APIKey: "xxxxxx"}
		j, _ := json.MarshalIndent(v1Profile, "", " ")
		os.WriteFile(filepath.Join(dir, "files-cli"), j, 0644)

		config := &files_sdk.Config{}
		profile := &Profiles{ConfigDir: dir}
		err := profile.Load(config, "")
		require.NoError(t, err)

		assert.Equal(t, "xxxxxx", config.APIKey)
	}
}

func TestCreateSession_InvalidPassword(t *testing.T) {
	assert := assert.New(t)

	sdkConfig, r := createRecorder("TestCreateSession_InvalidPassword")
	defer r.Stop()
	_, config := createTempConfig(sdkConfig)
	var err error
	stdOut := bytes.NewBufferString("")
	stdIn := &StubInput{inputs: []string{"testdomain", "\r", "", "testuser", "\r"}}
	config.Overrides = Overrides{Out: stdOut, In: stdIn}.Init()
	err = CreateSession(files_sdk.SessionCreateParams{Password: "badpassword"}, config)

	assert.Equal("testdomain", config.Current().Subdomain)
	assert.Equal("testuser", config.Current().Username)
	assert.Equal("Invalid username or password", err.(files_sdk.ResponseError).ErrorMessage)
}

func TestCreateSession_ValidPassword(t *testing.T) {
	assert := assert.New(t)

	sdkConfig, r := createRecorder("TestCreateSession_ValidPassword")
	defer r.Stop()
	_, config := createTempConfig(sdkConfig)
	var err error
	stdOut := bytes.NewBufferString("")
	stdIn := &StubInput{inputs: []string{"testdomain", "\r", "", "testuser", "\r"}}
	config.Overrides = Overrides{Out: stdOut, In: stdIn}
	err = CreateSession(files_sdk.SessionCreateParams{Password: "goodpassword"}, config)
	assert.NoError(err)
	assert.Equal("testdomain", config.Current().Subdomain)
	assert.Equal("testuser", config.Current().Username)
}

func TestCreateSession_SessionUnauthorizedError_TOTP(t *testing.T) {
	assert := assert.New(t)

	var params files_sdk.SessionCreateParams
	var err error

	stdOut := bytes.NewBufferString("")
	pipeInput("123456\n", func() {
		params, err = SessionUnauthorizedError(
			files_sdk.SessionCreateParams{Password: "password"},
			files_sdk.ResponseError{
				Type: "not-authenticated/two-factor-authentication-error",
				Data: files_sdk.Data{
					TwoFactorAuthenticationMethod: []string{"totp"},
				},
			},
			&Profiles{Overrides: Overrides{Out: stdOut, Timeout: time.Second * 5}},
		)
	})

	assert.Nil(err, "has no error")
	assert.Equal("\ntotp: ", stdOut.String(), "displays prompt")

	assert.Equal("password", params.Password, "retains password")
	assert.Equal("123456", params.Otp, "populates one-time password")
}

func TestConfig_CheckVersion(t *testing.T) {
	type args struct {
		version            string
		fetchLatestVersion func() (version.Version, bool)
		Profile
		installedViaBrew bool
	}
	tests := []struct {
		name       string
		args       args
		wantStdout string
		Profile
	}{
		{
			name: "Custom install with old version",
			args: args{
				version: "1.1.0",
				fetchLatestVersion: func() (version.Version, bool) {
					return version.Version{Major: 1, Minor: 2, Patch: 9}, true
				},
				installedViaBrew: false,
			},
			wantStdout: "files-cli version 1.1.0 is out of date. Latest version is 1.2.9\nDownload latest version from\nhttps://github.com/Files-com/files-cli/releases\n\n",
			Profile:    Profile{LastValidVersionCheck: time.Time{}},
		},
		{
			name: "brew install with old version",
			args: args{
				version: "1.1.0",
				fetchLatestVersion: func() (version.Version, bool) {
					return version.Version{Major: 1, Minor: 2, Patch: 9}, true
				},
				installedViaBrew: true,
			},
			wantStdout: "files-cli version 1.1.0 is out of date. Latest version is 1.2.9\nUpgrade via Homebrew\n\tbrew upgrade files-cli\n\n",
			Profile:    Profile{LastValidVersionCheck: time.Time{}},
		},
		{
			name: "already checked yesterday",
			args: args{
				version: "1.1.0",
				fetchLatestVersion: func() (version.Version, bool) {
					return version.Version{Major: 1, Minor: 2, Patch: 9}, true
				},
				installedViaBrew: true,
				Profile:          Profile{LastValidVersionCheck: time.Now().Add(-24 * time.Hour)},
			},
			wantStdout: "",
			Profile:    Profile{LastValidVersionCheck: time.Now().Add(-24 * time.Hour)},
		},
		{
			name: "already checked 3 days ago",
			args: args{
				version: "1.1.0",
				fetchLatestVersion: func() (version.Version, bool) {
					return version.Version{Major: 1, Minor: 2, Patch: 9}, true
				},
				installedViaBrew: true,
				Profile:          Profile{LastValidVersionCheck: time.Now().Add(-(24 * time.Hour) * 3)},
			},
			wantStdout: "files-cli version 1.1.0 is out of date. Latest version is 1.2.9\nUpgrade via Homebrew\n\tbrew upgrade files-cli\n\n",
			Profile:    Profile{LastValidVersionCheck: time.Now().Add(-(24 * time.Hour) * 3)},
		},
		{
			name: "out of date but not checking because within 48 hours of last good check",
			args: args{
				version: "1.1.0",
				fetchLatestVersion: func() (version.Version, bool) {
					return version.Version{Major: 1, Minor: 2, Patch: 9}, true
				},
				installedViaBrew: true,
				Profile:          Profile{LastValidVersionCheck: time.Now().Add(-24 * time.Hour)},
			},
			wantStdout: "",
			Profile:    Profile{LastValidVersionCheck: time.Now().Add(-24 * time.Hour)},
		},
		{
			name: "was out of date but client was upgraded",
			args: args{
				version: "1.2.9",
				fetchLatestVersion: func() (version.Version, bool) {
					return version.Version{Major: 1, Minor: 2, Patch: 9}, true
				},
				installedViaBrew: true,
				Profile:          Profile{LastValidVersionCheck: time.Now().Add(-(24 * time.Hour) * 3)},
			},
			wantStdout: "",
			Profile:    Profile{LastValidVersionCheck: time.Now()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdOut := bytes.NewBufferString("")
			tempDir, err := os.MkdirTemp(os.TempDir(), tt.name)
			require.NoError(t, err)
			profile := (&Profiles{ConfigDir: tempDir}).Init()
			profile.Profile = tt.name
			profile.Profiles[tt.name] = &tt.args.Profile
			profile.CheckVersion(tt.args.version, tt.args.fetchLatestVersion, tt.args.installedViaBrew, stdOut)
			assert.Equal(t, tt.wantStdout, stdOut.String())
			assert.Equal(t, tt.LastValidVersionCheck.Truncate(60*time.Second), tt.args.LastValidVersionCheck.Truncate(60*time.Second))
		})
	}
}
