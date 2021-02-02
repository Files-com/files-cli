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
	"sync"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/dnaeon/go-vcr/cassette"
	recorder "github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
)

func captureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	writer.Close()
	return <-out
}

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

func createTempConfig() (*os.File, Config) {
	_, err := os.Stat("tmp")
	if os.IsNotExist(err) {
		os.MkdirAll("tmp", 0755)
	}
	file, err := ioutil.TempFile("tmp", "file-cli-config-test")
	if err != nil {
		log.Fatal(err)
	}

	config := Config{}
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

	files_sdk.GlobalConfig.HttpClient = httpClient

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "X-Filesapi-Auth")
		return nil
	})
	return r
}

func TestCreateSession_InvalidPassword(t *testing.T) {
	assert := assert.New(t)

	r := createRecorder("TestCreateSession_InvalidPassword")
	defer r.Stop()
	file, config := createTempConfig()
	defer os.Remove(file.Name())
	var err error
	re := captureOutput(func() {
		pipeInput("testdomain\ntestuser\n", func() {
			err = CreateSession(files_sdk.SessionCreateParams{Password: "badpassword"}, config)
		})
	})

	assert.Equal("Subdomain: Username: ", re)
	assert.Equal("Invalid username or password", err.(files_sdk.ResponseError).ErrorMessage)
}

func TestCreateSession_ValidPassword(t *testing.T) {
	assert := assert.New(t)

	r := createRecorder("TestCreateSession_ValidPassword")
	defer r.Stop()
	file, config := createTempConfig()
	defer os.Remove(file.Name())
	var err error
	re := captureOutput(func() {
		pipeInput("testdomain\ntestuser\n", func() {
			err = CreateSession(files_sdk.SessionCreateParams{Password: "goodpassword"}, config)
		})
	})

	assert.Equal("Subdomain: Username: ", re)
	assert.Equal(nil, err)
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

	re := captureOutput(func() {
		params, err = SessionUnauthorizedError(
			files_sdk.SessionCreateParams{Password: "password"},
			files_sdk.ResponseError{
				Type: "not-authenticated/two-factor-authentication-error",
				Data: files_sdk.Data{
					TwoFactorAuthenticationMethod: []string{"u2f"},
					U2fSIgnRequests:               []files_sdk.U2fSignRequests{signRequest},
					PartialSessionId:              "123456",
				},
			},
		)
	})

	if err.Error() == "failed to find any devices" {
		assert.EqualError(err, "failed to find any devices")
		assert.Equal("\n", re)
	} else {
		assert.EqualError(err, "failed to get authentication response after 25 seconds", "Unplug u2f device")
		assert.Equal("\nDevice version: U2F_V2...................................................................................................", re)
	}

	assert.Equal("", params.Password, "clears password")
	assert.Equal("123456", params.PartialSessionId, "Uses PartialSessionId instead of password")
	assert.Equal("null", params.Otp, "Otp is set to null because of error")
}

func TestCreateSession_SessionUnauthorizedError_TOTP(t *testing.T) {
	assert := assert.New(t)

	var params files_sdk.SessionCreateParams
	var err error

	re := captureOutput(func() {
		pipeInput("123456\n", func() {
			params, err = SessionUnauthorizedError(
				files_sdk.SessionCreateParams{Password: "password"},
				files_sdk.ResponseError{
					Type: "not-authenticated/two-factor-authentication-error",
					Data: files_sdk.Data{
						TwoFactorAuthenticationMethod: []string{"totp"},
					},
				},
			)
		})
	})

	assert.Nil(err, "has no error")
	assert.Equal("\ntotp: ", re, "displays prompt")

	assert.Equal("password", params.Password, "retains password")
	assert.Equal("123456", params.Otp, "populates one-time password")
}
