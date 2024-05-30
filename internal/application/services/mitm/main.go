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

package InfiniteMITMApplicationMITMService

import (
	"crypto/tls"
	"crypto/x509"
	"embed"
	"fmt"
	"infinite-mitm/configs"
	handlers "infinite-mitm/internal/application/services/mitm/handlers"
	errors "infinite-mitm/pkg/modules/errors"
	"net/http"

	"github.com/elazarl/goproxy"
)

type CustomLogger struct{}

func (l CustomLogger) Printf(format string, v ...interface{}) {
	// Ignore goproxy logs
}

const overrideLogger = true

func InitializeServer(f embed.FS) (*http.Server, error) {
	CACert, err := f.ReadFile("cert/rootCA.pem")
	if err != nil {
		return nil, errors.Log(errors.ErrRootCertificateException, err.Error())
	}

	CAKey, err := f.ReadFile("cert/rootCA.key")
	if err != nil {
		return nil, errors.Log(errors.ErrRootCertificateException, err.Error())
	}

	cert, err := tls.X509KeyPair(CACert, CAKey)
	if err != nil {
		return nil, errors.Log(errors.ErrRootCertificateException, err.Error())
	}

	CACertPool := x509.NewCertPool()
	if !CACertPool.AppendCertsFromPEM(CACert) {
		return nil, errors.Log(errors.ErrRootCertificateException, "failed to add root CA certificate to pool")
	}

	goproxy.GoproxyCa = cert
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false

	if overrideLogger {
		proxy.Logger = CustomLogger{}
	}

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	for _, handler := range internalRequestHandlers() {
		proxy.OnRequest(handler.Match).DoFunc(handler.Fn)
	}

	for _, handler := range internalResponseHandlers() {
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

func internalRequestHandlers() []handlers.RequestHandlerStruct {
	handlersList := []handlers.RequestHandlerStruct{
	}

	handlersList = append(handlersList, handlers.HandleHaloWaypointRequests())
	return handlersList
}

func internalResponseHandlers() []handlers.ResponseHandlerStruct {
	handlersList := []handlers.ResponseHandlerStruct{
	}

	handlersList = append(handlersList, handlers.HandleHaloWaypointResponses())
	return handlersList
}
