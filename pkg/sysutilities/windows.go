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

//go:build windows
// +build windows

package sysutilities

import (
	"infinite-mitm/pkg/errors"
	"log"
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/sys/windows"
)

func isAdmin() bool {
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid,
	)

	if err != nil {
		return false
	}

	defer windows.FreeSid(sid)

	token := windows.Token(0)
	member, err := token.IsMember(sid)
	return err == nil && member
}

func runAsAdmin() {
	verb := "runas"
	exe, err := os.Executable()
	if err != nil {
		log.Fatalln(errors.ErrFatalException, err.Error())
	}

	cmd := exec.Command(exe)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	var showCmd int32 = 1 // SW_SHOWNORMAL
	verbPtr, _ := windows.UTF16PtrFromString(verb)
	exePtr, _ := windows.UTF16PtrFromString(exe)
	dirPtr, _ := windows.UTF16PtrFromString("")
	argPtr, _ := windows.UTF16PtrFromString("")

	err = windows.ShellExecute(0, verbPtr, exePtr, argPtr, dirPtr, showCmd)
	if err != nil {
		log.Fatalln(errors.ErrFatalException, err.Error())
	}

	os.Exit(0)
}
