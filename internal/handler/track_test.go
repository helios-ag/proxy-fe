package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestTrackHandler(t *testing.T) {
	db, mock := redismock.NewClientMock()

	// Expect the SAdd command to be called with specific arguments
	mock.ExpectSAdd("user:mock-uuid:articles", 123).SetVal(1)

	// Create a mock cookie reader and writer
	//cookies.Read = func(r *http.Request, name string) (string, error) {
	//	return "", http.ErrNoCookie
	//}
	//cookies.Write = func(w http.ResponseWriter, cookie http.Cookie) {
	//	http.SetCookie(w, &cookie)
	//}

	// Create a request payload
	requestPayload := RecordRequest{ID: 123}
	payloadBytes, _ := json.Marshal(requestPayload)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/track", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	TrackHandler(db, rr, req)

	// Check the status code
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Check if the cookie was set
	cookie := rr.Result().Cookies()[0]
	assert.Equal(t, "userId", cookie.Name)
	assert.Equal(t, "mock-uuid", cookie.Value)
	assert.True(t, cookie.HttpOnly)
	assert.WithinDuration(t, time.Now().Add(60*24*time.Minute), cookie.Expires, time.Minute)

	// Ensure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
