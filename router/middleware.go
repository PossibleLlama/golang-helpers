package router

import (
	"net/http"
	"time"

	"github.com/PossibleLlama/golang-helpers/logging"
	"github.com/PossibleLlama/golang-helpers/strings"
)

const (
	HEADER_AUTH   = "Authorization"
	JSON_ENCODING = "application/json; charset=utf-8"
)

// SetDefaultHeadersMiddleware Sets default headers
func SetDefaultHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", JSON_ENCODING)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
		w.Header().Set("Access-Control-Allow-Headers", HEADER_AUTH)
		next.ServeHTTP(w, r)
	})
}

// CheckTraceTokenMiddleware Checks for x-trace-token,
// and if it doesn't exist creates one
func CheckTraceTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.Header.Get(logging.TraceToken)) < 6 {
			token := strings.RandAlphabeticString(16)
			r.Header.Set(logging.TraceToken, token)
			w.Header().Set(logging.TraceToken, token)
		} else {
			w.Header().Set(logging.TraceToken, r.Header.Get(logging.TraceToken))
		}
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware Logs details related to a request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := logging.LoggingResponseWriter{
			ResponseWriter:   w,
			ResponseMetadata: &logging.LoggingResponseMetadata{},
		}

		next.ServeHTTP(&lw, r)
		logging.LogResponse(
			w.Header().Get(logging.TraceToken),
			r,
			lw.ResponseMetadata,
			time.Since(start)*time.Millisecond)
	})
}
