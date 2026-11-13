package wiktionary

import (
	"context"
	"net/http"
)

const WiktionaryURL = "https://en.wiktionary.org/api/rest_v1"

// DefaultUserAgent is sent when no custom User-Agent is configured.
// Wikimedia rejects requests with generic User-Agents (HTTP 403), see
// https://meta.wikimedia.org/wiki/User-Agent_policy. Applications should
// prefer WithUserAgent with their own name and contact information.
const DefaultUserAgent = "wiktionary-go (https://github.com/oskov/wiktionary-go)"

// WithUserAgent sets the User-Agent header on every request.
func WithUserAgent(userAgent string) ClientOption {
	return WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.Header.Set("User-Agent", userAgent)
		return nil
	})
}

// NewWiktionaryClient creates a client for the English Wiktionary REST API
// with DefaultUserAgent. Options passed in opts are applied afterwards, so
// WithUserAgent overrides the default.
func NewWiktionaryClient(opts ...ClientOption) (*ClientWithResponses, error) {
	return NewClientWithResponses(WiktionaryURL, append([]ClientOption{WithUserAgent(DefaultUserAgent)}, opts...)...)
}
