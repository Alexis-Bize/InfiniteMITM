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
	ErrRootCertificateException = errors.New("root certificate exception")
	// HTTP
	ErrHTTPForbidden = errors.New("forbidden")
	ErrHTTPUnauthorized = errors.New("unauthorized")
	ErrHTTPNotFound = errors.New("not found")
	ErrHTTPInternal = errors.New("internal error")
	ErrHTTPRequestException = errors.New("http request exception")
	// JSON
	ErrJSONUnmarshalException = errors.New("json unmarshal exception")
	// Miscellaneous
	ErrIOReadException = errors.New("io read exception")
	ErrFatalException = errors.New("fatal exception")
)

func (e *MITMError) Error() string {
	return fmt.Sprintf("message: %s; original error: %v", e.Message, e.Err)
}

func (e *MITMError) Log() {
	log.Println(e.Err.Error())
}

func (e *MITMError) Unwrap() error {
	return e.Err
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
