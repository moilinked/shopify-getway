package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port              string   `yaml:"port"`
	ShopifyAPIKey     string   `yaml:"shopify_api_key"`
	ShopifyAPISecret  string   `yaml:"shopify_api_secret"`
	ShopifyAPIVersion string   `yaml:"shopify_api_version"`
	WebhookBaseURL    string   `yaml:"webhook_base_url"`
	WebhookTopics     []string `yaml:"webhook_topics"`
	DebugAuth         bool     `yaml:"debug_auth"`
	LogLevel          string   `yaml:"log_level"`
}

func Load(path string) (Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return Config{}, err
	}

	cfg.Port = strings.TrimSpace(cfg.Port)
	if cfg.Port == "" {
		cfg.Port = "9998"
	}

	cfg.ShopifyAPIKey = strings.TrimSpace(cfg.ShopifyAPIKey)
	cfg.ShopifyAPISecret = strings.TrimSpace(cfg.ShopifyAPISecret)
	if cfg.ShopifyAPIKey == "" || cfg.ShopifyAPISecret == "" {
		return Config{}, errors.New("shopify_api_key and shopify_api_secret are required")
	}

	cfg.ShopifyAPIVersion = strings.TrimSpace(cfg.ShopifyAPIVersion)
	if cfg.ShopifyAPIVersion == "" {
		cfg.ShopifyAPIVersion = "2025-01"
	}

	cfg.WebhookBaseURL = strings.TrimSpace(cfg.WebhookBaseURL)
	if cfg.WebhookBaseURL == "" {
		return Config{}, errors.New("webhook_base_url is required")
	}
	baseURL, err := url.Parse(cfg.WebhookBaseURL)
	if err != nil {
		return Config{}, fmt.Errorf("parse webhook_base_url: %w", err)
	}
	if baseURL.Scheme != "https" && baseURL.Scheme != "http" {
		return Config{}, errors.New("webhook_base_url must start with http:// or https://")
	}

	cfg.LogLevel = strings.TrimSpace(strings.ToLower(cfg.LogLevel))
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	if len(cfg.WebhookTopics) == 0 {
		cfg.WebhookTopics = []string{"APP_UNINSTALLED"}
	}
	for i := range cfg.WebhookTopics {
		cfg.WebhookTopics[i] = strings.TrimSpace(strings.ToUpper(cfg.WebhookTopics[i]))
		if cfg.WebhookTopics[i] == "" {
			return Config{}, errors.New("webhook_topics cannot contain empty values")
		}
	}

	return cfg, nil
}
