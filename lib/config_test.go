package lib

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Files-com/files-cli/lib/version"

	"github.com/Files-com/files-sdk-go/v2/lib"

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

func createTempConfig() (*os.File, *Config) {
	_, err := os.Stat("tmp")
	if os.IsNotExist(err) {
		os.MkdirAll("tmp", 0755)
	}
	file, err := ioutil.TempFile("tmp", "file-cli-config-test")
	if err != nil {
		log.Fatal(err)
	}

	config := &Config{}
	config.configPathOverride = file.Name()
	return file, config
}

func createRecorder(fixture string) *recorder.Recorder {
	var r *recorder.Recorder
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
	httpClient := &http.Client{
		Transport: r,
	}
	files_sdk.GlobalConfig.Debug = lib.Bool(false)
	files_sdk.GlobalConfig.SetHttpClient(httpClient)

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "X-Filesapi-Auth")
		return nil
	})
	return r
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

func TestCreateSession_InvalidPassword(t *testing.T) {
	assert := assert.New(t)

	r := createRecorder("TestCreateSession_InvalidPassword")
	defer r.Stop()
	file, config := createTempConfig()
	defer os.Remove(file.Name())
	var err error
	stdOut := bytes.NewBufferString("")
	stdIn := &StubInput{inputs: []string{"testdomain", "\r", "", "testuser", "\r"}}
	config.Overrides = Overrides{Out: stdOut, In: stdIn}.Init()
	err = CreateSession(files_sdk.SessionCreateParams{Password: "badpassword"}, config)

	assert.Equal("testdomain", config.Subdomain)
	assert.Equal("testuser", config.Username)
	assert.Equal("Invalid username or password", err.(files_sdk.ResponseError).ErrorMessage)
}

func TestCreateSession_ValidPassword(t *testing.T) {
	assert := assert.New(t)

	r := createRecorder("TestCreateSession_ValidPassword")
	defer r.Stop()
	file, config := createTempConfig()
	defer os.Remove(file.Name())
	var err error
	stdOut := bytes.NewBufferString("")
	stdIn := &StubInput{inputs: []string{"testdomain", "\r", "", "testuser", "\r"}}
	config.Overrides = Overrides{Out: stdOut, In: stdIn}
	err = CreateSession(files_sdk.SessionCreateParams{Password: "goodpassword"}, config)
	assert.NoError(err)
	assert.Equal("testdomain", config.Subdomain)
	assert.Equal("testuser", config.Username)
}

func TestCreateSession_SessionUnauthorizedError_U2F(t *testing.T) {
	assert := assert.New(t)
	signRequest := files_sdk.U2fSignRequests{
		Challenge:   "taco",
		AppId:       "taco.com",
		SignRequest: files_sdk.SignRequest{KeyHandle: "xxxxx"},
	}

	var params files_sdk.SessionCreateParams
	var err error

	stdOut := bytes.NewBufferString("")
	params, err = SessionUnauthorizedError(
		files_sdk.SessionCreateParams{Password: "password"},
		files_sdk.ResponseError{
			Type:         "not-authenticated/two-factor-authentication-error",
			ErrorMessage: "2FA Authentication error: Token from U2F is required",
			Data: files_sdk.Data{
				TwoFactorAuthenticationMethod: []string{"u2f"},
				U2fSIgnRequests:               []files_sdk.U2fSignRequests{signRequest},
				PartialSessionId:              "123456",
			},
		},
		Config{Overrides: Overrides{Out: stdOut, Timeout: time.Second * 5}},
	)

	if err.Error() == "failed to find any devices" {
		assert.EqualError(err, "failed to find any devices")
		assert.Equal("Token from U2F is required\n", stdOut.String())
	} else {
		assert.EqualError(err, "failed to get authentication response after 25 seconds", "Unplug u2f device")
		assert.Contains(stdOut.String(), "Device version: U2F_V2")
	}

	assert.Equal("", params.Password, "clears password")
	assert.Equal("123456", params.PartialSessionId, "Uses PartialSessionId instead of password")
	assert.Equal("null", params.Otp, "Otp is set to null because of error")
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
			Config{Overrides: Overrides{Out: stdOut, Timeout: time.Second * 5}},
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
		Config
		installedViaBrew bool
	}
	tests := []struct {
		name       string
		args       args
		wantStdout string
		Config
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
			wantStdout: "files-cli version 1.1.0 is out of date\nDownload latest version from\nhttps://github.com/Files-com/files-cli/releases\n\n",
			Config:     Config{VersionOutOfDate: true, LastVersionCheck: time.Now().Local()},
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
			wantStdout: "files-cli version 1.1.0 is out of date\nUpgrade via Homebrew\n\tbrew upgrade files-cli\n\n",
			Config:     Config{VersionOutOfDate: true, LastVersionCheck: time.Now().Local()},
		},
		{
			name: "already checked yesterday",
			args: args{
				version: "1.1.0",
				fetchLatestVersion: func() (version.Version, bool) {
					return version.Version{Major: 1, Minor: 2, Patch: 9}, true
				},
				installedViaBrew: true,
				Config:           Config{LastVersionCheck: time.Now().Add(-24 * time.Hour)},
			},
			wantStdout: "",
			Config:     Config{LastVersionCheck: time.Now().Add(-24 * time.Hour)},
		},
		{
			name: "already checked 3 days ago",
			args: args{
				version: "1.1.0",
				fetchLatestVersion: func() (version.Version, bool) {
					return version.Version{Major: 1, Minor: 2, Patch: 9}, true
				},
				installedViaBrew: true,
				Config:           Config{LastVersionCheck: time.Now().Add(-(24 * time.Hour) * 3)},
			},
			wantStdout: "files-cli version 1.1.0 is out of date\nUpgrade via Homebrew\n\tbrew upgrade files-cli\n\n",
			Config:     Config{VersionOutOfDate: true, LastVersionCheck: time.Now().Local()},
		},
		{
			name: "already known out of date",
			args: args{
				version: "1.1.0",
				fetchLatestVersion: func() (version.Version, bool) {
					return version.Version{Major: 1, Minor: 2, Patch: 9}, true
				},
				installedViaBrew: true,
				Config:           Config{LastVersionCheck: time.Now().Add(-24 * time.Hour), VersionOutOfDate: true},
			},
			wantStdout: "",
			Config:     Config{VersionOutOfDate: true, LastVersionCheck: time.Now().Add(-24 * time.Hour)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdOut := bytes.NewBufferString("")
			tt.args.Config.CheckVersion(tt.args.version, tt.args.fetchLatestVersion, tt.args.installedViaBrew, stdOut)
			assert.Equal(t, tt.wantStdout, stdOut.String())
			assert.Equal(t, tt.Config.VersionOutOfDate, tt.args.Config.VersionOutOfDate)
			assert.Equal(t, tt.Config.LastVersionCheck.Truncate(60*time.Second), tt.args.Config.LastVersionCheck.Truncate(60*time.Second))
		})
	}
}
