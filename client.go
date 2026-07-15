package govapi

import (
	"context"
	"net/http"
)

const defaultBaseURL = "https://api.henrikdev.xyz"

const (
	PremiumWebhookEventMatch = MATCH
	PremiumWebhookEventMMR   = MMR
)

// New returns a HenrikDev API client. An empty API key leaves authorization
// unset for endpoints that permit anonymous requests.
func New(apiKey string, opts ...ClientOption) (*ClientWithResponses, error) {
	auth := WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		if apiKey != "" {
			req.Header.Set("Authorization", apiKey)
		}
		return nil
	})

	return NewClientWithResponses(defaultBaseURL, append([]ClientOption{auth}, opts...)...)
}

// Ptr returns a pointer to v for optional request fields and parameters.
func Ptr[T any](v T) *T {
	return &v
}
