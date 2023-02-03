package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/drakkan/sftpgo/v2/common"

	"github.com/drakkan/sftpgo/v2/util"
	"github.com/rs/zerolog"

	"github.com/drakkan/sftpgo/v2/logger"

	"github.com/drakkan/sftpgo/v2/config"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	ip_address "github.com/Files-com/files-sdk-go/v2/ipaddress"
	remote_server "github.com/Files-com/files-sdk-go/v2/remoteserver"
	"github.com/drakkan/sftpgo/v2/dataprovider"
	"github.com/drakkan/sftpgo/v2/kms"
	"github.com/drakkan/sftpgo/v2/service"
	"github.com/drakkan/sftpgo/v2/vfs"
	"github.com/sftpgo/sdk"
	"github.com/spf13/pflag"
)

const (
	logFilePathFlag         = "log-file-path"
	logVerboseFlag          = "log-verbose"
	logUTCTimeFlag          = "log-utc-time"
	appendToIpWhitelistFlag = "append-to-ip-whitelist"

	defaultLogMaxSize   = 10
	defaultLogMaxBackup = 5
	defaultLogMaxAge    = 28
	defaultLogCompress  = false
)

type AgentService struct {
	service.Service
	*files_sdk.Config
	ConfigPath string
	files_sdk.RemoteServerConfigurationFile
	portableFsProvider                 string
	permissions                        map[string][]string
	portableSFTPFingerprints           []string
	portablePublicKeys                 []string
	shutdown                           chan bool
	PortableSFTPEndpoint               string
	PortableSFTPPrefix                 string
	PortableSFTPDisableConcurrentReads bool
	PortableSFTPDBufferSize            int64
	PortableAllowedPatterns            []string
	PortableDeniedPatterns             []string
	PortableLogFile                    string
	PortableLogVerbose                 bool
	PortableLogUTCTime                 bool
	PortableSFTPFingerprints           []string
	ipWhitelist                        map[string]bool
	appendedIpWhitelist                []string
	sftpGoConfigPath                   string
	subdomain                          string
	context.Context
}

func (a *AgentService) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&a.subdomain, "subdomain", "", "")
	flags.StringVar(&a.ConfigPath, "config", "./files-agent-config.json", "")
	flags.StringVarP(&a.PortableLogFile, logFilePathFlag, "l", "files-agent.log", "Leave empty to disable logging")
	flags.BoolVarP(&a.PortableLogVerbose, logVerboseFlag, "v", false, "Enable verbose logs")
	flags.BoolVar(&a.PortableLogUTCTime, logUTCTimeFlag, false, "Use UTC time for logging")
	flags.StringArrayVar(&a.PortableAllowedPatterns, "allowed-patterns", []string{},
		`Allowed file patterns case insensitive.
The format is:
/dir::pattern1,pattern2.
For example: "/somedir::*.jpg,a*b?.png"`)
	flags.StringArrayVar(&a.PortableDeniedPatterns, "denied-patterns", []string{},
		`Denied file patterns case insensitive.
The format is:
/dir::pattern1,pattern2.
For example: "/somedir::*.jpg,a*b?.png"`)
	flags.StringVar(&a.PortableSFTPEndpoint, "sftp-endpoint", "", `SFTP endpoint as host:port for SFTP
provider`)
	flags.StringSliceVar(&a.PortableSFTPFingerprints, "sftp-fingerprints", []string{}, `SFTP fingerprints to verify remote host
key for SFTP provider`)
	flags.StringVar(&a.PortableSFTPPrefix, "sftp-prefix", "", `SFTP prefix allows restrict all
operations to a given path within the
remote SFTP server`)
	flags.Int64Var(&a.PortableSFTPDBufferSize, "sftp-buffer-size", 0, `The size of the buffer (in MB) to use
for transfers. By enabling buffering,
the reads and writes, from/to the
remote SFTP server, are split in
multiple concurrent requests and this
allows data to be transferred at a
faster rate, over high latency networks,
by overlapping round-trip times`)
	flags.StringSliceVar(&a.appendedIpWhitelist, "append-to-ip-whitelist", []string{"127.0.0.1"}, "Add additional IPs to whitelist")
}

