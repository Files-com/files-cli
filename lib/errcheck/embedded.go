package errcheck

import (
	files_sdk "github.com/Files-com/files-sdk-go/v3/file"
)

func CheckEmbeddedErrors(parent any) error {
	if jobFile, ok := parent.(files_sdk.JobFile); ok {
		if jobFile.Err != nil {
			return jobFile.Err
		}
	}
	return nil
}
