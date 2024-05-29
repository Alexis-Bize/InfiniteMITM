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
	"infinite-mitm/configs"
	InfiniteMITMApplicationMITMService "infinite-mitm/internal/application/services/mitm"
	ProxyModule "infinite-mitm/pkg/modules/proxy"
	"log"
)

func Start(f embed.FS) error {
	var err error

	server, err := InfiniteMITMApplicationMITMService.InitializeServer(f)
	if err != nil {
		return err
	}

	err = ProxyModule.ToggleProxy("on")
	if err != nil {
		return err
	}

	log.Printf("%s - Starting proxy server on port %d\n", configs.GetConfig().Name, configs.GetConfig().Proxy.Port)
	return server.ListenAndServe()
}
