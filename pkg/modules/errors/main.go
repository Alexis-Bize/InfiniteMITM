package errors

import (
	"errors"
	"log"
)

var (
	ErrProxy = errors.New("proxy error")
	ErrProxyToggleInvalidCommand = errors.New("invalid proxy toggle command")
)

func Handle(err error, message string) error {
	if message == "" {
		message = err.Error()
	}

	log.Println(message)
	return err
}

func Is(err error, expected error) bool {
	return err == expected
}