func (a *AgentService) Init(ctx context.Context, requirePaths bool) error {
	a.Context = ctx
	a.Config.Debug = true
	a.Config.SetLogger(logger.GetLogger())
	a.Config.Subdomain = a.subdomain
	profile := &Profiles{}
	err := profile.Load(&files_sdk.Config{}, "default")
	if err != nil {
		return err
	}
	if a.Config.Subdomain == "" && profile.Subdomain != "" {
		a.Config.Subdomain = profile.Subdomain
	}
	a.permissions = make(map[string][]string)
	a.shutdown = make(chan bool)
	a.ipWhitelist = make(map[string]bool)
	a.Service.PortableMode = 1

	err = a.InitPaths()
	if requirePaths {
		return err
	}
	return nil
}

func (a *AgentService) InitPaths() error {
	if util.IsFileInputValid(a.PortableLogFile) && !filepath.IsAbs(a.PortableLogFile) {
		var err error
		a.PortableLogFile, err = filepath.Abs(a.PortableLogFile)
		if err != nil {
			return err
		}
	}
	a.LogFilePath = a.PortableLogFile
	absPath, err := filepath.Abs(a.ConfigPath)
	if err != nil {
		return err
	}
	a.ConfigPath = absPath
	_, err = os.Stat(a.ConfigPath)
	if a.PortableLogFile == "" {
		dir, _ := filepath.Split(a.ConfigPath)
		a.PortableLogFile = filepath.Join(dir, "files-agent.log")
	}
	return err
}

func (a *AgentService) LoadConfig(ctx context.Context) error {
	a.Context = ctx
	a.Config.Debug = true

	logLevel := zerolog.DebugLevel
	if !a.LogVerbose {
		logLevel = zerolog.DebugLevel
	}

	logger.InitLogger(a.LogFilePath, a.LogMaxSize, a.LogMaxBackups, a.LogMaxAge, a.LogCompress, a.LogUTCTime, logLevel)
	logger.EnableConsoleLogger(logLevel)
	if a.LogFilePath == "" {
		return fmt.Errorf("log path is empty")
	}
	a.Config.SetLogger(logger.GetLogger())
	err := a.loadConfig()
	if err != nil {
		return err
	}
	err = a.updateCloudConfig(a.Context, "shutdown", "API check")
	if err != nil {
		return err
	}
	a.portableFsProvider = "osfs"
	fsProvider := sdk.GetProviderByName(a.portableFsProvider)
	if !filepath.IsAbs(a.Root) {
		if fsProvider == sdk.LocalFilesystemProvider {
			a.portableFsProvider, _ = filepath.Abs(a.Root)
		} else {
			a.portableFsProvider = os.TempDir()
		}
	}
	p, err := a.mapPermissions()
	if err != nil {
		return err
	}
	a.permissions["/"] = p

	portableSFTPPrivateKey := a.PrivateKey
	if portableSFTPPrivateKey == "" {
		return fmt.Errorf("missing private key")
	}

	a.portableSFTPFingerprints = append(a.portableSFTPFingerprints, a.PrivateKey)
	a.RemoteServerConfigurationFile.ServerHostKey = a.portableSFTPFingerprints[0]
	a.portablePublicKeys = append(a.portablePublicKeys, a.PublicKey)

	err = a.loadPublicIpAddress(ctx)
	if err != nil {
		return err
	}
	a.ConfigDir, a.ConfigFile, err = a.createServerTempConfig()
	if err != nil {
		return err
	}

	a.LogMaxSize = defaultLogMaxSize
	a.LogMaxBackups = defaultLogMaxBackup
	a.LogMaxAge = defaultLogMaxAge
	a.LogCompress = defaultLogCompress
	a.LogVerbose = a.PortableLogVerbose
	a.LogUTCTime = a.PortableLogUTCTime
	a.Service.PortableMode = 1
	a.PortableUser = dataprovider.User{
		BaseUser: sdk.BaseUser{
			Username:    "admin",
			PublicKeys:  a.portablePublicKeys,
			Permissions: a.permissions,
			HomeDir:     a.Root,
			Status:      1,
		},
		Filters: dataprovider.UserFilters{
			BaseUserFilters: sdk.BaseUserFilters{
				FilePatterns:   a.parsePatternsFilesFilters(),
				StartDirectory: "",
			},
		},
		FsConfig: vfs.Filesystem{
			Provider: sdk.GetProviderByName(a.portableFsProvider),
			SFTPConfig: vfs.SFTPFsConfig{
				BaseSFTPFsConfig: sdk.BaseSFTPFsConfig{
					Endpoint:                a.PortableSFTPEndpoint,
					Fingerprints:            a.portableSFTPFingerprints,
					Prefix:                  a.PortableSFTPPrefix,
					DisableCouncurrentReads: a.PortableSFTPDisableConcurrentReads,
					BufferSize:              a.PortableSFTPDBufferSize,
				},
				PrivateKey: kms.NewPlainSecret(portableSFTPPrivateKey),
			},
		},
	}

	return nil
}

