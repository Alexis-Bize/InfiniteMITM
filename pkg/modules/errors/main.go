package ErrorsModule

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
	ErrPromptException = errors.New("prompt exception")
	ErrWatcherException = errors.New("watcher exception")
	ErrIOReadException = errors.New("io read exception")
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
