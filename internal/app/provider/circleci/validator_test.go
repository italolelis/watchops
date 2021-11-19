package circleci_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/italolelis/watchops/internal/app/provider/circleci"
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
				r.Header.Add("Circleci-Signature", "")

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
				r.Header.Add("Circleci-Signature", "wrong")

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
					"/",
					strings.NewReader("a=1&b=2"),
				)
				r.Header.Add("Circleci-Signature", "v1=06fed32dbfee920451881ec80f4507536a8ae9cc696306e969a7a3a8f29940ec")

				return r
			},
			shouldFail: false,
		},
	}

	v := circleci.NewValidator("valid")

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
