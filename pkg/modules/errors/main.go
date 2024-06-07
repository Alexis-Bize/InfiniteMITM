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
	ErrProxyToggleInvalidCommand = errors.New("invalid proxy toggle command")
	// HTTP
	ErrHTTPForbidden = errors.New("forbidden")
	ErrHTTPUnauthorized = errors.New("unauthorized")
	ErrHTTPNotFound = errors.New("not found")
	ErrHTTPError = errors.New("http error")
	ErrHTTPRequestException = errors.New("http request exception")
	ErrHTTPBodyReadException = errors.New("http body read exception")
	// JSON
	ErrJSONUnmarshalException = errors.New("json unmarshal exception")
	// Application
	ErrAppException = errors.New("application exception")
	ErrInvalidYAMLException = errors.New("invalid yaml exception")
	ErrFatalException = errors.New("fatal exception")
	ErrWatcherException = errors.New("watcher exception")
	ErrIOReadException = errors.New("io read exception")
	ErrRootCertificateException = errors.New("root certificate exception")
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
