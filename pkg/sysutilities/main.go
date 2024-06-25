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
	"bytes"
	"embed"
	"fmt"
	"infinite-mitm/pkg/errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/atotto/clipboard"
	"github.com/ncruces/zenity"
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
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", url).Start()
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	}
}

func CopyToClipboard(data string) {
	clipboard.WriteAll(data)
}

func SaveToDisk(data []byte, outputDir string, filename string, contentType string) {
	defer func() {
		_ = recover();
	}()

	contentType = strings.Split(contentType, ";")[0]
	var mimeExtensions = map[string]string{
		"application/json":         "json",
		"application/xml":          "xml",
		"text/html":                "html",
		"text/plain":               "txt",
		"image/jpeg":               "jpeg",
		"image/jpg":                "jpg",
		"image/png":                "png",
		"image/gif":                "gif",
		"application/octet-stream": "bin",
	}

	extension := mimeExtensions[contentType]
	if extension == "" {
		extension = "bin"
	}

	filename = fmt.Sprintf("%s.%s", filename, extension)
	filePath, err := zenity.SelectFileSave(
		zenity.Title("Save body content"),
		zenity.Filename(filepath.Join(outputDir, filename)),
		zenity.ConfirmOverwrite(),
	)

	if err == nil {
		os.WriteFile(filePath, data, 0644)
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

func CheckForRootCertificate(certName string) (bool, *errors.MITMError) {
	var installed bool
	var mitmErr *errors.MITMError

	switch runtime.GOOS {
	case "windows":
		installed, mitmErr = checkForRootCertificateOnWindows(certName)
	case "darwin":
		installed, mitmErr = checkForRootCertificateOnDarwin(certName)
	default:
		installed, mitmErr = false, nil
	}

	if mitmErr != nil {
		return false, mitmErr
	}

	return installed, nil
}

func IsAdmin() bool {
	return isAdmin()
}

func RunAsAdmin() {
	runAsAdmin()
}

func InstallRootCertificate(f *embed.FS, certFilename string) {
	installRootCertificate(f, certFilename)
}

func checkForRootCertificateOnWindows(certName string) (bool, *errors.MITMError) {
	cmd := exec.Command("certutil", "-verifystore", "-user", "root", certName)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, errors.Create(errors.ErrProxyCertificateException, err.Error())
	}

	if strings.TrimSpace(out.String()) == "" {
		return false, errors.Create(errors.ErrProxyCertificateException, "missing root certificate")
	}

	return true, nil
}

func checkForRootCertificateOnDarwin(certName string) (bool, *errors.MITMError) {
	cmd := exec.Command("security", "find-certificate", "-a", "-c", certName)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, errors.Create(errors.ErrProxyCertificateException, err.Error())
	}

	if strings.TrimSpace(out.String()) == "" {
		return false, errors.Create(errors.ErrProxyCertificateException, "missing root certificate")
	}

	return true, nil
}
