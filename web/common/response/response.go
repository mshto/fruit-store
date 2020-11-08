package response

import (
	"encoding/json"
	"net/http"
)

const (
	jsonContentType = "application/json; charset=utf-8"
)

// EmptyResp EmptyResp
type EmptyResp struct {
}

// errorMsg struct that is encoded and sent to the user. According to connectwise conventions
type errorMsg struct {
	Error string `json:"error"`
}

// RenderFailedResponse writes message and statusCode to response
func RenderFailedResponse(w http.ResponseWriter, code int, err error) {
	b, err := json.Marshal(errorMsg{
		Error: err.Error(),
	})
	if err != nil {
		render(w, http.StatusInternalServerError, nil)
		return
	}

	render(w, code, b)
}

// RenderResponse writes message and statusCode to response
func RenderResponse(w http.ResponseWriter, code int, body interface{}) {
	b, err := json.Marshal(body)
	if err != nil {
		render(w, http.StatusInternalServerError, nil)
		return
	}

	render(w, code, b)
}

func render(w http.ResponseWriter, code int, b []byte) {
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(code)
	_, _ = w.Write(b) // nolint
}
