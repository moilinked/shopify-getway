package httpapi

import (
	"net/http"

	"shopify-gateway/internal/config"
	mdw "shopify-gateway/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(cfg config.Config) http.Handler {
	r := chi.NewRouter()
	r.Use(mdw.AllowAllCORS)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mdw.RequestLogger)
	r.Use(middleware.Recoverer)
	r.Use(mdw.ShopContextMiddleware())

	r.Get("/ping", HandlePing)

	r.Route("/admin", func(admin chi.Router) {
		admin.Use(mdw.ShopifySessionTokenMiddleware(cfg.ShopifyAPIKey, cfg.ShopifyAPISecret, cfg.DebugAuth))
		admin.Get("/ping", HandleAdminPing)
	})

	r.Route("/app", func(app chi.Router) {
		app.Use(mdw.ShopifyAppProxySignatureMiddleware(cfg.ShopifyAPISecret, cfg.DebugAuth))
		app.Get("/ping", HandleAppPing)
	})

	r.With(mdw.ShopifyWebhookMiddleware(cfg.ShopifyAPISecret, cfg.DebugAuth)).
		Post("/webhooks/shopify", HandleShopifyWebhook)

	return r
}
