package lib

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/Files-com/files-cli/lib/version"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/session"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
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
	SessionId                 string    `json:"session_id"`
	SessionExpiry             time.Time `json:"session_expiry"`
	LastValidVersionCheck     time.Time `json:"last_valid_version_check"`
	Subdomain                 string    `json:"subdomain"`
	Username                  string    `json:"username"`
	APIKey                    string    `json:"api_key"`
	Endpoint                  string    `json:"endpoint,omitempty"`
	configPathOverride        string
	files_sdk.Environment     `json:"environment"`
	ConcurrentConnectionLimit int      `json:"concurrent_connection_limit"`
	ResourceFormat            []string `json:"resource_format"`
}

func (p Profile) SetResourceFormat(cmd *cobra.Command, defaultFormat []string) []string {
	if len(p.ResourceFormat) > 0 && !cmd.Flags().Changed("format") {
		return p.ResourceFormat
	}
	return defaultFormat
}

type ResetConfig struct {
	Subdomain                 bool
	Username                  bool
	APIKey                    bool
	Endpoint                  bool
	Session                   bool
	VersionCheck              bool
	ConcurrentConnectionLimit bool
	ResourceFormat            bool
}

var SessionExpiry = time.Hour * 6
var CheckVersionEvery = time.Hour * 48

const CLICurrentVersionURL = "https://api.github.com/repos/Files-com/files-cli/releases/latest"

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
	if reset.ConcurrentConnectionLimit {
		p.Current().ConcurrentConnectionLimit = 0
	}
	if reset.ResourceFormat {
		p.Current().ResourceFormat = nil
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
		p.Save()
	}

	if profile != "" {
		p.Profile = profile
	} else {
		p.Profile = "default"
	}
	p.Environment = config.Environment
	p.Config = config

	p.SetOnConfig()
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
	if p.Config.EndpointOverride == "" {
		p.Config.EndpointOverride = p.Current().Endpoint
	} else {
		p.Current().Endpoint = p.Config.EndpointOverride
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
	return os.WriteFile(path, file, 0600)
}

func (p *Profiles) ValidSession() bool {
	return p.Current().SessionId != "" && !p.SessionExpired()
}

func (p *Profiles) SessionExpired() bool {
	return p.Current().SessionId != "" && time.Now().Local().After(p.Current().SessionExpiry)
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
				config.Logger.Printf("Versioning checking failed: %v", err.Error())
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
		req.Close = true
		resp, err := config.Do(req)
		if checkingFailed(err) {
			return version.Version{}, false
		}
		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if checkingFailed(err) {
			return version.Version{}, false
		}
		releases := make(map[string]interface{})
		err = json.Unmarshal(data, &releases)
		if checkingFailed(err) {
			return version.Version{}, false
		}
		tagNameTemp, okKey := releases["tag_name"]
		tagName, okString := tagNameTemp.(string)
		if okKey && okString {
			latestVersion, err := version.New(strings.Replace(tagName, "v", "", 1))
			if checkingFailed(err) {
				return version.Version{}, false
			}

			config.Logger.Printf("Latest version: %v", latestVersion)
			return latestVersion, true
		} else {
			checkingFailed(clierr.Errorf(clierr.ErrorCodeFatal, "failed to parse tag_name from releases - %v", releases))
			return version.Version{}, false
		}
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

// initConfig initializes the configuration for a set of profiles.
// It ensures that the configuration directory and file exist.
// If there's a permission issue with the configuration file, it renames the file and creates a new one.
func (p *Profiles) initConfig() error {
	// Determine the configuration root directory
	root, err := p.configRoot()
	if err != nil {
		return err
	}

	// Create the root directory if it doesn't exist
	if _, err := os.Stat(root); os.IsNotExist(err) {
		if err := os.MkdirAll(root, 0755); err != nil {
			return err
		}
	}

	path, err := p.configPath()
	if err != nil {
		return err
	}

	stat, err := os.Stat(path)
	var missingReadPerms bool
	if err == nil {
		missingReadPerms = stat.Mode().Perm()&0400 == 0
	}

	// Handle configuration file permissions and existence
	var createFile bool
	if os.IsPermission(err) || missingReadPerms {
		backupPath := fmt.Sprintf("%v-backup", path)
		if err := os.Rename(path, backupPath); err != nil {
			panic(err)
		}
		fmt.Fprintf(os.Stderr, "Warning: The config file '%v' has a permissions issue and cannot be read. To avoid data loss, it has been renamed to '%v'. A new config file will be created.\n", path, backupPath)
		createFile = true
	} else if os.IsNotExist(err) {
		createFile = true
	}

	// Create a new configuration file if it doesn't exist or was renamed
	if createFile {
		if err := os.WriteFile(path, []byte("{}"), 0600); err != nil {
			return err
		}
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
	var responseError files_sdk.ResponseError
	ok := errors.As(err, &responseError)
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

		return paramsSessionCreate, clierr.Errorf(clierr.ErrorCodeUsage, "%v is unsupported as login method", responseError.Data.TwoFactorAuthenticationMethod)
	}

	return paramsSessionCreate, err
}

func parseTermInput(text string) string {
	text = strings.ReplaceAll(text, "\r", "")
	text = strings.ReplaceAll(text, "\n", "") // Windows command prompt
	return text
}

func Reauthenicate(profile *Profiles) (err error) {
	if profile.Config.AdditionalHeaders == nil {
		profile.Config.AdditionalHeaders = make(map[string]string)
	}
	fmt.Fprintf(profile.Out, "Password or 2FA Code: ")
	profile.Config.AdditionalHeaders["X-Files-Reauthentication"], err = passwordPrompt()
	fmt.Fprintf(profile.Out, "\n")
	return
}

func passwordPrompt() (password string, err error) {
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return password, err
	}

	password = parseTermInput(string(bytePassword))

	return password, nil
}

