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
	"fmt"
	"infinite-mitm/configs"
	"infinite-mitm/pkg/errors"

	"infinite-mitm/pkg/github"

	"infinite-mitm/pkg/integrity"

	"github.com/Masterminds/semver/v3"
)

const (
	refBranch = "main"
)

func CheckForApplicationUpdate() (bool, string, *errors.MITMError) {
	releases, err := github.GetReleases()
	if err != nil {
		return false, "", err
	}

	latestRelease := (*releases)[0]
	newVersionAvailable, versionErr := isNewVersionAvailable(configs.GetConfig().Version, latestRelease.TagName)
	if versionErr != nil {
		return false, "", errors.Create(errors.ErrUpdaterException, versionErr.Error())
	} else if !newVersionAvailable {
		return false, "", nil
	}

	return true, latestRelease.TagName, nil
}

func CheckForIntegrityFileUpdate(filename string) (bool, *github.PublicFile, *errors.MITMError) {
	integrity, err := integrity.ReadIntegrityConfig()
	if err != nil {
		return false, nil, err
	}

	val, exists := integrity.Files[filename]
	if !exists {
		return false, nil, errors.Create(errors.ErrUpdaterException, fmt.Sprintf("filename %s not found in integrity.Files", filename))
	}

	pubFile, err := github.GetPublicFile(filename, refBranch)
	if err != nil {
		return false, nil, err
	}

	if pubFile.Sha != val.Sha1 {
		return true, pubFile, nil
	}

	return false, pubFile, nil
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