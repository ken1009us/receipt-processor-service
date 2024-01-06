package util

import (
	"encoding/json"
	"net/http"
)


type APIError struct {
	Message string `json:"message"`
	Code int `json:"code"`
}


func WriteError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(APIError{Message: message, Code: code})
}