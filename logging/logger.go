package logging

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
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
var commitSha string
var scmLink string

type loggerBuilder struct {
	version       string
	project       string
	service       string
	scmLinkToRepo string
	logLevel      zapcore.Level
}

func NewLoggerBuilder() *loggerBuilder {
	return &loggerBuilder{
		logLevel: zapcore.DebugLevel,
	}
}

func (b *loggerBuilder) WithVersion(v string) *loggerBuilder {
	b.version = v
	return b
}

func (b *loggerBuilder) WithProject(p string) *loggerBuilder {
	b.project = p
	return b
}

func (b *loggerBuilder) WithService(s string) *loggerBuilder {
	b.service = s
	return b
}

func (b *loggerBuilder) WithScmLinkToRepo(s string) *loggerBuilder {
	b.scmLinkToRepo = s
	return b
}

func (b *loggerBuilder) WithLogLevel(l zapcore.Level) *loggerBuilder {
	b.logLevel = l
	return b
}

func (b *loggerBuilder) Build() {
	commitSha = b.version
	scmLink = b.scmLinkToRepo
	initLogger(b.logLevel, b.project, b.service)
}

func initLogger(level zapcore.Level, project, service string) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(level)
	zapConfig.Encoding = "json"
	zapConfig.EncoderConfig.TimeKey = "time"
	zapConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	zapConfig.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	zapConfig.EncoderConfig.MessageKey = "message"
	zapConfig.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	zapConfig.DisableStacktrace = true
	zapConfig.OutputPaths = []string{"stdout"}
	zapConfig.EncoderConfig.EncodeCaller = remoteSourceCallerEncoder

	zapLogger, err := zapConfig.Build(zap.Fields(
		zap.String("project", project),
		zap.String("service", service),
		zap.String("version", commitSha),
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

func remoteSourceCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	var link string
	if link = linkOrEmpty(caller.TrimmedPath()); len(link) == 0 {
		zapcore.ShortCallerEncoder(caller, enc)
	} else {
		enc.AppendString(link + "#L" + strconv.Itoa(caller.Line))
	}
}

func linkOrEmpty(input string) string {
	if scmLink == "" || input == "" {
		return ""
	}
	withoutRepoName := strings.Join(strings.Split(input, "/")[1:], "/")
	withoutLineNum := strings.Split(withoutRepoName, ":")[0]

	var linkLocation string
	if commitSha == "" {
		linkLocation = "HEAD"
	} else {
		linkLocation = commitSha
	}
	return scmLink + "/blob/" + linkLocation + "/" + withoutLineNum
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
