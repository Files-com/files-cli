package lib

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/Files-com/files-cli/lib/version"
	"github.com/Files-com/files-sdk-go/v2/session"
	"golang.org/x/crypto/ssh/terminal"

	"fmt"
	"io"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
)

type Overrides struct {
	Out     io.Writer
	In      io.Reader
	Timeout time.Duration
}

func (o Overrides) Init() Overrides {
	o.Timeout = time.Second * 25
	return o
}

type Config struct {
	Overrides          `json:"-"`
	SessionId          string    `json:"session_id"`
	SessionExpiry      time.Time `json:"session_expiry"`
	LastVersionCheck   time.Time `json:"last_version_check"`
	VersionOutOfDate   bool      `json:"version_out_of_date"`
	Subdomain          string    `json:"subdomain"`
	Username           string    `json:"username"`
	APIKey             string    `json:"api_key"`
	Endpoint           string    `json:"endpoint,omitempty"`
	configPathOverride string
}

type ResetConfig struct {
	Subdomain    bool
	Username     bool
	APIKey       bool
	Endpoint     bool
	Session      bool
	VersionCheck bool
}

var SessionExpiry = time.Hour * 6
var CheckVersionEvery = time.Hour * 48

const CLICurrentVersionURL = "https://raw.githubusercontent.com/Files-com/files-cli/master/_VERSION"

func (c Config) ResetWith(reset ResetConfig) error {
	if reset.Subdomain {
		c.Subdomain = ""
	}
	if reset.Username {
		c.Username = ""
	}
	if reset.APIKey {
		c.APIKey = ""
	}
	if reset.Endpoint {
		c.Endpoint = ""
	}
	if reset.Session {
		c.SessionId = ""
	}
	if reset.VersionCheck {
		c.VersionOutOfDate = false
		c.LastVersionCheck = time.Now()
	}
	return c.Save()
}

func (c Config) Reset() error {
	return initConfig()
}

func (c *Config) Load() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	configRoot := filepath.Join(usr.HomeDir, ".config")
	configPath := filepath.Join(configRoot, "files-cli")
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		initConfig()
		return nil
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &c)

	if err != nil {
		return err
	}

	c.SetGlobal()
	return nil
}

func (c *Config) SetGlobal() {
	files_sdk.GlobalConfig.SessionId = c.SessionId
	files_sdk.GlobalConfig.Subdomain = c.Subdomain
	files_sdk.GlobalConfig.APIKey = c.APIKey
	files_sdk.GlobalConfig.Endpoint = c.Endpoint
}

func (c *Config) Save() error {
	file, _ := json.MarshalIndent(c, "", " ")
	path, err := c.configPath()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, file, 0644)
}

func (c *Config) ValidSession() bool {
	return c.SessionId != "" && !c.SessionExpired()
}

func (c *Config) SessionExpired() bool {
	return c.SessionId != "" && time.Now().Local().After(c.SessionExpiry)
}

func (c *Config) CheckVersion(versionString string, fetchLatestVersion func() (version.Version, bool), installedViaBrew bool, writer io.Writer) {
	defer c.Save()

	if time.Now().Local().Before(c.LastVersionCheck.Add(CheckVersionEvery)) {
		return
	}

	runningVersion, _ := version.New(versionString)

	if !c.VersionOutOfDate {
		latestVersion, ok := fetchLatestVersion()
		if !ok {
			return
		}
		c.LastVersionCheck = time.Now()

		if latestVersion.Equal(runningVersion) || latestVersion.Greater(runningVersion) {
			return
		}
		c.VersionOutOfDate = true
	}

	writer.Write([]byte(fmt.Sprintf("files-cli version %v is out of date\n", runningVersion)))
	if installedViaBrew {
		writer.Write([]byte(fmt.Sprintf("Upgrade via Homebrew\n\tbrew upgrade files-cli\n\n")))
		return
	}

	writer.Write([]byte(fmt.Sprintf("Download latest version from\nhttps://github.com/Files-com/files-cli/releases\n\n")))
}

func FetchLatestVersionNumber(parentCtx context.Context) func() (version.Version, bool) {
	return func() (version.Version, bool) {
		checkingFailed := func(err error) bool {
			if err != nil {
				files_sdk.GlobalConfig.Logger().Printf("Versioning checking failed: %v", err.Error())
				return true
			}
			return false
		}

		ctx, cancel := context.WithTimeout(parentCtx, 4*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, "GET", CLICurrentVersionURL, nil)
		if checkingFailed(err) {
			return version.Version{}, false
		}

		resp, err := files_sdk.GlobalConfig.GetHttpClient().Do(req)
		if checkingFailed(err) {
			return version.Version{}, false
		}
		data, err := ioutil.ReadAll(resp.Body)
		if checkingFailed(err) {
			return version.Version{}, false
		}
		latestVersion, err := version.New(string(data))
		if checkingFailed(err) {
			return version.Version{}, false
		}
		return latestVersion, true
	}
}

