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

package MITMApplicationMITMService

import (
	"infinite-mitm/configs"
	handlers "infinite-mitm/internal/application/services/mitm/handlers"
	pattern "infinite-mitm/internal/application/services/mitm/helpers/pattern"
	domains "infinite-mitm/internal/modules/domains"
	errors "infinite-mitm/pkg/modules/errors"
	utilities "infinite-mitm/pkg/modules/utilities"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/elazarl/goproxy"
	"gopkg.in/yaml.v2"
)

type YAMLRequestNode struct {
	Body    string `yaml:"body,omitempty"`
	Headers interface {} `yaml:"headers,omitempty"`
}

type YAMLResponseNode struct {
	Body    string `yaml:"body,omitempty"`
	Headers interface {} `yaml:"headers,omitempty"`
}

type YAMLNode struct {
	Path     string   `yaml:"path"`
	Methods  []string `yaml:"methods"`
	Request  YAMLRequestNode `yaml:"request,omitempty"`
	Response YAMLResponseNode `yaml:"response,omitempty"`
}

type YAML struct {
	Blobs     []YAMLNode `yaml:"blobs,omitempty"`
	Authoring []YAMLNode `yaml:"authoring,omitempty"`
	Discovery []YAMLNode `yaml:"discovery,omitempty"`
	HaloStats []YAMLNode `yaml:"stats,omitempty"`
	Settings  []YAMLNode `yaml:"settings,omitempty"`
	GameCMS   []YAMLNode `yaml:"gamecms,omitempty"`
	Economy   []YAMLNode `yaml:"economy,omitempty"`
}

func ReadClientMITMConfig() ([]handlers.RequestHandlerStruct, []handlers.ResponseHandlerStruct, error) {
	var clientRequestHandlers []handlers.RequestHandlerStruct
	var clientResponseHandlers []handlers.ResponseHandlerStruct

	filePath := filepath.Join(configs.GetConfig().Extra.ProjectDir, "mitm.yaml")
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return clientRequestHandlers, clientResponseHandlers, errors.Log(errors.ErrFatalException, err.Error())
	}

	var content YAML
	err = yaml.Unmarshal(yamlFile, &content)
	if err != nil {
		return clientRequestHandlers, clientResponseHandlers, errors.Log(errors.ErrFatalException, err.Error())
	}

	if len(content.Settings) > 0 {
		for _, v := range content.Settings {
			if v.Response != (YAMLResponseNode{}) {
				path := v.Path
				if path == "" || !strings.HasPrefix(path, "/") {
					continue
				}

				methods := v.Methods
				handler := func () handlers.ResponseHandlerStruct {
					target := pattern.Create(`(?i)` + regexp.QuoteMeta(domains.HaloWaypointSVCDomains.Settings + path))

					return handlers.ResponseHandlerStruct{
						Match: goproxy.UrlMatches(target),
						Fn: func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
							if !utilities.Contains(methods, resp.Request.Method) {
								return resp
							}

							matches := pattern.Match(target, resp.Request.URL.String())
							kv := utilities.InterfaceToMap(v.Response.Headers)
							for key, value := range kv {
								resp.Header.Set(key, pattern.ReplaceMatches(value, matches))
							}

							return resp
						},
					}
				}

				clientResponseHandlers = append(clientResponseHandlers, handler())
			}
		}
	}

	return clientRequestHandlers, clientResponseHandlers, nil
}
