package shortener

import (
	"crypto/rand"
	"errors"
	"fmt"
	"regexp"
)

const (
	codeAlphabet = "abcdefghijkmnopqrstuvwxyzABCDEFGHIJKMNOPQRSTUVWXYZ1234567890"
)

var ErrInvalidCode = errors.New("invalid code")

var codePattern = regexp.MustCompile("^[a-zA-Z0-9]+$")

func GenerateCode(lenght int) (string, error) {
	b := make([]byte, lenght)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	for i := range b {
		b[i] = codeAlphabet[int(b[i])%len(codeAlphabet)]
	}

	return string(b), nil
}

func ValidateCode(code string) error {
	if code == "" {
		return fmt.Errorf("%w: empty", ErrInvalidCode)
	}
	if !codePattern.MatchString(code) {
		return fmt.Errorf("%w: must be alphanumeric", ErrInvalidCode)
	}

	return nil
}
