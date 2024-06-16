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

package MITM

import (
	"embed"
	"net/http"
	"sync"

	configs "infinite-mitm/configs"
	application "infinite-mitm/internal/application"
	events "infinite-mitm/internal/application/events"
	mitm "infinite-mitm/internal/application/services/mitm"
	prompt "infinite-mitm/internal/application/services/prompt"
	kill "infinite-mitm/internal/application/services/signal/kill"
	networkView "infinite-mitm/internal/application/services/ui/views/network"
	"infinite-mitm/pkg/errors"
	proxy "infinite-mitm/pkg/proxy"
	"infinite-mitm/pkg/sysutilities"

	"github.com/gookit/event"
)

var mu sync.Mutex
var server *http.Server

func Start(f *embed.FS) *errors.MITMError {
	var mitmErr *errors.MITMError
	kill.Init()

	certInstalled := true
	if mitmErr := application.Init(f); mitmErr != nil {
		if errors.ErrProxyCertificateException != mitmErr.Unwrap() {
			return mitmErr
		}

		certInstalled = false
	}

	option, mitmErr := prompt.WelcomePrompt(certInstalled)
	if mitmErr != nil {
		return mitmErr
	}

	if prompt.Start.Is(option) {
		var wg sync.WaitGroup
		wg.Add(4)

		go func() {
			defer wg.Done()
			networkView.Create()
			sysutilities.KillProcess()
		}()

		go func() {
			defer wg.Done()
			enableProxy()
			startServer(f, false, &wg)
		}()

		go func() {
			defer wg.Done()
			mitm.WatchClientMITMConfig()
		}()

		wg.Wait()
	} else if prompt.InstallRootCertificate.Is(option) {
		sysutilities.OpenBrowser(configs.GetConfig().Repository + "/blob/main/docs/Install-Root-Certificate.md")
	} else if prompt.ForceKillProxy.Is(option) || prompt.Exit.Is(option) {
		if mitmErr := proxy.ToggleProxy("off"); mitmErr != nil {
			mitmErr.Log()
		}
	}

	sysutilities.KillProcess()
	return nil
}

func enableProxy() {
	kill.Register(func() {
		if mitmErr := proxy.ToggleProxy("off"); mitmErr != nil {
			mitmErr.Log()
		}
	})

	if mitmErr := proxy.ToggleProxy("on"); mitmErr != nil {
		mitmErr.Log()
	}
}

func startServer(f *embed.FS, isRestart bool, wg *sync.WaitGroup) {
	s, mitmErr := mitm.CreateServer(f)
	if mitmErr != nil {
		mitmErr.Log()
		return
	}

	server = s

	if !isRestart {
		event.On(events.RestartServer, event.ListenerFunc(func(e event.Event) error {
			restartServer(f, wg)
			return nil
		}))
	}

	if err := server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			errors.Create(errors.ErrProxyServerException, err.Error()).Log()
		}
	}
}

func shutdownServer() {
	if err := server.Close(); err != nil {
		errors.Create(errors.ErrProxyServerException, err.Error()).Log()
	}
}

func restartServer(f *embed.FS, wg *sync.WaitGroup) {
	mu.Lock()
	defer mu.Unlock()

	shutdownServer()
	wg.Add(1)

	go func() {
		defer wg.Done()
		startServer(f, true, wg)
	}()
}
