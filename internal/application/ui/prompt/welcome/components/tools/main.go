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

package MITMApplicationWelcomePromptUIToolsComponent

import (
	selectServersUI "infinite-mitm/internal/application/ui/tools/select-servers"
	"infinite-mitm/pkg/proxy"
	"infinite-mitm/pkg/smartcache"
	"infinite-mitm/pkg/theme"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

type PromptOption int

const (
	SelectServers PromptOption = iota
	ForceKillProxy
	ClearSmartCache
	GoBack
)

var optionToString = map[PromptOption]string{
	SelectServers:   "🌎 Select Servers",
	ForceKillProxy:  "🛑 Force Kill Proxy",
	ClearSmartCache: "🧹 Clear SmartCache",
	GoBack:          "← Go back",
}

var hasRootCertificateInstalled *bool

func (d PromptOption) String() string {
	return optionToString[d]
}

func (d PromptOption) Is(option string) bool {
	return d.String() == option
}

func Run() {
	var selected string
	var options []huh.Option[string]

	options = append(
		options,
		huh.NewOption(SelectServers.String(), SelectServers.String()),
		huh.NewOption(ForceKillProxy.String(), ForceKillProxy.String()),
		huh.NewOption(ClearSmartCache.String(), ClearSmartCache.String()),
		huh.NewOption(GoBack.String(), GoBack.String()),
	)

	huh.NewSelect[string]().
		Title("🛠️ Tools").
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
	} else if ForceKillProxy.Is(selected) {
		spinner.New().Title("Killing active proxy...").
			TitleStyle(lipgloss.NewStyle().
				Foreground(theme.ColorNormalFg)).
				Run()

		proxy.ToggleProxy("off")
	} else if SelectServers.Is(selected) {
		selectServersUI.Create()
	}
}
