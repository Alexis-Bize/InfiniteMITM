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

package integrity

import (
	"fmt"
	"infinite-mitm/configs"
	"infinite-mitm/pkg/errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const QOSServersFilename = "qosservers.json"

type YAML struct {
	Files map[string]FileContent `yaml:"files"`
	Version int `yaml:"version"`
}

type FileContent struct {
	Sha1 string `yaml:"sha1"`
}

const (
	IntegrityFilename = "integrity.yaml"
	IntegrityVersion = 1
)

var IntegrityFilepath = filepath.Join(configs.GetConfig().Extra.ProjectDir, IntegrityFilename)

func ReadIntegrityConfig() (YAML, *errors.MITMError) {
	yamlFile, err := os.ReadFile(IntegrityFilepath); if err != nil {
		return YAML{}, errors.Create(errors.ErrYAMLReadException, err.Error())
	}

	var content YAML
	if err = yaml.Unmarshal(yamlFile, &content); err != nil {
		return YAML{}, errors.Create(errors.ErrYAMLReadException, err.Error())
	}

	if content.Version != IntegrityVersion {
		return YAML{}, errors.Create(errors.ErrMITMYamlSchemaOutdatedException, fmt.Sprintf("your %s is outdated, please delete it and restart the application to fix this issue.", IntegrityFilename))
	}

	return content, nil
}

func WriteIntegrityFile(content YAML) {
	buffer, err := yaml.Marshal(content); if err == nil {
		os.WriteFile(IntegrityFilepath, buffer, 0644)
	}
}
