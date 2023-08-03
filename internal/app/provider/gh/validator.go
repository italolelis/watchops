package gh

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/v53/github"
)

type Validator struct {
	secret []byte
}

func NewValidator(secret string) *Validator {
	return &Validator{secret: []byte(secret)}
}

func (g *Validator) Validate(r *http.Request) error {
	payload, err := github.ValidatePayload(r, g.secret)
	if err != nil {
		return fmt.Errorf("failed to check the origin request signature: %w", err)
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(payload))
	return nil
}

func (g *Validator) IsSupported(incomingProvider string) bool {
	return sourceType == incomingProvider
}
