package lib

import (
	"context"
	"fmt"
	"io"

	"github.com/hashicorp/go-retryablehttp"

	flib "github.com/Files-com/files-sdk-go/v2/lib"
)

func ErrorWithOriginalResponse(err error, format string, logger retryablehttp.Logger) (interface{}, error) {
	originalResponse, ok := err.(flib.ErrorWithOriginalResponse)
	if ok {
		if format != "json" {
			logger.Printf("Original error: `%v`", err.Error())
			return nil, fmt.Errorf("the server returned an unexpected response can only format response as json (use `--format json`)")
		} else {
			logger.Printf("Recovering from original error: `%v`", err.Error())
			return originalResponse.OriginalResponse(), nil
		}
	} else {
		return nil, err
	}
}

func HandleResponse(ctx context.Context, i interface{}, errIn error, format string, fields string, usePager bool, stdout io.Writer, stderr io.Writer, logger retryablehttp.Logger) {
	if errIn != nil {
		originalResponse, err := ErrorWithOriginalResponse(errIn, format, logger)
		if err != nil {
			ClientError(ctx, err, stderr)
		} else {
			Format(ctx, originalResponse, "json", "", usePager, stdout)
		}
	} else {
		err := Format(ctx, i, format, fields, usePager, stdout)
		if err != nil {
			ClientError(ctx, err, stderr)
		}
	}
}
