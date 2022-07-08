package web

import (
	"bytes"
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func RequestLogger(logger zerolog.Logger) func(http.Handler) http.Handler {
	var bufferPool = sync.Pool{
		New: func() interface{} { return new(bytes.Buffer) },
	}
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(rw, r.ProtoMajor)
			start := time.Now()
			buf := bufferPool.Get().(*bytes.Buffer)
			buf.Reset()
			defer bufferPool.Put(buf)
			ww.Tee(buf)
			defer func() {
				statusCode := ww.Status()
				var ev *zerolog.Event
				msg := http.StatusText(statusCode)
				if statusCode >= 200 && statusCode < 400 {
					ev = logger.Info()
				} else if statusCode >= 400 && statusCode < 500 {
					ev = logger.Warn()
				} else {
					ev = logger.Error()
				}
				ev = ev.
					Str("request-id", middleware.GetReqID(r.Context())).
					Int("status", statusCode).
					Int("bytes", ww.BytesWritten()).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("query", r.URL.RawQuery).
					Str("ip", r.RemoteAddr).
					Str("user-agent", r.UserAgent()).
					Dur("latency", time.Since(start))
				//if statusCode < 200 || statusCode >= 400 {
				//	ev = ev.RawJSON("body", buf.Bytes())
				//}
				ev.Msg(msg)
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

type ctxKey int

const Authenticated ctxKey = ctxKey(0)

type Authenticator interface {
	Authenticate(r *http.Request) (any, error)
}

func Authentication(authenticator Authenticator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			identity, err := authenticator.Authenticate(r)
			if err != nil {
				ae := ErrHTTP{
					http.StatusUnauthorized,
					http.StatusText(http.StatusUnauthorized),
				}
				renderJson(w, ae.StatusCode, ae)
				return
			}
			ctx := context.WithValue(r.Context(), Authenticated, identity)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

/*
func getBearerTokenFromHeader(r *http.Request) (string, error) {
	const (
		headerAuthorization = "Authorization"
		headerPrefixBearer  = "BEARER"
	)
	bearer := r.Header.Get(headerAuthorization)
	size := len(headerPrefixBearer) + 1
	if len(bearer) > size && strings.ToUpper(bearer[0:size-1]) == headerPrefixBearer {
		return bearer[size:], nil
	}
	return "", err
}
*/
