package middleware

import (
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/hlog"
)

// AccessLogger logs information about incoming HTTP requests.
func AccessLogger(r *http.Request, status, size int, dur time.Duration) {
	reqID := chimiddleware.GetReqID(r.Context())

	hlog.FromRequest(r).
		Info().
		Dur("duration_ms", dur).
		Str("host", r.Host).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("remote_addr", r.RemoteAddr).
		Str("request_id", reqID).
		Int("size", size).
		Int("status", status).
		Msg("handle request")
}
