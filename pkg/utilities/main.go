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

package utilities

import (
	"fmt"
	errors "infinite-mitm/pkg/errors"
	"log"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"syscall"

	"github.com/charmbracelet/huh/spinner"
	"golang.org/x/term"
)

func Contains(slice []string, item string) bool {
	defer func() {
		_ = recover();
	}()

	for _, v := range slice {
		if v == item {
			return true
		}
	}

	return false
}

func InterfaceToMap(i interface{}) map[string]string {
	defer func() {
		_ = recover();
	}()

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

func GetHomeDirectory() (string, *errors.MITMError) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
		return "", nil
	}

	return home, nil
}

func GetTerminalSize() (int, int) {
	const defaultWidth = 80
	const defaultHeight = 25

	fd := int(os.Stdin.Fd())
	width, height, err := term.GetSize(fd)
	if err != nil {
		return defaultWidth, defaultHeight
	}

	return width, height
}

func ClearTerminal() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func WrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	var sb strings.Builder
	for len(text) > width {
		cut := width
		for cut > 0 && text[cut] != ' ' && text[cut] != '\n' {
			cut--
		}

		if cut == 0 {
			cut = width
		}

		sb.WriteString(text[:cut] + "\n")
		text = text[cut:]
		for len(text) > 0 && (text[0] == ' ' || text[0] == '\n') {
			text = text[1:]
		}
	}

	sb.WriteString(text)
	return sb.String()
}

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

func KillProcess() {
	process, err := os.FindProcess(os.Getpid())
	if err == nil {
		process.Signal(syscall.SIGINT)
	}
}
