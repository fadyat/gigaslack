package main

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"googlesheets-slackbot-golang/cmd/config"
	"googlesheets-slackbot-golang/internal"
	"googlesheets-slackbot-golang/internal/service"
	"log"
	"net/http"
)

var (
	Version = "dev"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("failed to create logger: ", err)
	}

	zap.ReplaceGlobals(logger)
	googleSpreadsheets := service.NewGoogleSpreadsheets(&cfg.Google)
	slackHandler := internal.NewSlackHandler(&cfg.Slack, googleSpreadsheets, zap.L())

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		WriteTimeout: cfg.Server.WriteTimeout,
		ReadTimeout:  cfg.Server.ReadTimeout,
		Handler:      slackHandler,
	}

	zap.L().Info("starting server", zap.String("version", Version), zap.String("addr", server.Addr))
	if e := server.ListenAndServe(); e != nil && !errors.Is(e, http.ErrServerClosed) {
		log.Fatal("failed to start server: ", e)
	}
}
