package opsgenie

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/italolelis/watchops/internal/app/provider"
)

const (
	sourceType = "opsgenie"
)

type opsgenieEvent struct {
	Action string `json:"action"`
	Alert  struct {
		ID        string `json:"alertId"`
		UpdatedAt int64  `json:"updatedAt"`
	} `json:"alert"`
	EscalationID  string    `json:"escalationId"`
	EsclationTime time.Time `json:"escalationTime"`
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

	event := provider.Event{
		EventType: e.Action,
		Signature: provider.GenerateSignature(payload),
		Source:    sourceType,
		Metadata:  payload,
		MsgID:     headers["msg_id"][0],
	}

	switch e.Action {
	case "Create",
		"Acknowledge",
		"UnAcknowledge",
		"AddTeam",
		"AddRecipient",
		"AddNote",
		"AddTags",
		"RemoveTags",
		"Close",
		"AssignOwnership",
		"TakeOwnership",
		"Delete",
		"UpdatePriority",
		"UpdateDescription",
		"UpdateMessage":
		event.ID = e.Alert.ID
		event.TimeCreated = time.Unix(e.Alert.UpdatedAt, 0)
	case "Escalate":
		event.ID = e.EscalationID
		event.TimeCreated = e.EsclationTime
	default:
		return provider.Event{}, fmt.Errorf("unsupported event type %s", e.Action)
	}

	return event, nil
}
