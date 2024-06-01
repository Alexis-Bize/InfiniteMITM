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

	"github.com/charmbracelet/huh"
)

func Welcome() (string, error) {
	var option string
	err := huh.NewSelect[string]().
		Title(fmt.Sprintf("%s - %s", configs.GetConfig().Name, configs.GetConfig().Version)).
		Options(
			huh.NewOption(Start.String(), Start.String()),
			huh.NewOption(InstallRootCertificate.String(), InstallRootCertificate.String()),
			huh.NewOption(ForceKillProxy.String(), ForceKillProxy.String()),
			huh.NewOption(Exit.String(), Exit.String()),
		).
		Value(&option).
		Run()

	return option, err
}
