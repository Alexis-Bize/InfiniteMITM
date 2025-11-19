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

package mitm

import (
	"fmt"
	"infinite-mitm/configs"
	"infinite-mitm/pkg/domains"
	"infinite-mitm/pkg/errors"
	"infinite-mitm/pkg/smartcache"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type TrafficDisplay string
type TrafficOptions struct {
	TrafficDisplay TrafficDisplay
}

type YAMLOptions struct {
	SmartCache     smartcache.SmartCacheYAMLOptions `yaml:"smart_cache"`
	TrafficDisplay TrafficDisplay `yaml:"traffic_display"`
}

type YAML struct {
	Domains domains.YAMLDomains `yaml:"domains"`
	Options YAMLOptions `yaml:"options"`
	Version int `yaml:"version"`
}

const (
	MITMFilename = "mitm.yaml"
	MITMVersion = 1
)

const (
	TrafficAll        TrafficDisplay = "all"
	TrafficOverrides  TrafficDisplay = "overrides"
	TrafficSmartCache TrafficDisplay = "smart_cached"
	TrafficSilent     TrafficDisplay = "silent"
)

var MITMFilepath = filepath.Join(configs.GetConfig().Extra.ProjectDir, MITMFilename)

func ReadClientMITMConfig() (YAML, *errors.MITMError) {
	yamlFile, err := os.ReadFile(MITMFilepath); if err != nil {
		return YAML{}, errors.Create(errors.ErrYAMLReadException, err.Error())
	}

	var content YAML
	if err = yaml.Unmarshal(yamlFile, &content); err != nil {
		return YAML{}, errors.Create(errors.ErrYAMLReadException, err.Error())
	}

	if content.Version != MITMVersion {
		return YAML{}, errors.Create(errors.ErrMITMYamlSchemaOutdatedException, fmt.Sprintf("your %s is outdated, please delete it and restart the application to fix this issue.", MITMFilename))
	}

	return content, nil
}

func WriteMITMFile(content YAML) {
	buffer, err := yaml.Marshal(content); if err == nil {
		os.WriteFile(MITMFilepath, buffer, 0644)
	}
}
