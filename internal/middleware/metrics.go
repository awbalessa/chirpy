package middleware

import (
	"net/http"
	"sync/atomic"
)

func MetricsInc(counter *atomic.Int32, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counter.Add(1)
		next.ServeHTTP(w, r)
	})
}
