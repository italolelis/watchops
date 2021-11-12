package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/italolelis/watchops/internal/app/provider"
	"github.com/italolelis/watchops/internal/app/wh"
	"github.com/italolelis/watchops/internal/pkg/log"
)

type WebhookHandler struct {
	c wh.Connector
}

// NewWebhookHandler creates a new content handler.
func NewWebhookHandler(c wh.Connector) *WebhookHandler {
	return &WebhookHandler{c: c}
}

func (h *WebhookHandler) Routes(validators ...provider.SignatureValidator) http.Handler {
	r := chi.NewRouter()
	r.With(ValidateSignature(validators...)).
		Post("/", h.HandleWebhook)

	return r
}

// SaveContent responsible to receive the callback from a webhook.
func (h *WebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	logger := log.WithContext(r.Context())

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorw("failed to read payload", "err", err)
		http.Error(w, "failed to read payload", http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	if err := h.c.Write(r.Context(), payload, r.Header); err != nil {
		logger.Errorw("failed to publish message", "err", err)
		http.Error(w, "failed to publish message", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
