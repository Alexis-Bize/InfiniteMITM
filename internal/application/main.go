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
	"infinite-mitm/pkg/mitm"
	"infinite-mitm/pkg/resources"
	"infinite-mitm/pkg/sysutilities"
	"infinite-mitm/pkg/theme"
	"infinite-mitm/pkg/updater"
	"os"
	"runtime"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

func Init() *errors.MITMError {
	var mitmErr *errors.MITMError

	spinner.New().Title("Looking for updates...").
		TitleStyle(lipgloss.NewStyle().
		Foreground(theme.ColorNormalFg)).
		Run()

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
			os.Exit(1)
		}
	}

	if runtime.GOOS == "windows" {
		spinner.New().Title("Verifying admin privileges...").
			TitleStyle(lipgloss.NewStyle().
				Foreground(theme.ColorNormalFg)).
				Run()

		if !sysutilities.IsAdmin() {
			sysutilities.RunAsAdmin()
		}
	}

	spinner.New().Title("Verifying local assets integrity...").
		TitleStyle(lipgloss.NewStyle().
		Foreground(theme.ColorNormalFg)).
		Run()

	mitmErr = resources.CreateRootAssets(); if mitmErr != nil {
		return mitmErr
	}

	_, mitmErr = mitm.ReadClientMITMConfig()
	if mitmErr != nil && errors.ErrMITMYamlSchemaOutdatedException == mitmErr.Unwrap() {
		return nil
	}

	spinner.New().Title("Looking for root certificate...").
		TitleStyle(lipgloss.NewStyle().
		Foreground(theme.ColorNormalFg)).
		Run()

	_, mitmErr = sysutilities.CheckForRootCertificate(configs.GetConfig().Proxy.Certificate.Name);
	if mitmErr != nil {
		return mitmErr
	}

	return nil
}
