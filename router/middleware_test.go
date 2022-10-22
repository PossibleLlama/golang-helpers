package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PossibleLlama/golang-helpers/strings"

	"github.com/stretchr/testify/assert"
)

func TestSetDefaultHeadersMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	actualRes := httptest.NewRecorder()
	testHandler := SetDefaultHeadersMiddleware(nextHandler)
	testHandler.ServeHTTP(actualRes, httptest.NewRequest("GET", "http://test", nil))

	assert.NotEmpty(t, actualRes.Header().Get("Content-Type"))
	assert.NotEmpty(t, actualRes.Header().Get("Access-Control-Allow-Origin"))
	assert.NotEmpty(t, actualRes.Header().Get("Access-Control-Allow-Methods"))
	assert.NotEmpty(t, actualRes.Header().Get("Access-Control-Allow-Headers"))

	assert.Equal(t, "application/json; charset=utf-8", actualRes.Header().Get("Content-Type"))
	assert.Equal(t, "*", actualRes.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, actualRes.Header().Get("Access-Control-Allow-Methods"), "GET")
	assert.Contains(t, actualRes.Header().Get("Access-Control-Allow-Methods"), "POST")
	assert.Contains(t, actualRes.Header().Get("Access-Control-Allow-Methods"), "PUT")
	assert.Contains(t, actualRes.Header().Get("Access-Control-Allow-Methods"), "DELETE")
	assert.Equal(t, "Authorization", actualRes.Header().Get("Access-Control-Allow-Headers"))
}

func TestCheckTraceTokenMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	actualReq := httptest.NewRequest("GET", "http://test", nil)
	actualRes := httptest.NewRecorder()
	testHandler := CheckTraceTokenMiddleware(nextHandler)
	testHandler.ServeHTTP(actualRes, actualReq)

	assert.NotEmpty(t, actualReq.Header.Get("X-Trace-Token"))
	assert.NotEmpty(t, actualRes.Header().Get("X-Trace-Token"))
	assert.Equal(t, actualReq.Header.Get("X-Trace-Token"), actualRes.Header().Get("X-Trace-Token"))
	assert.Len(t, actualReq.Header.Get("X-Trace-Token"), 16)
}

func TestCheckTraceTokenMiddlewareWithProvided2Char(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	actualReq := httptest.NewRequest("GET", "http://test", nil)
	actualReq.Header.Add("X-Trace-Token", strings.RandAlphabeticString(2))
	actualRes := httptest.NewRecorder()
	testHandler := CheckTraceTokenMiddleware(nextHandler)
	testHandler.ServeHTTP(actualRes, actualReq)

	assert.NotEmpty(t, actualReq.Header.Get("X-Trace-Token"))
	assert.NotEmpty(t, actualRes.Header().Get("X-Trace-Token"))
	assert.Equal(t, actualReq.Header.Get("X-Trace-Token"), actualRes.Header().Get("X-Trace-Token"))
	assert.Len(t, actualReq.Header.Get("X-Trace-Token"), 16)
}

func TestCheckTraceTokenMiddlewareWithProvided10Char(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	actualReq := httptest.NewRequest("GET", "http://test", nil)
	actualReq.Header.Add("X-Trace-Token", strings.RandAlphabeticString(10))
	actualRes := httptest.NewRecorder()
	testHandler := CheckTraceTokenMiddleware(nextHandler)
	testHandler.ServeHTTP(actualRes, actualReq)

	assert.NotEmpty(t, actualReq.Header.Get("X-Trace-Token"))
	assert.NotEmpty(t, actualRes.Header().Get("X-Trace-Token"))
	assert.Equal(t, actualReq.Header.Get("X-Trace-Token"), actualRes.Header().Get("X-Trace-Token"))
	assert.Len(t, actualReq.Header.Get("X-Trace-Token"), 10)
}
