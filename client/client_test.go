package wiktionary

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const definitionContentType = `application/json; charset=utf-8; profile="https://www.mediawiki.org/wiki/Specs/definition/0.8.0"`

const definitionFixture = `{
  "en": [
    {
      "partOfSpeech": "Noun",
      "language": "English",
      "definitions": [
        {
          "definition": "Something that is representative of all such things in a group.",
          "parsedExamples": [{"example": "This is an example."}],
          "examples": ["This is an example."]
        },
        {
          "definition": "Something that serves to illustrate a rule."
        }
      ]
    }
  ],
  "fr": [
    {
      "partOfSpeech": "Noun",
      "language": "French",
      "definitions": [{"definition": "example"}]
    }
  ]
}`

// newTestServer starts a server that records the last request and replies
// with the given status, content type and body.
func newTestServer(t *testing.T, status int, contentType, body string) (*httptest.Server, *http.Request) {
	t.Helper()
	var last http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		last = *r.Clone(context.Background())
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return srv, &last
}

func TestNewWiktionaryClientUserAgent(t *testing.T) {
	tests := []struct {
		name string
		opts []ClientOption
		want string
	}{
		{name: "default", want: DefaultUserAgent},
		{name: "custom", opts: []ClientOption{WithUserAgent("my-app/1.0")}, want: "my-app/1.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, last := newTestServer(t, http.StatusOK, definitionContentType, definitionFixture)

			// Point the client at the test server instead of WiktionaryURL.
			opts := append([]ClientOption{WithBaseURL(srv.URL)}, tt.opts...)
			c, err := NewWiktionaryClient(opts...)
			if err != nil {
				t.Fatalf("NewWiktionaryClient: %v", err)
			}
			if _, err := c.GetPageDefinitionTermWithResponse(context.Background(), "example", nil); err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if got := last.Header.Get("User-Agent"); got != tt.want {
				t.Errorf("User-Agent = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetPageDefinitionTermRequestPath(t *testing.T) {
	srv, last := newTestServer(t, http.StatusOK, definitionContentType, definitionFixture)

	c, err := NewClientWithResponses(srv.URL, WithUserAgent("test"))
	if err != nil {
		t.Fatalf("NewClientWithResponses: %v", err)
	}
	if _, err := c.GetPageDefinitionTermWithResponse(context.Background(), "hello world", nil); err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if last.Method != http.MethodGet {
		t.Errorf("method = %s, want GET", last.Method)
	}
	if got, want := last.URL.Path, "/page/definition/hello world"; got != want {
		t.Errorf("path = %q, want %q", got, want)
	}
}

func TestGetParsedResponse(t *testing.T) {
	srv, _ := newTestServer(t, http.StatusOK, definitionContentType, definitionFixture)

	c, err := NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("NewClientWithResponses: %v", err)
	}
	rsp, err := c.GetPageDefinitionTermWithResponse(context.Background(), "example", nil)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	parsed, err := rsp.GetParsedResponse()
	if err != nil {
		t.Fatalf("GetParsedResponse: %v", err)
	}
	if len(parsed) != 2 {
		t.Fatalf("got %d languages, want 2", len(parsed))
	}

	en := parsed["en"]
	if len(en) != 1 {
		t.Fatalf("got %d English terms, want 1", len(en))
	}
	term := en[0]
	if term.PartOfSpeech != "Noun" || term.Language != "English" {
		t.Errorf("term = %+v, want Noun/English", term)
	}
	if len(term.Definitions) != 2 {
		t.Fatalf("got %d definitions, want 2", len(term.Definitions))
	}
	first := term.Definitions[0]
	if first.Definition != "Something that is representative of all such things in a group." {
		t.Errorf("definition = %q", first.Definition)
	}
	if len(first.ParsedExamples) != 1 || first.ParsedExamples[0].Example != "This is an example." {
		t.Errorf("parsedExamples = %+v", first.ParsedExamples)
	}
	if len(first.Examples) != 1 || first.Examples[0] != "This is an example." {
		t.Errorf("examples = %+v", first.Examples)
	}
	if second := term.Definitions[1]; second.ParsedExamples != nil || second.Examples != nil {
		t.Errorf("expected no examples for second definition, got %+v", second)
	}
}

func TestGetParsedResponseErrors(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		contentType string
		body        string
	}{
		{
			name:        "not found",
			status:      http.StatusNotFound,
			contentType: "application/problem+json",
			body:        `{"type":"https://mediawiki.org/wiki/HyperSwitch/errors/not_found","title":"Not found.","status":404}`,
		},
		{
			name:        "forbidden",
			status:      http.StatusForbidden,
			contentType: "text/plain",
			body:        "Please set a user-agent",
		},
		{
			name:        "invalid json",
			status:      http.StatusOK,
			contentType: "text/plain",
			body:        "not json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, _ := newTestServer(t, tt.status, tt.contentType, tt.body)

			c, err := NewClientWithResponses(srv.URL)
			if err != nil {
				t.Fatalf("NewClientWithResponses: %v", err)
			}
			rsp, err := c.GetPageDefinitionTermWithResponse(context.Background(), "example", nil)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if _, err := rsp.GetParsedResponse(); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

// TestLiveGetPageDefinitionTerm hits the real Wiktionary API.
// Run with: WIKTIONARY_LIVE_TEST=1 go test ./client -run Live
func TestLiveGetPageDefinitionTerm(t *testing.T) {
	if os.Getenv("WIKTIONARY_LIVE_TEST") == "" {
		t.Skip("set WIKTIONARY_LIVE_TEST=1 to run against the live API")
	}

	c, err := NewWiktionaryClient()
	if err != nil {
		t.Fatalf("NewWiktionaryClient: %v", err)
	}
	rsp, err := c.GetPageDefinitionTermWithResponse(context.Background(), "example", nil)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	parsed, err := rsp.GetParsedResponse()
	if err != nil {
		t.Fatalf("GetParsedResponse: %v (body: %.200s)", err, rsp.Body)
	}
	if len(parsed["en"]) == 0 {
		t.Fatal("expected English definitions for \"example\"")
	}
}
