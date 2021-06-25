package cmd

import (
	"testing"

	clib "github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestExternalEventsCreateSuccess(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestExternalEventsCreate")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	ExternalEventsInit()
	str := clib.CaptureOutput(func() {
		out, err := callCmd(ExternalEvents, config, []string{"create", "--status", "success", "--body", "this is a success test"})
		assert.NoError(err)
		assert.Equal("", out)
	})
	event := files_sdk.ExternalEvent{}
	event.UnmarshalJSON([]byte(str))
	assert.Equal("this is a success test", event.Body)
	assert.Equal("client_log", event.EventType)
	assert.Equal("success", event.Status)
}

func TestExternalEventsCreateError(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestExternalEventsCreateError")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	ExternalEventsInit()
	str := clib.CaptureOutput(func() {
		out, err := callCmd(ExternalEvents, config, []string{"create", "--status", "error", "--body", "this is a error test"})
		assert.NoError(err)
		assert.Equal("", out)
	})
	event := files_sdk.ExternalEvent{}
	err = event.UnmarshalJSON([]byte(str))
	assert.NoError(err)
	assert.Equal("this is a error test", event.Body)
	assert.Equal("client_log", event.EventType)
	assert.Equal("error", event.Status)
}

func TestExternalEventsCreateNoStatus(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestExternalEventsCreateNoStatus")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	ExternalEventsInit()
	str := clib.CaptureOutput(func() {
		out, err := callCmd(ExternalEvents, config, []string{"create", "--status", "taco", "--body", "this is a error test"})
		assert.NoError(err)
		assert.Equal("", out)
	})
	event := files_sdk.ExternalEvent{}
	err = event.UnmarshalJSON([]byte(str))
	assert.Error(err)
	assert.Contains(str, "missing required field: ExternalEventCreateParams{}.Status")
}
