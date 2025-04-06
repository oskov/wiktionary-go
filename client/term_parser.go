package wiktionary

import (
	"encoding/json"
	"fmt"
)

type ParsedPageDefinitionTerm struct {
	PartOfSpeech string       `json:"partOfSpeech"`
	Language     string       `json:"language"`
	Definitions  []Definition `json:"definitions"`
}

type TermDefinition struct {
	Definition     string          `json:"definition"`
	ParsedExamples []ParsedExample `json:"parsedExamples,omitempty"`
	Examples       []string        `json:"examples,omitempty"`
}

type ParsedExample struct {
	Example string `json:"example"`
}

// GetPageDefinitionTermResponse is a response from the GetPageDefinitionTerm endpoint. The openapi scheme seems incorrect for this endpoint,
// so this method tries to parse the response body. to provide a meaningful response.
func (r GetPageDefinitionTermResponse) GetParsedResponse() (map[string][]ParsedPageDefinitionTerm, error) {
	if r.StatusCode() != 200 {
		return nil, fmt.Errorf("bad status code: %d", r.StatusCode())
	}
	var parsedResponse map[string][]ParsedPageDefinitionTerm
	if err := json.Unmarshal(r.Body, &parsedResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return parsedResponse, nil
}
