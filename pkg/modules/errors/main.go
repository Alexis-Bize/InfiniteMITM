package ErrorsModule

import (
	"errors"
	"log"
)

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

func Log(err error, message string) error {
	if message == "" {
		message = err.Error()
	}

	log.Println(message)
	return err
}

func Is(err error, expected error) bool {
	return err == expected
}