func (a *AgentService) Start(_ bool) error {
	err := a.LoadConfig(a.Context)
	if err != nil {
		return err
	}
	logger.Debug("files-cli", "", "AgentService.Start")
	a.Config.SetLogger(logger.GetLogger())
	err = a.Service.StartPortableMode(int(a.Port), -1, -1, []string{}, false,
		false, "", "", "", "")
	if err == nil {
		go a.afterStart()
		logger.Debug("files-cli", "", "AgentService.Start finished")
	} else {
		logger.Debug("files-cli", "", "AgentService.Start err: %v", err)
		buf := make([]byte, 1<<16)
		runtime.Stack(buf, true)
		logger.Debug("files-cli", "", "%s", buf)
		return err
	}

	return a.Service.Error
}

func (a *AgentService) Stop() {
	a.updateCloudConfig(a.Context, "shutdown", "stopping")
	a.Service.Stop()
}

func (a *AgentService) Wait() {
	logger.Debug("files-cli", "", "AgentService.Wait")
	a.Service.Wait()
	a.updateCloudConfig(a.Context, "shutdown", "done waiting")
	os.Exit(1)
}

func (a *AgentService) ServiceArgs() []string {
	args := []string{
		"agent", "start",
		"--config", a.ConfigPath,
		fmt.Sprintf("--%v", logFilePathFlag), a.LogFilePath,
		fmt.Sprintf("--%v", logVerboseFlag), fmt.Sprintf("%v", a.PortableLogVerbose),
		fmt.Sprintf("--%v", logUTCTimeFlag), fmt.Sprintf("%v", a.PortableLogUTCTime),
		fmt.Sprintf("--%v", appendToIpWhitelistFlag), strings.Join(a.appendedIpWhitelist, ","),
	}

	if len(a.PortableAllowedPatterns) > 0 {
		args = append(args, "--allowed-patterns", strings.Join(a.PortableAllowedPatterns, ","))
	}
	return args
}

func (a *AgentService) parsePatternsFilesFilters() []sdk.PatternsFilter {
	var patterns []sdk.PatternsFilter
	for _, val := range a.PortableAllowedPatterns {
		p, exts := getPatternsFilterValues(strings.TrimSpace(val))
		if p != "" {
			patterns = append(patterns, sdk.PatternsFilter{
				Path:            path.Clean(p),
				AllowedPatterns: exts,
				DeniedPatterns:  []string{},
			})
		}
	}
	for _, val := range a.PortableDeniedPatterns {
		p, exts := getPatternsFilterValues(strings.TrimSpace(val))
		if p != "" {
			found := false
			for index, e := range patterns {
				if path.Clean(e.Path) == path.Clean(p) {
					patterns[index].DeniedPatterns = append(patterns[index].DeniedPatterns, exts...)
					found = true
					break
				}
			}
			if !found {
				patterns = append(patterns, sdk.PatternsFilter{
					Path:            path.Clean(p),
					AllowedPatterns: []string{},
					DeniedPatterns:  exts,
				})
			}
		}
	}
	return patterns
}

