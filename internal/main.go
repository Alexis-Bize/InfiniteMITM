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

package InfiniteMITM

import (
	"embed"
	mitm "infinite-mitm/internal/application/services/mitm"
	prompt "infinite-mitm/internal/application/services/prompt"
	networkTable "infinite-mitm/internal/application/services/ui/network"
	errors "infinite-mitm/pkg/modules/errors"
	proxy "infinite-mitm/pkg/modules/proxy"
	"os"
)

func Start(f embed.FS) error {
	var err error

	option, err := prompt.Welcome()
	if err != nil {
		os.Exit(0)
	}

	if prompt.StartProxyServer.Is(option) {
		server, err := mitm.InitializeServer(f)
		if err != nil {
			return err
		}

		err = proxy.ToggleProxy("on")
		if err != nil {
			return err
		}

		go func() {
			program := networkTable.Create()
			if _, err := program.Run(); err != nil {
				errors.Log(errors.ErrFatalException, err.Error())
				os.Exit(1)
			}
		}()

		if err := server.ListenAndServe(); err != nil {
			errors.Log(errors.ErrFatalException, err.Error())
			os.Exit(1)
		}
	}

	if prompt.ForceKillProxy.Is(option) {
		proxy.ToggleProxy("off")
	}

	return nil
}
