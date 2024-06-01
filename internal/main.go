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
	application "infinite-mitm/internal/application"
	mitm "infinite-mitm/internal/application/services/mitm"
	prompt "infinite-mitm/internal/application/services/prompt"
	signal "infinite-mitm/internal/application/services/signal"
	networkTable "infinite-mitm/internal/application/services/ui/network"
	errors "infinite-mitm/pkg/modules/errors"
	proxy "infinite-mitm/pkg/modules/proxy"
)

func Start(f *embed.FS) error {
	if err := application.CreateRootAssets(f); err != nil {
		return err
	}

	option, err := prompt.Welcome()
	if err != nil {
		return err
	}

	if prompt.Start.Is(option) {
		startServer(f)
	} else if prompt.ForceKillProxy.Is(option) || prompt.Exit.Is(option) {
		signal.Stop()
	}

	return nil
}

func startProxy() {
	if err := proxy.ToggleProxy("on"); err != nil {
		errors.Create(errors.ErrFatalException, err.Error()).Log()
		signal.Stop()
	}
}

func createNetworkView() {
	program := networkTable.Create()
	if _, err := program.Run(); err != nil {
		errors.Create(errors.ErrFatalException, err.Error()).Log()
		signal.Stop()
	}
}

func startServer(f *embed.FS) error {
	server, err := mitm.InitializeServer(f)
	if err != nil {
		return err
	}

	go startProxy()
	go createNetworkView()

	if err := server.ListenAndServe(); err != nil {
		errors.Create(errors.ErrFatalException, err.Error()).Log()
		signal.Stop()
	}

	return nil
}