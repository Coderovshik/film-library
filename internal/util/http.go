package util

import (
	"encoding/json"
	"errors"
	"io"
	"log"
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

func Unauthorized(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
}

func Forbidden(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}

func OK(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func BadRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func JSON(w http.ResponseWriter, r *http.Request, statusCode int, obj any) {
	w.WriteHeader(statusCode)
	w.Header().Set("content-type", "application/json")
	jsonBytes, _ := json.Marshal(obj)
	w.Write(jsonBytes)
}

func BindJSON(w http.ResponseWriter, r *http.Request, object any) bool {
	if err := json.NewDecoder(r.Body).Decode(object); err != nil {
		log.Printf("ERROR: failed to decode request body err=%s\n", err.Error())

		var sErr *json.SyntaxError
		if errors.As(err, &sErr) {
			BadRequest(w, r)
			w.Header().Set("content-type", "text/plain")
			w.Write([]byte("invalid json"))
			return false
		}

		if errors.Is(err, io.EOF) {
			BadRequest(w, r)
			w.Header().Set("content-type", "text/plain")
			w.Write([]byte("empty request"))
			return false
		}

		var utErr *json.UnmarshalTypeError
		if errors.As(err, &utErr) {
			BadRequest(w, r)
			w.Header().Set("content-type", "text/plain")
			w.Write([]byte("incorrect request typing"))
			return false
		}

		InternalServerError(w, r)
		return false
	}

	return true
}
