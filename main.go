package main

import (
	app "proxy/internal"
	"proxy/internal/config"
)

func main() {
	appConfig := config.GetConfig()

	application := &app.App{}
	application.Initialize(appConfig)
	application.Run(":3000")
}
