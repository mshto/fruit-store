package middleware

import (
	"net/http"
)

// WithCORSMiddleware middleware to check product client and site IDs
type WithCORSMiddleware struct {
}

// NewWithCORSMiddleware new EntitlementMiddleware
func NewWithCORSMiddleware() *WithCORSMiddleware {
	return &WithCORSMiddleware{}
}

func (s *WithCORSMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}

	// Stop here for a Preflighted OPTIONS request.
	if r.Method == "OPTIONS" {
		return
	}

	next(w, r)
}
