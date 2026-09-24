package lib

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type failingWriter struct{}

var errOutputFailed = errors.New("synthetic output failure")

func (failingWriter) Write([]byte) (int, error) {
	return 0, errOutputFailed
}

func TestJSONEnvelopeIterReportsOutputFailure(t *testing.T) {
	err := JSONEnvelopeIter(&cursorSliceIter{}, nil, nil, false, "", failingWriter{})

	require.ErrorIs(t, err, errOutputFailed)
}
