package util

import (
	"net/http"
	"strings"
	"testing"
)

type MockRequestObject struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

type MockResponseWriter struct {
	statusCode int
}

func (mrw *MockResponseWriter) WriteHeader(statusCode int) {
	mrw.statusCode = statusCode
}

func (mrw *MockResponseWriter) Header() http.Header {
	return http.Header{}
}

func (mrw *MockResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func TestBindJSON(t *testing.T) {
	const validJSON = `
		{
			"id": 1,
			"message": "hello"
		}`
	w := &MockResponseWriter{}
	r, _ := http.NewRequest("", "", strings.NewReader(validJSON))

	expReq := MockRequestObject{
		ID:      1,
		Message: "hello",
	}

	var req MockRequestObject
	ok := BindJSON(w, r, &req)

	if !(ok && req.ID == expReq.ID && req.Message == expReq.Message) {
		t.Errorf("Expected %+v, got %+v", expReq, req)
	}

	const emptyJSON = ""
	r, _ = http.NewRequest("", "", strings.NewReader(emptyJSON))

	expReq = MockRequestObject{}

	req = MockRequestObject{}
	ok = BindJSON(w, r, &req)

	if !(!ok && req.ID == expReq.ID && req.Message == expReq.Message && w.statusCode == http.StatusBadRequest) {
		t.Errorf("Expected %+v, got %d: %+v", expReq, w.statusCode, req)
	}

	const invalidJSON = "{]"
	r, _ = http.NewRequest("", "", strings.NewReader(invalidJSON))

	expReq = MockRequestObject{}

	req = MockRequestObject{}
	ok = BindJSON(w, r, &req)

	if !(!ok && req.ID == expReq.ID && req.Message == expReq.Message && w.statusCode == http.StatusBadRequest) {
		t.Errorf("Expected %+v, got %d: %+v", expReq, w.statusCode, req)
	}

	const wrongTypesJSON = `
		{
			"id": "1",
			"message": "hello"
		}`
	r, _ = http.NewRequest("", "", strings.NewReader(wrongTypesJSON))

	expReq = MockRequestObject{
		Message: "hello",
	}

	req = MockRequestObject{}
	ok = BindJSON(w, r, &req)

	if !(!ok && req.ID == expReq.ID && req.Message == expReq.Message && w.statusCode == http.StatusBadRequest) {
		t.Errorf("Expected %+v, got %d: %+v", expReq, w.statusCode, req)
	}
}
