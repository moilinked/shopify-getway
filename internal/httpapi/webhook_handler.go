package httpapi

import (
	"io"
	"net/http"
	"strings"

	"shopify-gateway/internal/logger"
	mw "shopify-gateway/internal/middleware"
)

func HandleShopifyWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		ErrBadRequest.Respond(w, "invalid request body")
		return
	}

	shop, _ := mw.ShopFromContext(r.Context())
	topic := strings.TrimSpace(r.Header.Get("X-Shopify-Topic"))
	webhookID := strings.TrimSpace(r.Header.Get("X-Shopify-Webhook-Id"))

	logger.Log.Info().
		Str("shop", shop).
		Str("topic", topic).
		Str("webhook_id", webhookID).
		Int("bytes", len(body)).
		Str("body", string(body)).
		Msg("shopify webhook received")

	RespondOK(w, map[string]any{
		"ok": true,
	})
}
