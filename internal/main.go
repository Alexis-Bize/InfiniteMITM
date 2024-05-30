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
	"infinite-mitm/configs"
	mitm "infinite-mitm/internal/application/services/mitm"
	prompt "infinite-mitm/internal/application/services/prompt"
	signal "infinite-mitm/internal/application/services/signal"
	networkTable "infinite-mitm/internal/application/services/ui/network"
	errors "infinite-mitm/pkg/modules/errors"
	proxy "infinite-mitm/pkg/modules/proxy"
	"io"
	"os"
	"path/filepath"
)

func Start(f *embed.FS) error {
	createRootAssets(f)

	option, err := prompt.Welcome()
	if err != nil {
		return err
	}

	if prompt.StartProxyServer.Is(option) {
		server, err := mitm.InitializeServer(f)
		if err != nil {
			return err
		}

		if err = proxy.ToggleProxy("on"); err != nil {
			return err
		}

		go func() {
			program := networkTable.Create()
			if _, err := program.Run(); err != nil {
				errors.Log(errors.ErrFatalException, err.Error())
				signal.Stop()
			}
		}()

		if err := server.ListenAndServe(); err != nil {
			errors.Log(errors.ErrFatalException, err.Error())
			signal.Stop()
		}

		return nil
	}

	if prompt.ForceKillProxy.Is(option) || prompt.Exit.Is(option) {
		signal.Stop()
		return nil
	}

	return nil
}

func createRootAssets(f *embed.FS) {
	projectDir := configs.GetConfig().Extra.ProjectDir
	sourceMITMTemplate := "assets/resource/shared/templates/mitm.yaml"
	outputMITMTemplate := filepath.Join(projectDir, filepath.Base(sourceMITMTemplate))

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		errors.Log(errors.ErrFatalException, err.Error())
		os.Exit(1)
	}

	if _, err := os.Stat(outputMITMTemplate); os.IsNotExist(err) {
		templateFile, err := f.Open(sourceMITMTemplate)
		if err != nil {
			errors.Log(errors.ErrFatalException, err.Error())
			os.Exit(1)
		}
		defer templateFile.Close()

		destinationFile, err := os.Create(outputMITMTemplate)
		if err != nil {
			errors.Log(errors.ErrFatalException, err.Error())
			os.Exit(1)
		}
		defer destinationFile.Close()

		if _, err = io.Copy(destinationFile, templateFile); err != nil {
			errors.Log(errors.ErrFatalException, err.Error())
			os.Exit(1)
		}

		if err = destinationFile.Sync(); err != nil {
			errors.Log(errors.ErrFatalException, err.Error())
			os.Exit(1)
		}
	}
}