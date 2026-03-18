package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"shopify-gateway/internal/logger"
)

func ShopifyWebhookMiddleware(apiSecret string, debugAuth bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestPath := r.URL.Path
			signature := strings.TrimSpace(r.Header.Get("X-Shopify-Hmac-Sha256"))
			if signature == "" {
				if debugAuth {
					logger.Log.Debug().Str("path", requestPath).Msg("shopify_webhook: missing hmac header")
				}
				http.Error(w, "missing webhook signature", http.StatusUnauthorized)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Log.Error().Str("path", requestPath).Err(err).Msg("shopify_webhook: read body failed")
				http.Error(w, "invalid request body", http.StatusBadRequest)
				return
			}
			_ = r.Body.Close()

			mac := hmac.New(sha256.New, []byte(apiSecret))
			_, _ = mac.Write(body)
			expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
			if !hmac.Equal([]byte(signature), []byte(expected)) {
				if debugAuth {
					logger.Log.Warn().Str("path", requestPath).
						Str("got_prefix", shortBase64(signature)).
						Str("expected_prefix", shortBase64(expected)).
						Msg("shopify_webhook: signature mismatch")
				}
				http.Error(w, "invalid webhook signature", http.StatusUnauthorized)
				return
			}

			if debugAuth {
				logger.Log.Info().Str("path", requestPath).Int("bytes", len(body)).Msg("shopify_webhook: verified")
			}

			ctx := r.Context()
			if shop := strings.TrimSpace(r.Header.Get("X-Shopify-Shop-Domain")); shop != "" {
				if shopDomain, err := extractShopDomain(shop); err == nil {
					ctx = reconcileShopContext(ctx, shopDomain)
				}
			}

			r.Body = io.NopCloser(bytes.NewReader(body))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func shortBase64(s string) string {
	if len(s) <= 8 {
		return s
	}
	return s[:8]
}
