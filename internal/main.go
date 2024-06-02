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
	networkTable "infinite-mitm/internal/application/services/ui/network"
	errors "infinite-mitm/pkg/modules/errors"
	proxy "infinite-mitm/pkg/modules/proxy"
	utilities "infinite-mitm/pkg/modules/utilities"

	"github.com/gookit/event"
)

var server *http.Server

func Start(f *embed.FS) error {
	kill.Init()

	if err := application.CreateRootAssets(f); err != nil {
		return err
	}

	installed, err := application.CheckForRootCertificate()
	if err != nil {
		return err
	}

	option, err := prompt.Welcome(installed)
	if err != nil {
		return err
	}

	if prompt.Start.Is(option) {
		var wg sync.WaitGroup
		wg.Add(4)
		
		go func() {
			defer wg.Done()
			createNetworkView()
		}()

		go func() {
			defer wg.Done()
			startServer(f, false, &wg)
		}()

		go func() {
			defer wg.Done()
			enableProxy()
		}()

		go func() {
			defer wg.Done()
			mitm.WatchClientMITMConfig()
		}()

		wg.Wait()
	} else if prompt.InstallRootCertificate.Is(option) {
		utilities.OpenBrowser(configs.GetConfig().Repository + "/blob/main/docs/Install-Root-Certificate.md")
	} else if prompt.ForceKillProxy.Is(option) || prompt.Exit.Is(option) {
		proxy.ToggleProxy("off")
	}

	return nil
}

func enableProxy() {
	if err := proxy.ToggleProxy("on"); err != nil {
		errors.Create(errors.ErrFatalException, err.Error()).Log()
		return
	}
	
	kill.Register(func() {
		proxy.ToggleProxy("off")
	})
}

func startServer(f *embed.FS, isRestart bool, wg *sync.WaitGroup) {
	s, err := mitm.InitializeServer(f)
	if err != nil {
		errors.Create(errors.ErrFatalException, err.Error()).Log()
		return
	}

	server = s

	if !isRestart {
		event.On(events.RestartServer, event.ListenerFunc(func(e event.Event) error {
			return restartServer(f, wg)
		}))
	}

	if isRestart {
		wg.Done()
	}

	if err := server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			errors.Create(errors.ErrFatalException, err.Error()).Log()
		}
	}
}

func createNetworkView() {
	program := networkTable.Create()
	if _, err := program.Run(); err != nil {
		errors.Create(errors.ErrFatalException, err.Error()).Log()
	}
}

func restartServer(f *embed.FS, wg *sync.WaitGroup) error {
	if server != nil {
		if err := server.Close(); err != nil {
			return err
		}
	}

	wg.Add(1)
	go startServer(f, true, wg)
	return nil
}
