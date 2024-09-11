package handler

import (
	"github.com/MicahParks/recaptcha"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
)

func CaptchaHandler(rdb *redis.Client, w http.ResponseWriter, r *http.Request) {
	verifier := recaptcha.NewVerifierV3("mySecret", recaptcha.VerifierV3Options{})

	http.MaxBytesReader(w, r.Body, 50000)

	err := r.ParseForm()
	if err != nil {
		log.Printf("Failed to parse form: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	frontendToken := r.Form.Get("g-recaptcha-response")
	remoteAddr := r.RemoteAddr

	response, err := verifier.Verify(r.Context(), frontendToken, remoteAddr)
	if err != nil {
		log.Printf("Failed to verify reCAPTCHA: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("reCAPTCHA V3 response: %#v", response)

	err = response.Check(recaptcha.V3ResponseCheckOptions{
		Action:   []string{"submit"},
		Hostname: []string{"example.com"},
		Score:    0.5,
	})
	if err != nil {
		log.Printf("Failed check for reCAPTCHA response: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
