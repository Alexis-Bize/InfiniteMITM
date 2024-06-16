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

package sysutilities

import (
	"infinite-mitm/pkg/errors"
	"log"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/charmbracelet/huh/spinner"
	"golang.org/x/term"
)

func GetHomeDirectory() (string, *errors.MITMError) {
	home, err := os.UserHomeDir(); if err != nil {
		log.Fatalln(errors.Create(errors.ErrFatalException, err.Error()))
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

func IsAdmin() bool {
	return isAdmin()
}

func RunAsAdmin() {
	runAsAdmin()
}
