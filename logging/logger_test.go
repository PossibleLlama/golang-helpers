package logging

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

const (
	responseBody    = "foo bar"
	responseBodyLen = 7
)

func TestLoggerResponse(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)

	// Set defaults via init
	InitLogger("v0.1.2", "proj", "svc")
	// Replace actual logger with dummy
	globalLogger = observedLogger

	req, err := http.NewRequest(http.MethodGet, "http://testing/", nil)
	assert.Nil(t, err)

	lrm := &LoggingResponseMetadata{
		Status: 0,
		Size:   0,
	}
	lwr := &LoggingResponseWriter{
		ResponseMetadata: lrm,
		ResponseWriter:   httptest.NewRecorder(),
	}
	lwr.WriteHeader(http.StatusPaymentRequired)
	bytes, err := lwr.Write([]byte(responseBody))

	assert.Nil(t, err)
	assert.Equal(t, http.StatusPaymentRequired, lwr.ResponseMetadata.Status)
	assert.Equal(t, responseBodyLen, lwr.ResponseMetadata.Size)
	assert.Equal(t, responseBodyLen, bytes)

	LogResponse("token", req, lrm, 0)

	assert.NotEmpty(t, observedLogs.Len(), "At least one log should have been output")
	assert.Len(t, observedLogs.All(), 1, "Should be exactly 1 log")

	currentLog := observedLogs.All()[0]
	assert.Equal(t, currentLog.Level, zapcore.InfoLevel)
	assert.Equal(t, currentLog.Message, "finished request")
}
