package lib

import (
	"context"
	"io"

	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/hashicorp/go-retryablehttp"
)

func ErrorWithOriginalResponse(err error, logger retryablehttp.Logger) (interface{}, error) {
	originalResponse, ok := err.(flib.ErrorWithOriginalResponse)
	if ok {
		logger.Printf("Recovering from original error: `%v`", originalResponse.Error())

		return originalResponse.OriginalResponse(), nil
	} else {
		return nil, err
	}
}

func HandleResponse(ctx context.Context, profile *Profiles, i interface{}, errIn error, format []string, fields []string, usePager bool, stdout io.Writer, stderr io.Writer, logger retryablehttp.Logger) error {
	var err error
	if errIn != nil {
		var originalResponse interface{}
		originalResponse, err = ErrorWithOriginalResponse(errIn, logger)
		if err == nil {
			err = Format(ctx, originalResponse, format, fields, usePager, stdout)
		}
	} else {
		err = Format(ctx, i, format, fields, usePager, stdout)
	}

	return ClientError(profile, err, stderr)
}
