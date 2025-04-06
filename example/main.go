package main

import (
	"context"
	"log"

	"github.com/oskov/wiktionary-go/client"
)

func main() {
	client, err := wiktionary.NewClientWithResponses(wiktionary.WiktionaryURL)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	rsp, err := client.GetPageDefinitionTermWithResponse(context.Background(), "example", nil)
	if err != nil {
		log.Fatalf("Failed to get page definition: %v", err)
	}
	if rsp.StatusCode() != 200 {
		log.Fatalf("Unexpected status code: %d", rsp.StatusCode())
	}
	parsedResponse, err := rsp.GetParsedResponse()
	if err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}
	log.Printf("Response: %+v", parsedResponse)
}
