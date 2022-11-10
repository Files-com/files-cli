package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

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
	files_sdk.Config
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
	flags.StringVarP(&a.PortableLogFile, logFilePathFlag, "l", "", "Leave empty to disable logging")
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
	err := a.loadConfig()
	if err != nil {
		return err
	}
	err = a.updateCloudConfig(a.Context, "shutdown")
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
	a.permissions = make(map[string][]string)
	a.permissions["/"] = p

	portableSFTPPrivateKey := a.PrivateKey
	if portableSFTPPrivateKey == "" {
		return fmt.Errorf("missing private key")
	}

	a.portableSFTPFingerprints = append(a.portableSFTPFingerprints, a.PrivateKey)
	a.portablePublicKeys = append(a.portablePublicKeys, a.PublicKey)

	err = a.loadPublicIpAddress(ctx, a.Config)
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

	a.LogFilePath = a.PortableLogFile
	a.LogMaxSize = defaultLogMaxSize
	a.LogMaxBackups = defaultLogMaxBackup
	a.LogMaxAge = defaultLogMaxAge
	a.LogCompress = defaultLogCompress
	a.LogVerbose = a.PortableLogVerbose
	a.LogUTCTime = a.PortableLogUTCTime
	a.shutdown = make(chan bool)
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

	go func() {
		<-a.shutdown
		a.updateCloudConfig(a.Context, "shutdown")
	}()

	return nil
}

func (a *AgentService) Start() error {
	a.Logger().Printf("%v", a.Service)
	if err := a.Service.StartPortableMode(int(a.Port), 0, 0, []string{}, false,
		false, "", "", "", ""); err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(a.Context, time.Second*30)
			defer cancel()
			for {
				if ctx.Err() != nil && a.Service.Error == nil {
					break
				}
				time.Sleep(time.Second * 5)
				err := a.updateCloudConfig(a.Context, "running")
				if err == nil {
					break
				}
			}
		}()
		a.Service.Wait()
		if a.Service.Error == nil {
			a.updateCloudConfig(a.Context, "running")
			os.Exit(0)
		}
	}
	a.updateCloudConfig(a.Context, "shutdown")
	os.Exit(1)
	return nil
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

func (a *AgentService) loadPublicIpAddress(ctx context.Context, clientConfig files_sdk.Config) (err error) {
	client := ip_address.Client{Config: clientConfig}
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
	a.Logger().Printf("Ip whitelist: %v", safeListTmp)
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
	a.Logger().Printf("Config Path: %v", configTempFile.Name())
	return dir, file, config.LoadConfig(dir, file)
}

func (a *AgentService) updateCloudConfig(ctx context.Context, status string) error {
	client := remote_server.Client{Config: a.Config}
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
	a.Logger().Printf("Response: %v", params)
	_, err := client.ConfigurationFile(ctx, params)
	if err != nil {
		a.Logger().Printf("Response: %v", err)
	}
	return err
}
