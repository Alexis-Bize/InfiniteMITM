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

package proxy

import (
	"fmt"
	"infinite-mitm/configs"
	"infinite-mitm/pkg/errors"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var proxyHost = configs.GetConfig().Proxy.Host
var proxyPort = strconv.Itoa(configs.GetConfig().Proxy.Port)

const darwinFallbackService = "Wi-Fi"

var darwinActiveService = darwinFallbackService

func ToggleProxy(command string) *errors.MITMError {
	if command != "on" && command != "off" {
		return errors.Create(errors.ErrProxyToggleInvalidCommand, "invalid command")
	}

	if command == "on" {
		if err := enableProxy(); err != nil {
			return errors.Create(errors.ErrProxy, err.Error())
		}
	} else {
		if err := disableProxy(); err != nil {
			return errors.Create(errors.ErrProxy, err.Error())
		}
	}

	return nil
}

func enableProxy() error {
	switch runtime.GOOS {
	case "windows":
		return enableProxyWindows()
	case "darwin":
		return enableProxyDarwin()
	}

	return nil
}

func disableProxy() error {
	switch runtime.GOOS {
	case "windows":
		return disableProxyWindows()
	case "darwin":
		return disableProxyDarwin()
	}

	return nil
}

func enableProxyWindows() error {
	proxyArg := fmt.Sprintf("proxy-server=\"http=%s:%s;https=%s:%s\"", proxyHost, proxyPort, proxyHost, proxyPort)
	return exec.Command("netsh", "winhttp", "set", "proxy", proxyArg, "\"<-loopback>\"").Run()
}

func enableProxyDarwin() error {
	darwinActiveService = detectDarwinActiveService()
	service := darwinActiveService

	commands := [][]string{
		{"networksetup", "-setwebproxy", service, proxyHost, proxyPort},
		{"networksetup", "-setsecurewebproxy", service, proxyHost, proxyPort},
		{"networksetup", "-setwebproxystate", service, "on"},
		{"networksetup", "-setsecurewebproxystate", service, "on"},
	}

	for _, args := range commands {
		if err := exec.Command(args[0], args[1:]...).Run(); err != nil {
			return err
		}
	}

	return nil
}

func disableProxyWindows() error {
	return exec.Command("netsh", "winhttp", "reset", "proxy").Run()
}

func disableProxyDarwin() error {
	service := darwinActiveService

	commands := [][]string{
		{"networksetup", "-setwebproxystate", service, "off"},
		{"networksetup", "-setsecurewebproxystate", service, "off"},
	}

	for _, args := range commands {
		if err := exec.Command(args[0], args[1:]...).Run(); err != nil {
			return err
		}
	}

	return nil
}

func detectDarwinActiveService() string {
	iface := defaultDarwinInterface()
	if iface == "" {
		return darwinFallbackService
	}

	out, err := exec.Command("networksetup", "-listnetworkserviceorder").Output()
	if err != nil {
		return darwinFallbackService
	}

	devSuffix := "Device: " + iface + ")"
	lines := strings.Split(string(out), "\n")

	for i := 1; i < len(lines); i++ {
		portLine := strings.TrimSpace(lines[i])
		if !strings.HasPrefix(portLine, "(Hardware Port:") || !strings.HasSuffix(portLine, devSuffix) {
			continue
		}

		serviceLine := strings.TrimSpace(lines[i-1])
		if idx := strings.Index(serviceLine, ") "); idx > 0 {
			name := strings.TrimSpace(serviceLine[idx+2:])
			if name != "" {
				return name
			}
		}
	}

	return darwinFallbackService
}

func defaultDarwinInterface() string {
	out, err := exec.Command("route", "-n", "get", "default").Output()
	if err != nil {
		return ""
	}

	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "interface:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[len(parts)-1]
			}
		}
	}

	return ""
}
