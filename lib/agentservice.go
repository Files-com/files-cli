package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

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
	logFilePathFlag = "log-file-path"
	logVerboseFlag  = "log-verbose"
	logUTCTimeFlag  = "log-utc-time"

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
	ipWhitelist                        []string
	context.Context
}

func (a *AgentService) AddFlags(flags *pflag.FlagSet) {
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
	flags.StringSliceVar(&a.ipWhitelist, "append-to-ip-whitelist", []string{"127.0.0.1"}, "Add additional IPs to whitelist")
}

func (a *AgentService) Init(ctx context.Context) error {
	a.Context = ctx
	var d bool
	d = true
	a.Config.Debug = &d
	a.Config.SetLogger(logger.GetLogger())
	a.permissions = make(map[string][]string)
	a.shutdown = make(chan bool)
	if util.IsFileInputValid(a.PortableLogFile) && !filepath.IsAbs(a.PortableLogFile) {
		var err error
		a.PortableLogFile, err = filepath.Abs(a.PortableLogFile)
		if err != nil {
			return nil
		}
	}
	if a.PortableLogFile == "" {
		exePath, err := os.Executable()
		if err != nil {
			return err
		}
		exeDir, _ := filepath.Split(exePath)
		a.PortableLogFile = filepath.Join(exeDir, "files-agent.log")
	}
	a.LogFilePath = a.PortableLogFile
	a.Service.PortableMode = 1
	return nil
}

func (a *AgentService) LoadConfig(ctx context.Context) error {
	a.Context = ctx
	var d bool
	d = true
	a.Config.Debug = &d

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
	whiteListFile, err := a.createServerTempWhiteList()
	if err != nil {
		return err
	}
	a.ConfigDir, a.ConfigFile, err = a.createServerTempConfig(whiteListFile)
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
	go func() {
		<-a.shutdown
		a.updateCloudConfig(a.Context, "shutdown", "shutdown chan received")
	}()
	err = a.Service.StartPortableMode(int(a.Port), 0, 0, []string{}, false,
		false, "", "", "", "")
	if err == nil {
		go func() {
			a.Config.SetLogger(logger.GetLogger())
			ctx, cancel := context.WithTimeout(a.Context, time.Second*30)
			defer cancel()
			for {
				if ctx.Err() != nil || a.Service.Error != nil {
					break
				}

				innerErr := a.updateCloudConfig(a.Context, "running", "after start loop")
				if innerErr == nil {
					break
				}
				time.Sleep(time.Second * 5)
			}
		}()
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
			exts := []string{}
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
	a.RemoteServerConfigurationFile.Id = int64(mapForId["id"].(float64))
	if a.ConfigVersion != "1" {
		return fmt.Errorf("agent upgrade required - `Your current version of the files-agent incompatible and requires an update.`")
	}
	return err
}

func (a *AgentService) loadPublicIpAddress(ctx context.Context) (err error) {
	client := ip_address.Client{Config: *a.Config}
	iter, err := client.GetReserved(ctx, files_sdk.IpAddressGetReservedParams{})
	if err != nil {
		return
	}
	for iter.Next() {
		a.ipWhitelist = append(a.ipWhitelist, iter.PublicIpAddress().IpAddress)
	}

	return
}

func (a *AgentService) createServerTempWhiteList() (file *os.File, err error) {
	safeListTmp := make(map[string]interface{})
	safeListTmp["addresses"] = a.ipWhitelist
	file, err = os.CreateTemp("", "whitelist-*.json")
	if err != nil {
		return
	}
	safeListBytes, err := json.Marshal(safeListTmp)
	logger.Debug("files-cli", "", "Ip whitelist: %v", safeListTmp)
	if err != nil {
		return
	}
	_, err = file.Write(safeListBytes)
	if err != nil {
		return
	}
	err = file.Close()
	return
}

func (a *AgentService) createServerTempConfig(whiteList *os.File) (string, string, error) {
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
	dir, file := filepath.Split(configTempFile.Name())
	logger.Debug("files-cli", "", "Config Path: %v", configTempFile.Name())
	return dir, file, config.LoadConfig(dir, file)
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
		return err
	}

	logger.Debug("files-cli", "", "Cloud Configuration Update Response: %v", newConfig)

	return nil
}