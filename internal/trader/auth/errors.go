package auth

import "fmt"

type KeySetParseError struct {
	err error
}

func (k *KeySetParseError) Error() string {
	return fmt.Sprintf("error parsing key set from JWT provider: %s", k.err.Error())
}

type TokenValidationError struct {
	err error
}

func (m *TokenValidationError) Error() string {
	return fmt.Sprintf("error validating token: %s", m.err.Error())
}

type InvalidFormatError struct {
	field string
}

func (i *InvalidFormatError) Error() string {
	return fmt.Sprintf("error parsing field %s from JWT token", i.field)
}
