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
	"infinite-mitm/configs"
	"infinite-mitm/pkg/errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const QOSServersFilename = "qosservers.json"

type YAML struct {
	Files map[string]FileContent `yaml:"files"`
}

type FileContent struct {
	Sha1 string `yaml:"sha1"`
}

const (
	IntegrityConfigFilename = "integrity.yaml"
)

var IntegrityConfigFilepath = filepath.Join(configs.GetConfig().Extra.ProjectDir, IntegrityConfigFilename)

func ReadIntegrityConfig() (YAML, *errors.MITMError) {
	yamlFile, err := os.ReadFile(IntegrityConfigFilepath); if err != nil {
		return YAML{}, errors.Create(errors.ErrYAMLReadException, err.Error())
	}

	var content YAML
	if err = yaml.Unmarshal(yamlFile, &content); err != nil {
		return YAML{}, errors.Create(errors.ErrYAMLReadException, err.Error())
	}

	return content, nil
}

func WriteIntegrityConfig(content YAML) {
	buffer, err := yaml.Marshal(content); if err == nil {
		os.WriteFile(IntegrityConfigFilepath, buffer, 0644)
	}
}
