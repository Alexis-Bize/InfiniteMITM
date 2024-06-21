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

package MITMApplicationWelcomePromptUI

import (
	"fmt"
	"infinite-mitm/configs"
	sscg "infinite-mitm/internal/application/ui/tools/sscg"
	"infinite-mitm/pkg/errors"
	"infinite-mitm/pkg/proxy"
	"infinite-mitm/pkg/smartcache"
	"infinite-mitm/pkg/sysutilities"
	"infinite-mitm/pkg/theme"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

type PromptOption int

const (
	// Welcome
	Start PromptOption = iota
	InstallRootCertificate
	Tools
	Credits
	Exit
	// Tools
	PingServers
	ForceKillProxy
	ClearSmartCache
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
	Tools:                  "🛠️ Tools [new]",
	Credits:                "🤝 Show Credits",
	Exit:                   "👋 Exit",
	// Tools
	PingServers:            "🌎 Ping Servers",
	ForceKillProxy:         "🛑 Force Kill Proxy",
	ClearSmartCache:        "🧹 Clear SmartCache",
	// Credits
	Author:                 "Made by: Zeny IC",
	Supporter:              "Supporter: Grunt.API",
	GitHub:                 "Source code: GitHub",
	// Generic
	GoBack:                 "← Go back",
}

var hasRootCertificateInstalled *bool

func (d PromptOption) String() string {
	return optionToString[d]
}

func (d PromptOption) Is(option string) bool {
	return d.String() == option
}

func WelcomePrompt(rootCertificateInstalled bool) (string, *errors.MITMError) {
	hasRootCertificateInstalled = &rootCertificateInstalled

	var selected string
	var options []huh.Option[string]

	if rootCertificateInstalled {
		options = append(options, huh.NewOption(Start.String(), Start.String()))
	} else {
		options = append(options, huh.NewOption(InstallRootCertificate.String(), InstallRootCertificate.String()))
	}

	options = append(
		options,
		huh.NewOption(Tools.String(), Tools.String()),
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

	if Tools.Is(selected) {
		return showTools()
	} else if Credits.Is(selected) {
		return showCredits()
	}

	return selected, nil
}

func showTools() (string, *errors.MITMError) {
	var selected string
	var options []huh.Option[string]

	options = append(
		options,
		huh.NewOption(PingServers.String(), PingServers.String()),
		huh.NewOption(ForceKillProxy.String(), ForceKillProxy.String()),
		huh.NewOption(ClearSmartCache.String(), ClearSmartCache.String()),
		huh.NewOption(GoBack.String(), GoBack.String()),
	)

	huh.NewSelect[string]().
		Title("Tools:").
		Options(options...).
		Value(&selected).
		WithTheme(theme.ThemeMITM()).
		Run()

	if ClearSmartCache.Is(selected) {
		spinner.New().Title("Clearing local cached files...").
		TitleStyle(lipgloss.NewStyle().
			Foreground(theme.ColorNormalFg)).
			Run()

		smartcache.Flush()
		return showTools()
	}

	if ForceKillProxy.Is(selected) {
		spinner.New().Title("Killing active proxy...").
			TitleStyle(lipgloss.NewStyle().
				Foreground(theme.ColorNormalFg)).
				Run()

		proxy.ToggleProxy("off")
		return showTools()
	}

	if PingServers.Is(selected) {
		sscg.Create()
		return showTools()
	}

	return WelcomePrompt(*hasRootCertificateInstalled)
}

func showCredits() (string, *errors.MITMError) {
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
		sysutilities.OpenBrowser("https://x.com/zeny_ic")
	case Supporter.String():
		sysutilities.OpenBrowser("https://x.com/gruntdotapi")
	case GitHub.String():
		sysutilities.OpenBrowser(configs.GetConfig().Repository)
	}

	return WelcomePrompt(*hasRootCertificateInstalled)
}
