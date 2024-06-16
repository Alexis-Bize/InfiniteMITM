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
	"embed"
	"fmt"
	"infinite-mitm/configs"
	mitm "infinite-mitm/internal/application/services/mitm"
	"infinite-mitm/pkg/errors"
	"infinite-mitm/pkg/resources"
	"infinite-mitm/pkg/sysutilities"
	"infinite-mitm/pkg/updater"

	"github.com/charmbracelet/huh"
)

func Init(f *embed.FS) *errors.MITMError {
	var mitmErr *errors.MITMError

	updateAvailable, latest, _ := updater.CheckForUpdates()
	if updateAvailable {
		var ignoreUpdate bool
		huh.NewConfirm().
			Title(fmt.Sprintf("🔥 A new version (%s) is available; would you like to download it?", latest)).
			Affirmative("Later").
			Negative("Yes please!").
			Value(&ignoreUpdate).
			Run()

		if !ignoreUpdate {
			sysutilities.OpenBrowser(fmt.Sprintf(configs.GetConfig().Repository + "/releases/tag/%s", latest))
			return nil
		}
	}

	if !sysutilities.IsAdmin() {
		sysutilities.RunAsAdmin()
	}

	mitmErr = resources.CreateRootAssets(f); if mitmErr != nil {
		return mitmErr
	}

	_, mitmErr = resources.CheckForRootCertificate(); if mitmErr != nil {
		return mitmErr
	}

	_, mitmErr = mitm.ReadClientMITMConfig()
	if mitmErr != nil && errors.ErrMITMYamlSchemaOutdatedException == mitmErr.Unwrap() {
		return nil
	}

	return mitmErr
}
