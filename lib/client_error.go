package lib

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/lib/errors"
)

// CliClientError is a wrapper around errors that are specific to sdk responses.
func CliClientError(profile *Profiles, err error, out ...io.Writer) error {
	clientErr := clientError(profile, err, out...)
	if clientErr == nil {
		return nil
	}
	return clierr.New(clierr.ErrorCodeFatal, clientErr)
}

func clientError(profile *Profiles, err error, out ...io.Writer) error {
	if len(out) == 0 {
		out = append(out, os.Stdout)
	}
	if err == nil {
		return nil
	}
	responseError, ok := errors.As[files_sdk.ResponseError](err, files_sdk.ResponseError{})

	if ok && responseError.Type == "not-authorized/reauthentication-needed-action" {
		fmt.Fprintf(out[0], "You are authenicated via a session ID (as opposed to an API key), we require that you provide the session user’s password/2FA again for certain types of requests where we want to add an additional level of security. We call this process Reauthentication. \n")
		fmt.Fprintf(out[0], "Use --reauthentication flag to be prompted for password/2FA authentication\n")
		path, err := os.Executable()
		if err != nil {
			path = os.Args[0]
		}
		filepath.Base(path)
		fmt.Fprintf(out[0], "\n\t%v %v --reauthentication\n", filepath.Base(path), strings.Join(os.Args[1:len(os.Args)], " "))
		return err
	}

	if ok && responseError.Type == "not-authorized/authentication-required" && profile.ValidSession() {
		fmt.Fprintf(out[0], "Your session is invalid. Please login with:\n")
		fmt.Fprintf(out[0], "\tfiles-cli login\n")
	}

	if ok && responseError.Type == "not-authorized/authentication-required" && profile.Config.GetAPIKey() != "" {
		fmt.Fprintf(out[0], "Set the correct api-key with:\n")
		fmt.Fprintf(out[0], "\tfiles-cli config set --api-key=\n")
	}

	return err
}
