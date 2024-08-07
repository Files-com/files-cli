package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/Files-com/files-cli/lib/version"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/dnaeon/go-vcr/cassette"
	recorder "github.com/dnaeon/go-vcr/recorder"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func createRecorder(fixture string) (sdkConfig files_sdk.Config, r *recorder.Recorder) {
	var err error
	if os.Getenv("GITLAB") != "" {
		fmt.Println("using ModeReplaying")
		r, err = recorder.NewAsMode(filepath.Join("fixtures", fixture), recorder.ModeReplaying, nil)
	} else {
		r, err = recorder.New(filepath.Join("fixtures", fixture))
	}
	if err != nil {
		panic(err)
	}
	sdkConfig = files_sdk.Config{}.Init().SetCustomClient(&http.Client{
		Transport: r,
	})

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "X-Filesapi-Auth")
		return nil
	})
	return
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
		profile.Save()
	}

	t.Log("creates custom profile")
	{
		config := &files_sdk.Config{}
		config.APIKey = "xxxxxxxxxx"
		profile := &Profiles{ConfigDir: dir}
		profile.Load(config, "custom")
		profile.Save()
	}

	t.Log("creates custom profile with environment")
	{
		config := &files_sdk.Config{}
		config.APIKey = "zzzzzzzzzzz"
		config.Environment = files_sdk.Staging
		profile := &Profiles{ConfigDir: dir}
		profile.Load(config, "staging")
		profile.Save()
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

	t.Log("updates custom profile with endpoint")
	{
		config := &files_sdk.Config{APIKey: "yyyyyyyy", EndpointOverride: "http://localhost:8080"}
		profile := &Profiles{ConfigDir: dir}
		profile.Load(config, "with-custom-endpoint")
		profile.Save()
		assert.Equal(t, "yyyyyyyy", config.APIKey)
		assert.Equal(t, "http://localhost:8080", config.Endpoint())

		// Reloading works
		config = &files_sdk.Config{}
		profile = &Profiles{ConfigDir: dir}
		profile.Load(config, "with-custom-endpoint")
		assert.Equal(t, "yyyyyyyy", config.APIKey)
		assert.Equal(t, "http://localhost:8080", config.Endpoint())
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

func TestCreateSession_CustomDomain(t *testing.T) {
	assert := assert.New(t)

	sdkConfig, r := createRecorder("TestCreateSession_CustomDomain")
	defer r.Stop()
	_, config := createTempConfig(&sdkConfig)
	var err error
	stdOut := bytes.NewBufferString("")
	stdIn := &StubInput{inputs: []string{"testdomain.com", "\r", "", "testuser", "\r"}}
	config.Overrides = Overrides{Out: stdOut, In: stdIn}.Init()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err = CreateSession(ctx, files_sdk.SessionCreateParams{Password: "badpassword"}, config)
	assert.NoError(ctx.Err())
	assert.Error(err)
	assert.Equal("testdomain.com", config.Current().Endpoint)
	assert.Equal("testuser", config.Current().Username)
}

func TestCreateSession_InvalidPassword(t *testing.T) {
	assert := assert.New(t)

	sdkConfig, r := createRecorder("TestCreateSession_InvalidPassword")
	defer r.Stop()
	_, config := createTempConfig(&sdkConfig)
	var err error
	stdOut := bytes.NewBufferString("")
	stdIn := &StubInput{inputs: []string{"testdomain", "\r", "", "testuser", "\r"}}
	config.Overrides = Overrides{Out: stdOut, In: stdIn}.Init()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err = CreateSession(ctx, files_sdk.SessionCreateParams{Password: "badpassword"}, config)
	assert.NoError(ctx.Err())
	assert.Equal("testdomain", config.Current().Subdomain)
	assert.Equal("testuser", config.Current().Username)
	assert.Equal("", config.Current().SessionId)
	assert.Equal(time.Time{}, config.Current().SessionExpiry)
	var resErr files_sdk.ResponseError
	require.True(t, errors.As(err, &resErr))
	assert.Equal("Invalid username or password", resErr.ErrorMessage)
}

func TestCreateSession_ValidPassword(t *testing.T) {
	assert := assert.New(t)

	sdkConfig, r := createRecorder("TestCreateSession_ValidPassword")
	defer r.Stop()
	_, config := createTempConfig(&sdkConfig)
	var err error
	stdOut := bytes.NewBufferString("")
	stdIn := &StubInput{inputs: []string{"testdomain", "\r", "", "testuser", "\r"}}
	config.Overrides = Overrides{Out: stdOut, In: stdIn}
	parentCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	for {
		ctx, cancel := context.WithTimeout(parentCtx, time.Second*5)
		err = CreateSession(ctx, files_sdk.SessionCreateParams{Password: "goodpassword"}, config)
		if ctx.Err() == nil {
			cancel()
			break
		}
		if parentCtx.Err() != nil {
			cancel()
			break
		}
	}

	assert.NoError(parentCtx.Err())
	assert.NoError(err)
	assert.Equal("testdomain", config.Current().Subdomain)
	assert.Equal("testuser", config.Current().Username)
	assert.NotEqual("", config.Current().SessionId)
	assert.NotEqual(time.Time{}, config.Current().SessionExpiry)
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