func InstalledViaBrew() bool {
	e, err := os.Executable()
	if err != nil {
		return false
	}
	symPath, err := os.Readlink(e)
	if err != nil {
		return false
	}

	if strings.Contains(symPath, "../Cellar/files-cli") {
		return true
	}

	return false
}

func (c *Config) configPath() (string, error) {
	if c.configPathOverride != "" {
		return c.configPathOverride, nil
	}
	root, err := configRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "files-cli"), nil
}

func configRoot() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".config"), nil
}

func initConfig() error {
	root, err := configRoot()
	if err != nil {
		return err
	}
	_, err = os.Stat(root)
	if os.IsNotExist(err) {
		os.MkdirAll(root, 0600)
	}

	path, err := (&Config{}).configPath()
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	f.Write([]byte("{}"))
	f.Close()
	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func SessionUnauthorizedError(paramsSessionCreate files_sdk.SessionCreateParams, err error, config Config) (files_sdk.SessionCreateParams, error) {
	responseError, ok := err.(files_sdk.ResponseError)
	if !ok {
		return paramsSessionCreate, err
	}
	switch eType := responseError.Type; eType {
	case "not-authorized/authentication-required":
		return paramsSessionCreate, err
	case "bad-request/invalid-username-or-password":
		return paramsSessionCreate, err
	case "not-authenticated/two-factor-authentication-error":
		fmt.Fprintf(config.Out, "%v\n", strings.Replace(responseError.ErrorMessage, "2FA Authentication error: ", "", 1))

		if contains(responseError.Data.TwoFactorAuthenticationMethod, "u2f") {
			return U2fResponse(paramsSessionCreate, responseError, config)
		}

		if contains(responseError.Data.TwoFactorAuthenticationMethod, "yubi") {
			return YubiResponse(paramsSessionCreate, responseError, config.Out)
		}

		if contains(responseError.Data.TwoFactorAuthenticationMethod, "totp") {
			return TotpResponse(paramsSessionCreate, config.Out)
		}

		if contains(responseError.Data.TwoFactorAuthenticationMethod, "sms") {
			return SmsResponse(paramsSessionCreate, config.Out)
		}

		return paramsSessionCreate, fmt.Errorf("%v is unsupported as login method", responseError.Data.TwoFactorAuthenticationMethod)
	}

	return paramsSessionCreate, err
}

func parseTermInput(text string) string {
	text = strings.Replace(text, "\r", "", -1)
	text = strings.Replace(text, "\n", "", -1) // Windows command prompt
	return text
}

func CreateSession(paramsSessionCreate files_sdk.SessionCreateParams, config *Config) error {
	var err error
	config.Subdomain, err = PromptUserWithPretext("Subdomain: %s", config.Subdomain, *config)
	if err != nil {
		return err
	}

	userNameDisplay := "Username: %s"
	if paramsSessionCreate.Username != "" {
		config.Username, err = PromptUserWithPretext(userNameDisplay, paramsSessionCreate.Username, *config)
	} else {
		config.Username, err = PromptUserWithPretext(userNameDisplay, config.Username, *config)
	}
	paramsSessionCreate.Username = config.Username

	if paramsSessionCreate.Password == "" {
		fmt.Fprintf(config.Out, "Password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}

		paramsSessionCreate.Password = parseTermInput(string(bytePassword))
		fmt.Fprintf(config.Out, "\n")
	}

	config.SetGlobal()
	files_sdk.GlobalConfig.SessionId = ""
	client := session.Client{Config: files_sdk.GlobalConfig}

	result, err := client.Create(context.TODO(), paramsSessionCreate)

	if err != nil {
		otpSessionCreate, err := SessionUnauthorizedError(paramsSessionCreate, err, *config)
		if err == nil {
			result, err = client.Create(context.TODO(), otpSessionCreate)
		}

		if err != nil {
			return err
		}
	}
	config.SessionId = result.Id
	config.SessionExpiry = time.Now().Local().Add(SessionExpiry)

	err = config.Save()
	if err != nil {
		return err
	}
	config.SetGlobal()
	return nil
}
