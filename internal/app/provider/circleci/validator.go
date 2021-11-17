package opsgenie

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
)

var (
	// ErrEmptyToken is used the the incoming request has not token set
	ErrEmptyToken = errors.New("the token can't be empty")
	// ErrInvalidToken is used the the incoming request has an invalid token
	ErrInvalidToken = errors.New("the token is invalid")
	// ErrReadBody is used the the incoming request can't be read
	ErrReadBody = errors.New("failed to read the request body")
)

type Validator struct {
	secret []byte
}

func NewValidator(secret string) *Validator {
	return &Validator{secret: []byte(secret)}
}

func (v *Validator) Validate(r *http.Request) error {
	var expectedSignature string = "v1="

	h := hmac.New(sha256.New, []byte(v.secret))

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return ErrReadBody
	}

	h.Write(b)

	expectedSignature += hex.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(r.Header.Get(signatureHeader)), []byte(expectedSignature)) {
		return ErrInvalidToken
	}

	return nil
}

func (g *Validator) IsSupported(incomingProvider string) bool {
	return sourceType == incomingProvider
}
