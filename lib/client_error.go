package lib

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
)

func ClientError(ctx context.Context, profile *Profiles, err error, out ...io.Writer) error {
	if len(out) == 0 {
		out = append(out, os.Stdout)
	}
	if err == nil {
		return nil
	}
	responseError, ok := err.(files_sdk.ResponseError)

	if ok && responseError.Type == "not-authorized/reauthentication-needed-action" {
		fmt.Fprintf(out[0], "You are authenicated via a session ID (as opposed to an API key), we require that you provide the session user’s password/2FA again for certain types of requests where we want to add an additional level of security. We call this process Reauthentication. \n")
		fmt.Fprintf(out[0], "Use --reauthentication flag to be prompted for password/2FA authentication\n")
		path, err := os.Executable()
		if err != nil {
			path = os.Args[0]
		}
		filepath.Base(path)
		fmt.Fprintf(out[0], "\n\t%v %v --reauthentication\n", filepath.Base(path), strings.Join(os.Args[1:len(os.Args)], " "))
		exit(ctx)
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

	if ok {
		fmt.Fprintf(out[0], "%v\n", responseError)
	} else {
		fmt.Fprintf(out[0], "%v\n", err)
	}

	exit(ctx)
	return err
}

func exit(ctx context.Context) {
	if ctx.Value("testing") != nil && !ctx.Value("testing").(bool) {
		os.Exit(1)
	}
}
