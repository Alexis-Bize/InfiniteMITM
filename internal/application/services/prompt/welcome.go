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
	errors "infinite-mitm/pkg/modules/errors"

	"github.com/charmbracelet/huh"
)

func Welcome(rootCertificateInstalled bool) (string, *errors.MITMError) {
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
		huh.NewOption(Exit.String(), Exit.String()),
	)

	err := huh.NewSelect[string]().
		Title(fmt.Sprintf("%s - %s", configs.GetConfig().Name, configs.GetConfig().Version)).
		Options(options...).
		Value(&selected).
		Run()

	if err != nil {
		return "", errors.Create(errors.ErrAppException, err.Error())
	}

	return selected, nil
}
