# go-valorant-api

Typed Go client for HenrikDev's unofficial VALORANT API 4.9.0, including the
premium webhook endpoints.

The client is generated from the checked-in OpenAPI snapshot, so request
parameters, response models, and supported endpoints stay aligned with the API
contract.

## Install

```bash
go get github.com/timmmFH/go-valorant-api/v2
```

## Use

```go
package main

import (
	"context"
	"log"

	govapi "github.com/timmmFH/go-valorant-api/v2"
)

func main() {
	client, err := govapi.New("your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	response, err := client.GetAccountV2WithResponse(
		context.Background(),
		"Henrik",
		"dev",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	if response.JSON200 == nil {
		log.Fatalf("HenrikDev returned %s: %s", response.Status(), response.Body)
	}

	log.Printf("PUUID: %s", response.JSON200.Data.PUUID)
}
```

Every operation has a `WithResponse` method that returns the decoded response,
the raw body, and the underlying `http.Response`. Use `WithHTTPClient` or
`WithBaseURL` when constructing the client to customize transport behavior or
target a test server.

## Premium webhooks

Premium operations use the same client and API key:

```go
events := []govapi.PremiumWebhookEvent{govapi.PremiumWebhookEventMatch}

response, err := client.AddWebhookUserWithResponse(
	context.Background(),
	govapi.PremiumWebhookUserAddRequest{
		PUUID:  govapi.Ptr("player-puuid"),
		Events: &events,
	},
)
```

The current HenrikDev OpenAPI contract does not define a JSON schema for the
successful webhook settings and update responses. Those bodies remain available
through `response.Body` until HenrikDev publishes the missing schemas.

## Regenerate

```bash
go generate ./...
```

The snapshot is normalized to OpenAPI 3.0 because the generator does not yet
support HenrikDev's nullable OpenAPI 3.1 forms.

This project is unofficial and is not endorsed by Riot Games.
