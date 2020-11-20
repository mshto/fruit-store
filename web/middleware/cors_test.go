package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestACORSMiddleware(t *testing.T) {
	type payload struct {
		origin string
		method string
	}

	tc := []struct {
		name string
		payload
	}{
		{
			name: "CORS middleware with success",
			payload: payload{
				origin: "origin",
				method: http.MethodGet,
			},
		},
		{
			name: "CORS middleware option with success",
			payload: payload{
				origin: "",
				method: http.MethodOptions,
			},
		},
	}

	for _, test := range tc {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req, _ := http.NewRequest(test.payload.method, "url", nil)
			req.Header.Set("Origin", test.payload.origin)

			rw := httptest.NewRecorder()

			NewWithCORSMiddleware().ServeHTTP(rw, req, func(rw http.ResponseWriter, r *http.Request) {})
		})
	}
}
