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

package InfiniteMITMApplicationService

import (
	"crypto/tls"
	"crypto/x509"
	"embed"
	"fmt"
	"infinite-mitm/configs"
	InfiniteMITMApplicationServiceModule "infinite-mitm/internal/application/services/mitm/modules"
	"log"
	"net/http"

	"gopkg.in/elazarl/goproxy.v1"
)

func StartServer(f embed.FS, userHandlers []InfiniteMITMApplicationServiceModule.HandlerStruct) {
	CACert, err := f.ReadFile("cert/rootCA.pem")
	if err != nil {
		log.Fatalf("failed to read root CA certificate: %v", err)
	}

	CAKey, err := f.ReadFile("cert/rootCA.key")
	if err != nil {
		log.Fatalf("failed to read root CA key: %v", err)
	}

	cert, err := tls.X509KeyPair(CACert, CAKey)
	if err != nil {
		log.Fatalf("failed to parse root CA certificate: %v", err)
	}

	CACertPool := x509.NewCertPool()
	if !CACertPool.AppendCertsFromPEM(CACert) {
		log.Fatalf("failed to add root CA certificate to pool")
	}

	goproxy.GoproxyCa = cert
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	for _, handler := range userHandlers {
		proxy.OnRequest(handler.Match).DoFunc(handler.Fn)
	}

	for _, handler := range internalHandlers() {
		proxy.OnRequest(handler.Match).DoFunc(handler.Fn)
	}

	server := &http.Server{
		Addr:    ":8888",
		Handler: proxy,
		TLSConfig: &tls.Config{
			RootCAs:            CACertPool,
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		},
	}

	fmt.Printf("%s - Starting proxy server on port 8888\n", configs.GetConfig().Name)
	log.Fatal(server.ListenAndServe())
}

func internalHandlers() []InfiniteMITMApplicationServiceModule.HandlerStruct {
	handlers := []InfiniteMITMApplicationServiceModule.HandlerStruct{
		InfiniteMITMApplicationServiceModule.PatchBookmarkedFilms(),
	}

	return handlers
}
