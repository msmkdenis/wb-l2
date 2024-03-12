package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// RequestLogger represents request logger middleware.
type (
	RequestLogger struct{}

	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// NewRequestLogger returns a new instance of RequestLogger.
func NewRequestLogger() *RequestLogger {
	return &RequestLogger{}
}

// Write implements the http.ResponseWriter interface.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader implements the http.ResponseWriter interface.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// RequestLogger logs each HTTP request.
func (r *RequestLogger) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI

		method := r.Method

		responseData := &responseData{}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		w = &lw

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		slog.Info("request_logger",
			slog.String("URI", uri),
			slog.String("method", method),
			slog.Duration("duration", duration),
			slog.Int("response_code", responseData.status),
			slog.Int("response_body_size", responseData.size))
	})
}
