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
	"bytes"
	"fmt"
	"infinite-mitm/configs"
	handlers "infinite-mitm/internal/application/services/mitm/handlers"
	pattern "infinite-mitm/internal/application/services/mitm/helpers/pattern"
	domains "infinite-mitm/internal/modules/domains"
	errors "infinite-mitm/pkg/modules/errors"
	utilities "infinite-mitm/pkg/modules/utilities"
	request "infinite-mitm/pkg/modules/utilities/request"
	"io"
	"net/http"
	"net/url"
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
	Root      []YAMLNode `yaml:"root,omitempty"`
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
		return clientRequestHandlers, clientResponseHandlers, errors.Create(errors.ErrFatalException, err.Error())
	}

	var content YAML
	err = yaml.Unmarshal(yamlFile, &content)
	if err != nil {
		return clientRequestHandlers, clientResponseHandlers, errors.Create(errors.ErrFatalException, err.Error())
	}

	for _, v := range content.Root {
		if v.Request != (YAMLRequestNode{}) {
			handler := createRequestHandler(domains.HaloWaypointSVCDomains.Root, v)
			clientRequestHandlers = append(clientRequestHandlers, handler)
		}

		if v.Response != (YAMLResponseNode{}) {
			handler := createResponseHandler(domains.HaloWaypointSVCDomains.Root, v)
			clientResponseHandlers = append(clientResponseHandlers, handler)
		}
	}

	for _, v := range content.Blobs {
		if v.Request != (YAMLRequestNode{}) {
			handler := createRequestHandler(domains.HaloWaypointSVCDomains.Blobs, v)
			clientRequestHandlers = append(clientRequestHandlers, handler)
		}

		if v.Response != (YAMLResponseNode{}) {
			handler := createResponseHandler(domains.HaloWaypointSVCDomains.Blobs, v)
			clientResponseHandlers = append(clientResponseHandlers, handler)
		}
	}

	for _, v := range content.Authoring {
		if v.Request != (YAMLRequestNode{}) {
			handler := createRequestHandler(domains.HaloWaypointSVCDomains.Authoring, v)
			clientRequestHandlers = append(clientRequestHandlers, handler)
		}

		if v.Response != (YAMLResponseNode{}) {
			handler := createResponseHandler(domains.HaloWaypointSVCDomains.Authoring, v)
			clientResponseHandlers = append(clientResponseHandlers, handler)
		}
	}

	for _, v := range content.Discovery {
		if v.Request != (YAMLRequestNode{}) {
			handler := createRequestHandler(domains.HaloWaypointSVCDomains.Discovery, v)
			clientRequestHandlers = append(clientRequestHandlers, handler)
		}

		if v.Response != (YAMLResponseNode{}) {
			handler := createResponseHandler(domains.HaloWaypointSVCDomains.Discovery, v)
			clientResponseHandlers = append(clientResponseHandlers, handler)
		}
	}

	for _, v := range content.HaloStats {
		if v.Request != (YAMLRequestNode{}) {
			handler := createRequestHandler(domains.HaloWaypointSVCDomains.HaloStats, v)
			clientRequestHandlers = append(clientRequestHandlers, handler)
		}

		if v.Response != (YAMLResponseNode{}) {
			handler := createResponseHandler(domains.HaloWaypointSVCDomains.HaloStats, v)
			clientResponseHandlers = append(clientResponseHandlers, handler)
		}
	}

	for _, v := range content.Settings {
		if v.Request != (YAMLRequestNode{}) {
			handler := createRequestHandler(domains.HaloWaypointSVCDomains.Settings, v)
			clientRequestHandlers = append(clientRequestHandlers, handler)
		}

		if v.Response != (YAMLResponseNode{}) {
			handler := createResponseHandler(domains.HaloWaypointSVCDomains.Settings, v)
			clientResponseHandlers = append(clientResponseHandlers, handler)
		}
	}

	for _, v := range content.GameCMS {
		if v.Request != (YAMLRequestNode{}) {
			handler := createRequestHandler(domains.HaloWaypointSVCDomains.GameCMS, v)
			clientRequestHandlers = append(clientRequestHandlers, handler)
		}

		if v.Response != (YAMLResponseNode{}) {
			handler := createResponseHandler(domains.HaloWaypointSVCDomains.GameCMS, v)
			clientResponseHandlers = append(clientResponseHandlers, handler)
		}
	}

	for _, v := range content.Economy {
		if v.Request != (YAMLRequestNode{}) {
			handler := createRequestHandler(domains.HaloWaypointSVCDomains.Economy, v)
			clientRequestHandlers = append(clientRequestHandlers, handler)
		}

		if v.Response != (YAMLResponseNode{}) {
			handler := createResponseHandler(domains.HaloWaypointSVCDomains.Economy, v)
			clientResponseHandlers = append(clientResponseHandlers, handler)
		}
	}

	return clientRequestHandlers, clientResponseHandlers, nil
}

