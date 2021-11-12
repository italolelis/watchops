package opsgenie

import (
	"encoding/json"
	"fmt"

	"github.com/italolelis/watchops/internal/app/provider"
)

const (
	sourceType = "opsgenie"
)

type opsgenieEvent struct {
	Action string `json:"action"`
}

// Parser is the GitHub webhook parser.
type Parser struct{}

// GetName returns the parser name.
func (p *Parser) GetName() string {
	return sourceType
}

// Parse parses the incoming webhook.
func (p *Parser) Parse(headers map[string][]string, payload []byte) (provider.Event, error) {
	var e opsgenieEvent
	if err := json.Unmarshal(payload, &e); err != nil {
		return provider.Event{}, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return provider.Event{
		EventType: e.Action,
		Signature: "",
		Source:    sourceType,
		Metadata:  payload,
	}, nil
}
