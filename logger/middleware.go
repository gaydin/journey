package logger

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"log/slog"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

const headerRequestID = "X-Request-ID"

type loggerKey struct{}

func Middleware(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(headerRequestID)
		if id == "" {
			id = uuid.New().String()
		}

		copyLogger := logger.With(
			slog.String("id", id),
			slog.String("url", r.URL.String()),
			slog.String("method", r.Method),
			slog.String("remote_ip", r.RemoteAddr),
		)
		r = r.WithContext(
			context.WithValue(r.Context(), loggerKey{}, copyLogger))

		rec := statusRecorder{w, 200}
		start := time.Now()
		next.ServeHTTP(&rec, r)
		since := time.Since(start).String()

		copyLogger = copyLogger.With(
			slog.Int("status", rec.status),
			slog.String("latency", since),
		)

		n := rec.status
		switch {
		case n >= 500:
			copyLogger.Error("server error")
		case n >= 400:
			copyLogger.Warn("Client error")
		case n >= 300:
			copyLogger.Info("Redirection")
		default:
			copyLogger.Info("Success")
		}
	})
}

func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}

func FromRequest(r *http.Request) *slog.Logger {
	return FromContext(r.Context())
}

func FromContext(ctx context.Context) *slog.Logger {
	return ctx.Value(loggerKey{}).(*slog.Logger)
}
