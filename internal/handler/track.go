package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io"
	"net/http"
	"proxy/internal/cookies"
	"proxy/internal/util"
	"time"
)

type RecordRequest struct {
	ID int `json:"id"`
}

func TrackHandler(rdb *redis.Client, w http.ResponseWriter, r *http.Request) {
	uuid, err := cookies.Read(r, "userId")
	if err != nil {
		uuid = util.Uuid()
	}
	if err != nil {
		cookie := &http.Cookie{
			Name:     "userId",
			Value:    uuid,
			HttpOnly: true,
			Expires:  time.Now().Add(60 * 24 * time.Minute),
		}
		cookies.Write(w, *cookie)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Decode the JSON payload into the Item struct
	var request RecordRequest
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rdb.SAdd(context.TODO(), fmt.Sprintf("user:%s:articles", uuid), request.ID)

	w.WriteHeader(http.StatusCreated)
}
