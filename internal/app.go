package app

import (
	"fmt"
	"log"
	"net/http"
	"proxy/api/authors"
	"proxy/api/posts"
	"proxy/internal/config"
	"proxy/internal/controllers"
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
}

func (a *App) Initialize(config *config.Config) {
	recaptcha.Init(recaptchaSecretKey)

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			config.DB.Host,
			config.DB.Port,
		),
	})
	pc := &posts.Client{
		HttpClient: &http.Client{},
		PostsUrl:   config.PostsUrl,
	}
	ac := &authors.Client{
		HttpClient: &http.Client{},
		AuthorsUrl: config.UsersUrl,
	}
	a.Router = mux.NewRouter()
	hf := controllers.NewControllerFactory(rdb, pc, ac)
	a.setRouters(hf)
}

func (a *App) setRouters(hf *controllers.ControllerFactory) {
	a.Get("/articles", a.handleRequest(hf.GetArticleController().GetArticles))
	a.Get("/articles/{id:[0-9]+}", a.handleRequest(hf.GetArticleController().GetArticle))
	a.Get("/authors", a.handleRequest(hf.GetAuthorArticlesController().GetAuthorArticles))
	a.Post("/track", a.handleRequest(rateLimitMiddleware(hf.GetTrackController().PostTrack)))
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

type RequestHandlerFunction func(w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}
}

func rateLimitMiddleware(next RequestHandlerFunction) RequestHandlerFunction {
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			recaptchaMiddleware(next)(w, r)
			return
		}
		next(w, r)
	}
}

func recaptchaMiddleware(next RequestHandlerFunction) RequestHandlerFunction {
	return func(w http.ResponseWriter, r *http.Request) {
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

		next(w, r)
	}
}
