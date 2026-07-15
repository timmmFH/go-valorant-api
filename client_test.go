package govapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientRequests(t *testing.T) {
	t.Run("premier affinity and authorization", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.Method != http.MethodGet {
				t.Errorf("method = %s, want GET", req.Method)
			}
			if req.URL.Path != "/valorant/v1/premier/Lazies/PONY" {
				t.Errorf("path = %s", req.URL.Path)
			}
			if got := req.URL.Query().Get("affinity"); got != "eu" {
				t.Errorf("affinity = %q, want eu", got)
			}
			if got := req.Header.Get("Authorization"); got != "test-key" {
				t.Errorf("Authorization = %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"status":200,"data":{"id":"team-id","name":"Lazies","tag":"PONY"}}`))
		}))
		defer server.Close()

		client, err := New("test-key", WithBaseURL(server.URL))
		if err != nil {
			t.Fatal(err)
		}

		response, err := client.PremierByNameWithResponse(
			context.Background(),
			"Lazies",
			"PONY",
			&PremierByNameParams{Affinity: Ptr("eu")},
		)
		if err != nil {
			t.Fatal(err)
		}
		if response.JSON200 == nil || response.JSON200.Data.ID != "team-id" {
			t.Fatalf("unexpected response: %#v", response)
		}
	})

	t.Run("premium webhook body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.Method != http.MethodPost {
				t.Errorf("method = %s, want POST", req.Method)
			}
			if req.URL.Path != "/public/v1/premium/webhook/users" {
				t.Errorf("path = %s", req.URL.Path)
			}

			var body PremiumWebhookUserAddRequest
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				t.Errorf("decode request: %v", err)
				return
			}
			if body.PUUID == nil || *body.PUUID != "player-puuid" {
				t.Errorf("PUUID = %v", body.PUUID)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"data":{"success":true,"user":{"id":"tracked","puuid":"player-puuid","region":"eu","enabled":true,"events":["MATCH"],"created_at":1,"updated_at":1}}}`))
		}))
		defer server.Close()

		client, err := New("premium-key", WithBaseURL(server.URL))
		if err != nil {
			t.Fatal(err)
		}
		events := []PremiumWebhookEvent{PremiumWebhookEventMatch}

		response, err := client.AddWebhookUserWithResponse(
			context.Background(),
			PremiumWebhookUserAddRequest{
				PUUID:  Ptr("player-puuid"),
				Events: &events,
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		if response.JSON201 == nil || !response.JSON201.Data.Success {
			t.Fatalf("unexpected response: %#v", response)
		}
	})
}
