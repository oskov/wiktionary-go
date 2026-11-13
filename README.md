# wiktionary-go

Wiktionary-Go is a Go client library for interacting with the Wiktionary API. The client is generated from the OpenAPI schema using the [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) library.

## Features

- Automatically generated Go client for the Wiktionary API.
- Based on the OpenAPI schema provided by Wiktionary.
- Easy-to-use and type-safe API client.

## Prerequisites

- Go 1.27 or later.

## Installation

```bash
go get github.com/oskov/wiktionary-go/client
```

## Usage

```go
import wiktionary "github.com/oskov/wiktionary-go/client"

client, err := wiktionary.NewWiktionaryClient(
	wiktionary.WithUserAgent("my-app/1.0 (https://example.com; me@example.com)"),
)
if err != nil {
	log.Fatal(err)
}
rsp, err := client.GetPageDefinitionTermWithResponse(ctx, "example", nil)
if err != nil {
	log.Fatal(err)
}
definitions, err := rsp.GetParsedResponse()
```

Wikimedia rejects requests without a descriptive `User-Agent` (HTTP 403), see the
[User-Agent policy](https://meta.wikimedia.org/wiki/User-Agent_policy).
`NewWiktionaryClient` sets a default one; pass `WithUserAgent` to identify your application.

See [example/main.go](./example/main.go) for a complete program.

## Generating the client

The client is generated using the OpenAPI schema and the oapi-codegen library. To regenerate the client:

1. Ensure the OpenAPI schema is downloaded to ./api/wiktionary/openapi.json. You can use the provided script:
```bash
./download-openapi.sh
```
2. Run the go generate command to generate the client:

```bash
go generate ./...
```