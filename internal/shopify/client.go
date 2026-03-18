package shopify

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Khan/genqlient/graphql"
)

var sharedTransport http.RoundTripper = &http.Transport{
	MaxIdleConns:        100,
	MaxIdleConnsPerHost: 10,
	IdleConnTimeout:     90 * time.Second,
}

type authTransport struct {
	token string
	base  http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-Shopify-Access-Token", t.token)
	req.Header.Set("Content-Type", "application/json")
	return t.base.RoundTrip(req)
}

func NewGQLClient(shopDomain, accessToken, apiVersion string) graphql.Client {
	endpoint := fmt.Sprintf(
		"https://%s/admin/api/%s/graphql.json",
		strings.TrimSpace(shopDomain),
		strings.TrimSpace(apiVersion),
	)
	httpClient := &http.Client{
		Timeout:   30 * time.Second,
		Transport: &authTransport{token: accessToken, base: sharedTransport},
	}
	return graphql.NewClient(endpoint, httpClient)
}
