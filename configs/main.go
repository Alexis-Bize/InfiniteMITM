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

package configs

import (
	"embed"
	utilities "infinite-mitm/pkg/modules/utilities"
	"log"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed application.yaml
var f embed.FS
var config Config

type Config struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
	Author      string `yaml:"author"`
	Repository  string `yaml:"repository"`
	Proxy struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"proxy"`
	Certificate struct {
		Name string `yaml:"name"`
	} `yaml:"certificate"`
	Extra struct {
		ProjectDir string
	}
}

func GetConfig() *Config {
	if config != (Config{}) {
		return &config
	}

	yamlFile, err := f.ReadFile("application.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalln(err)
	}

	home, err := utilities.GetHomeDirectory()
	if err != nil {
		log.Fatalln(err)
	}

	config.Extra.ProjectDir = path.Join(home, strings.Replace(config.Name, " ", "", -1))
	return &config
}
