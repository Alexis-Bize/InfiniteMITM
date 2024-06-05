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

package main

import (
	"embed"
	MITM "infinite-mitm/internal"
	MITMApplicationUIServiceTestTable "infinite-mitm/internal/application/services/ui/test"
	"log"
)

//go:generate goversioninfo -icon=assets/resources/windows/icon_256x256.ico
//go:embed assets/resources/shared/templates/*
//go:embed cert/InfiniteMITMRootCA.pem
//go:embed cert/InfiniteMITMRootCA.key
var f embed.FS

const debug = true

func main() {
	if debug {
		MITMApplicationUIServiceTestTable.Create()
		return
	}

	err := MITM.Start(&f)
	if err != nil {
		log.Println(err)
	}
}
