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

package MITMApplicationWelcomePromptUICreditsComponent

import (
	"infinite-mitm/configs"
	"infinite-mitm/pkg/errors"
	"infinite-mitm/pkg/sysutilities"
	"infinite-mitm/pkg/theme"

	"github.com/charmbracelet/huh"
)

type PromptOption int

const (
	Author PromptOption = iota
	Supporter
	GitHub
	GoBack
)

var optionToString = map[PromptOption]string{
	Author:    "Made by: Zeny IC",
	Supporter: "Supporter: Grunt.API",
	GitHub:    "Source code: GitHub",
	GoBack:    "← Go back",
}

func (d PromptOption) String() string {
	return optionToString[d]
}

func (d PromptOption) Is(option string) bool {
	return d.String() == option
}

func Run() *errors.MITMError {
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
		return Run()
	case Supporter.String():
		sysutilities.OpenBrowser("https://x.com/gruntdotapi")
		return Run()
	case GitHub.String():
		sysutilities.OpenBrowser(configs.GetConfig().Repository)
		return Run()
	}

	return nil
}
