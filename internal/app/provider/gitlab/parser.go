package gitlab

import (
	"encoding/json"
	"fmt"

	"github.com/italolelis/fourkeys/internal/app/provider"
)

const (
	sourceType = "gitlab"
)

var (
	supportedTypes = []string{"push", "merge_request",
		"note", "tag_push", "issue",
		"pipeline", "job", "deployment",
		"build"}
)

type gitlabEvent struct {
	ObjectKind string `json:"object_kind"`
}

// Parser is the GitHub webhook parser.
type Parser struct{}

// GetName returns the parser name.
func (p *Parser) GetName() string {
	return sourceType
}

// Parse parses the incoming webhook.
func (p *Parser) Parse(headers map[string][]string, payload []byte) (provider.Event, error) {
	var e gitlabEvent
	if err := json.Unmarshal(payload, &e); err != nil {
		return provider.Event{}, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return provider.Event{
		EventType: e.ObjectKind,
		Signature: "",
		Source:    sourceType,
		Metadata:  payload,
	}, nil
}