func createPattern(domain domains.Domain, path string) *regexp.Regexp {
	if path == "" || !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	re := pattern.Create(`(?i)` + regexp.QuoteMeta(domain + path))
	return re
}

func createRequestHandler(domain domains.Domain, node YAMLNode) handlers.RequestHandlerStruct {
	target := createPattern(domain, node.Path)
	return handlers.RequestHandlerStruct{
		Match: goproxy.UrlMatches(target),
		Fn: func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			if !utilities.Contains(node.Methods, req.Method) {
				return req, nil
			}

			body := node.Request.Body
			matches := pattern.Match(target, req.URL.String())

			if body != "" && req.Method != "GET" && req.Method != "DELETE" {
				buffer, err := readBodyFile(body, req.Header, matches)

				if err != nil {
					errors.Create(errors.ErrIOReadException, "invalid request body, skipping").Log()
				} else if buffer != nil {
					bufferLength := len(buffer)
					req.Body = io.NopCloser(bytes.NewBuffer(buffer))
					req.ContentLength = int64(bufferLength)
					req.Header.Set("Content-Length", fmt.Sprintf("%d", bufferLength))
				}
			}

			
			kv := utilities.InterfaceToMap(node.Request.Headers)
			for key, value := range kv {
				fmt.Println(key)
				req.Header.Set(key, pattern.ReplaceParameters(pattern.ReplaceMatches(value, matches)))
			}

			return req, nil
		},
	}
}

func createResponseHandler(domain domains.Domain, node YAMLNode) handlers.ResponseHandlerStruct {
	target := createPattern(domain, node.Path)
	return handlers.ResponseHandlerStruct{
		Match: goproxy.UrlMatches(target),
		Fn: func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			if !utilities.Contains(node.Methods, resp.Request.Method) {
				return resp
			}

			body := node.Response.Body
			matches := pattern.Match(target, resp.Request.URL.String())

			if body != "" {
				buffer, err := readBodyFile(body, resp.Header, matches)
				if err != nil {
					errors.Create(errors.ErrIOReadException, "invalid response body, skipping").Log()
				} else if buffer != nil {
					bufferLength := len(buffer)
					resp.Body = io.NopCloser(bytes.NewBuffer(buffer))
					resp.ContentLength = int64(bufferLength)
					resp.Header.Set("Content-Length", fmt.Sprintf("%d", bufferLength))
				}
			}
			
			kv := utilities.InterfaceToMap(node.Response.Headers)
			for key, value := range kv {
				resp.Header.Set(key, pattern.ReplaceParameters(pattern.ReplaceMatches(value, matches)))
			}

			return resp
		},
	}
}

func readBodyFile(body string, headers http.Header, matches []string) ([]byte, error) {
	body = pattern.ReplaceParameters(pattern.ReplaceMatches(body, matches))

	var buffer []byte
	var err error

	if isURL(body) {
		buffer, err = request.Send("GET", body, nil, request.HeadersToMap(headers))
	} else if isSystemPath(body) {
		buffer, err = os.ReadFile(body)
	}


	return buffer, err
}

func isURL(str string) bool {
	_, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}

	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func isSystemPath(str string) bool {
	if strings.HasPrefix(str, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return false
		}

		str = filepath.Join(homeDir, str[2:])
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z]:\\|^\/`, str)
	return matched || filepath.IsAbs(str)
}

