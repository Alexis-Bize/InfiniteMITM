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
	"fmt"
	"net/http"
	"runtime"
	"sync"

	"infinite-mitm/configs"
	mitmApplication "infinite-mitm/internal/application"
	embedFS "infinite-mitm/internal/application/embed"
	eventsService "infinite-mitm/internal/application/services/events"
	mitmService "infinite-mitm/internal/application/services/mitm"
	killService "infinite-mitm/internal/application/services/signal/kill"
	networkUI "infinite-mitm/internal/application/ui/network"
	welcomePromptUI "infinite-mitm/internal/application/ui/prompt/welcome"
	"infinite-mitm/pkg/errors"
	"infinite-mitm/pkg/proxy"
	"infinite-mitm/pkg/sysutilities"

	"github.com/gookit/event"
)

var (
	server *http.Server
	restartMutex sync.Mutex
	once sync.Once
)

func Start(f *embed.FS, omitInit bool) *errors.MITMError {
	var mitmErr *errors.MITMError

	certInstalled := true
	sysutilities.ClearTerminal()

	if !omitInit {
		embedFS.Set(f)

		if mitmErr := mitmApplication.Init(); mitmErr != nil {
			if errors.ErrProxyCertificateException != mitmErr.Unwrap() {
				return mitmErr
			}

			certInstalled = false
		}
	}

	option, mitmErr := welcomePromptUI.Run(certInstalled)
	if mitmErr != nil {
		return mitmErr
	}

	if welcomePromptUI.Start.Is(option) {
		var wg sync.WaitGroup
		var stopChan = make(chan struct{})

		wg.Add(3)

		go func() {
			defer wg.Done()
			networkUI.Create()
			stopServer()
			close(stopChan)
		}()

		go func() {
			defer wg.Done()
			mitmService.WatchClientMITMConfig(stopChan)
		}()

		go func() {
			defer wg.Done()

			once.Do(func() {
				killService.Register(func() {
					disableProxy()
				})

				event.On(eventsService.RestartServer, event.ListenerFunc(func(e event.Event) error {
					restartServer(f, &wg)
					return nil
				}))
			})

			startServer(f)
		}()

		wg.Wait()
		return Start(f, true)
	}

	if welcomePromptUI.InstallRootCertificate.Is(option) {
		if runtime.GOOS == "windows" {
			sysutilities.InstallRootCertificate(f, fmt.Sprintf("%s.cer", configs.GetConfig().Proxy.Certificate.Name))
			return Start(f, false)
		}

		sysutilities.OpenBrowser(configs.GetConfig().Repository + "/blob/main/cert/" + fmt.Sprintf("%s.pem", configs.GetConfig().Proxy.Certificate.Name))
	} else if welcomePromptUI.Exit.Is(option) {
		if mitmErr := proxy.ToggleProxy("off"); mitmErr != nil {
			mitmErr.Log()
		}
	}

	sysutilities.KillProcess()
	return nil
}

func startServer(f *embed.FS) {
	enableProxy()

	s, mitmErr := mitmService.CreateServer(f)
	if mitmErr != nil {
		event.MustFire(eventsService.ProxyStatusMessage, event.M{"details": mitmErr.String()})
		return
	}

	server = s
	if err := server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			event.MustFire(eventsService.ProxyStatusMessage, event.M{"details": err.Error()})
		}
	}
}

func restartServer(f *embed.FS, wg *sync.WaitGroup) {
	restartMutex.Lock()
	defer restartMutex.Unlock()

	wg.Add(1)
	stopServer()

	go func() {
		defer wg.Done()
		startServer(f)
	}()
}

func stopServer() {
	disableProxy()

	if err := server.Close(); err != nil {
		errors.Create(errors.ErrProxyServerException, err.Error()).Log()
	}
}

func enableProxy() {
	if mitmErr := proxy.ToggleProxy("on"); mitmErr != nil {
		mitmErr.Log()
	}
}

func disableProxy() {
	if mitmErr := proxy.ToggleProxy("off"); mitmErr != nil {
		mitmErr.Log()
	}
}
