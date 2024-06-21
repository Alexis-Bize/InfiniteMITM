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
	events "infinite-mitm/internal/application/services/events"
	handlers "infinite-mitm/internal/application/services/mitm/handlers"
	helpers "infinite-mitm/internal/application/services/mitm/helpers"
	context "infinite-mitm/internal/application/services/mitm/modules/context"
	traffic "infinite-mitm/internal/application/services/mitm/modules/traffic"
	"infinite-mitm/pkg/domains"
	"infinite-mitm/pkg/errors"
	"infinite-mitm/pkg/pattern"
	"infinite-mitm/pkg/request"
	"infinite-mitm/pkg/smartcache"
	"infinite-mitm/pkg/utilities"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/fsnotify/fsnotify"
	"github.com/gookit/event"
	"gopkg.in/yaml.v2"
)

type YAMLOptions struct {
	SmartCache     smartcache.SmartCacheYAMLOptions `yaml:"smart_cache"`
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
		if v.Request.Body != "" || len(v.Request.Headers) != 0 || len(v.Request.Before.Commands) != 0 {
			clientRequestHandlers = append(clientRequestHandlers, createRequestHandler(domain, v))
		}

		if v.Response.Body != "" || len(v.Response.Headers) != 0 || len(v.Response.Before.Commands) != 0 {
			clientResponseHandlers = append(clientResponseHandlers, createResponseHandler(domain, v))
		}
	}

	return clientRequestHandlers, clientResponseHandlers
}

func createRequestHandler(domain domains.DomainType, node domains.YAMLDomainNode) handlers.RequestHandlerStruct {
	target := pattern.Create(domain, node.Path)
	return handlers.RequestHandlerStruct{
		Match: helpers.MatchRequestUrl(target),
		Fn: func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			if !utilities.Contains(node.Methods, req.Method) {
				return req, nil
			}

			body := node.Request.Body
			matches := pattern.Match(target, request.StripPort(req.URL.String()))
			beforeHandlers := node.Request.Before

			kv := utilities.InterfaceToMap(node.Request.Headers)
			for key, value := range kv {
				req.Header.Set(key, pattern.ReplaceParameters(pattern.ReplaceMatches(value, matches)))
			}

			beforeCommands := createCommands(beforeHandlers.Commands, matches)
			runCommands(beforeCommands)

			if body != "" {
				buffer, err := readBodyFile(body, matches, req.Header)
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
	target := pattern.Create(domain, node.Path)
	return handlers.ResponseHandlerStruct{
		Match: helpers.MatchResponseUrl(target),
		Fn: func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			if !utilities.Contains(node.Methods, resp.Request.Method) {
				return resp
			}

			body := node.Response.Body
			matches := pattern.Match(target, request.StripPort(resp.Request.URL.String()))
			beforeHandlers := node.Response.Before

			kv := utilities.InterfaceToMap(node.Response.Headers)
			for key, value := range kv {
				resp.Header.Set(key, pattern.ReplaceParameters(pattern.ReplaceMatches(value, matches)))
			}

			beforeCommands := createCommands(beforeHandlers.Commands, matches)
			runCommands(beforeCommands)

			if body != "" {
				buffer, err := readBodyFile(body, matches, resp.Request.Header)
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

func isUrl(str string) bool {
	_, err := url.ParseRequestURI(str); if err != nil {
		return false
	}

	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func createCommands(commands []domains.YAMLDomainTrafficRunCommand, matches []string) []string {
	var commandsList []string

	if len(commands) != 0 {
		for _, commands := range commands {
			var runList []string
			for _, run := range commands.Run {
				replace := pattern.ReplaceParameters(pattern.ReplaceMatches(run, matches))
				if !isUrl(replace) {
					replace = filepath.Clean(replace)
				}

				runList = append(runList, replace)
			}

			commandsList = append(commandsList, strings.Join(runList, " "))
		}
	}

	return commandsList
}

func runCommands(commands []string) {
	if len(commands) == 0 {
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for i, cmd := range commands {
			args := strings.Split(cmd, " ")
			length := len(args)
			if length == 0 {
				continue
			}

			var c *exec.Cmd
			if length == 1 {
				c = exec.Command(args[0])
			} else if length >= 2 {
				c = exec.Command(args[0], args[1:]...)
			}

			if err := c.Run(); err != nil {
				event.MustFire(events.ProxyStatusMessage, event.M{"details": fmt.Sprintf("command #%d encountered an error: %s", i + 1, err.Error())})
			}
		}
	}()

	wg.Wait()
}

func readBodyFile(body string, matches []string, header http.Header) ([]byte, error) {
	str := pattern.ReplaceParameters(pattern.ReplaceMatches(body, matches))

	if isUrl(str) {
		buffer, mitmErr := request.Send("GET", str, nil, header);
		if mitmErr != nil {
			return nil, fmt.Errorf(mitmErr.Message)
		}

		return buffer, nil
	}

	buffer, err := os.ReadFile(filepath.Clean(str));
	if err != nil {
		return nil, err
	}

	return buffer, nil
}