func getPatternsFilterValues(value string) (string, []string) {
	if strings.Contains(value, "::") {
		dirExts := strings.Split(value, "::")
		if len(dirExts) > 1 {
			dir := strings.TrimSpace(dirExts[0])
			var exts []string
			for _, e := range strings.Split(dirExts[1], ",") {
				cleanedExt := strings.TrimSpace(e)
				if cleanedExt != "" {
					exts = append(exts, cleanedExt)
				}
			}
			if dir != "" && len(exts) > 0 {
				return dir, exts
			}
		}
	}
	return "", nil
}

func (a *AgentService) mapPermissions() ([]string, error) {
	switch a.PermissionSet {
	case "read_only":
		return []string{"list", "download"}, nil
	case "write_only":
		return []string{"list", "upload", "overwrite", "delete", "rename", "create_dirs"}, nil
	case "read_write":
		return []string{"list", "download", "upload", "overwrite", "delete", "rename", "create_dirs"}, nil
	default:
		return []string{}, fmt.Errorf("invalid or missing permissions: %v", a.PermissionSet)
	}
}

func (a *AgentService) afterStart() {
	a.pingServer()
	for range time.Tick(24 * time.Hour) {
		err := a.whitelistRefresh()
		if err != nil {
			logger.Info("files-cli", "", "Unable to refresh files.com whitelist - %v", err.Error())
		}
	}
}

func (a *AgentService) pingServer() {
	logger.Debug("files-cli", "", "Contacting Files.com to attempt an agent connection")
	a.Config.SetLogger(logger.GetLogger())
	ctx, cancel := context.WithTimeout(a.Context, time.Minute*6)
	defer cancel()
	attemptCount := 0
	for {
		if ctx.Err() != nil || a.Service.Error != nil {
			break
		}
		requestCtx, requestCancel := context.WithTimeout(ctx, time.Minute)
		err := a.updateCloudConfig(requestCtx, "running", "after start loop")
		requestCancel()

		if err == nil {
			break
		}
		attemptCount += 1
		time.Sleep(time.Second * time.Duration(5*attemptCount))
	}
}

func (a *AgentService) whitelistRefresh() error {
	err := a.loadPublicIpAddress(a.Context)
	if err != nil {
		return err
	}
	_, err = a.whitelistFile()
	if err != nil {
		return err
	}
	return a.reloadConfig()
}

func (a *AgentService) loadConfig() error {
	absPath, err := filepath.Abs(a.ConfigPath)
	if err != nil {
		return err
	}
	a.ConfigPath = absPath
	configBytes, err := os.ReadFile(a.ConfigPath)
	if err != nil {
		return err
	}
	mapForId := make(map[string]interface{})
	err = json.Unmarshal(configBytes, &mapForId)
	if err != nil {
		return err
	}
	err = json.Unmarshal(configBytes, &a.RemoteServerConfigurationFile)
	if err != nil {
		return err
	}
	if a.RemoteServerConfigurationFile.Port == 0 {
		a.RemoteServerConfigurationFile.Port = 58550
	}
	if a.Root == "" {
		a.Root, err = os.UserHomeDir()
		if err != nil {
			return err
		}
	}
	a.RemoteServerConfigurationFile.Id = int64(mapForId["id"].(float64))
	if a.ConfigVersion != "1" {
		return fmt.Errorf("agent upgrade required - `Your current version of the files-agent incompatible and requires an update.`")
	}
	return err
}

func (a *AgentService) loadPublicIpAddress(ctx context.Context) (err error) {
	for _, ip := range a.appendedIpWhitelist {
		a.ipWhitelist[ip] = true
	}

	client := ip_address.Client{Config: *a.Config}
	iter, err := client.GetReserved(ctx, files_sdk.IpAddressGetReservedParams{})
	if err != nil {
		return
	}
	for iter.Next() {
		if iter.Err() != nil {
			return iter.Err()
		}
		a.ipWhitelist[iter.PublicIpAddress().IpAddress] = true
	}

	if iter.Err() != nil {
		return iter.Err()
	}

	if a.Config.Subdomain != "" {
		a.Config.RootPath()
		url, err := url.Parse(a.Config.Endpoint)
		if err != nil {
			logger.Debug("files-cli", "", "Unable to parse url %v", a.Config.Endpoint)
			return err
		}
		ips, err := net.LookupIP(url.Hostname())
		if err != nil {
			return err
		}

		for _, ip := range ips {
			a.ipWhitelist[ip.String()] = true
		}
	}

	return
}

