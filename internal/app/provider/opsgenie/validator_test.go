package opsgenie

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateToken(t *testing.T) {
	cases := []struct {
		name       string
		token      func() *http.Request
		shouldFail bool
	}{
		{
			name: "empty header",
			token: func() *http.Request {
				r, _ := http.NewRequestWithContext(
					context.Background(),
					http.MethodGet,
					"/", nil,
				)

				return r
			},
			shouldFail: true,
		},
		{
			name: "empty token",
			token: func() *http.Request {
				r, _ := http.NewRequestWithContext(
					context.Background(),
					http.MethodGet,
					"/", nil,
				)
				r.Header.Add("X-TOKEN", "")

				return r
			},
			shouldFail: true,
		},
		{
			name: "invalid token",
			token: func() *http.Request {
				r, _ := http.NewRequestWithContext(
					context.Background(),
					http.MethodGet,
					"/", nil,
				)
				r.Header.Add("X-TOKEN", "wrong")

				return r
			},
			shouldFail: true,
		},
		{
			name: "correct token",
			token: func() *http.Request {
				r, _ := http.NewRequestWithContext(
					context.Background(),
					http.MethodGet,
					"/", nil,
				)
				r.Header.Add("X-TOKEN", "valid")

				return r
			},
			shouldFail: false,
		},
	}

	v := NewValidator("valid")

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			err := v.Validate(c.token())
			if c.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
