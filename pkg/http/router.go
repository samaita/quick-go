package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Module is an interface for application modules
type Module interface {
	SetupRoutes(r *gin.Engine)
}

// NewRouter initializes and returns a Gin router with the provided modules
func NewRouter(modules ...Module) *gin.Engine {
	r := gin.Default()

	for _, module := range modules {
		module.SetupRoutes(r)
	}

	return r
}

// RunHTTPServer sets up and starts the HTTP server
func RunHTTPServer(r *gin.Engine, port string) {
	addr := fmt.Sprintf("%s", port)
	fmt.Printf("Server is running on %s...\n", addr)
	r.Run(addr)
}
