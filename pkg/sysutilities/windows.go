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
	"embed"
	"fmt"
	"infinite-mitm/pkg/errors"
	"log"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

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

func installRootCertificate(f *embed.FS, certFilename string) {
	certData, err := f.ReadFile(fmt.Sprintf("cert/%s", certFilename))
	if err != nil {
		log.Fatalln(errors.ErrFatalException, err.Error())
	}

	rootStore, err := windows.UTF16PtrFromString("ROOT")
	if err != nil {
		log.Fatalln(errors.ErrFatalException, err.Error())
	}

	certStore, err := windows.CertOpenStore(
		windows.CERT_STORE_PROV_SYSTEM, 0, 0,
		windows.CERT_SYSTEM_STORE_CURRENT_USER,
		uintptr(unsafe.Pointer(rootStore)),
	)

	if err != nil {
		log.Fatalln(errors.ErrFatalException, err.Error())
	}

	defer windows.CertCloseStore(certStore, 0)
	cert, err := windows.CertCreateCertificateContext(
		windows.X509_ASN_ENCODING|windows.PKCS_7_ASN_ENCODING,
		&certData[0],
		uint32(len(certData)),
	)

	if err != nil {
		log.Fatalln(errors.ErrFatalException, err.Error())
	}

	defer windows.CertFreeCertificateContext(cert)
	err = windows.CertAddCertificateContextToStore(certStore, cert, windows.CERT_STORE_ADD_ALWAYS, nil)
	if err != nil {
		log.Fatalln(errors.ErrFatalException, err.Error())
	}
}
