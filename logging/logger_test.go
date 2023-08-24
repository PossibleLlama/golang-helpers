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

func TestLoggerLogs(t *testing.T) {
	var tests = []struct {
		name  string
		f     func()
		level zapcore.Level
	}{
		{
			name:  "info logging",
			f:     func() { LogInfo("a", "b") },
			level: zapcore.InfoLevel,
		}, {
			name:  "warn logging",
			f:     func() { LogWarn("a", "b") },
			level: zapcore.WarnLevel,
		}, {
			name:  "error logging",
			f:     func() { LogError("a", "b") },
			level: zapcore.ErrorLevel,
		},
	}

	for _, testItem := range tests {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)

	// Set defaults via init
	InitLogger("v0.1.2", "proj", "svc")
	// Replace actual logger with dummy
	globalLogger = observedLogger

		t.Run(testItem.name, func(t *testing.T) {
	assert.Empty(t, observedLogs.Len(), "Before logging, there should be no lines")

			testItem.f()

	assert.NotEmpty(t, observedLogs.Len(), "At least one log should have been output")
	assert.Len(t, observedLogs.All(), 1, "There should be exactly one log per ran test")
	currentLog := observedLogs.All()[0]

	assert.NotEmpty(t, currentLog.Message, "Output from logger should not be empty")

			assert.Equal(t, currentLog.Level, testItem.level, "Output from logger should be at the correct level")
	assert.Equal(t, currentLog.Message, "b", "Output from logger should contain line")
		})
	}
}
