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

package domains

import (
	"fmt"
	"strings"
)

type DomainType = string
type Domains struct {
	Root      DomainType
	Blobs     DomainType
	Authoring DomainType
	Discovery DomainType
	HaloStats DomainType
	Settings  DomainType
	GameCMS   DomainType
	Economy   DomainType
	Lobby     DomainType
	Skill     DomainType
}

type YAMLDomains struct {
	Root      []YAMLDomainNode `yaml:"root"`
	Blobs     []YAMLDomainNode `yaml:"blobs"`
	Authoring []YAMLDomainNode `yaml:"authoring"`
	Discovery []YAMLDomainNode `yaml:"discovery"`
	HaloStats []YAMLDomainNode `yaml:"stats"`
	Settings  []YAMLDomainNode `yaml:"settings"`
	GameCMS   []YAMLDomainNode `yaml:"gamecms"`
	Economy   []YAMLDomainNode `yaml:"economy"`
	Lobby     []YAMLDomainNode `yaml:"lobby"`
	Skill     []YAMLDomainNode `yaml:"skill"`
}

type YAMLDomainNode struct {
	Path     string   `yaml:"path"`
	Methods  []string `yaml:"methods"`
	Request  YAMLDomainRequestNode `yaml:"request"`
	Response YAMLDomainResponseNode `yaml:"response"`
}

type YAMLDomainRequestNode struct {
	Body    string `yaml:"body"`
	Headers map[string]string `yaml:"headers"`
	Before  YAMLDomainTrafficCommands `yaml:"before"`
}

type YAMLDomainTrafficCommands struct {
	Commands []YAMLDomainTrafficRunCommand `yaml:"commands"`
}

type YAMLDomainTrafficRunCommand struct {
	Run []string `yaml:"run"`
}

type YAMLDomainResponseNode struct {
	Body       string `yaml:"body"`
	Headers    map[string]string `yaml:"headers"`
	StatusCode int `yaml:"code"`
	Before     YAMLDomainTrafficCommands `yaml:"before"`
}

type YAMLContentDomainPair struct {
	Content []YAMLDomainNode
	Domain  DomainType
}

type YAMLContentDomainPairs []YAMLContentDomainPair

const (
	Root      DomainType  = "svc.halowaypoint.com:443"
	Blobs     DomainType  = "blobs-infiniteugc.svc.halowaypoint.com:443"
	Authoring DomainType  = "authoring-infiniteugc.svc.halowaypoint.com:443"
	Discovery DomainType  = "discovery-infiniteugc.svc.halowaypoint.com:443"
	HaloStats DomainType  = "halostats.svc.halowaypoint.com:443"
	Settings  DomainType  = "settings.svc.halowaypoint.com:443"
	GameCMS   DomainType  = "gamecms-hacs.svc.halowaypoint.com:443"
	Economy   DomainType  = "economy.svc.halowaypoint.com:443"
	Lobby     DomainType  = "lobby-hi.svc.halowaypoint.com:443"
	Skill     DomainType  = "skill.svc.halowaypoint.com:443"
)

var HaloWaypointSVCDomains = Domains{
	Root:      Root,
	Blobs:     Blobs,
	Authoring: Authoring,
	Discovery: Discovery,
	HaloStats: HaloStats,
	Settings:  Settings,
	GameCMS:   GameCMS,
	Economy:   Economy,
	Lobby:     Lobby,
	Skill:     Skill,
}

var SmartCachableHostnames = []string{
	DomainToHostname(Blobs),
	DomainToHostname(Authoring),
	DomainToHostname(Discovery),
	DomainToHostname(GameCMS),
	DomainToHostname(HaloStats),
	DomainToHostname(Skill),
}

func GetYAMLContentDomainPairs(yamlDomains YAMLDomains) YAMLContentDomainPairs {
	pairs := YAMLContentDomainPairs{
		{Content: yamlDomains.Root, Domain: HaloWaypointSVCDomains.Root},
		{Content: yamlDomains.Blobs, Domain: HaloWaypointSVCDomains.Blobs},
		{Content: yamlDomains.Authoring, Domain: HaloWaypointSVCDomains.Authoring},
		{Content: yamlDomains.Discovery, Domain: HaloWaypointSVCDomains.Discovery},
		{Content: yamlDomains.HaloStats, Domain: HaloWaypointSVCDomains.HaloStats},
		{Content: yamlDomains.Settings, Domain: HaloWaypointSVCDomains.Settings},
		{Content: yamlDomains.GameCMS, Domain: HaloWaypointSVCDomains.GameCMS},
		{Content: yamlDomains.Economy, Domain: HaloWaypointSVCDomains.Economy},
		{Content: yamlDomains.Lobby, Domain: HaloWaypointSVCDomains.Lobby},
		{Content: yamlDomains.Skill, Domain: HaloWaypointSVCDomains.Skill},
	}

	return pairs
}

func DomainToHostname(domain DomainType) string {
	return strings.Split(domain, ":")[0]
}

func DomainToBaseUrl(domain DomainType) string {
	return fmt.Sprintf("https://%s", DomainToHostname(domain))
}
