package middleware

import (
	"net/http"
	"sync"
)

// CacheStatusWriter wraps http.ResponseWriter to track cache status
type CacheStatusWriter struct {
	http.ResponseWriter
	cacheStatus string
	mu          sync.Mutex
}

// SetCacheStatus sets the cache status (HIT, MISS, SKIP)
func (w *CacheStatusWriter) SetCacheStatus(status string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.cacheStatus = status
}

// WriteHeader adds the X-Cache-Status header before writing the response
func (w *CacheStatusWriter) WriteHeader(statusCode int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.cacheStatus != "" {
		w.ResponseWriter.Header().Set("X-Cache-Status", w.cacheStatus)
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

// CacheStatusMiddleware wraps handlers to add cache status tracking
func CacheStatusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csw := &CacheStatusWriter{ResponseWriter: w, cacheStatus: "SKIP"}
		next.ServeHTTP(csw, r)
	})
}
