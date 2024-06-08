package ProxyModule

import (
	"fmt"
	"infinite-mitm/configs"
	errors "infinite-mitm/pkg/modules/errors"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var proxyHost = configs.GetConfig().Proxy.Host
var proxyPort = strconv.Itoa(configs.GetConfig().Proxy.Port)

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
	commands := []string{
		fmt.Sprintf("networksetup -setwebproxy Wi-Fi %s %s", proxyHost, proxyPort),
		fmt.Sprintf("networksetup -setsecurewebproxy Wi-Fi %s %s", proxyHost, proxyPort),
		"networksetup -setwebproxystate Wi-Fi on",
		"networksetup -setsecurewebproxystate Wi-Fi on",
	}

	for _, cmd := range commands {
		parts := strings.Fields(cmd)
		if err := exec.Command(parts[0], parts[1:]...).Run(); err != nil {
			return err
		}
	}

	return nil
}

func disableProxyWindows() error {
	return exec.Command("netsh", "winhttp", "reset", "proxy").Run()
}

func disableProxyDarwin() error {
	commands := []string{
		"networksetup -setwebproxystate Wi-Fi off",
		"networksetup -setsecurewebproxystate Wi-Fi off",
	}

	for _, cmd := range commands {
		parts := strings.Fields(cmd)
		if err := exec.Command(parts[0], parts[1:]...).Run(); err != nil {
			return err
		}
	}

	return nil
}
