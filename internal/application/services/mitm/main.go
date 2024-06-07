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
	context "infinite-mitm/internal/application/services/mitm/helpers/context"
	errors "infinite-mitm/pkg/modules/errors"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/google/uuid"
	"github.com/gookit/event"
)

type emptyLogger struct{}

var certName = configs.GetConfig().Proxy.Certificate.Name

func CreateServer(f *embed.FS) (*http.Server, *errors.MITMError) {
	CACert, err := f.ReadFile(fmt.Sprintf("cert/%s.pem", certName))
	if err != nil {
		return nil, errors.Create(errors.ErrRootCertificateException, err.Error())
	}

	CAKey, err := f.ReadFile(fmt.Sprintf("cert/%s.key", certName))
	if err != nil {
		return nil, errors.Create(errors.ErrRootCertificateException, err.Error())
	}

	cert, err := tls.X509KeyPair(CACert, CAKey)
	if err != nil {
		return nil, errors.Create(errors.ErrRootCertificateException, err.Error())
	}

	CACertPool := x509.NewCertPool()
	if !CACertPool.AppendCertsFromPEM(CACert) {
		return nil, errors.Create(errors.ErrRootCertificateException, "failed to add root CA certificate to pool")
	}

	goproxy.GoproxyCa = cert
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false
	proxy.Logger = emptyLogger{}

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest()

	content, mitmErr := ReadClientMITMConfig()
	if mitmErr != nil {
		return nil, mitmErr
	}

	clientRequestHandlers, clientResponseHandlers := CreateClientMITMHandlers(content)
	clientRequestHandlersCount := len(clientRequestHandlers)
	clientResponseHandlersCount := len(clientResponseHandlers)
	totalClientHandlersCount := clientRequestHandlersCount + clientResponseHandlersCount

	event.MustFire(events.ProxyStatusMessage, event.M{
		"details": fmt.Sprintf("[%s] found %d override(s); %d request(s) and %d response(s)", YAMLFilename, totalClientHandlersCount, clientRequestHandlersCount, clientResponseHandlersCount),
	})

	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request,*http.Response) {
		customCtx := context.ContextHandler(ctx)
		customCtx.SetUserData("uuid", uuid.New().String())
		customCtx.SetUserData("proxified", map[string]bool{"req": false, "resp": false})
		return req, nil
	})

	for _, handler := range clientRequestHandlers {
		proxy.OnRequest(handler.Match).DoFunc(handler.Fn)
	}

	for _, handler := range internalRequestHandlers() {
		proxy.OnRequest(handler.Match).DoFunc(handler.Fn)
	}

	for _, handler := range clientResponseHandlers {
		proxy.OnResponse(handler.Match).DoFunc(handler.Fn)
	}

	for _, handler := range internalResponseHandlers() {
		proxy.OnResponse(handler.Match).DoFunc(handler.Fn)
	}

	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		return resp
	})

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

func internalRequestHandlers() []handlers.RequestHandlerStruct {
	handlersList := []handlers.RequestHandlerStruct{
	}

	handlersList = append(handlersList, handlers.HandleRootRequests())
	return handlersList
}

func internalResponseHandlers() []handlers.ResponseHandlerStruct {
	handlersList := []handlers.ResponseHandlerStruct{
	}

	handlersList = append(handlersList, handlers.HandleRootResponses())
	return handlersList
}

func (l emptyLogger) Printf(format string, v ...interface{}) {
	// Ignore goproxy logs
}
