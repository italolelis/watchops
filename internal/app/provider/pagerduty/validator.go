package pagerduty

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	webhookSignaturePrefix = "v1="
	webhookSignatureHeader = "X-PagerDuty-Signature"
	webhookBodyReaderLimit = 2 * 1024 * 1024 // 2MB
)

var (
	// ErrNoValidSignatures is returned when a webhook is not properly signed
	// with the expected signature. When receiving this error, it is reccommended
	// that the server return HTTP 403 to prevent redelivery.
	ErrNoValidSignatures = errors.New("invalid webhook signature")
	// ErrMalformedHeader is returned when the *http.Request is missing the
	// X-PagerDuty-Signature header. When receiving this error, it is recommended
	// that the server return HTTP 400 to prevent redelivery.
	ErrMalformedHeader = errors.New("X-PagerDuty-Signature header is either missing or malformed")
	// ErrMalformedBody is returned when the *http.Request body is either
	// missing or malformed. When receiving this error, it's recommended that the
	// server return HTTP 400 to prevent redelivery.
	ErrMalformedBody = errors.New("HTTP request body is either empty or malformed")
)

type Validator struct {
	secret []byte
}

func NewValidator(secret string) *Validator {
	return &Validator{secret: []byte(secret)}
}

func (v *Validator) Validate(r *http.Request) error {
	h := r.Header.Get(webhookSignatureHeader)
	if len(h) == 0 {
		return ErrMalformedHeader
	}

	orb := r.Body

	b, err := ioutil.ReadAll(io.LimitReader(r.Body, webhookBodyReaderLimit))
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	defer func() { _ = orb.Close() }()
	r.Body = ioutil.NopCloser(bytes.NewReader(b))

	if len(b) == 0 {
		return ErrMalformedBody
	}

	sigs := v.extractPayloadSignatures(h)
	if len(sigs) == 0 {
		return ErrMalformedHeader
	}

	s := v.calculateSignature(b, v.secret)

	for _, sig := range sigs {
		if hmac.Equal(s, sig) {
			return nil
		}
	}

	return ErrNoValidSignatures
}

func (g *Validator) IsSupported(incomingProvider string) bool {
	return sourceType == incomingProvider
}

func (g *Validator) extractPayloadSignatures(s string) [][]byte {
	var sigs [][]byte

	for _, sv := range strings.Split(s, ",") {
		// Ignore any signatures that are not the initial v1 version.
		if !strings.HasPrefix(sv, webhookSignaturePrefix) {
			continue
		}

		sig, err := hex.DecodeString(strings.TrimPrefix(sv, webhookSignaturePrefix))
		if err != nil {
			continue
		}

		sigs = append(sigs, sig)
	}

	return sigs
}

func (g *Validator) calculateSignature(payload []byte, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write(payload)
	return mac.Sum(nil)
}
