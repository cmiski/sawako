package middleware

import (
	"log"
	"net/http"
	"time"
)

// custom response writer wrapper to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// middleware function to log incoming requests and their latency
func Logging(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // default status code, will be updated if WriteHeader is called
		}

		next.ServeHTTP(rw, r)

		latency := time.Since(start)

		log.Printf(
			"method=%s path=%s status=%d latency=%s",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			latency,
		)
	})

}
