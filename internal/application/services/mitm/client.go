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
	context "infinite-mitm/internal/application/services/mitm/modules/context"
	traffic "infinite-mitm/internal/application/services/mitm/modules/traffic"
	domains "infinite-mitm/pkg/domains"
	errors "infinite-mitm/pkg/errors"
	pattern "infinite-mitm/pkg/pattern"
	request "infinite-mitm/pkg/request"
	cache "infinite-mitm/pkg/smartcache"
	utilities "infinite-mitm/pkg/utilities"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/fsnotify/fsnotify"
	"github.com/gookit/event"
	"gopkg.in/yaml.v2"
)

type SmartCacheYAMLOptions struct {
	Enabled  bool
	Strategy cache.StrategyType
}

type YAMLOptions struct {
	SmartCache     SmartCacheYAMLOptions `yaml:"smart_cache"`
	TrafficDisplay traffic.TrafficDisplay `yaml:"traffic_display"`
}

type YAML struct {
	Domains domains.YAMLDomains `yaml:"domains"`
	Options YAMLOptions `yaml:"options"`
	Version int `yaml:"version"`
}

const YAMLFilename = "mitm.yaml"
const YAMLVersion = 1

func WatchClientMITMConfig() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		event.MustFire(events.ProxyStatusMessage, event.M{
			"details": errors.Create(errors.ErrWatcherException, err.Error()).String(),
		})

		return
	}
	defer watcher.Close()

	filePath := filepath.Join(configs.GetConfig().Extra.ProjectDir, YAMLFilename)
	err = watcher.Add(filePath)
	if err != nil {
		event.MustFire(events.ProxyStatusMessage, event.M{
			"details": errors.Create(errors.ErrWatcherException, err.Error()).String(),
		})

		return
	}

	for {
		select {
		case watchEvent, ok := <-watcher.Events:
			if ok && watchEvent.Op&fsnotify.Write == fsnotify.Write {
				event.MustFire(events.ProxyStatusMessage, event.M{
					"details": fmt.Sprintf("[%s] changes detected; restarting proxy server...", YAMLFilename),
				})

				time.Sleep(time.Second * 1)
				event.MustFire(events.RestartServer, event.M{})
			}
		case err, ok := <-watcher.Errors:
			if ok && err != nil {
				event.MustFire(events.ProxyStatusMessage, event.M{
					"details": errors.Create(errors.ErrWatcherException, err.Error()).String(),
				})
			}
		}
	}
}

func ReadClientMITMConfig() (YAML, *errors.MITMError) {
	filePath := filepath.Join(configs.GetConfig().Extra.ProjectDir, YAMLFilename)
	yamlFile, err := os.ReadFile(filePath); if err != nil {
		log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
		return YAML{}, nil
	}

	var content YAML
	if err = yaml.Unmarshal(yamlFile, &content); err != nil {
		return YAML{}, errors.Create(errors.ErrYAMLReadException, err.Error())
	}

	if content.Version != YAMLVersion {
		return YAML{}, errors.Create(errors.ErrMITMYamlSchemaOutdatedException, fmt.Sprintf("your %s is outdated, all configs will be ignored", YAMLFilename))
	}

	return content, nil
}

func CreateClientMITMHandlers(yaml YAML) ([]handlers.RequestHandlerStruct, []handlers.ResponseHandlerStruct) {
	var clientRequestHandlers []handlers.RequestHandlerStruct
	var clientResponseHandlers []handlers.ResponseHandlerStruct

	for _, pair := range domains.GetYAMLContentDomainPairs(yaml.Domains) {
		reqHandlers, respHandlers := processNodes(pair.Content, pair.Domain)
		clientRequestHandlers = append(clientRequestHandlers, reqHandlers...)
		clientResponseHandlers = append(clientResponseHandlers, respHandlers...)
	}

	return clientRequestHandlers, clientResponseHandlers
}

func processNodes(contentList []domains.YAMLDomainNode, domain domains.DomainType) ([]handlers.RequestHandlerStruct, []handlers.ResponseHandlerStruct) {
	var clientRequestHandlers []handlers.RequestHandlerStruct
	var clientResponseHandlers []handlers.ResponseHandlerStruct
	for _, v := range contentList {
		if v.Request != (domains.YAMLDomainRequestNode{}) {
			clientRequestHandlers = append(clientRequestHandlers, createRequestHandler(domain, v))
		}

		if v.Response != (domains.YAMLDomainResponseNode{}) {
			clientResponseHandlers = append(clientResponseHandlers, createResponseHandler(domain, v))
		}
	}

	return clientRequestHandlers, clientResponseHandlers
}

func createPattern(domain domains.DomainType, path string) *regexp.Regexp {
	if path == "" || !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	re := pattern.Create(`(?i)` + regexp.QuoteMeta(domain) + path)
	return re
}

func createRequestHandler(domain domains.DomainType, node domains.YAMLDomainNode) handlers.RequestHandlerStruct {
	target := createPattern(domain, node.Path)
	return handlers.RequestHandlerStruct{
		Match: goproxy.UrlMatches(target),
		Fn: func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			if !utilities.Contains(node.Methods, req.Method) {
				return req, nil
			}

			body := node.Request.Body
			matches := pattern.Match(target, req.URL.String())

			if body != "" && req.Method != http.MethodGet && req.Method != http.MethodDelete {
				buffer, err := readBodyFile(body, req.Header, matches)

				if err != nil {
					mitmErr := errors.Create(errors.ErrIOReadException, fmt.Sprintf("invalid request body for %s; %s", body, err.Error()))
					event.MustFire(events.ProxyStatusMessage, event.M{"details": mitmErr.String()})
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

func createResponseHandler(domain domains.DomainType, node domains.YAMLDomainNode) handlers.ResponseHandlerStruct {
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
					mitmErr := errors.Create(errors.ErrIOReadException, fmt.Sprintf("invalid response body for %s; %s", body, err.Error()))
					event.MustFire(events.ProxyStatusMessage, event.M{"details": mitmErr.String()})
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
	str := pattern.ReplaceParameters(pattern.ReplaceMatches(body, matches))

	if isURL(str) {
		buffer, mitmErr := request.Send("GET", body, nil, request.HeadersToMap(headers));
		if mitmErr != nil {
			return nil, fmt.Errorf(mitmErr.Message)
		}

		return buffer, nil
	}

	buffer, err := os.ReadFile(cleanSystemPath(str));
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func isURL(str string) bool {
	_, err := url.ParseRequestURI(str); if err != nil {
		return false
	}

	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func cleanSystemPath(str string) string {
	if strings.HasPrefix(str, "~/") {
		home, err := utilities.GetHomeDirectory();
		if err != nil {
			return str
		}

		str = filepath.Join(home, str[2:])
	}

	return filepath.Clean(str)
}
