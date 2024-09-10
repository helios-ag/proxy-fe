package main

import (
	app "proxy/internal"
	"proxy/internal/config"
	"time"
)

const (
	cacheAuthorArticles   = 2 * time.Minute
	detailedArticledCache = 5 * time.Minute
	cacheAuthorsList      = 5 * time.Minute
)

func main() {
	appConfig := config.GetConfig()

	application := &app.App{}
	application.Initialize(appConfig)
	application.Run(":3000")
}
