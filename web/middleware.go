package web

import (
	"bytes"
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
