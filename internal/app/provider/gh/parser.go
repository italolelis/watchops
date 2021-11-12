package gh

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/italolelis/watchops/internal/app/provider"
)

const (
	sourceType        = "github"
	eventTypeHeader   = "X-Github-Event"
	ghSignatureHeader = "X-Hub-Signature-256"
)

type ghMetadata struct {
	Action      string       `json:"action,omitempty"`
	Push        *push        `json:"push,omitempty"`
	PullRequest *pullRequest `json:"pull_request,omitempty"`
	Repository  *repo        `json:"repository,omitempty"`
	Release     *release     `json:"release,omitempty"`
	Review      *review      `json:"review,omitempty"`
	Sender      *user        `json:"sender,omitempty"`
	Deployment  *deployment  `json:"deployment,omitempty"`
}

type push struct {
	ID  string `json:"id"`
	Ref string `json:"ref"`
}

type pullRequest struct {
	ID        int64  `json:"id"`
	Number    int    `json:"number"`
	Title     string `json:"title"`
	State     string `json:"state"`
	Merged    bool   `json:"merged"`
	Repo      *repo  `json:"repo,omitempty"`
	MergedBy  *user  `json:"merged_by,omitempty"`
	Comments  int    `json:"comments"`
	Commits   int    `json:"commits"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
	HeadRef   string `json:"head_ref"`
	BaseRef   string `json:"base_ref"`
}

type release struct {
	Name       string `json:"name"`
	Draft      bool   `json:"draft"`
	PreRelease bool   `json:"prerelease"`
	Author     *user  `json:"author,omitempty"`
}

type review struct {
	User  *user  `json:"user,omitempty"`
	State string `json:"state"`
}

type deployment struct {
	ID          int64            `json:"id"`
	Environment string           `json:"environment"`
	Status      deploymentStatus `json:"status"`
}

type deploymentStatus struct {
	State string `json:"state"`
}

type user struct {
	Login string `json:"login"`
	URL   string `json:"url"`
}

type repo struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	URL      string `json:"url"`
}

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
		Signature: headers[ghSignatureHeader][0],
		Source:    sourceType,
	}

	var metadata ghMetadata

	switch e := rawEvent.(type) {
	case *github.PushEvent:
		if e.GetHeadCommit() == nil {
			return provider.Event{}, ErrorNoHeadCommit
		}

		event.TimeCreated = e.GetHeadCommit().GetTimestamp().Time
		event.ID = e.GetHeadCommit().GetID()
		metadata = ghMetadata{
			Sender: convertUser(e.GetSender()),
			Push: &push{
				ID:  e.GetHeadCommit().GetID(),
				Ref: e.GetRef(),
			},
			Repository: &repo{
				Name:     e.GetRepo().GetName(),
				FullName: e.GetRepo().GetFullName(),
				Private:  e.GetRepo().GetPrivate(),
				URL:      e.GetRepo().GetURL(),
			},
		}
	case *github.PullRequestEvent:
		event.TimeCreated = e.GetPullRequest().GetUpdatedAt()
		event.ID = strconv.Itoa(e.GetNumber())
		metadata = ghMetadata{
			Action:      e.GetAction(),
			Sender:      convertUser(e.GetSender()),
			Repository:  convertRepo(e.GetRepo()),
			PullRequest: convertPullRequest(e.GetPullRequest()),
		}
	case *github.PullRequestReviewEvent:
		event.TimeCreated = e.GetReview().GetSubmittedAt()
		event.ID = strconv.Itoa(int(e.GetReview().GetID()))
		metadata = ghMetadata{
			Action:      e.GetAction(),
			Sender:      convertUser(e.GetSender()),
			Review:      convertReview(e.GetReview()),
			Repository:  convertRepo(e.GetRepo()),
			PullRequest: convertPullRequest(e.GetPullRequest()),
		}
	case *github.PullRequestReviewCommentEvent:
		event.TimeCreated = e.GetComment().GetUpdatedAt()
		event.ID = strconv.Itoa(int(e.GetComment().GetID()))
		metadata = ghMetadata{
			Action:      e.GetAction(),
			Sender:      convertUser(e.GetSender()),
			Repository:  convertRepo(e.GetRepo()),
			PullRequest: convertPullRequest(e.GetPullRequest()),
		}
	case *github.DeploymentEvent:
		event.TimeCreated = e.GetDeployment().GetUpdatedAt().Time
		event.ID = strconv.Itoa(int(e.GetDeployment().GetID()))
		metadata = ghMetadata{
			Deployment: convertDeployment(e.GetDeployment()),
			Sender:     convertUser(e.GetSender()),
			Repository: convertRepo(e.GetRepo()),
		}
	case *github.DeploymentStatusEvent:
		event.TimeCreated = e.GetDeploymentStatus().GetUpdatedAt().Time
		event.ID = strconv.Itoa(int(e.GetDeploymentStatus().GetID()))
		metadata = ghMetadata{
			Deployment: convertDeploymentStatus(e),
			Sender:     convertUser(e.GetSender()),
			Repository: convertRepo(e.GetRepo()),
		}
	case *github.ReleaseEvent:
		if !e.GetRelease().GetCreatedAt().Time.IsZero() {
			event.TimeCreated = e.GetRelease().GetCreatedAt().Time
		} else if !e.GetRelease().GetPublishedAt().Time.IsZero() {
			event.TimeCreated = e.GetRelease().GetPublishedAt().Time
		}

		event.ID = strconv.Itoa(int(e.GetRelease().GetID()))

		metadata = ghMetadata{
			Action:     e.GetAction(),
			Sender:     convertUser(e.GetSender()),
			Repository: convertRepo(e.GetRepo()),
			Release:    convertRelease(e.GetRelease()),
		}
	default:
		return provider.Event{}, fmt.Errorf("unsupported event type %s", eventType)
	}

	rawMetadata, err := json.Marshal(metadata)
	if err != nil {
		return provider.Event{}, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	event.Metadata = rawMetadata

	return event, nil
}

func convertPullRequest(p *github.PullRequest) *pullRequest {
	return &pullRequest{
		ID:        p.GetID(),
		Number:    p.GetNumber(),
		State:     p.GetState(),
		Title:     p.GetTitle(),
		Comments:  p.GetComments(),
		Commits:   p.GetCommits(),
		Merged:    p.GetMerged(),
		Additions: p.GetAdditions(),
		Deletions: p.GetDeletions(),
		Repo:      convertRepo(p.Base.Repo),
		MergedBy:  convertUser(p.MergedBy),
		HeadRef:   p.GetHead().GetRef(),
		BaseRef:   p.GetBase().GetRef(),
	}
}

func convertDeployment(e *github.Deployment) *deployment {
	return &deployment{
		ID:          e.GetID(),
		Environment: e.GetEnvironment(),
	}
}

func convertDeploymentStatus(e *github.DeploymentStatusEvent) *deployment {
	return &deployment{
		ID:          e.GetDeployment().GetID(),
		Environment: e.GetDeployment().GetEnvironment(),
		Status: deploymentStatus{
			State: e.GetDeploymentStatus().GetState(),
		},
	}
}

func convertReview(r *github.PullRequestReview) *review {
	return &review{
		User:  convertUser(r.User),
		State: r.GetState(),
	}
}

func convertRelease(r *github.RepositoryRelease) *release {
	return &release{
		Name:       r.GetName(),
		Author:     convertUser(r.Author),
		Draft:      r.GetDraft(),
		PreRelease: r.GetPrerelease(),
	}
}

func convertRepo(r *github.Repository) *repo {
	return &repo{
		Name:     r.GetName(),
		FullName: r.GetFullName(),
		Private:  r.GetPrivate(),
		URL:      r.GetURL(),
	}
}

func convertUser(u *github.User) *user {
	return &user{
		Login: u.GetLogin(),
		URL:   u.GetURL(),
	}
}
