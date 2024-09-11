package cookies

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWrite(t *testing.T) {
	tests := []struct {
		name        string
		cookie      http.Cookie
		expectError error
	}{
		{
			name: "Valid cookie",
			cookie: http.Cookie{
				Name:  "test",
				Value: "testValue",
			},
			expectError: nil,
		},
		{
			name: "Cookie value too long",
			cookie: http.Cookie{
				Name:  "test",
				Value: string(make([]byte, 4096)), // Create a string with 4096 bytes
			},
			expectError: ErrValueTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := Write(w, tt.cookie)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("Expected error %v, got %v", tt.expectError, err)
			}
		})
	}
}

func TestRead(t *testing.T) {
	tests := []struct {
		name        string
		cookie      *http.Cookie
		expectValue string
		expectError error
	}{
		{
			name: "Valid cookie",
			cookie: &http.Cookie{
				Name:  "test",
				Value: base64.URLEncoding.EncodeToString([]byte("testValue")),
			},
			expectValue: "testValue",
			expectError: nil,
		},
		{
			name:        "No cookie",
			cookie:      nil,
			expectValue: "",
			expectError: http.ErrNoCookie,
		},
		{
			name: "Invalid base64 value",
			cookie: &http.Cookie{
				Name:  "test",
				Value: "invalidBase64",
			},
			expectValue: "",
			expectError: ErrInvalidValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.cookie != nil {
				r.AddCookie(tt.cookie)
			}

			value, err := Read(r, "test")
			if value != tt.expectValue {
				t.Errorf("Expected value %s, got %s", tt.expectValue, value)
			}
			if !errors.Is(err, tt.expectError) {
				t.Errorf("Expected error %v, got %v", tt.expectError, err)
			}
		})
	}
}
