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
	events "infinite-mitm/internal/application/events"
	handlers "infinite-mitm/internal/application/services/mitm/handlers"
	context "infinite-mitm/internal/application/services/mitm/helpers/context"
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
	"github.com/fsnotify/fsnotify"
	"github.com/gookit/event"
	"gopkg.in/yaml.v2"
)

type YAMLRequestNode struct {
	Body    string `yaml:"body,omitempty"`
	Headers interface{} `yaml:"headers,omitempty"`
}

type YAMLResponseNode struct {
	Body       string `yaml:"body,omitempty"`
	Headers    interface{} `yaml:"headers,omitempty"`
	StatusCode int `yaml:"code,omitempty"`
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
	Lobby     []YAMLNode `yaml:"lobby,omitempty"`
}

type ContentDomainPair struct {
	content []YAMLNode
	domain  domains.Domain
}

const YAMLFilename = "mitm.yaml"

func WatchClientMITMConfig() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		errors.Create(errors.ErrWatcherException, err.Error()).Log()
		return
	}
	defer watcher.Close()

	filePath := filepath.Join(configs.GetConfig().Extra.ProjectDir, YAMLFilename)
	err = watcher.Add(filePath)
	if err != nil {
		errors.Create(errors.ErrWatcherException, err.Error()).Log()
		return
	}

	go func() {
		for {
			select {
			case watchEvent, ok := <-watcher.Events:
				if ok && watchEvent.Op&fsnotify.Write == fsnotify.Write {
					event.MustFire(events.RestartServer, map[string]any{})
				}
			case err, ok := <-watcher.Errors:
				if ok && err != nil {
					errors.Create(errors.ErrWatcherException, err.Error()).Log()
				}
			}
		}
	}()

	<-make(chan struct{})
}

func ReadClientMITMConfig() ([]handlers.RequestHandlerStruct, []handlers.ResponseHandlerStruct, error) {
	var clientRequestHandlers []handlers.RequestHandlerStruct
	var clientResponseHandlers []handlers.ResponseHandlerStruct

	filePath := filepath.Join(configs.GetConfig().Extra.ProjectDir, YAMLFilename)
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return clientRequestHandlers, clientResponseHandlers, errors.Create(errors.ErrFatalException, err.Error())
	}

	var content YAML
	err = yaml.Unmarshal(yamlFile, &content)
	if err != nil {
		return clientRequestHandlers, clientResponseHandlers, errors.Create(errors.ErrFatalException, err.Error())
	}

	contentDomainPairs := []ContentDomainPair{
		{content: content.Root, domain: domains.HaloWaypointSVCDomains.Root},
		{content: content.Blobs, domain: domains.HaloWaypointSVCDomains.Blobs},
		{content: content.Authoring, domain: domains.HaloWaypointSVCDomains.Authoring},
		{content: content.Discovery, domain: domains.HaloWaypointSVCDomains.Discovery},
		{content: content.HaloStats, domain: domains.HaloWaypointSVCDomains.HaloStats},
		{content: content.Settings, domain: domains.HaloWaypointSVCDomains.Settings},
		{content: content.GameCMS, domain: domains.HaloWaypointSVCDomains.GameCMS},
		{content: content.Economy, domain: domains.HaloWaypointSVCDomains.Economy},
		{content: content.Lobby, domain: domains.HaloWaypointSVCDomains.Lobby},
	}

	for _, pair := range contentDomainPairs {
		reqHandlers, respHandlers := processNodes(pair.content, pair.domain)
		clientRequestHandlers = append(clientRequestHandlers, reqHandlers...)
		clientResponseHandlers = append(clientResponseHandlers, respHandlers...)
	}

	return clientRequestHandlers, clientResponseHandlers, nil
}

func processNodes(contentList []YAMLNode, domain domains.Domain) ([]handlers.RequestHandlerStruct, []handlers.ResponseHandlerStruct) {
	var clientRequestHandlers []handlers.RequestHandlerStruct
	var clientResponseHandlers []handlers.ResponseHandlerStruct

	for _, v := range contentList {
		if v.Request != (YAMLRequestNode{}) {
			handler := createRequestHandler(domain, v)
			clientRequestHandlers = append(clientRequestHandlers, handler)
		}

		if v.Response != (YAMLResponseNode{}) {
			handler := createResponseHandler(domain, v)
			clientResponseHandlers = append(clientResponseHandlers, handler)
		}
	}

	return clientRequestHandlers, clientResponseHandlers
}

func createPattern(domain domains.Domain, path string) *regexp.Regexp {
	if path == "" || !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	re := pattern.Create(`(?i)` + regexp.QuoteMeta(domain) + path)
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
				req.Header.Set(key, pattern.ReplaceParameters(pattern.ReplaceMatches(value, matches)))
			}

			customCtx := context.ContextHandler(ctx)
			proxified := customCtx.GetUserData("proxified")

			if proxified != nil {
				pr := proxified.(map[string]bool)
				pr["req"] = true
				customCtx.SetUserData("proxified", pr)
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
					errors.Create(errors.ErrIOReadException, fmt.Sprintf("invalid response body for %s", body)).Log()
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

			if node.Response.StatusCode != 0 {
				resp.Status = http.StatusText(node.Response.StatusCode)
				resp.StatusCode = node.Response.StatusCode
			}

			customCtx := context.ContextHandler(ctx)
			proxified := customCtx.GetUserData("proxified")

			if proxified != nil {
				pr := proxified.(map[string]bool)
				pr["resp"] = true
				customCtx.SetUserData("proxified", pr)
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
		home, err := utilities.GetHomeDirectory()
		if err != nil {
			return false
		}

		str = filepath.Join(home, str[2:])
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z]:\\|^\/`, str)
	return matched || filepath.IsAbs(str)
}

