package logging

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	ENV_NAME = "APP_ENV"
	ENV_DEV  = "dev"
)

const (
	// #nosec G101 -- Not a secret
	TraceToken = "X-Trace-Token"
)

var globalLogger *zap.Logger

func InitLogger(version, project, service string) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	zapConfig.Encoding = "json"
	zapConfig.EncoderConfig.TimeKey = "time"
	zapConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	zapConfig.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	zapConfig.EncoderConfig.MessageKey = "message"
	zapConfig.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	zapConfig.DisableStacktrace = true
	zapConfig.OutputPaths = []string{"stdout"}

	zapLogger, err := zapConfig.Build(zap.Fields(
		zap.String("project", project),
		zap.String("service", service),
		zap.String("version", version),
		zap.String("environment", getEnv()),
	), zap.AddCallerSkip(1))

	if err != nil {
		fmt.Println("failed to initialize logger:", err)
		os.Exit(1)
	}
	defer zapLogger.Sync()

	globalLogger = zapLogger
}

func getEnv() string {
	env := os.Getenv(ENV_NAME)
	if env == "" {
		env = ENV_DEV
	}
	return env
}

func withTrace(token string) *zap.Logger {
	return globalLogger.With(
		zap.String(TraceToken, token))
}

func LogResponse(token string, req *http.Request, lw *LoggingResponseMetadata, duration time.Duration) {
	withTrace(token).With(
		zap.String("http_method", req.Method),
		zap.String("http_host", req.Host),
		zap.String("http_path", req.URL.String()),
		zap.String("http_query", req.URL.Query().Encode()),
		zap.String("http_remote", req.RemoteAddr),
		zap.Any("http_headers", headersToArray(req.Header.Clone())),
		zap.Duration("http_duration", duration),
		zap.Int("http_status", lw.Status),
		zap.Int("http_size", lw.Size),
	).Info("finished request")
}

func LogDebug(token, msg string) {
	withTrace(token).Debug(msg)
}

func LogInfo(token, msg string) {
	withTrace(token).Info(msg)
}

func LogWarn(token, msg string) {
	withTrace(token).Warn(msg)
}

func LogError(token, msg string) {
	withTrace(token).Error(msg)
}

func headersToArray(h http.Header) map[string][]string {
	h.Del(TraceToken)
	hasher := sha512.New()
	hasher.Write([]byte(h.Get("Authorization")))
	h.Set("Authorization-Hash", base64.URLEncoding.EncodeToString(hasher.Sum(nil)))
	h.Set("Authorization-Len", strconv.Itoa(len(h.Get("Authorization"))))
	h.Del("Authorization")
	return h
}

type (
	// LoggingResponseMetadata wrapper for holding response metadata
	LoggingResponseMetadata struct {
		Status int
		Size   int
	}

	// LoggingResponseWriter wrapper around the default response writer to enable
	// the logging of the response metadata
	LoggingResponseWriter struct {
		http.ResponseWriter
		ResponseMetadata *LoggingResponseMetadata
	}
)

func (w *LoggingResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.ResponseMetadata.Size = size
	return size, err
}

func (w *LoggingResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.ResponseMetadata.Status = statusCode
}
