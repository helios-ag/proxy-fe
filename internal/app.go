package app

import (
	"fmt"
	"log"
	"net/http"
	"proxy/internal/config"
	"proxy/internal/handler"
	"proxy/internal/util"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"

	"github.com/dpapathanasiou/go-recaptcha"
)

const (
	// Public key for recapthca must be set on frontend
	recaptchaSecretKey = "your-recaptcha-secret-key"
	// values used
	limit = 1
	burst = 10
)

var limiter = rate.NewLimiter(limit, burst)

type App struct {
	Router *mux.Router
	RDB    *redis.Client
}

func (a *App) Initialize(config *config.Config) {
	recaptcha.Init(recaptchaSecretKey)

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			config.DB.Host,
			config.DB.Port,
		),
	})

	a.RDB = rdb
	a.Router = mux.NewRouter()
	a.setRouters()
}

func (a *App) setRouters() {
	a.Get("/articles", a.handleRequest(handler.ArticlesHandler))
	a.Get("/articles/{id:[0-9]+}", a.handleRequest(handler.ArticleHandler))
	a.Get("/authors", a.handleRequest(handler.AuthorArticlesHandler))
	a.Post("/track", a.handleRequest(rateLimitMiddleware(handler.TrackHandler)))
}

func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

type RequestHandlerFunction func(rdb *redis.Client, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.RDB, w, r)
	}
}

func rateLimitMiddleware(next RequestHandlerFunction) RequestHandlerFunction {
	return func(rdb *redis.Client, w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			recaptchaMiddleware(next)(rdb, w, r)
			return
		}
		next(rdb, w, r)
	}
}

func recaptchaMiddleware(next RequestHandlerFunction) RequestHandlerFunction {
	return func(rdb *redis.Client, w http.ResponseWriter, r *http.Request) {
		recaptchaResponse := r.FormValue("g-recaptcha-response")
		if recaptchaResponse == "" {
			util.Error(w, http.StatusBadRequest, "Missing recaptcha-response")
			return
		}

		remoteIP := r.RemoteAddr
		if _, err := recaptcha.Confirm(remoteIP, recaptchaResponse); err != nil {
			util.Error(w, http.StatusBadRequest, "Recaptcha rejected")
			return
		}

		next(rdb, w, r)
	}
}
