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

package Utilities

import (
	"fmt"
	errors "infinite-mitm/pkg/modules/errors"
	"os"
	"os/exec"
	"reflect"
	"runtime"

	"github.com/charmbracelet/huh/spinner"
)

func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}

	return false
}

func InterfaceToMap(i interface{}) map[string]string {
	v := reflect.ValueOf(i)
	output := make(map[string]string)

	if i == nil {
		return output
	}

	for _, key := range v.MapKeys() {
		ckey := fmt.Sprintf("%v", key.Interface())
		cvalue := fmt.Sprintf("%v", v.MapIndex(key).Interface())
		output[ckey] = cvalue
	}

	return output
}

func GetHomeDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Create(errors.ErrFatalException, err.Error())
	}

	return home, nil
}

// https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func OpenBrowser(url string) {
	spinner.New().Title("Attempting to open your browser...").Run()

	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", url).Start()
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	}
}