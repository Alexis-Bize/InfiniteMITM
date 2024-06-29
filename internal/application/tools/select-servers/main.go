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

package MITMApplicationSelectServersTool

import (
	"encoding/json"
	"fmt"
	embedFS "infinite-mitm/internal/application/embed"
	"infinite-mitm/pkg/domains"
	"infinite-mitm/pkg/errors"
	"infinite-mitm/pkg/mitm"
	"infinite-mitm/pkg/pattern"
	"infinite-mitm/pkg/resources"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

type QOSServers []QOSServer
type QOSServer struct {
	Region    string `json:"region"`
	ServerURL string `json:"serverUrl"`
}

type PingResult struct {
	Region    string
	ServerURL string
	Ping      int
}

type MinMax struct {
	Min PingResult
	Max PingResult
}

const PingErrorValue = -1
const qosserversFilename = "qosservers.json"

var (
	once    sync.Once
	servers QOSServers
)

func GetQOSServers() QOSServers {
	once.Do(func() {
		serverFilepath := path.Join("assets/resources/shared", qosserversFilename)
		qosservers, err := embedFS.Get().ReadFile(serverFilepath); if err != nil {
			log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
		}

		if err := json.Unmarshal([]byte(qosservers), &servers); err != nil {
			log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
		}
	})

	return servers
}

func GetPingTime(serverURL string) (*probing.Statistics, *errors.MITMError) {
	pinger, err := probing.NewPinger(serverURL)
	if runtime.GOOS == "windows" {
		pinger.SetPrivileged(true)
	}

	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.Create(errors.ErrPingFailedException, err.Error())
	}

	pinger.Count = 1
	pinger.Timeout = time.Second * 5

	err = pinger.Run()
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.Create(errors.ErrPingFailedException, err.Error())
	}

	return pinger.Statistics(), nil
}

func GetMinMax(pairs []PingResult) MinMax {
	minPing := int(^uint(0) >> 1)
	maxPing := 0
	minRegion := ""
	maxRegion := ""
	minServerURL := ""
	maxServerURL := ""

	for _, pair := range pairs {
		if pair.Ping < minPing {
			minRegion = pair.Region
			minServerURL = pair.ServerURL
			minPing = pair.Ping
		}

		if pair.Ping > maxPing {
			maxRegion = pair.Region
			maxServerURL = pair.ServerURL
			maxPing = pair.Ping
		}
	}

	return MinMax{
		Min: PingResult{Region: minRegion, ServerURL: minServerURL, Ping: minPing},
		Max: PingResult{Region: maxRegion, ServerURL: maxServerURL, Ping: maxPing},
	}
}

func ReadLocalQOSServers() QOSServers {
	serverFilepath := filepath.Join(resources.GetDirPaths()[resources.JsonDirKey], qosserversFilename)
	qosservers, err := os.ReadFile(serverFilepath); if err != nil {
		return GetQOSServers()
	}

	var localServers QOSServers
	if err = json.Unmarshal(qosservers, &localServers); err != nil {
		return GetQOSServers()
	}

	return localServers
}

func WriteLocalQOSServers(servers QOSServers) {
	serverFilepath := filepath.Join(resources.GetDirPaths()[resources.JsonDirKey], qosserversFilename)
	buffer, err := json.Marshal(servers); if err == nil {
		os.WriteFile(serverFilepath, buffer, 0644)
	}
}

func WriteMITMFile() {
	filePathPattern := path.Join(pattern.MITMDirPrefix, resources.ResourcesDirKey, resources.JsonDirKey, qosserversFilename)
	content, mitmErr := mitm.ReadClientMITMConfig(); if mitmErr != nil {
		return
	}

	var updatedLobbyDomain = []domains.YAMLDomainNode{}
	for _, config := range content.Domains.Lobby {
		if !strings.HasSuffix(config.Path, "/qosservers") {
			updatedLobbyDomain = append(updatedLobbyDomain, config)
		}
	}

	updatedLobbyDomain = append(updatedLobbyDomain, domains.YAMLDomainNode{
		Path: "/titles/:title/qosservers",
		Methods: []string{"GET"},
		Response: domains.YAMLDomainResponseNode{
			Body: filePathPattern,
			Headers: map[string]string{"content-type": pattern.MatchParameters.JSON},
			StatusCode: 200,
		},
	})

	content.Domains.Lobby = updatedLobbyDomain
	mitm.WriteClientMITMConfig(content)
}
