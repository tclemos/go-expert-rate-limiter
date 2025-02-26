package main

import (
	"github.com/tclemos/go-expert-rate-limiter/config"
	"github.com/tclemos/go-expert-rate-limiter/internal/webserver"
)

func main() {
	var webServerConfig config.WebServerConfig
	if err := config.LoadConfig(".env", &webServerConfig); err != nil {
		panic(err)
	}

	webserver.Start(webServerConfig)
}