func (a *AgentService) whitelistFile() (file *os.File, err error) {
	if common.Config.WhiteListFile == "" {
		file, err = os.CreateTemp("", "whitelist-*.json")
		common.Config.WhiteListFile = file.Name()
	} else {
		file, err = os.OpenFile(common.Config.WhiteListFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	}
	if err != nil {
		return
	}

	err = a.saveWhitelist(file)
	return
}

func (a *AgentService) saveWhitelist(file *os.File) error {
	safeListTmp := make(map[string][]string)
	var addresses []string
	for k := range a.ipWhitelist {
		addresses = append(addresses, k)
	}
	safeListTmp["addresses"] = addresses
	safeListBytes, err := json.Marshal(safeListTmp)
	if err != nil {
		return err
	}
	logger.Debug("files-cli", "", "Ip whitelist: %v", safeListTmp)
	_, err = file.Write(safeListBytes)
	defer file.Close()
	return err
}

func (a *AgentService) createServerTempConfig() (string, string, error) {
	whiteList, err := a.whitelistFile()
	if err != nil {
		return "", "", err
	}
	tmpConfig := make(map[string]interface{})
	tmpConfig["common"] = map[string]string{"whitelist_file": whiteList.Name()}
	configTempFile, err := os.CreateTemp("", "config-*.json")

	if err != nil {
		return "", "", err
	}
	tmpConfigBytes, err := json.Marshal(tmpConfig)
	if err != nil {
		return "", "", err
	}
	_, err = configTempFile.Write(tmpConfigBytes)
	if err != nil {
		return "", "", err
	}
	err = configTempFile.Close()
	if err != nil {
		return "", "", err
	}
	a.sftpGoConfigPath = configTempFile.Name()
	logger.Debug("files-cli", "", "Config Path: %v", a.sftpGoConfigPath)
	dir, file := filepath.Split(a.sftpGoConfigPath)
	return dir, file, a.reloadConfig()
}

func (a *AgentService) reloadConfig() error {
	dir, file := filepath.Split(a.sftpGoConfigPath)
	err := config.LoadConfig(dir, file)
	if err != nil {
		return err
	}

	fmt.Println(common.Config)

	return common.Reload()
}

func (a *AgentService) updateCloudConfig(ctx context.Context, status string, source string) error {
	client := remote_server.Client{Config: *a.Config}
	params := files_sdk.RemoteServerConfigurationFileParams{}

	params.Status = status
	params.Id = a.Id
	params.Port = a.Port
	params.Root = a.Root
	params.PermissionSet = a.PermissionSet
	params.ApiToken = a.ApiToken
	params.Hostname = a.Hostname
	params.ConfigVersion = a.ConfigVersion
	params.PrivateKey = a.PrivateKey
	params.PublicKey = a.PublicKey
	logger.Debug("files-cli", "", "Update Cloud Configuration - source (%v) : %v", source, params)
	newConfig, err := client.ConfigurationFile(ctx, params)
	if err != nil {
		logger.Debug("files-cli", "", "Cloud Configuration Update Error - source (%v): %v", source, err)
		remoteAgentModel, findErr := client.Find(ctx, files_sdk.RemoteServerFindParams{Id: params.Id})
		if findErr == nil {
			logger.Debug("files-cli", "", "Server set attributes - hostname: %v, port: %v, disabled: %v, root: %v",
				remoteAgentModel.Hostname, remoteAgentModel.Port, remoteAgentModel.Disabled, remoteAgentModel.FilesAgentRoot,
			)
		}
		return err
	}

	logger.Debug("files-cli", "", "Cloud Configuration Update Response - hostname: %v, port: %v, status: %v, root: %v",
		newConfig.Hostname, newConfig.Port, newConfig.Status, newConfig.Root)

	return nil
}
