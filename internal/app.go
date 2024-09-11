package app

import (
	"fmt"
	"log"
	"net/http"
	"proxy/internal/config"
	"proxy/internal/handler"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	RDB    *redis.Client
}

func (a *App) Initialize(config *config.Config) {

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
	a.Post("/authors", a.handleRequest(handler.AuthorArticlesHandler))
	a.Post("/captcha", a.handleRequest(handler.CaptchaHandler))
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
