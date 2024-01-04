package main

import (
	"log"

	"github.com/samaita/quick-go/config"
	authHandler "github.com/samaita/quick-go/internal/auth/handler"
	healthHandler "github.com/samaita/quick-go/internal/health/handler"
	"github.com/samaita/quick-go/pkg/http"
)

var (
	cfg *config.Config
)

func main() {
	var (
		err error
	)

	// init config
	cfg, err = config.LoadConfig("")
	if err != nil {
		log.Fatalln(err)
	}

	moduleHealth := &healthHandler.HealthHandler{}
	moduleAuth := &authHandler.AuthHandler{}

	// Create a Gin router using the NewRouter function from the router package
	r := http.NewRouter(moduleHealth, moduleAuth)

	// Run the HTTP server
	http.RunHTTPServer(r, cfg.App.HTTP.Port)

}
