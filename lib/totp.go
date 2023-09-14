package lib

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
)

func TotpResponse(paramsSessionCreate files_sdk.SessionCreateParams, out io.Writer) (files_sdk.SessionCreateParams, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintf(out, "totp: ")
	text, _ := reader.ReadString('\n')
	paramsSessionCreate.Otp = strings.Replace(text, "\n", "", -1)
	return paramsSessionCreate, nil
}
