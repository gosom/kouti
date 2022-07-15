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

func IdentityFromRequestContext(r *http.Request) any {
	return r.Context().Value(Authenticated)
}

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

type Authorizator interface {
	Authorize(identity any, r *http.Request) error
}

func Authentication(authenticator Authenticator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			identity := 2
			/* TODO uncomment me
			identity, err := authenticator.Authenticate(r)
			if err != nil {
				ae := ErrHTTP{
					http.StatusUnauthorized,
					http.StatusText(http.StatusUnauthorized),
					nil,
				}
				renderJson(w, ae.StatusCode, ae)
				return
			}
			*/
			ctx := context.WithValue(r.Context(), Authenticated, identity)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Authorization(authorizator Authorizator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			identity := IdentityFromRequestContext(r)
			if err := authorizator.Authorize(identity, r); err != nil {
				ae := ErrHTTP{
					http.StatusForbidden,
					http.StatusText(http.StatusForbidden),
					nil,
				}
				renderJson(w, ae.StatusCode, ae)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
