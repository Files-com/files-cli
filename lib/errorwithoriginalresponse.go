package lib

import (
	"context"
	"io"

	"github.com/hashicorp/go-retryablehttp"

	flib "github.com/Files-com/files-sdk-go/v2/lib"
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

func HandleResponse(ctx context.Context, i interface{}, errIn error, format string, fields string, usePager bool, stdout io.Writer, stderr io.Writer, logger retryablehttp.Logger) {
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
	if err != nil {
		ClientError(ctx, err, stderr)
	}
}