// Copyright 2024 Alexis Bize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package github

import (
	"encoding/json"
	"fmt"
	"infinite-mitm/configs"
	"infinite-mitm/pkg/errors"
	"net/http"
	"net/url"
	"strings"
)

type Release struct {
	TagName string `json:"tag_name"`
}

type Releases []Release

type PublicFile struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Sha string `json:"sha"`
	Size int `json:"size"`
	Url string `json:"url"`
	HtmlUrl string `json:"html_url"`
	GitUrl string `json:"git_url"`
	DownloadUrl string `json:"download_url"`
	Type string `json:"type"`
	Content string `json:"content"`
	Encoding string `json:"encoding"`
	Links Links `json:"_links"`
}

type Links struct {
	Self string `json:"self"`
	Git string `json:"git"`
	Html string `json:"html"`
}

func GetReleases() (*Releases, *errors.MITMError) {
	baseUrl, err := getGitHubAPIBaseUrl()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/releases", baseUrl)
	resp, httpErr := http.Get(url)
	if httpErr != nil {
		return  nil, errors.Create(errors.ErrHTTPRequestException, httpErr.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return  nil, errors.Create(errors.ErrHTTPCodeError, fmt.Sprintf("status code: %d", resp.StatusCode))
	}

	var releases Releases
	if decodeErr := json.NewDecoder(resp.Body).Decode(&releases); decodeErr != nil {
		return  nil, errors.Create(errors.ErrJSONReadException, decodeErr.Error())
	}

	return &releases, nil
}

func GetPublicFile(filename string) (*PublicFile, *errors.MITMError) {
	baseUrl, err := getGitHubAPIBaseUrl()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/contents/public/%s", baseUrl, filename)
	resp, httpErr := http.Get(url)
	if httpErr != nil {
		return  nil, errors.Create(errors.ErrHTTPRequestException, httpErr.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Create(errors.ErrHTTPCodeError, fmt.Sprintf("status code: %d", resp.StatusCode))
	}

	var publicFile PublicFile
	if decodeErr := json.NewDecoder(resp.Body).Decode(&publicFile); decodeErr != nil {
		return  nil, errors.Create(errors.ErrJSONReadException, decodeErr.Error())
	}

	return  &publicFile, nil
}

func getGitHubAPIBaseUrl() (string, *errors.MITMError) {
	owner, repo := extractOwnerAndRepo()
	if owner == "" || repo == "" {
		return "", errors.Create(errors.ErrFatalException, "could not extract repository details")
	}

	return fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo), nil
}

func extractOwnerAndRepo() (string, string) {
	githubURL := configs.GetConfig().Repository
	u, err := url.Parse(githubURL)
	if err != nil {
		return "", ""
	}

	path := strings.TrimPrefix(u.Path, "/")
	components := strings.Split(path, "/")
	if len(components) != 2 {
		return "", ""
	}

	return components[0], components[1]
}
