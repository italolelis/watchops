package provider

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// UnkownTypeError represents a unkown incoming event type from the provider.
type UnkownTypeError struct {
	Type string
}

func (u *UnkownTypeError) Error() string {
	return fmt.Sprintf("unsupported event type %s", u.Type)
}

// Parser holds the methods for a webhook parser.
type Parser interface {
	GetName() string
	Parse(headers map[string][]string, payload []byte) (Event, error)
}

type SignatureValidator interface {
	Validate(r *http.Request) error
	IsSupported(incomingProvider string) bool
}

// Writer represents the writer repository.
type Writer interface {
	Add(context.Context, Event) error
	Close() error
}

// Event represents a generic data event.
type Event struct {
	EventType   string          `json:"event_type"`
	ID          string          `json:"id"`
	Metadata    json.RawMessage `json:"metadata"`
	TimeCreated time.Time       `json:"time_created"`
	Signature   string          `json:"signature"`
	MsgID       string          `json:"msg_id"`
	Source      string          `json:"source"`
}

func GenerateSignature(value []byte) string {
	hasher := sha1.New()
	hasher.Write(value)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

// ParserRegistry is a registry of parsers.
type ParserRegistry struct {
	parsers map[string]Parser
}

// NewParserRegistry creates a new parser registry.
func NewParserRegistry() *ParserRegistry {
	return &ParserRegistry{
		parsers: make(map[string]Parser),
	}
}

// Get returns the parser for the given name.
func (pr *ParserRegistry) Get(name string) Parser {
	return pr.parsers[name]
}

// Register registers a parser.
func (pr *ParserRegistry) Register(parser Parser) *ParserRegistry {
	pr.parsers[parser.GetName()] = parser

	return pr
}
