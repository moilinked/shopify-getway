package shopify

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/Khan/genqlient/graphql"
)

func TestWebhookRegistrarRequiresBaseURL(t *testing.T) {
	t.Parallel()

	r := NewWebhookRegistrar("2025-01", "", []string{"APP_UNINSTALLED"})
	err := r.EnsureShopSubscriptions(context.Background(), "demo.myshopify.com", "shpat_test")
	if err == nil {
		t.Fatal("expected error when webhook_base_url is missing")
	}
}

func TestWebhookRegistrarCreatesMissingTopics(t *testing.T) {
	t.Parallel()

	client := &fakeGraphQLClient{}
	r := NewWebhookRegistrar("2025-01", "https://example.com", []string{
		"APP_UNINSTALLED",
		"CUSTOMERS_DATA_REQUEST",
	})
	r.newClient = func(shopDomain, accessToken, apiVersion string) graphql.Client {
		return client
	}

	err := r.EnsureShopSubscriptions(context.Background(), "demo.myshopify.com", "shpat_test")
	if err != nil {
		t.Fatalf("EnsureShopSubscriptions: %v", err)
	}

	wantCreates := []string{"APP_UNINSTALLED", "CUSTOMERS_DATA_REQUEST"}
	if !reflect.DeepEqual(client.createdTopics, wantCreates) {
		t.Fatalf("created topics: got %v, want %v", client.createdTopics, wantCreates)
	}
	if len(client.deletedIDs) != 0 {
		t.Fatalf("deleted ids: got %v, want none", client.deletedIDs)
	}
}

func TestWebhookRegistrarReplacesStaleSubscriptions(t *testing.T) {
	t.Parallel()

	client := &fakeGraphQLClient{
		listNodes: []map[string]any{
			{"id": "1", "topic": "APP_UNINSTALLED", "uri": "https://old.example.com/webhooks/shopify"},
			{"id": "2", "topic": "CUSTOMERS_DATA_REQUEST", "uri": "https://example.com/webhooks/shopify"},
			{"id": "3", "topic": "CUSTOMERS_DATA_REQUEST", "uri": "https://old.example.com/webhooks/shopify"},
		},
	}
	r := NewWebhookRegistrar("2025-01", "https://example.com", []string{
		"APP_UNINSTALLED",
		"CUSTOMERS_DATA_REQUEST",
	})
	r.newClient = func(shopDomain, accessToken, apiVersion string) graphql.Client {
		return client
	}

	err := r.EnsureShopSubscriptions(context.Background(), "demo.myshopify.com", "shpat_test")
	if err != nil {
		t.Fatalf("EnsureShopSubscriptions: %v", err)
	}

	if !reflect.DeepEqual(client.deletedIDs, []string{"1", "3"}) {
		t.Fatalf("deleted ids: got %v, want %v", client.deletedIDs, []string{"1", "3"})
	}
	if !reflect.DeepEqual(client.createdTopics, []string{"APP_UNINSTALLED"}) {
		t.Fatalf("created topics: got %v, want %v", client.createdTopics, []string{"APP_UNINSTALLED"})
	}
}

type fakeGraphQLClient struct {
	listNodes     []map[string]any
	createdTopics []string
	deletedIDs    []string
	createUserErr error
	deleteUserErr error
}

func (f *fakeGraphQLClient) MakeRequest(_ context.Context, req *graphql.Request, resp *graphql.Response) error {
	switch req.OpName {
	case "WebhookSubscriptionsByTopic":
		return setResponseData(resp.Data, map[string]any{
			"webhookSubscriptions": map[string]any{
				"nodes": f.listNodes,
			},
		})
	case "WebhookSubscriptionCreate":
		var vars struct {
			Topic               string `json:"topic"`
			WebhookSubscription struct {
				Uri *string `json:"uri"`
			} `json:"webhookSubscription"`
		}
		if err := decodeVariables(req.Variables, &vars); err != nil {
			return err
		}
		f.createdTopics = append(f.createdTopics, vars.Topic)
		userErrors := []map[string]any{}
		if f.createUserErr != nil {
			userErrors = append(userErrors, map[string]any{
				"field":   []string{"topic"},
				"message": f.createUserErr.Error(),
			})
		}
		return setResponseData(resp.Data, map[string]any{
			"webhookSubscriptionCreate": map[string]any{
				"webhookSubscription": map[string]any{
					"id":    "new-id",
					"topic": vars.Topic,
					"uri":   derefString(vars.WebhookSubscription.Uri),
				},
				"userErrors": userErrors,
			},
		})
	case "WebhookSubscriptionDelete":
		var vars struct {
			ID string `json:"id"`
		}
		if err := decodeVariables(req.Variables, &vars); err != nil {
			return err
		}
		f.deletedIDs = append(f.deletedIDs, vars.ID)
		userErrors := []map[string]any{}
		if f.deleteUserErr != nil {
			userErrors = append(userErrors, map[string]any{
				"field":   []string{"id"},
				"message": f.deleteUserErr.Error(),
			})
		}
		return setResponseData(resp.Data, map[string]any{
			"webhookSubscriptionDelete": map[string]any{
				"deletedWebhookSubscriptionId": vars.ID,
				"userErrors":                   userErrors,
			},
		})
	default:
		return errors.New("unexpected op: " + req.OpName)
	}
}

func setResponseData(target any, payload any) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, target)
}

func decodeVariables(input any, target any) error {
	raw, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, target)
}

func derefString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
