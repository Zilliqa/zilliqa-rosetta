package server

import (
	"log"
	"net/http"
	"time"
)

// LoggerMiddleware is a simple logger middleware that prints the requests in
// an ad-hoc fashion to the stdlib's log.
func LoggerMiddleware(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}
