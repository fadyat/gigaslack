package main

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"googlesheets-slackbot-golang/cmd/config"
	"googlesheets-slackbot-golang/internal"
	"googlesheets-slackbot-golang/internal/service"
	"log"
	"net/http"
)

var (
	Version = "dev"
)

func initLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("failed to create logger: ", err)
	}

	zap.ReplaceGlobals(logger)
}

func main() {
	initLogger()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	creds, err := google.CredentialsFromJSON(
		context.Background(), cfg.Google.Credentials, sheets.SpreadsheetsReadonlyScope,
	)
	if err != nil {
		log.Fatal("failed to get credentials from json: ", err)
	}

	googleSpreadsheets, err := service.NewGoogleSpreadsheets(&cfg.Google, creds)
	if err != nil {
		log.Fatal("failed to create google spreadsheets service: ", err)
	}

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
