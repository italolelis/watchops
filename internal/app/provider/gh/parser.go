package gh

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/italolelis/watchops/internal/app/provider"
)

const (
	sourceType      = "github"
	eventTypeHeader = "X-Github-Event"
	signatureHeader = "X-Hub-Signature-256"
)

var (
	// ErrorUnknownType is used when the given type is not supported.
	ErrorUnknownType = errors.New("unknown event type")

	// ErrorNoHeadCommit is used when the the event has no head commit.
	ErrorNoHeadCommit = errors.New("no head commit provided in the event")
)

// Parser is the GitHub webhook parser.
type Parser struct{}

// GetName returns the parser name.
func (p *Parser) GetName() string {
	return sourceType
}

// Parse parses the incoming webhook.
// nolint: exhaustivestruct, funlen
func (p *Parser) Parse(headers map[string][]string, payload []byte) (provider.Event, error) {
	eventType := headers[eventTypeHeader][0]

	rawEvent, err := github.ParseWebHook(eventType, payload)
	if err != nil {
		return provider.Event{}, fmt.Errorf("could not parse webhook: %w", err)
	}

	event := provider.Event{
		EventType: eventType,
		Signature: headers[signatureHeader][0],
		Source:    sourceType,
		Metadata:  payload,
		MsgID:     headers["msg_id"][0],
	}

	switch e := rawEvent.(type) {
	case *github.PushEvent:
		if e.GetHeadCommit() == nil {
			return provider.Event{}, ErrorNoHeadCommit
		}

		event.TimeCreated = e.GetHeadCommit().GetTimestamp().Time
		event.ID = e.GetHeadCommit().GetID()
	case *github.PullRequestEvent:
		event.TimeCreated = e.GetPullRequest().GetUpdatedAt()
		event.ID = strconv.Itoa(e.GetNumber())
	case *github.PullRequestReviewEvent:
		event.TimeCreated = e.GetReview().GetSubmittedAt()
		event.ID = strconv.FormatInt(e.GetReview().GetID(), 2)
	case *github.PullRequestReviewCommentEvent:
		event.TimeCreated = e.GetComment().GetUpdatedAt()
		event.ID = strconv.FormatInt(e.GetComment().GetID(), 2)
	case *github.DeploymentEvent:
		event.TimeCreated = e.GetDeployment().GetUpdatedAt().Time
		event.ID = strconv.FormatInt(e.GetDeployment().GetID(), 2)
	case *github.DeploymentStatusEvent:
		event.TimeCreated = e.GetDeploymentStatus().GetUpdatedAt().Time
		event.ID = strconv.FormatInt(e.GetDeploymentStatus().GetID(), 2)
	case *github.IssuesEvent:
		event.TimeCreated = e.GetIssue().GetUpdatedAt()
		event.ID = e.GetRepo().GetName() + "/" + strconv.Itoa(e.GetIssue().GetNumber())
	case *github.IssueCommentEvent:
		event.TimeCreated = e.GetComment().GetUpdatedAt()
		event.ID = strconv.FormatInt(e.GetComment().GetID(), 2)
	case *github.CheckRunEvent:
		if !e.GetCheckRun().GetCompletedAt().Time.IsZero() {
			event.TimeCreated = e.GetCheckRun().GetCompletedAt().Time
		} else if !e.GetCheckRun().GetStartedAt().Time.IsZero() {
			event.TimeCreated = e.GetCheckRun().GetStartedAt().Time
		}
		event.ID = strconv.FormatInt(e.GetCheckRun().GetID(), 2)
	case *github.CheckSuiteEvent:
		if !e.GetCheckSuite().GetApp().GetUpdatedAt().IsZero() {
			event.TimeCreated = e.GetCheckSuite().GetApp().GetUpdatedAt()
		} else if !e.GetCheckSuite().GetApp().GetCreatedAt().IsZero() {
			event.TimeCreated = e.GetCheckSuite().GetApp().GetCreatedAt()
		}
		event.ID = strconv.FormatInt(e.GetCheckSuite().GetID(), 2)
	case *github.StatusEvent:
		event.TimeCreated = e.GetUpdatedAt().Time
		event.ID = strconv.FormatInt(e.GetID(), 2)
	case *github.ReleaseEvent:
		if !e.GetRelease().GetCreatedAt().Time.IsZero() {
			event.TimeCreated = e.GetRelease().GetCreatedAt().Time
		} else if !e.GetRelease().GetPublishedAt().Time.IsZero() {
			event.TimeCreated = e.GetRelease().GetPublishedAt().Time
		}

		event.ID = strconv.FormatInt(e.GetRelease().GetID(), 2)
	default:
		return provider.Event{}, fmt.Errorf("unsupported event type %s", eventType)
	}

	return event, nil
}
