package logging

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	responseBody    = "foo bar"
	responseBodyLen = 7
)

func TestLog(t *testing.T) {
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
	LogDebug("token", "debug message")
	LogInfo("token", "info message")
	LogWarn("token", "warn message")
	LogError("token", "error message")
}
