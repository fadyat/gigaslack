package internal

import (
	"go.uber.org/zap"
	"net/http"
)

type writer struct{}

func (writer) WriteWithLogs(w http.ResponseWriter, body string, lg *zap.Logger) {
	_, err := w.Write([]byte(body))
	if err != nil {
		lg.Info("failed to write response body", zap.Error(err))
	}
}
