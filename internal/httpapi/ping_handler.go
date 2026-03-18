package httpapi

import (
	"net/http"

	mw "shopify-gateway/internal/middleware"
)

func HandlePing(w http.ResponseWriter, r *http.Request) {
	shop, _ := mw.ShopFromContext(r.Context())
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":   true,
		"shop": shop,
	})
}

func HandleAdminPing(w http.ResponseWriter, r *http.Request) {
	shop, _ := mw.ShopFromContext(r.Context())
	claims, _ := mw.ClaimsFromContext(r.Context())
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":     true,
		"shop":   shop,
		"claims": claims,
	})
}

func HandleAppPing(w http.ResponseWriter, r *http.Request) {
	shop, _ := mw.ShopFromContext(r.Context())
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":   true,
		"shop": shop,
	})
}
