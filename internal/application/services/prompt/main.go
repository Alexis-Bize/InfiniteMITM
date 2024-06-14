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

package MITMApplicationPromptService

import (
	"fmt"
	"infinite-mitm/configs"
	theme "infinite-mitm/internal/application/services/ui/theme"
	errors "infinite-mitm/pkg/errors"
	"infinite-mitm/pkg/smartcache"
	utilities "infinite-mitm/pkg/utilities"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

type PromptOption int

const (
	// Welcome
	Start PromptOption = iota
	InstallRootCertificate
	ForceKillProxy
	ClearSmartCache
	Credits
	Exit
	// Credits
	Author
	Supporter
	GitHub
	// Generic
	GoBack
)

var optionToString = map[PromptOption]string{
	// Welcome
	Start:                  "🔒 Start Proxy Server",
	InstallRootCertificate: "🔐 Install Root Certificate",
	ForceKillProxy:         "🛑 Force Kill Proxy",
	ClearSmartCache:        "🧹 Clear SmartCache",
	Credits:                "🤝 Credits",
	Exit:                   "👋 Exit",
	// Credits
	Author:                 "Made by: Zeny IC",
	Supporter:              "Supporter: Grunt.API",
	GitHub:                 "Source code: GitHub",
	// Generic
	GoBack:                 "← Go back",
}

func (d PromptOption) String() string {
	return optionToString[d]
}

func (d PromptOption) Is(option string) bool {
	return d.String() == option
}

func WelcomePrompt(rootCertificateInstalled bool) (string, *errors.MITMError) {
	var selected string
	var options []huh.Option[string]

	if rootCertificateInstalled {
		options = append(options, huh.NewOption(Start.String(), Start.String()))
	} else {
		options = append(options, huh.NewOption(InstallRootCertificate.String(), InstallRootCertificate.String()))
	}

	options = append(
		options,
		huh.NewOption(ForceKillProxy.String(), ForceKillProxy.String()),
		huh.NewOption(ClearSmartCache.String(), ClearSmartCache.String()),
		huh.NewOption(Credits.String(), Credits.String()),
		huh.NewOption(Exit.String(), Exit.String()),
	)

	err := huh.NewSelect[string]().
		Title(fmt.Sprintf("%s - %s", configs.GetConfig().Name, configs.GetConfig().Version)).
		Options(options...).
		Value(&selected).
		WithTheme(theme.ThemeMITM()).
		Run()

	if err != nil {
		return "", errors.Create(errors.ErrPromptException, err.Error())
	}

	if Credits.Is(selected) {
		return showCredits(rootCertificateInstalled)
	} else if ClearSmartCache.Is(selected) {
		spinner.New().Title("Clearing local cached files...").Run()
		smartcache.Flush()
		return WelcomePrompt(rootCertificateInstalled)
	}

	return selected, nil
}

func showCredits(rootCertificateInstalled bool) (string, *errors.MITMError) {
	var selected string
	var options []huh.Option[string]

	options = append(
		options,
		huh.NewOption(Author.String(), Author.String()),
		huh.NewOption(Supporter.String(), Supporter.String()),
		huh.NewOption(GitHub.String(), GitHub.String()),
		huh.NewOption(GoBack.String(), GoBack.String()),
	)

	huh.NewSelect[string]().
		Title("Credits:").
		Options(options...).
		Value(&selected).
		WithTheme(theme.ThemeMITM()).
		Run()

	switch selected {
	case Author.String():
		utilities.OpenBrowser("https://x.com/zeny_ic")
	case Supporter.String():
		utilities.OpenBrowser("https://x.com/gruntdotapi")
	case GitHub.String():
		utilities.OpenBrowser(configs.GetConfig().Repository)
	}

	return WelcomePrompt(rootCertificateInstalled)
}
