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

package updater

import (
	"encoding/json"
	"fmt"
	"infinite-mitm/configs"
	"infinite-mitm/pkg/errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/Masterminds/semver/v3"
)

type Release struct {
	TagName string `json:"tag_name"`
}

type Releases []Release

func CheckForUpdates() (bool, string, *errors.MITMError) {
	owner, repo := extractOwnerAndRepo()
	if owner == "" || repo == "" {
		return false, "", errors.Create(errors.ErrUpdaterException, "could not extract repository details")
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		return false, "", errors.Create(errors.ErrHTTPRequestException, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, "", errors.Create(errors.ErrHTTPCodeError, fmt.Sprintf("status code: %d", resp.StatusCode))
	}

	var releases Releases
	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		return false, "", errors.Create(errors.ErrJSONReadException, err.Error())
	}

	latestRelease := releases[0]
	newVersionAvailable, err := isNewVersionAvailable(configs.GetConfig().Version, latestRelease.TagName)
	if err != nil {
		return false, "", errors.Create(errors.ErrUpdaterException, err.Error())
	} else if !newVersionAvailable {
		return false, "", nil
	}

	return true, latestRelease.TagName, nil
}

func isNewVersionAvailable(currentVersion, latestVersion string) (bool, error) {
	cv, err := semver.NewVersion(currentVersion)
	if err != nil {
		return false, err
	}

	lv, err := semver.NewVersion(latestVersion)
	if err != nil {
		return false, err
	}

	return lv.Compare(cv) == 1, nil
}

func extractOwnerAndRepo() (string, string) {
	githubUrl := configs.GetConfig().Repository
	u, err := url.Parse(githubUrl)
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
