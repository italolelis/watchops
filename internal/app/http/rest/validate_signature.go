package rest

import (
	"net/http"
	"strings"

	"github.com/italolelis/watchops/internal/app/provider"
	"github.com/italolelis/watchops/internal/pkg/log"
)

// ValidateSignature middleware to validate the signature header.
func ValidateSignature(validators ...provider.SignatureValidator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			incomingProvider := getSource(r.Header)

			logger := log.WithContext(r.Context()).Named("request_validator").With("provider", incomingProvider)
			logger.Debug("validating request")

			var validated bool

			for _, v := range validators {
				pLogger := logger
				pLogger.Debug("checking validator for incoming provider")

				if !v.IsSupported(incomingProvider) {
					continue
				}

				pLogger.Debug("incoming provider supported")

				if err := v.Validate(r); err != nil {
					pLogger.Debugw("failed to validate the request for provider", "err", err)
					continue
				}

				validated = true
				break
			}

			defer r.Body.Close()

			if !validated {
				http.Error(w, "failed to check the origin request signature", http.StatusUnauthorized)
				return
			}

			logger.Debug("request validated")
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func getSource(headers map[string][]string) string {
	if _, ok := headers["X-Gitlab-Event"]; ok {
		return "gitlab"
	}

	if _, ok := headers["Circleci-Event-Type"]; ok {
		return "circleci"
	}

	if _, ok := headers["X-Pagerduty-Signature"]; ok {
		return "pagerduty"
	}

	source := strings.TrimSpace(strings.Split(headers["User-Agent"][0], "/")[0])
	switch source {
	case "GitHub-Hookshot":
		return "github"
	case "Opsgenie Http Client":
		return "opsgenie"
	default:
		return ""
	}
}
