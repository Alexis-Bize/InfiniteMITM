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
	credits "infinite-mitm/internal/application/ui/prompt/welcome/components/credits"
	tools "infinite-mitm/internal/application/ui/prompt/welcome/components/tools"
	"infinite-mitm/pkg/errors"
	"infinite-mitm/pkg/theme"

	"github.com/charmbracelet/huh"
)

type PromptOption int

const (
	// Welcome
	Start PromptOption = iota
	InstallRootCertificate
	Tools
	Credits
	Exit
)

var optionToString = map[PromptOption]string{
	// Welcome
	Start:                  "🔒 Start Proxy Server",
	InstallRootCertificate: "🔐 Install Root Certificate",
	Tools:                  "🧰 Show Tools",
	Credits:                "🤝 Show Credits",
	Exit:                   "👋 Quit",
}

func (d PromptOption) String() string {
	return optionToString[d]
}

func (d PromptOption) Is(option string) bool {
	return d.String() == option
}

func Run(rootCertificateInstalled bool) (string, *errors.MITMError) {
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
		tools.Run()
		return Run(rootCertificateInstalled)
	} else if Credits.Is(selected) {
		credits.Run()
		return Run(rootCertificateInstalled)
	}

	return selected, nil
}
