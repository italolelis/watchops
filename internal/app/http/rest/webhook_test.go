package rest_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/italolelis/watchops/internal/app/http/rest"
	"github.com/italolelis/watchops/internal/app/wh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type simpleConnector struct{ shoudlError bool }

func (c *simpleConnector) Write(ctx context.Context, payload []byte, headers map[string][]string) error {
	if c.shoudlError {
		return fmt.Errorf("connector error")
	}

	return nil
}

type wrongBody struct{}

func (wrongBody) Read(p []byte) (n int, err error) {
	return 0, errors.New(`everything is broken /o\`)
}

func TestWebhookHandler(t *testing.T) {
	cases := []struct {
		name       string
		token      func() *http.Request
		connector  wh.Connector
		body       io.Reader
		statusCode int
	}{
		{
			name:       "valid incoming request",
			connector:  &simpleConnector{},
			body:       strings.NewReader(`{"foo": "bar"}`),
			statusCode: http.StatusNoContent,
		},
		{
			name:       "failed connector",
			connector:  &simpleConnector{true},
			body:       strings.NewReader(`{"foo": "bar"}`),
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "empty body",
			connector:  &simpleConnector{true},
			body:       wrongBody{},
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "invalid body",
			connector:  &simpleConnector{true},
			body:       strings.NewReader(``),
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/webhooks", c.body)
			require.NoError(t, err)

			wh := rest.NewWebhookHandler(c.connector)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(wh.HandleWebhook)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, c.statusCode, rr.Code)
		})
	}

}
