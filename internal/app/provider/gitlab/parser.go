package gitlab

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/italolelis/watchops/internal/app/provider"
)

const (
	sourceType = "gitlab"
)

type gitlabEvent struct {
	ObjectKind       string           `json:"object_kind"`
	CheckoutSHA      string           `json:"checkout_sha"`
	Commits          []commit         `json:"commits"`
	ObjectAttributes objectAttributes `json:"object_attributes"`
	BuildID          string           `json:"build_id"`
	DeploymentID     string           `json:"deployment_id"`
	StatusChangedAt  time.Time        `json:"status_changed_at"`
	BuildFinishedAt  time.Time        `json:"build_finished_at"`
	BuildStartedAt   time.Time        `json:"build_started_at"`
	BuildCreatedAt   time.Time        `json:"build_created_at"`
}

func (g gitlabEvent) BuildTime() time.Time {
	if !g.StatusChangedAt.IsZero() {
		return g.StatusChangedAt
	} else if !g.BuildStartedAt.IsZero() {
		return g.BuildStartedAt
	} else {
		return g.BuildCreatedAt
	}
}

type commit struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
}

type objectAttributes struct {
	ID         string    `json:"id"`
	UpdatedAt  time.Time `json:"updated_at"`
	FinishedAt time.Time `json:"finished_at"`
	CreatedAt  time.Time `json:"created_at"`
}

func (o objectAttributes) Time() time.Time {
	if !o.UpdatedAt.IsZero() {
		return o.UpdatedAt
	} else if !o.FinishedAt.IsZero() {
		return o.FinishedAt
	} else {
		return o.CreatedAt
	}
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

	event := provider.Event{
		EventType: e.ObjectKind,
		Signature: provider.GenerateSignature(payload),
		Source:    sourceType,
		Metadata:  payload,
		MsgID:     headers["msg_id"][0],
	}

	switch e.ObjectKind {
	case "push",
		"tag_push":
		event.ID = e.CheckoutSHA
		for _, c := range e.Commits {
			if c.ID == event.ID {
				i, err := strconv.ParseInt(c.Timestamp, 10, 64)
				if err != nil {
					return provider.Event{}, errors.New("failed to parse timestamp from push event")
				}

				event.TimeCreated = time.Unix(i, 0)
			}
		}

	case "merge_request",
		"note",
		"issue",
		"pipeline":
		event.ID = e.ObjectAttributes.ID
		event.TimeCreated = e.ObjectAttributes.Time()

	case "job":
		event.ID = e.BuildID
		event.TimeCreated = e.ObjectAttributes.Time()
	case "deployment":
		event.ID = e.DeploymentID
		event.TimeCreated = e.StatusChangedAt
	case "build":
		event.ID = e.BuildID
		event.TimeCreated = e.BuildTime()
	default:
		return provider.Event{}, &provider.UnkownTypeError{Type: e.ObjectKind}
	}

	if event.TimeCreated.IsZero() {
		rawTime := headers["publish_time"][0]
		publishTime, err := strconv.ParseInt(rawTime, 10, 64)
		if err != nil {
			return provider.Event{}, errors.New("unable to parse time created event")
		}

		event.TimeCreated = time.Unix(publishTime, 0)
	}

	return event, nil
}
