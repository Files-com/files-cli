package lib

import (
	"context"
	"encoding/json"
	"net/http"
	"os/user"
	"path/filepath"
	"reflect"
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

type Profiles struct {
	Profiles              map[string]*Profile `json:"profiles"`
	*files_sdk.Config     `json:"-"`
	Overrides             `json:"-"`
	Profile               string `json:"-"`
	files_sdk.Environment `json:"-"`
	ConfigDir             string `json:"-"`
}

type Profile struct {
	SessionId             string    `json:"session_id"`
	SessionExpiry         time.Time `json:"session_expiry"`
	LastValidVersionCheck time.Time `json:"last_valid_version_check"`
	Subdomain             string    `json:"subdomain"`
	Username              string    `json:"username"`
	APIKey                string    `json:"api_key"`
	Endpoint              string    `json:"endpoint,omitempty"`
	configPathOverride    string
	files_sdk.Environment `json:"environment"`
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

func (p *Profiles) Current() *Profile {
	env, ok := p.Profiles[p.Profile]
	if !ok {
		p.Profiles[p.Profile] = &Profile{Environment: p.Environment, APIKey: p.APIKey}
		return p.Current()
	}
	return env
}

func (p *Profiles) ResetWith(reset ResetConfig) error {
	if reset.Subdomain {
		p.Current().Subdomain = ""
	}
	if reset.Username {
		p.Current().Username = ""
	}
	if reset.APIKey {
		p.Current().APIKey = ""
	}
	if reset.Endpoint {
		p.Current().Endpoint = ""
	}
	if reset.Session {
		p.Current().SessionId = ""
	}
	if reset.VersionCheck {
		p.Current().LastValidVersionCheck = time.Now()
	}
	return p.Save()
}

func (p *Profiles) Reset() error {
	return p.initConfig()
}

func (p *Profiles) Init() *Profiles {
	p.Profiles = make(map[string]*Profile)
	return p
}

func (p *Profiles) Load(config *files_sdk.Config, profile string) error {
	p.Init()
	p.initConfig()

	data, err := p.read()
	if err != nil {
		return err
	}

	v1Profile := Profile{}

	err = json.Unmarshal(data, &v1Profile)
	if err != nil {
		return err
	}

	if reflect.ValueOf(v1Profile).IsZero() {
		err = json.Unmarshal(data, &p)
		if err != nil {
			return err
		}
	} else {
		p.Profiles["default"] = &v1Profile
	}

	if profile != "" {
		p.Profile = profile
	} else {
		p.Profile = "default"
	}
	p.Environment = config.Environment
	p.Config = config

	p.SetOnConfig()
	p.Save()
	return nil
}

func (p *Profiles) read() (b []byte, err error) {
	var path string
	path, err = p.configPath()
	if err != nil {
		return b, err
	}
	return os.ReadFile(path)
}

func (p *Profiles) SetOnConfig() {
	p.Config.SessionId = p.Current().SessionId
	if p.Config.Endpoint == "" {
		p.Config.Endpoint = p.Current().Endpoint
	} else {
		p.Current().Endpoint = p.Config.Endpoint
	}
	if p.Config.Subdomain == "" {
		p.Config.Subdomain = p.Current().Subdomain
	} else {
		p.Current().Subdomain = p.Config.Subdomain
	}
	if p.Config.APIKey == "" {
		p.Config.APIKey = p.Current().APIKey
	} else {
		p.Current().APIKey = p.Config.APIKey
	}
	p.Config.Environment = p.Current().Environment
}

func (p *Profiles) Save() error {
	file, _ := json.MarshalIndent(p, "", " ")
	path, err := p.configPath()
	if err != nil {
		return err
	}
	return os.WriteFile(path, file, 0644)
}

func (p *Profiles) ValidSession() bool {
	return p.SessionId != "" && !p.SessionExpired()
}

func (p *Profiles) SessionExpired() bool {
	return p.SessionId != "" && time.Now().Local().After(p.Current().SessionExpiry)
}

func (p *Profiles) CheckVersion(versionString string, fetchLatestVersion func() (version.Version, bool), installedViaBrew bool, writer io.Writer) {
	defer p.Save()

	if time.Now().Local().Before(p.Current().LastValidVersionCheck.Add(CheckVersionEvery)) {
		return
	}

	runningVersion, _ := version.New(versionString)

	latestVersion, ok := fetchLatestVersion()
	if !ok {
		return
	}

	if latestVersion.Equal(runningVersion) || latestVersion.Greater(runningVersion) {
		p.Current().LastValidVersionCheck = time.Now()
		return
	}

	writer.Write([]byte(fmt.Sprintf("files-cli version %v is out of date. Latest version is %v\n", runningVersion, latestVersion)))
	if installedViaBrew {
		writer.Write([]byte(fmt.Sprintf("Upgrade via Homebrew\n\tbrew upgrade files-cli\n\n")))
		return
	}

	writer.Write([]byte(fmt.Sprintf("Download latest version from\nhttps://github.com/Files-com/files-cli/releases\n\n")))
}

func FetchLatestVersionNumber(config files_sdk.Config, parentCtx context.Context) func() (version.Version, bool) {
	return func() (version.Version, bool) {
		checkingFailed := func(err error) bool {
			if err != nil {
				config.Logger().Printf("Versioning checking failed: %v", err.Error())
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

		resp, err := config.GetHttpClient().Do(req)
		if checkingFailed(err) {
			return version.Version{}, false
		}
		data, err := io.ReadAll(resp.Body)
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

func (p *Profiles) configPath() (string, error) {
	root, err := p.configRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "files-cli"), nil
}

func (p *Profiles) configRoot() (string, error) {
	if p.ConfigDir != "" {
		return p.ConfigDir, nil
	}
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".config"), nil
}

func (p *Profiles) initConfig() error {
	root, err := p.configRoot()
	if err != nil {
		return err
	}
	_, err = os.Stat(root)
	if os.IsNotExist(err) {
		os.MkdirAll(root, 0600)
	}

	path, err := p.configPath()
	if err != nil {
		return err
	}

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		f.Write([]byte("{}"))
		f.Close()
	}
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

func SessionUnauthorizedError(paramsSessionCreate files_sdk.SessionCreateParams, err error, config *Profiles) (files_sdk.SessionCreateParams, error) {
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

func CreateSession(paramsSessionCreate files_sdk.SessionCreateParams, profile *Profiles) error {
	var err error
	profile.Current().Subdomain, err = PromptUserWithPretext("Subdomain: %s", profile.Current().Subdomain, profile)
	if err != nil {
		return err
	}

	userNameDisplay := "Username: %s"
	if paramsSessionCreate.Username != "" {
		profile.Current().Username, err = PromptUserWithPretext(userNameDisplay, paramsSessionCreate.Username, profile)
	} else {
		profile.Current().Username, err = PromptUserWithPretext(userNameDisplay, profile.Current().Username, profile)
	}
	paramsSessionCreate.Username = profile.Current().Username

	if paramsSessionCreate.Password == "" {
		fmt.Fprintf(profile.Out, "Password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}

		paramsSessionCreate.Password = parseTermInput(string(bytePassword))
		fmt.Fprintf(profile.Out, "\n")
	}

	profile.SetOnConfig()
	profile.SessionId = ""
	client := session.Client{Config: *profile.Config}

	result, err := client.Create(context.TODO(), paramsSessionCreate)

	if err != nil {
		otpSessionCreate, err := SessionUnauthorizedError(paramsSessionCreate, err, profile)
		if err == nil {
			result, err = client.Create(context.TODO(), otpSessionCreate)
		}

		if err != nil {
			return err
		}
	}
	profile.SessionId = result.Id
	profile.Current().SessionExpiry = time.Now().Local().Add(SessionExpiry)

	err = profile.Save()
	if err != nil {
		return err
	}
	profile.SetOnConfig()
	return nil
}
