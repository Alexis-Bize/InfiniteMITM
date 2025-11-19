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

package MITMApplication

import (
	"fmt"
	"infinite-mitm/configs"
	"infinite-mitm/pkg/errors"
	"infinite-mitm/pkg/integrity"
	"infinite-mitm/pkg/mitm"
	"infinite-mitm/pkg/resources"
	"infinite-mitm/pkg/spinner"
	"infinite-mitm/pkg/sysutilities"
	"infinite-mitm/pkg/updater"
	"os"
	"path/filepath"
	"runtime"

	"github.com/charmbracelet/huh"
)

func Init() *errors.MITMError {
	var mitmErr *errors.MITMError

	spinner.Run("Looking for application update...")
	applicationUpdateAvailable, latest, _ := updater.CheckForApplicationUpdate()
	if applicationUpdateAvailable {
		var ignoreUpdate bool
		huh.NewConfirm().
			Title(fmt.Sprintf("🔥 A new version (%s) is available; would you like to download it?", latest)).
			Affirmative("Later").
			Negative("Yes please!").
			Value(&ignoreUpdate).
			Run()

		if !ignoreUpdate {
			spinner.Run("Attempting to open your browser...")
			sysutilities.OpenBrowser(fmt.Sprintf(configs.GetConfig().Repository + "/releases/tag/%s", latest))
			os.Exit(0)
		}
	}

	if runtime.GOOS == "windows" {
		spinner.Run("Verifying admin privileges...")
		if !sysutilities.IsAdmin() {
			sysutilities.RunAsAdmin()
		}
	}

	spinner.Run("Looking for root certificate...")
	_, mitmErr = sysutilities.CheckForRootCertificate(configs.GetConfig().Proxy.Certificate.Name);
	if mitmErr != nil {
		return mitmErr
	}

	spinner.Run("Verifying local assets integrity...")
	mitmErr = resources.CreateRootAssets(); if mitmErr != nil {
		return mitmErr
	}

	spinner.Run("Looking for server list update...")
	serverListUpdateAvailable, qosServersFile, _ := updater.CheckForIntegrityFileUpdate(integrity.QOSServersFilename)
	if serverListUpdateAvailable {
		spinner.Run("Updating existing server list with latest version...")
		data, mitmErr := sysutilities.DownloadFile(qosServersFile.DownloadUrl)

		if mitmErr == nil {
			key := resources.GetDirPaths()[resources.PubDirKey]
			os.WriteFile(filepath.Join(key, integrity.QOSServersFilename), data, 0644)

			content, mitmErr := integrity.ReadIntegrityConfig();
			if mitmErr == nil {
				content.Files[integrity.QOSServersFilename] = integrity.FileContent{
					Sha1: qosServersFile.Sha,
				}

				integrity.WriteIntegrityFile(content)

				huh.NewConfirm().
					Title("⚠️ The server list has been updated; you may need to reselect your preferred ones.").
					Affirmative("Sounds good!").
					Negative("").
					Run()
			}
		}
	}

	_, mitmErr = mitm.ReadClientMITMConfig()
	if mitmErr != nil {
		return mitmErr
	}

	return nil
}
