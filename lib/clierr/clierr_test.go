package clierr

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCliErrorDedupesJoinedErrorMessages(t *testing.T) {
	err := errors.New("server temporarily unavailable")
	cliErr := New(ErrorCodeFatal, errors.Join(err, err, errors.New("server temporarily unavailable")))

	assert.Equal(t, "server temporarily unavailable - status (7)", cliErr.Error())
}

func TestCliErrorKeepsUniqueJoinedErrorMessages(t *testing.T) {
	cliErr := New(ErrorCodeFatal, errors.Join(
		errors.New("first failure"),
		errors.New("second failure"),
		errors.New("first failure"),
	))

	assert.Equal(t, "first failure\nsecond failure - status (7)", cliErr.Error())
}

func TestCliErrorDedupesWrappedJoinedErrorMessages(t *testing.T) {
	err := fmt.Errorf("request failed: %w", errors.Join(
		errors.New("server temporarily unavailable"),
		errors.New("server temporarily unavailable"),
	))
	cliErr := New(ErrorCodeFatal, err)

	assert.Equal(t, "request failed: server temporarily unavailable - status (7)", cliErr.Error())
}

func TestCliErrorKeepsSingleWrappedErrorContext(t *testing.T) {
	err := fmt.Errorf("open you will never find me: %w", errors.New("not found"))
	cliErr := New(ErrorCodeFatal, err)

	assert.Equal(t, "open you will never find me: not found - status (7)", cliErr.Error())
}

func TestCliErrorStillUnwrapsOriginalError(t *testing.T) {
	err := errors.New("server temporarily unavailable")
	cliErr := New(ErrorCodeFatal, errors.Join(err, errors.New("server temporarily unavailable")))

	assert.ErrorIs(t, cliErr, err)
}
