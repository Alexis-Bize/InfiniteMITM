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

package resources

import (
	"bytes"
	"embed"
	"infinite-mitm/configs"
	"infinite-mitm/pkg/errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/huh/spinner"
)

var certName = configs.GetConfig().Proxy.Certificate.Name
var projectDir = configs.GetConfig().Extra.ProjectDir
var resourcesPath = filepath.Join(projectDir, "resources")
var resourcesDirPaths = map[string]string{
	"root": resourcesPath,
	"ugc":              filepath.Join(resourcesPath, "ugc"),
	"ugc-maps":         filepath.Join(resourcesPath, "ugc", "maps"),
	"ugc-egv":          filepath.Join(resourcesPath, "ugc", "enginegamevariants"),
	"ugc-egv-cgui":     filepath.Join(resourcesPath, "ugc", "enginegamevariants", "cgui-markups"),
	"bin":              filepath.Join(resourcesPath, "bin"),
	"bin-flags":        filepath.Join(resourcesPath, "bin", "flags"),
	"tools":            filepath.Join(resourcesPath, "tools"),
	"tools-ivt":        filepath.Join(resourcesPath, "tools", "InfiniteVariantTool"),
	"json":             filepath.Join(resourcesPath, "json"),
	"smartcache":       filepath.Join(resourcesPath, "cache"),
	"downloads":        filepath.Join(resourcesPath, "downloads"),
}

func GetSmartCacheDirPath() string {
	return resourcesDirPaths["smartcache"]
}

func GetDownloadsDirPath() string {
	return resourcesDirPaths["downloads"]
}

func CreateRootAssets(f *embed.FS) *errors.MITMError {
	spinner.New().Title("Checking local assets integrity...").Run()

	sourceMITMTemplate := "assets/resources/shared/templates/mitm.yaml"
	outputMITMTemplate := filepath.Join(projectDir, filepath.Base(sourceMITMTemplate))

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
		return nil
	}

	if _, err := os.Stat(outputMITMTemplate); os.IsNotExist(err) {
		templateFile, err := f.Open(sourceMITMTemplate)
		if err != nil {
			log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
			return nil
		}
		defer templateFile.Close()

		destinationFile, err := os.Create(outputMITMTemplate)
		if err != nil {
			log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
			return nil
		}
		defer destinationFile.Close()

		if _, err = io.Copy(destinationFile, templateFile); err != nil {
			log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
			return nil
		}

		if err = destinationFile.Sync(); err != nil {
			log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
			return nil
		}
	}

	for _, v := range resourcesDirPaths {
		os.MkdirAll(v, 0755)
	}

	return nil
}

func CheckForRootCertificate() (bool, *errors.MITMError) {
	spinner.New().Title("Checking root certificate...").Run()

	var installed bool
	var mitmErr *errors.MITMError

	switch runtime.GOOS {
	case "windows":
		installed, mitmErr = checkForRootCertificateOnWindows()
	case "darwin":
		installed, mitmErr = checkForRootCertificateOnDarwin()
	default:
		installed, mitmErr = false, nil
	}

	if mitmErr != nil {
		return false, mitmErr
	}

	return installed, nil
}

func checkForRootCertificateOnWindows() (bool, *errors.MITMError) {
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

		return false, errors.Create(errors.ErrProxyCertificateException, err.Error())
	}

	if strings.TrimSpace(out.String()) == "" {
		return false, errors.Create(errors.ErrProxyCertificateException, "missing root certificate")
	}

	return true, nil
}

func checkForRootCertificateOnDarwin() (bool, *errors.MITMError) {
	cmd := exec.Command("security", "find-certificate", "-a", "-c", certName)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, errors.Create(errors.ErrProxyCertificateException, err.Error())
	}

	if strings.TrimSpace(out.String()) == "" {
		return false, errors.Create(errors.ErrProxyCertificateException, "missing root certificate")
	}

	return true, nil
}
