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

package errors

import (
	"errors"
	"fmt"
	"log"
)

type MITMError struct {
	Message string
	Err     error
}

var (
	// Proxy
	ErrProxy = errors.New("proxy error")
	ErrProxyServerException = errors.New("proxy server exception")
	ErrProxyToggleInvalidCommand = errors.New("invalid proxy toggle command")
	ErrProxyCertificateException = errors.New("proxy certificate exception")
	// HTTP
	ErrHTTPCodeError = errors.New("http code error")
	ErrHTTPRequestException = errors.New("http request exception")
	ErrHTTPBodyReadException = errors.New("http body read exception")
	// JSON/YAML
	ErrJSONReadException = errors.New("json read exception")
	ErrYAMLReadException = errors.New("yaml read exception")
	// Application
	ErrFatalException = errors.New("fatal exception")
	ErrUpdaterException = errors.New("updater exception")
	ErrPromptException = errors.New("prompt exception")
	ErrWatcherException = errors.New("watcher exception")
	ErrIOReadException = errors.New("io read exception")
	ErrMITMYamlSchemaOutdatedException = errors.New("mitm yaml schema outdated exception")
	ErrPingFailedException = errors.New("ping failed exception")
)

func (e *MITMError) Error() string {
	return fmt.Sprintf("message: %s; original error: %v", e.Message, e.Err)
}

func (e *MITMError) Log() {
	Log(e)
}

func (e *MITMError) String() string {
	return String(e)
}

func (e *MITMError) Unwrap() error {
	return e.Err
}

func String(e *MITMError) string {
	return fmt.Sprintf("[error] message: %s; original error: %v", e.Message, e.Err)
}

func Log(e *MITMError) {
	log.Println(String(e))
}

func Create(err error, message string) *MITMError {
	if message == "" {
		message = err.Error()
	}

	return &MITMError{
		Message: message,
		Err:     err,
	}
}
