package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmizerany/assert"
)

func TestRenderFailedResponse(t *testing.T) {
	type expected struct {
		code int
	}

	type payload struct {
		statusCode int
		err        error
	}

	tc := []struct {
		name string
		payload
		expected
	}{
		{
			name: "Render failed response with success",
			payload: payload{
				statusCode: http.StatusBadRequest,
				err:        errors.New("not found"),
			},
			expected: expected{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, test := range tc {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			rw := httptest.NewRecorder()

			RenderFailedResponse(rw, test.payload.statusCode, test.payload.err)
			assert.Equal(t, test.expected.code, rw.Code)
		})
	}
}

func TestRenderResponse(t *testing.T) {
	type expected struct {
		code int
	}

	type payload struct {
		statusCode int
		body       interface{}
	}

	tc := []struct {
		name string
		payload
		expected
	}{
		{
			name: "Render failed response with success",
			payload: payload{
				statusCode: http.StatusOK,
				body:       "{}",
			},
			expected: expected{
				code: http.StatusOK,
			},
		},
		{
			name: "Render failed response with fail",
			payload: payload{
				statusCode: http.StatusOK,
				body:       make(chan int),
			},
			expected: expected{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, test := range tc {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			rw := httptest.NewRecorder()

			RenderResponse(rw, test.payload.statusCode, test.payload.body)
			assert.Equal(t, test.expected.code, rw.Code)
		})
	}
}