func ValidateDomain(domain string) bool {
	_, err := net.LookupHost(domain)
	if err != nil {
		return false
	}

	return true
}

func CreateSession(ctx context.Context, paramsSessionCreate files_sdk.SessionCreateParams, profile *Profiles) error {
	var err error
	profile.Current().Subdomain, err = PromptUserWithPretext(ctx, "Site: %s", lib.DefaultString(profile.Current().Subdomain, profile.Current().Endpoint), "Enter your site subdomain (e.g. mysite) or custom domain (e.g. myfiles.customdomain.com)", profile)
	if err != nil {
		return err
	}

	_, err = url.Parse(files_sdk.Config{Subdomain: profile.Current().Subdomain}.Endpoint())
	if err != nil {
		customDomain := profile.Current().Subdomain
		_, err = url.Parse(files_sdk.Config{EndpointOverride: customDomain}.Endpoint())
		if err != nil {
			return clierr.Errorf(clierr.ErrorCodeFatal, "invalid domain or subdomain: %v", profile.Current().Subdomain)
		}
		profile.Current().Subdomain = ""
		profile.Current().Endpoint = customDomain
	}

	if strings.Contains(profile.Current().Subdomain, ".") {
		profile.Current().Endpoint = profile.Current().Subdomain
		profile.Current().Subdomain = ""
	}

	userNameDisplay := "Username: %s"
	if paramsSessionCreate.Username != "" {
		profile.Current().Username, err = PromptUserWithPretext(ctx, userNameDisplay, paramsSessionCreate.Username, "", profile)
		if err != nil {
			return err
		}
	} else {
		profile.Current().Username, err = PromptUserWithPretext(ctx, userNameDisplay, profile.Current().Username, "", profile)
		if err != nil {
			return err
		}
	}
	paramsSessionCreate.Username = profile.Current().Username

	if paramsSessionCreate.Password == "" {
		fmt.Fprintf(profile.Out, "Password: ")
		paramsSessionCreate.Password, err = passwordPrompt()
		fmt.Fprintf(profile.Out, "\n")
		if err != nil {
			return err
		}
	}

	profile.SetOnConfig()
	result, err := createSession(*profile.Config, paramsSessionCreate)
	if err != nil {
		var x509Err x509.HostnameError
		if errors.As(err, &x509Err) {
			customDomain := profile.Current().Subdomain
			profile.Current().Subdomain = ""
			profile.SessionId = ""
			profile.Current().Endpoint = customDomain
			profile.SetOnConfig()
			result, err = createSession(*profile.Config, paramsSessionCreate)
		}
		var respErr files_sdk.ResponseError
		if errors.As(err, &respErr) {
			if respErr.Type == "not-authenticated/lockout-region-mismatch" {
				profile.Current().Endpoint = respErr.Data.Host
				profile.SessionId = ""
				profile.SetOnConfig()
				result, err = createSession(*profile.Config, paramsSessionCreate)
			}
		}
	}
	if err != nil {
		otpSessionCreate, err := SessionUnauthorizedError(paramsSessionCreate, err, profile)
		profile.Config.APIKey = ""
		if err == nil {
			result, err = createSession(*profile.Config, otpSessionCreate)
		}

		if err != nil {
			return err
		}
	}
	profile.Current().SessionId = result.Id
	profile.Current().SessionExpiry = time.Now().Local().Add(SessionExpiry)

	err = profile.Save()
	if err != nil {
		return err
	}
	profile.SetOnConfig()
	return nil
}

func createSession(config files_sdk.Config, params files_sdk.SessionCreateParams) (files_sdk.Session, error) {
	return (&session.Client{Config: config}).Create(params, files_sdk.RequestOption(func(r *http.Request) error {
		r.Header.Del("X-FilesAPI-Auth")
		r.Header.Del("X-FilesAPI-Key")
		return nil
	}))
}
