package proxy

import (
	"fmt"
	"infinite-mitm/configs"
	"infinite-mitm/pkg/modules/errors"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

var proxyHost = configs.GetConfig().Proxy.Host
var proxyPort = strconv.Itoa(configs.GetConfig().Proxy.Port)

func ToggleProxy(command string) error {
	err := toggle(command)
	if err != nil {
		return errors.Handle(errors.ErrProxy, err.Error())
	}

	return nil
}

func toggle(command string) error {
	if command != "on" && command != "off" {
		return errors.ErrProxyToggleInvalidCommand
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-c
		disableProxy()
	}()

	switch runtime.GOOS {
	case "windows":
		if command == "on" {
			return enableProxyWindows()
		} else {
			return disableProxyWindows()
		}
	case "darwin":
		if command == "on" {
			return enableProxyDarwin()
		} else {
			return disableProxyDarwin()
		}
	}

	return nil
}

func enableProxyWindows() error {
	return exec.Command("netsh", "winhttp", "set", "proxy", proxyHost+":"+proxyPort).Run()
}

func disableProxyWindows() error {
	return exec.Command("netsh", "winhttp", "reset", "proxy").Run()
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

func disableProxy() {
	switch runtime.GOOS {
	case "windows":
		disableProxyWindows()
	case "darwin":
		disableProxyDarwin()
	}
}
