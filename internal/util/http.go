package util

import (
	"encoding/json"
	"net/http"
)

const (
	ErrorTypeValidation = "Validation"
	ErrorTypeConflict   = "Conflict"
)

type ErrorMessage struct {
	ErrorType string `json:"errorType"`
	Body      string `json:"body"`
}

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func OK(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func JSON(w http.ResponseWriter, r *http.Request, statusCode int, obj any) {
	w.WriteHeader(statusCode)
	w.Header().Set("content-type", "application/json")
	jsonBytes, _ := json.Marshal(obj)
	w.Write(jsonBytes)
}
