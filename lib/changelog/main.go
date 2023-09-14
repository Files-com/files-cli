package changelog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/Files-com/files-cli/lib/version"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

const (
	Releases = "https://api.github.com/repos/Files-com/files-cli/releases"
	Tags     = "https://api.github.com/repos/Files-com/files-cli/tags"
)

func ParseTag(tag string) string {
	if tag == "" {
		return tag
	}

	re := regexp.MustCompile("[0-9/.]+")
	result := re.FindAllString(tag, -1)
	if len(result) > 0 {
		return fmt.Sprintf("v%v", result[0])
	}

	return tag
}

func GetLog(ctx context.Context, cmd *cobra.Command, config files_sdk.Config, tag string, prepend string) error {
	if tag == "" {
		tag = "latest"
	} else {
		var err error
		tag, err = url.JoinPath("tags", tag)
		if err != nil {
			return err
		}
	}

	path, err := url.JoinPath(Releases, tag)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return err
	}
	resp, err := config.Do(req)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var release Release
	err = json.Unmarshal(data, &release)
	if err != nil {
		return err
	}

	if release.Body == "" {
		return nil
	}

	body := strings.Replace(release.Body, "## Changelog", fmt.Sprintf("## %v - %v", release.TagName, release.CreatedAt), 1)

	if prepend != "" {
		body = prepend + body
	}

	out, err := glamour.Render(body, "dark")

	fmt.Fprintf(cmd.OutOrStdout(), out)
	return err
}

func GetAllTags(ctx context.Context, config files_sdk.Config, min string, max string) ([]string, error) {
	var versions []string
	req, err := http.NewRequestWithContext(ctx, "GET", Tags, nil)

	if err != nil {
		return versions, err
	}
	resp, err := config.Do(req)
	if err != nil {
		return versions, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return versions, err
	}
	var tags []Tag
	err = json.Unmarshal(data, &tags)
	if err != nil {
		return versions, err
	}

	minVersion, err := version.New(min)
	if err != nil {
		return versions, err
	}

	maxVersion, err := version.New(max)
	if err != nil {
		return versions, err
	}

	for _, tag := range tags {
		tagVersion, err := version.New(tag.Name)
		if err != nil {
			return versions, err
		}

		if maxVersion.Greater(tagVersion) {
			continue
		}

		if tagVersion.Greater(minVersion) {
			break
		}

		versions = append(versions, fmt.Sprintf("v%v", tagVersion.String()))
	}

	return versions, nil
}

type Release struct {
	Url             string    `json:"url"`
	HtmlUrl         string    `json:"html_url"`
	AssetsUrl       string    `json:"assets_url"`
	UploadUrl       string    `json:"upload_url"`
	TarballUrl      string    `json:"tarball_url"`
	ZipballUrl      string    `json:"zipball_url"`
	DiscussionUrl   string    `json:"discussion_url"`
	Id              int       `json:"id"`
	NodeId          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Body            string    `json:"body"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Author          struct {
		Login             string `json:"login"`
		Id                int    `json:"id"`
		NodeId            string `json:"node_id"`
		AvatarUrl         string `json:"avatar_url"`
		GravatarId        string `json:"gravatar_id"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Assets []struct {
		Url                string    `json:"url"`
		BrowserDownloadUrl string    `json:"browser_download_url"`
		Id                 int       `json:"id"`
		NodeId             string    `json:"node_id"`
		Name               string    `json:"name"`
		Label              string    `json:"label"`
		State              string    `json:"state"`
		ContentType        string    `json:"content_type"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		Uploader           struct {
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
	} `json:"assets"`
}

type Tag struct {
	Name       string `json:"name"`
	ZipballUrl string `json:"zipball_url"`
	TarballUrl string `json:"tarball_url"`
	Commit     struct {
		Sha string `json:"sha"`
		Url string `json:"url"`
	} `json:"commit"`
	NodeId string `json:"node_id"`
}
