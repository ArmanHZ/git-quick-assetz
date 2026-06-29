package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Release struct {
	URL             string  `json:"url"`
	AssetsURL       string  `json:"assets_url"`
	UploadURL       string  `json:"upload_url"`
	HTMLURL         string  `json:"html_url"`
	ID              int64   `json:"id"`
	NodeID          string  `json:"node_id"`
	TagName         string  `json:"tag_name"`
	TargetCommitish string  `json:"target_commitish"`
	Name            string  `json:"name"`
	Draft           bool    `json:"draft"`
	Immutable       bool    `json:"immutable"`
	Prerelease      bool    `json:"prerelease"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	PublishedAt     string  `json:"published_at"`
	Assets          []Asset `json:"assets"`
}

type Asset struct {
	URL                string `json:"url"`
	ID                 int64  `json:"id"`
	NodeID             string `json:"node_id"`
	Name               string `json:"name"`
	Label              string `json:"label"`
	ContentType        string `json:"content_type"`
	State              string `json:"state"`
	Size               int64  `json:"size"`
	Digest             string `json:"digest"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func GetReleases(repoURL string) ([]Release, error) {
	ownerAndRepo, err := ExtractOwnerAndRepoNames(repoURL)
	if err != nil {
		return nil, err
	}

	gitAPIReqURL := fmt.Sprintf(`https://api.github.com/repos/%s/releases`, ownerAndRepo)

	resp, err := http.Get(gitAPIReqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var releases []Release

	err = json.Unmarshal(body, &releases)
	if err != nil {
		return nil, err
	}

	return releases, err
}

// URL list has file name and url info
// TODO: Refactor with better names.
func DownloadAssets(urlList [][]string, dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	for _, url := range urlList {
		resp, err := http.Get(url[1])
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		filename := url[0]
		dst := filepath.Join(dir, filename)

		file, err := os.Create(dst)
		if err != nil {
			return err
		}

		_, err = io.Copy(file, resp.Body)
		file.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
