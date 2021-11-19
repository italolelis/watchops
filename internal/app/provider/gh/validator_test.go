package gh_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/italolelis/watchops/internal/app/provider/gh"
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
					http.MethodPost,
					"/",
					strings.NewReader("a=1&b=2"),
				)
				r.Header.Add("Content-Type", "application/json")

				return r
			},
			shouldFail: true,
		},
		{
			name: "empty body",
			token: func() *http.Request {
				r, _ := http.NewRequestWithContext(
					context.Background(),
					http.MethodPost,
					"/",
					http.NoBody,
				)

				return r
			},
			shouldFail: true,
		},
		{
			name: "invalid header",
			token: func() *http.Request {
				r, _ := http.NewRequestWithContext(
					context.Background(),
					http.MethodGet,
					"/", nil,
				)
				r.Header.Add("X-Hub-Signature", "")

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
				r.Header.Add("X-Hub-Signature", "wrong")

				return r
			},
			shouldFail: true,
		},
		{
			name: "correct token",
			token: func() *http.Request {
				body := map[string]interface{}{
					"test": "test",
				}
				rawBody, _ := json.Marshal(body)

				r, _ := http.NewRequestWithContext(
					context.Background(),
					http.MethodGet,
					"/",
					bytes.NewReader(rawBody),
				)
				r.Header.Add("Content-Type", "application/json")
				r.Header.Add("X-Hub-Signature", "sha256=92b251de78498d9b43dc0cd194b63b91cbf6d5ff8b2e3359f00d66bfbcc8efa6")

				return r
			},
			shouldFail: false,
		},
	}

	v := gh.NewValidator("valid")

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
