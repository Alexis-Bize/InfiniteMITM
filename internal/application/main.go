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
	"infinite-mitm/configs"
	errors "infinite-mitm/pkg/modules/errors"
	Utilities "infinite-mitm/pkg/modules/utilities"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/huh/spinner"
)

var certName = "InfiniteMITM"

func CreateRootAssets(f *embed.FS) error {
	spinner.New().Title("Checking local assets integrity...").Run()
	
	projectDir := configs.GetConfig().Extra.ProjectDir
	sourceMITMTemplate := "assets/resource/shared/templates/mitm.yaml"
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
	cmd := exec.Command("certutil", "-verifystore", "root", certName)
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

func checkForRootCertificateOnDarwin() (bool, error) {
	home, err := Utilities.GetHomeDirectory()
	if err != nil {
		return false, err
	}
	
	keychainPath := home + "/Library/Keychains/login.keychain-db"
	cmd := exec.Command("security", "find-certificate", "-a", "-c", certName, keychainPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return false, errors.Create(errors.ErrRootCertificateException, err.Error())
	}

	if strings.TrimSpace(out.String()) == "" {
		return false, nil
	}

	return true, nil
}
