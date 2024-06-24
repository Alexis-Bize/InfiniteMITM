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
	"infinite-mitm/configs"
	embedFS "infinite-mitm/internal/application/embed"
	"infinite-mitm/pkg/errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	SmartCacheDirKey = "smartcache"
	DownloadsDirKey  = "downloads"
	ResourcesDirKey  = "resources"
	UGCDirKey        = "ugc"
	UGCMapsDirKey    = "ugc-maps"
	UGCEGVDirKey     = "ugc-egv"
	UGCEGVCGUIDirKey = "ugc-egv-cgui"
	BinDirKey        = "bin"
	BinFlagsDirKey   = "bin-flags"
	ToolsDirKey      = "tools"
	ToolsIVTDirKey   = "tools-ivt"
	JsonDirKey       = "json"
)

var projectDir = configs.GetConfig().Extra.ProjectDir
var dirPaths map[string]string

func init() {
	var resourcesPath = filepath.Join(projectDir, ResourcesDirKey)
	dirPaths = map[string]string{
		// ~/InfiniteMITM
		SmartCacheDirKey:  filepath.Join(projectDir, "cache"),
		DownloadsDirKey:   filepath.Join(projectDir, "downloads"),
		ResourcesDirKey:   filepath.Join(projectDir, "resources"),
		// ~/InfiniteMITM/resources/ugc
		UGCDirKey:         filepath.Join(resourcesPath, "ugc"),
		UGCMapsDirKey:     filepath.Join(resourcesPath, "ugc", "maps"),
		UGCEGVDirKey:      filepath.Join(resourcesPath, "ugc", "enginegamevariants"),
		UGCEGVCGUIDirKey:  filepath.Join(resourcesPath, "ugc", "enginegamevariants", "cgui-markups"),
		// ~/InfiniteMITM/resources/bin
		BinDirKey:         filepath.Join(resourcesPath, "bin"),
		BinFlagsDirKey:    filepath.Join(resourcesPath, "bin", "flags"),
		// ~/InfiniteMITM/resources/tools
		ToolsDirKey:       filepath.Join(resourcesPath, "tools"),
		ToolsIVTDirKey:    filepath.Join(resourcesPath, "tools", "InfiniteVariantTool"),
		// ~/InfiniteMITM/resources/json
		JsonDirKey:        filepath.Join(resourcesPath, "json"),
	}
}

func GetRootPath() string {
	return projectDir
}

func GetDirPaths() map[string]string {
	return dirPaths
}

func GetResourcesPath() string {
	return dirPaths[ResourcesDirKey]
}

func GetSmartCacheDirPath() string {
	return dirPaths[SmartCacheDirKey]
}

func GetDownloadsDirPath() string {
	return dirPaths[DownloadsDirKey]
}

func CreateRootAssets() *errors.MITMError {
	sourceMITMTemplate := "assets/resources/shared/mitm.yaml"
	outputMITMTemplate := filepath.Join(projectDir, filepath.Base(sourceMITMTemplate))

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
	}

	if _, err := os.Stat(outputMITMTemplate); os.IsNotExist(err) {
		templateFile, err := embedFS.Get().Open(sourceMITMTemplate)
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
