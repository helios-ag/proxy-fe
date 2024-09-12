Project uses devcontainer (can be opened with vscode)

Makefile available for tests and build
Inside dev container run `go run main.go`

Api available at `http://localhost:3000`

Available endpoints:
`GET /articles`

`GET /articles/{id:[0-9]+}`

`GET /authors`

a.Post("/track", a.handleRequest(rateLimitMiddleware(handler.TrackHandler)))

Check `internal/app.go` for some details

TODO:

* Proper DI
* Add e2e tests (i.e. cute)