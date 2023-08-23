package logging

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	responseBody    = "foo bar"
	responseBodyLen = 7
)

func TestLoggerResponse(t *testing.T) {
	InitLogger("v0.1.2", "proj", "svc")
	req, err := http.NewRequest("GET", "http://testing/", nil)
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

	assert.Equal(t, "", captureOutput(func() {
	LogDebug("token", "debug message")
	}), "Output from void logger should be empty")
}

func TestLoggerLogs(t *testing.T) {
	var tests = []struct {
		level    string
		f        func()
		contains []string
	}{
		{
			level: "debug",
			f:     func() { LogDebug("token", "debug message") },
			contains: []string{
				"\"" + TraceToken + "\":\"token\"",
				"\"message\":\"debug message\"",
			},
		}, {
			level: "info",
			f:     func() { LogInfo("token", "info message") },
			contains: []string{
				"\"" + TraceToken + "\":\"token\"",
				"\"message\":\"info message\"",
			},
		}, {
			level: "warn",
			f:     func() { LogWarn("token", "warn message") },
			contains: []string{
				"\"" + TraceToken + "\":\"token\"",
				"\"message\":\"warn message\"",
			},
		}, {
			level: "error",
			f:     func() { LogError("token", "error message") },
			contains: []string{
				"\"" + TraceToken + "\":\"token\"",
				"\"message\":\"error message\"",
			},
		},
	}

	InitLogger("v0.1.2", "proj", "svc")

	for _, testItem := range tests {
		t.Run(testItem.level, func(t *testing.T) {
			output := captureOutput(func() {
				testItem.f()
			})
			assert.NotEmpty(t, output, "Output from logger should not be empty")
			for _, e := range testItem.contains {
				assert.Contains(t, output, "\"level\":\""+testItem.level+"\"")
				assert.Contains(t, output, e, "Output from logger should contain line")
			}
			assert.Contains(t, output, "\"version\":\"v0.1.2\"")
			assert.Contains(t, output, "\"project\":\"proj\"")
			assert.Contains(t, output, "\"service\":\"svc\"")
			assert.Contains(t, output, "\"environment\":\"dev\"")
			assert.Contains(t, output, "\"time\":\"") // Check that time has a value, but don't check value
		})
	}
}

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}
