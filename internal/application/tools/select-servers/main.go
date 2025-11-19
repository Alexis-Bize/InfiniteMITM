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
	"infinite-mitm/pkg/domains"
	"infinite-mitm/pkg/errors"
	"infinite-mitm/pkg/integrity"
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

const PingErrorValue = -1

var (
	once    sync.Once
	servers QOSServers
)

func GetQOSServers() QOSServers {
	once.Do(func() {
		key := resources.GetDirPaths()[resources.PubDirKey]
		serverFilepath := path.Join(key, integrity.QOSServersFilename)
		qosservers, err := os.ReadFile(serverFilepath); if err != nil {
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
		return nil, errors.Create(errors.ErrPingFailedException, err.Error())
	}

	pinger.Count = 3
	pinger.Timeout = time.Second * 5

	err = pinger.Run()
	if err != nil {
		return nil, errors.Create(errors.ErrPingFailedException, err.Error())
	}

	return pinger.Statistics(), nil
}

func ReadLocalQOSServers() QOSServers {
	serverFilepath := filepath.Join(resources.GetDirPaths()[resources.JsonDirKey], integrity.QOSServersFilename)
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
	serverFilepath := filepath.Join(resources.GetDirPaths()[resources.JsonDirKey], integrity.QOSServersFilename)
	buffer, err := json.Marshal(servers); if err == nil {
		os.WriteFile(serverFilepath, buffer, 0644)
	}
}

func WriteMITMFile() {
	filePathPattern := path.Join(pattern.MITMDirPrefix, resources.ResourcesDirKey, resources.JsonDirKey, integrity.QOSServersFilename)
	content, err := mitm.ReadClientMITMConfig(); if err != nil {
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
	mitm.WriteMITMFile(content)
}
