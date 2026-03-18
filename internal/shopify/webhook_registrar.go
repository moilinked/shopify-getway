package shopify

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/Khan/genqlient/graphql"
)

var defaultWebhookTopics = []WebhookSubscriptionTopic{
	WebhookSubscriptionTopicAppUninstalled,
}

type webhookClientFactory func(shopDomain, accessToken, apiVersion string) graphql.Client

type WebhookRegistrar struct {
	APIVersion string
	BaseURL    string
	Topics     []WebhookSubscriptionTopic

	newClient webhookClientFactory
}

func NewWebhookRegistrar(apiVersion, baseURL string, topics []string) *WebhookRegistrar {
	normalizedTopics := normalizeTopics(topics)
	if len(normalizedTopics) == 0 {
		normalizedTopics = append([]WebhookSubscriptionTopic(nil), defaultWebhookTopics...)
	}

	return &WebhookRegistrar{
		APIVersion: strings.TrimSpace(apiVersion),
		BaseURL:    strings.TrimSpace(baseURL),
		Topics:     normalizedTopics,
		newClient:  NewGQLClient,
	}
}

func (r *WebhookRegistrar) CallbackURL() string {
	callbackURL, _ := r.callbackURL()
	return callbackURL
}

func (r *WebhookRegistrar) EnsureShopSubscriptions(ctx context.Context, shopDomain, accessToken string) error {
	if r == nil {
		return nil
	}
	shopDomain = strings.TrimSpace(shopDomain)
	accessToken = strings.TrimSpace(accessToken)
	if shopDomain == "" || accessToken == "" {
		return errors.New("shop domain and access token are required")
	}

	callbackURL, err := r.callbackURL()
	if err != nil {
		return err
	}

	client := r.newClient(shopDomain, accessToken, r.APIVersion)
	existing, err := r.listSubscriptions(ctx, client)
	if err != nil {
		return err
	}

	byTopic := make(map[WebhookSubscriptionTopic][]WebhookSubscriptionsByTopicWebhookSubscriptionsWebhookSubscriptionConnectionNodesWebhookSubscription, len(existing))
	for _, subscription := range existing {
		byTopic[subscription.Topic] = append(byTopic[subscription.Topic], subscription)
	}

	for _, topic := range r.Topics {
		subscriptions := byTopic[topic]
		hasTarget := false
		for _, subscription := range subscriptions {
			if subscription.Uri == callbackURL && !hasTarget {
				hasTarget = true
				continue
			}
			if err := r.deleteSubscription(ctx, client, subscription.Id); err != nil {
				return fmt.Errorf("delete webhook subscription for %s: %w", topic, err)
			}
		}
		if hasTarget {
			continue
		}
		if err := r.createSubscription(ctx, client, topic, callbackURL); err != nil {
			return fmt.Errorf("create webhook subscription for %s: %w", topic, err)
		}
	}

	return nil
}

func (r *WebhookRegistrar) callbackURL() (string, error) {
	if r == nil || strings.TrimSpace(r.BaseURL) == "" {
		return "", errors.New("webhook_base_url is required to register Shopify webhooks")
	}

	base, err := url.Parse(strings.TrimSpace(r.BaseURL))
	if err != nil {
		return "", fmt.Errorf("parse webhook_base_url: %w", err)
	}
	if base.Scheme != "https" && base.Scheme != "http" {
		return "", errors.New("webhook_base_url must start with http:// or https://")
	}

	base.Path = strings.TrimRight(base.Path, "/") + "/webhooks/shopify"
	base.RawQuery = ""
	base.Fragment = ""
	return base.String(), nil
}

func (r *WebhookRegistrar) listSubscriptions(ctx context.Context, client graphql.Client) ([]WebhookSubscriptionsByTopicWebhookSubscriptionsWebhookSubscriptionConnectionNodesWebhookSubscription, error) {
	resp, err := WebhookSubscriptionsByTopic(ctx, client, r.Topics)
	if err != nil {
		return nil, err
	}

	return resp.WebhookSubscriptions.Nodes, nil
}

func (r *WebhookRegistrar) createSubscription(ctx context.Context, client graphql.Client, topic WebhookSubscriptionTopic, callbackURL string) error {
	resp, err := WebhookSubscriptionCreate(ctx, client, topic, WebhookSubscriptionInput{
		Format: ptrWebhookSubscriptionFormat(WebhookSubscriptionFormatJson),
		Uri:    &callbackURL,
	})
	if err != nil {
		return err
	}

	if resp.WebhookSubscriptionCreate == nil {
		return errors.New("missing webhookSubscriptionCreate payload")
	}
	return firstCreateUserError(resp.WebhookSubscriptionCreate.UserErrors)
}

func (r *WebhookRegistrar) deleteSubscription(ctx context.Context, client graphql.Client, id string) error {
	resp, err := WebhookSubscriptionDelete(ctx, client, id)
	if err != nil {
		return err
	}

	if resp.WebhookSubscriptionDelete == nil {
		return errors.New("missing webhookSubscriptionDelete payload")
	}
	return firstDeleteUserError(resp.WebhookSubscriptionDelete.UserErrors)
}

func normalizeTopics(topics []string) []WebhookSubscriptionTopic {
	if len(topics) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(topics))
	normalized := make([]WebhookSubscriptionTopic, 0, len(topics))
	for _, topic := range topics {
		topic = strings.TrimSpace(strings.ToUpper(topic))
		if topic == "" {
			continue
		}
		if _, ok := seen[topic]; ok {
			continue
		}
		seen[topic] = struct{}{}
		normalized = append(normalized, WebhookSubscriptionTopic(topic))
	}
	return normalized
}

func firstCreateUserError(userErrors []WebhookSubscriptionCreateWebhookSubscriptionCreateWebhookSubscriptionCreatePayloadUserErrorsUserError) error {
	if len(userErrors) == 0 {
		return nil
	}

	first := userErrors[0]
	if field := first.GetField(); len(field) > 0 {
		return fmt.Errorf("%s: %s", strings.Join(field, "."), first.GetMessage())
	}
	return errors.New(first.GetMessage())
}

func firstDeleteUserError(userErrors []WebhookSubscriptionDeleteWebhookSubscriptionDeleteWebhookSubscriptionDeletePayloadUserErrorsUserError) error {
	if len(userErrors) == 0 {
		return nil
	}

	first := userErrors[0]
	if field := first.GetField(); len(field) > 0 {
		return fmt.Errorf("%s: %s", strings.Join(field, "."), first.GetMessage())
	}
	return errors.New(first.GetMessage())
}

func ptrWebhookSubscriptionFormat(v WebhookSubscriptionFormat) *WebhookSubscriptionFormat {
	return &v
}
