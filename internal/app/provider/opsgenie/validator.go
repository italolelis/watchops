package opsgenie

import (
	"bytes"
	"errors"
	"net/http"
)

const tokenHeader = "X-TOKEN"

var (
	// ErrEmptyToken is used the the incoming request has not token set
	ErrEmptyToken = errors.New("the token can't be empty")
	// ErrInvalidToken is used the the incoming request has an invalid token
	ErrInvalidToken = errors.New("the token is invalid")
)

type Validator struct {
	secret []byte
}

func NewValidator(secret string) *Validator {
	return &Validator{secret: []byte(secret)}
}

func (v *Validator) Validate(r *http.Request) error {
	token := []byte(r.Header.Get(tokenHeader))
	if bytes.Equal(token, []byte("")) {
		return ErrEmptyToken
	}

	if !bytes.Equal(token, v.secret) {
		return ErrInvalidToken
	}

	return nil
}

func (g *Validator) IsSupported(incomingProvider string) bool {
	return "opsgenie" == incomingProvider
}
