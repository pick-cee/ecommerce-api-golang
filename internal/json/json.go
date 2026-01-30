package json

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message,omitempty"`
}

func Write(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func Read(r *http.Request, data any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return  decoder.Decode(data)
}

func WriteError(w http.ResponseWriter, status int, err error, message string) {
	Write(w, status, ErrorResponse{
		Error: err.Error(),
		Message: message,
	})
}