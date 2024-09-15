package util

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJSON(t *testing.T) {
	rr := httptest.NewRecorder()

	payload := map[string]string{"message": "hello"}

	JSON(rr, http.StatusOK, payload)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("controllers returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("controllers returned wrong content type: got %v want %v", contentType, "application/json")
	}

	expected, _ := json.Marshal(payload)
	if strings.TrimSpace(rr.Body.String()) != string(expected) {
		t.Errorf("controllers returned unexpected body: got %v want %v", rr.Body.String(), string(expected))
	}
}

func TestError(t *testing.T) {
	rr := httptest.NewRecorder()

	errorMessage := "something went wrong"

	Error(rr, http.StatusBadRequest, errorMessage)

	// Check the status code.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("controllers returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("controllers returned wrong content type: got %v want %v", contentType, "application/json")
	}

	expected, _ := json.Marshal(map[string]string{"error": errorMessage})
	if strings.TrimSpace(rr.Body.String()) != string(expected) {
		t.Errorf("controllers returned unexpected body: got %v want %v", rr.Body.String(), string(expected))
	}
}
