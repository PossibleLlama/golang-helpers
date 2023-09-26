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

	NewLoggerBuilder().Build()
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

		NewLoggerBuilder().Build()
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
	NewLoggerBuilder().Build()

	assert.True(t, globalLogger.Level().Enabled(zapcore.DebugLevel), "Logger should produce debug messages")
	assert.True(t, globalLogger.Level().Enabled(zapcore.ErrorLevel), "Logger should produce error messages")
	assert.True(t, globalLogger.Level().Enabled(zapcore.FatalLevel), "Logger should produce fatal messages")
}

func TestQuietLogGivesLimitedLogs(t *testing.T) {
	NewLoggerBuilder().WithLogLevel(zap.PanicLevel).Build()

	assert.False(t, globalLogger.Level().Enabled(zapcore.DebugLevel), "Quiet logger should not produce debug messages")
	assert.False(t, globalLogger.Level().Enabled(zapcore.ErrorLevel), "Quiet logger should not produce error messages")
	assert.True(t, globalLogger.Level().Enabled(zapcore.FatalLevel), "Quiet logger should produce fatal messages")
}

func TestLinkOrEmpty(t *testing.T) {
	var tests = []struct {
		name    string
		commit  string
		scmLink string
		input   string
		expect  string
	}{
		{
			name:    "Empty gives empty",
			commit:  "",
			scmLink: "",
			input:   "",
			expect:  "",
		}, {
			name:    "Empty input gives empty",
			commit:  "",
			scmLink: "repoName",
			input:   "",
			expect:  "",
		}, {
			name:    "Scm and single file path input gives output to HEAD",
			commit:  "",
			scmLink: "repoName",
			input:   "repoName",
			expect:  "repoName/blob/HEAD",
		}, {
			name:    "Scm and multiple file path input gives output to HEAD",
			commit:  "",
			scmLink: "repoName",
			input:   "/repoName/bar",
			expect:  "repoName/blob/HEAD/bar",
		}, {
			name:    "Scm and multiple deep file path input gives output to HEAD",
			commit:  "",
			scmLink: "repoName",
			input:   "repoName/bar/bang/buzz",
			expect:  "repoName/blob/HEAD/bar/bang/buzz",
		}, {
			name:    "Scm, commit and single file path input gives output to commit",
			commit:  "abc123",
			scmLink: "repoName",
			input:   "repoName",
			expect:  "repoName/blob/abc123",
		}, {
			name:    "Scm, commit and multiple file path input gives output to commit",
			commit:  "abc123",
			scmLink: "repoName",
			input:   "repoName/bar",
			expect:  "repoName/blob/abc123/bar",
		}, {
			name:    "Scm and multiple file path with line input gives output to HEAD without line",
			commit:  "",
			scmLink: "repoName",
			input:   "repoName/bar:123",
			expect:  "repoName/blob/HEAD/bar",
		}, {
			name:    "Scm and file path with prefix gives output without prefix",
			commit:  "",
			scmLink: "repoName",
			input:   "prefix/not/to/output/repoName/bar/buzz",
			expect:  "repoName/blob/HEAD/bar/buzz",
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			commitSha = testItem.commit
			scmLink = testItem.scmLink
			actual := linkOrEmpty(testItem.input)

			assert.Equal(t, testItem.expect, actual)
		})
	}
}
