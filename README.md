Project uses devcontainer (can be opened with vscode)

Makefile available for tests and build
Inside dev container run `go run main.go`

Api available at `http://localhost:3000`

Available endpoints:

`GET /articles` - will return list of articles, if userId present in cookies, will return viewed articles 

`GET /articles?author=1` - will return list of articles by author with id 1

`GET /articles/{id:[0-9]+}` - will return single article

`GET /authors` - returns list of authors

`POST /track` with payload `{ id: [0-9]+}` - will record viewed article, endpoint protected with rate limiter and google recaptcha

Check `internal/app.go` for some details, like cache time, recaptcha key and etc



TODO:

* Proper DI
* Add more tests
* Use better rate limiting lib
* Add e2e tests (i.e. cute)