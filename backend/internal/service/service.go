package service

import (
	"errors"
	"strings"
)

var ErrValidation = errors.New("service: validation failed")

func required(values ...string) error {
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return ErrValidation
		}
	}
	return nil
}
