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
	"crypto/tls"
	"crypto/x509"
	"embed"
	"fmt"
	"infinite-mitm/configs"
	events "infinite-mitm/internal/application/events"
	handlers "infinite-mitm/internal/application/services/mitm/handlers"
	context "infinite-mitm/internal/application/services/mitm/modules/context"
	traffic "infinite-mitm/internal/application/services/mitm/modules/traffic"
	"infinite-mitm/pkg/domains"
	"infinite-mitm/pkg/errors"
	"infinite-mitm/pkg/smartcache"
	"net/http"
	"regexp"

	"github.com/elazarl/goproxy"
	"github.com/google/uuid"
	"github.com/gookit/event"
)

type emptyLogger struct{}

var certName = configs.GetConfig().Proxy.Certificate.Name
var smartCache *smartcache.SmartCache

func CreateServer(f *embed.FS) (*http.Server, *errors.MITMError) {
	var clientRequestHandlers []handlers.RequestHandlerStruct
	var clientResponseHandlers []handlers.ResponseHandlerStruct

	CACert, err := f.ReadFile(fmt.Sprintf("cert/%s.pem", certName)); if err != nil {
		return nil, errors.Create(errors.ErrProxyCertificateException, err.Error())
	}

	CAKey, err := f.ReadFile(fmt.Sprintf("cert/%s.key", certName)); if err != nil {
		return nil, errors.Create(errors.ErrProxyCertificateException, err.Error())
	}

	cert, err := tls.X509KeyPair(CACert, CAKey); if err != nil {
		return nil, errors.Create(errors.ErrProxyCertificateException, err.Error())
	}

	CACertPool := x509.NewCertPool(); if !CACertPool.AppendCertsFromPEM(CACert) {
		return nil, errors.Create(errors.ErrProxyCertificateException, "failed to add root CA certificate to pool")
	}

	goproxy.GoproxyCa = cert
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false
	proxy.Logger = emptyLogger{}

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	content, mitmErr := ReadClientMITMConfig(); if mitmErr != nil {
		event.MustFire(events.ProxyStatusMessage, event.M{"details": mitmErr.String()})
	}

	if mitmErr == nil {
		clientRequestHandlers, clientResponseHandlers = CreateClientMITMHandlers(content)
		clientActiveRequestHandlersCount := len(clientRequestHandlers)
		clientActiveResponseHandlersCount := len(clientResponseHandlers)
		totalClientHandlersCount := len(clientRequestHandlers) + len(clientResponseHandlers)

		domainText := "handler"
		if totalClientHandlersCount == 0 || totalClientHandlersCount > 1 {
			domainText += "s"
		}

		requestText := "request"
		if clientActiveRequestHandlersCount == 0 || clientActiveRequestHandlersCount > 1 {
			requestText += "s"
		}

		responseText := "response"
		if clientActiveResponseHandlersCount == 0 || clientActiveResponseHandlersCount > 1 {
			responseText += "s"
		}

		event.MustFire(events.ProxyStatusMessage, event.M{
			"details": fmt.Sprintf(
				"[%s] found %d %s; %d %s and %d %s",
				YAMLFilename,
				totalClientHandlersCount,
				domainText,
				clientActiveRequestHandlersCount,
				requestText,
				clientActiveResponseHandlersCount,
				responseText,
			),
		})
	}

	var rootCondition = goproxy.UrlMatches(regexp.MustCompile(`(?i)` + regexp.QuoteMeta(domains.HaloWaypointSVCDomains.Root)))
	smartCacheEnabled := content.Options.SmartCache.Enabled

	if !smartCacheEnabled {
		smartCache = nil
	} else if smartCache == nil {
		smartCache = smartcache.New(smartcache.StrategyType(content.Options.SmartCache.Strategy))
	}

	proxy.OnRequest(rootCondition).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		customCtx := context.ContextHandler(ctx)
		customCtx.SetUserData("uuid", uuid.New().String())
		customCtx.SetUserData("proxified", map[string]bool{"req": false, "resp": false})

		if smartCacheEnabled && smartcache.IsURLSmartCachable(req.URL.String(), req.Method) {
			customCtx.SetUserData("cache", smartCache)
		}

		return req, nil
	})

	for _, handler := range clientRequestHandlers {
		proxy.OnRequest(handler.Match).DoFunc(handler.Fn)
	}

	for _, handler := range internalRequestHandlers(content.Options) {
		proxy.OnRequest(handler.Match).DoFunc(handler.Fn)
	}

	for _, handler := range clientResponseHandlers {
		proxy.OnResponse(handler.Match).DoFunc(handler.Fn)
	}

	for _, handler := range internalResponseHandlers(content.Options) {
		proxy.OnResponse(handler.Match).DoFunc(handler.Fn)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", configs.GetConfig().Proxy.Port),
		Handler: proxy,
		TLSConfig: &tls.Config{
			RootCAs:            CACertPool,
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		},
	}

	return server, nil
}

func internalRequestHandlers(options YAMLOptions) []handlers.RequestHandlerStruct {
	handlersList := []handlers.RequestHandlerStruct{}
	handlersList = append(handlersList, handlers.HandleRootRequests(traffic.TrafficOptions{
		TrafficDisplay: options.TrafficDisplay,
	}))

	return handlersList
}

func internalResponseHandlers(options YAMLOptions) []handlers.ResponseHandlerStruct {
	handlersList := []handlers.ResponseHandlerStruct{}
	handlersList = append(handlersList, handlers.HandleRootResponses(traffic.TrafficOptions{
		TrafficDisplay: options.TrafficDisplay,
	}))

	return handlersList
}

func (l emptyLogger) Printf(format string, v ...interface{}) {
	// Ignore goproxy logs
}
