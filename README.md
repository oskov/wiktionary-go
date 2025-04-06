# wiktionary-go

Wiktionary-Go is a Go client library for interacting with the Wiktionary API. The client is generated from the OpenAPI schema using the [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) library.

## Features

- Automatically generated Go client for the Wiktionary API.
- Based on the OpenAPI schema provided by Wiktionary.
- Easy-to-use and type-safe API client.

## Prerequisites

- Go 1.24.1 or later.

## Installation

Clone the repository:

```bash
go get github.com/oskov/wiktionary-go
```

## Usage

See ./example/main.go

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