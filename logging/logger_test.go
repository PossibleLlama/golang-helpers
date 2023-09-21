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
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
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
			name:  "debug logging",
			f:     func() { LogDebug("a", "b") },
			level: zapcore.DebugLevel,
		}, {
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
		observedZapCore, observedLogs := observer.New(zap.DebugLevel)
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

func TestLogEnablesAllLevelsOfLogs(t *testing.T) {
	// Set defaults via init
	InitLogger("v0.1.2", "proj", "svc")

	assert.True(t, globalLogger.Level().Enabled(zapcore.DebugLevel), "Logger should produce debug messages")
	assert.True(t, globalLogger.Level().Enabled(zapcore.ErrorLevel), "Logger should produce error messages")
	assert.True(t, globalLogger.Level().Enabled(zapcore.FatalLevel), "Logger should produce fatal messages")
}

func TestQuietLogGivesLimitedLogs(t *testing.T) {
	// Set defaults via init
	InitQuietLogger()

	assert.False(t, globalLogger.Level().Enabled(zapcore.DebugLevel), "Quiet logger should not produce debug messages")
	assert.False(t, globalLogger.Level().Enabled(zapcore.ErrorLevel), "Quiet logger should not produce error messages")
	assert.True(t, globalLogger.Level().Enabled(zapcore.FatalLevel), "Quiet logger should produce fatal messages")
}

func TestGithubLinkOrEmpty(t *testing.T) {
	var tests = []struct {
		name   string
		commit string
		input  string
		expect string
	}{
		{
			name:   "Empty gives empty",
			commit: "",
			input:  "",
			expect: "",
		}, {
			name:   "Foo gives empty",
			commit: "",
			input:  "foo",
			expect: "",
		}, {
			name:   "github gives empty",
			commit: "",
			input:  "github",
			expect: "",
		}, {
			name:   "foo/github gives empty",
			commit: "",
			input:  "foo/github",
			expect: "",
		}, {
			name:   "foo/bar/github gives empty",
			commit: "",
			input:  "foo/bar/github",
			expect: "",
		}, {
			name:   "/foo/github gives empty",
			commit: "",
			input:  "/foo/github",
			expect: "",
		}, {
			name:   "/foo/bar/github gives empty",
			commit: "",
			input:  "/foo/bar/github",
			expect: "",
		}, {
			name:   "github/org gives empty",
			commit: "",
			input:  "github/foo",
			expect: "",
		}, {
			name:   "github/org/repo gives empty",
			commit: "",
			input:  "github/org/repo",
			expect: "",
		}, {
			name:   "github/org/repo/file without commit gives link to main branch",
			commit: "",
			input:  "github/org/repo/file",
			expect: "https://github.com/org/repo/blob/main/file",
		}, {
			name:   "github/org/repo/file with commit gives link to specific commit",
			commit: "abc",
			input:  "github/org/repo/file",
			expect: "https://github.com/org/repo/blob/abc/file",
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			commitSha = testItem.commit
			actual := githubLinkOrEmpty(testItem.input)

			assert.Equal(t, testItem.expect, actual)
		})
	}
}
