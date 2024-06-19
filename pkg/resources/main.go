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
	"embed"
	"infinite-mitm/configs"
	"infinite-mitm/pkg/errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

var projectDir string
var resourcesPath string
var dirPaths map[string]string

func init() {
	projectDir = configs.GetConfig().Extra.ProjectDir
	resourcesPath = filepath.Join(projectDir, "resources")
	dirPaths = map[string]string{
		// ~/InfiniteMITM
		"smartcache":       filepath.Join(projectDir, "cache"),
		"downloads":        filepath.Join(projectDir, "downloads"),
		// ~/InfiniteMITM/resources
		"ugc":              filepath.Join(resourcesPath, "ugc"),
		"ugc-maps":         filepath.Join(resourcesPath, "ugc", "maps"),
		"ugc-egv":          filepath.Join(resourcesPath, "ugc", "enginegamevariants"),
		"ugc-egv-cgui":     filepath.Join(resourcesPath, "ugc", "enginegamevariants", "cgui-markups"),
		"bin":              filepath.Join(resourcesPath, "bin"),
		"bin-flags":        filepath.Join(resourcesPath, "bin", "flags"),
		"tools":            filepath.Join(resourcesPath, "tools"),
		"tools-ivt":        filepath.Join(resourcesPath, "tools", "InfiniteVariantTool"),
		"json":             filepath.Join(resourcesPath, "json"),
	}
}

func GetDirPaths() map[string]string {
	return dirPaths
}

func GetSmartCacheDirPath() string {
	return dirPaths["smartcache"]
}

func GetDownloadsDirPath() string {
	return dirPaths["downloads"]
}

func CreateRootAssets(f *embed.FS) *errors.MITMError {
	sourceMITMTemplate := "assets/resources/shared/templates/mitm.yaml"
	outputMITMTemplate := filepath.Join(projectDir, filepath.Base(sourceMITMTemplate))

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
	}

	if _, err := os.Stat(outputMITMTemplate); os.IsNotExist(err) {
		templateFile, err := f.Open(sourceMITMTemplate)
		if err != nil {
			log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
		}
		defer templateFile.Close()

		destinationFile, err := os.Create(outputMITMTemplate)
		if err != nil {
			log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
		}
		defer destinationFile.Close()

		if _, err = io.Copy(destinationFile, templateFile); err != nil {
			log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
		}

		if err = destinationFile.Sync(); err != nil {
			log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
		}
	}

	for _, v := range dirPaths {
		os.MkdirAll(v, 0755)
	}

	return nil
}
