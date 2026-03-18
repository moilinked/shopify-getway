package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	mw "shopify-gateway/internal/middleware"
	"shopify-gateway/internal/shopify"
)

type WebhookRegistrationHandler struct {
	Registrar *shopify.WebhookRegistrar
}

type registerWebhooksRequest struct {
	Shop        string `json:"shop"`
	AccessToken string `json:"accessToken"`
}

func (h *WebhookRegistrationHandler) Register(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.Registrar == nil {
		ErrInternal.Respond(w, "webhook registrar is not configured")
		return
	}

	var req registerWebhooksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrBadRequest.Respond(w, "invalid request body")
		return
	}

	shop := strings.TrimSpace(req.Shop)
	if ctxShop, ok := mw.ShopFromContext(r.Context()); ok {
		shop = ctxShop
	}
	if shop == "" {
		ErrBadRequest.Respond(w, "shop is required")
		return
	}
	if strings.TrimSpace(req.AccessToken) == "" {
		ErrBadRequest.Respond(w, "accessToken is required")
		return
	}

	if err := h.Registrar.EnsureShopSubscriptions(r.Context(), shop, req.AccessToken); err != nil {
		ErrInternal.Respond(w, err.Error())
		return
	}

	RespondOK(w, map[string]any{
		"ok":           true,
		"shop":         shop,
		"callback_url": h.Registrar.CallbackURL(),
		"topics":       h.Registrar.Topics,
	})
}
