package internal

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
	"googlesheets-slackbot-golang/cmd/config"
	"googlesheets-slackbot-golang/internal/service"
	"io"
	"net/http"
)

type SlackHandler struct {
	cfg                *config.Slack
	slackAPI           *slack.Client
	googleSpreadsheets *service.GoogleSpreadsheets
	lg                 *zap.Logger
	writer             writer
}

func NewSlackHandler(
	cfg *config.Slack,
	googleSpreadsheets *service.GoogleSpreadsheets,
	lg *zap.Logger,
) *SlackHandler {

	return &SlackHandler{
		cfg:                cfg,
		slackAPI:           slack.New(cfg.BotToken),
		googleSpreadsheets: googleSpreadsheets,
		lg:                 lg,
		writer:             writer{},
	}
}

func (ss *SlackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		ss.lg.Info("method not allowed")
		return
	}

	ss.verifyMiddleware(ss.slashCommandHandler)(w, r)
}

func (ss *SlackHandler) verifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ss.lg.Info("failed to read request body", zap.Error(err))
			return
		}

		sv, err := slack.NewSecretsVerifier(r.Header, ss.cfg.SigningSecret)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ss.lg.Info("failed to create secrets verifier", zap.Error(err))
			return
		}

		if _, err = sv.Write(body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ss.lg.Info("failed to write request body to secrets verifier", zap.Error(err))
			return
		}

		if err = sv.Ensure(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			ss.lg.Info("failed to verify request signature", zap.Error(err))
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))
		next(w, r)
	}
}

func (ss *SlackHandler) slashCommandHandler(w http.ResponseWriter, r *http.Request) {
	cmd, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ss.lg.Info("failed to parse slash command", zap.Error(err))
		return
	}

	user, err := ss.slackAPI.GetUserProfile(&slack.GetUserProfileParameters{
		UserID: cmd.UserID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ss.lg.Info("failed to get user profile", zap.Error(err))
		return
	}

	spreadsheetData, err := ss.googleSpreadsheets.TakeByValue(user.Email)

	switch {
	case err == nil:
	case errors.Is(err, service.ErrValueNotFound):
		w.WriteHeader(http.StatusOK)
		ss.writer.WriteWithLogs(w, "You are not in the spreadsheet :(", ss.lg)
		ss.lg.Info("user not found in spreadsheet", zap.String("user", user.Email))
		return
	case errors.Is(err, service.ErrHeadersNotFound), errors.Is(err, service.ErrSearchColumnNotFound), errors.Is(err, service.ErrTakeColumnNotFound):
		w.WriteHeader(http.StatusOK)
		ss.writer.WriteWithLogs(w, "Some of the table columns are changed, please contact the administrator", ss.lg)
		ss.lg.Info("failed to take values from spreadsheet", zap.Error(err))
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		ss.lg.Info("failed to take values from spreadsheet", zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	ss.writer.WriteWithLogs(w, fmt.Sprintf("Hello, %s!\n\n%s\n", user.Email, ss.cfg.Custom.SuccessMsg), ss.lg)
	ss.writer.WriteWithLogs(w, fmt.Sprintf("%s\n", spreadsheetData), ss.lg)
	ss.lg.Info("successfully processed slash command", zap.String("user", user.Email))
}
