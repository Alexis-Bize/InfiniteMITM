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
	"bytes"
	"embed"
	"fmt"
	"infinite-mitm/configs"
	errors "infinite-mitm/pkg/modules/errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

var certName = configs.GetConfig().Proxy.Certificate.Name

func CreateRootAssets(f *embed.FS) error {
	spinner.New().Title("Checking local assets integrity...").Run()

	projectDir := configs.GetConfig().Extra.ProjectDir
	resourcesPath := filepath.Join(projectDir, "resources")

	sourceMITMTemplate := "assets/resources/shared/templates/mitm.yaml"
	outputMITMTemplate := filepath.Join(projectDir, filepath.Base(sourceMITMTemplate))

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return errors.Create(errors.ErrFatalException, err.Error())
	}

	if _, err := os.Stat(outputMITMTemplate); os.IsNotExist(err) {
		templateFile, err := f.Open(sourceMITMTemplate)
		if err != nil {
			return errors.Create(errors.ErrFatalException, err.Error())
		}
		defer templateFile.Close()

		destinationFile, err := os.Create(outputMITMTemplate)
		if err != nil {
			return errors.Create(errors.ErrFatalException, err.Error())
		}
		defer destinationFile.Close()

		if _, err = io.Copy(destinationFile, templateFile); err != nil {
			return errors.Create(errors.ErrFatalException, err.Error())
		}

		if err = destinationFile.Sync(); err != nil {
			return errors.Create(errors.ErrFatalException, err.Error())
		}
	}

	if _, err := os.Stat(resourcesPath); os.IsNotExist(err) {
		previewDir := projectDir

		if runtime.GOOS == "window" {
			previewDir = strings.ReplaceAll(previewDir, "/", "\\")
		}

		var ignoreResourcesCreation bool
		err := huh.NewConfirm().
			Title("Would you like to create a resources directory to help you organize your work?").
			Description(fmt.Sprintf("The directory will be created under \"%s\"", previewDir)).
			Affirmative("No, thanks. I know what I'm doing.").
			Negative("Yes please!").
			Value(&ignoreResourcesCreation).
			Run()

		if err == nil && !ignoreResourcesCreation {
			spinner.New().Title("Creating resources...").Run()

			os.MkdirAll(resourcesPath, 0755)
			os.MkdirAll(filepath.Join(resourcesPath, "ugc", "maps"), 0755)
			os.MkdirAll(filepath.Join(resourcesPath, "ugc", "enginegamevariants", "cgui-markups"), 0755)
			os.MkdirAll(filepath.Join(resourcesPath, "bin", "flags"), 0755)
			os.MkdirAll(filepath.Join(resourcesPath, "tools", "InfiniteVariantToolCLI"), 0755)
			os.MkdirAll(filepath.Join(resourcesPath, "json"), 0755)
		}
	}

	return nil
}

func CheckForRootCertificate() (bool, error) {
	spinner.New().Title("Checking root certificate...").Run()

	var installed bool
	var err error

	switch runtime.GOOS {
	case "windows":
		installed, err = checkForRootCertificateOnWindows()
	case "darwin":
		installed, err = checkForRootCertificateOnDarwin()
	default:
		return false, nil
	}

	return installed, err
}

func checkForRootCertificateOnWindows() (bool, error) {
	cmd := exec.Command("certutil", "-verifystore", "-user", "root", certName)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 2148073489 {
				return false, nil
			}
		}

		return false, errors.Create(errors.ErrRootCertificateException, err.Error())
	}

	if strings.TrimSpace(out.String()) == "" {
		return false, nil
	}

	return true, nil
}

func checkForRootCertificateOnDarwin() (bool, error) {
	cmd := exec.Command("security", "find-certificate", "-a", "-c", certName)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, errors.Create(errors.ErrRootCertificateException, err.Error())
	}

	if strings.TrimSpace(out.String()) == "" {
		return false, nil
	}

	return true, nil
}
