package main

import (
	"context"
	"flag"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"shopify-gateway/internal/config"
	"shopify-gateway/internal/httpapi"
	"shopify-gateway/internal/logger"
	"shopify-gateway/internal/shopify"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("load config")
	}

	if err := logger.Init(cfg.LogLevel); err != nil {
		logger.Log.Fatal().Err(err).Msg("init logger")
	}
	defer logger.Close()

	registrar := shopify.NewWebhookRegistrar(cfg.ShopifyAPIVersion, cfg.WebhookBaseURL, cfg.WebhookTopics)

	router := httpapi.NewRouter(cfg, registrar)
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Log.Info().Str("addr", ":"+cfg.Port).Msgf("listening on http://localhost:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal().Err(err).Msg("server error")
		}
	}()

	<-ctx.Done()
	logger.Log.Info().Msg("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error().Err(err).Msg("server shutdown error")
	}
}
