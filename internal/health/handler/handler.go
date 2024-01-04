// pkg/modules/auth/auth_handler.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler is an implementation of the Module interface for health check
type HealthHandler struct{}

// SetupRoutes sets up authentication-related routes
func (a *HealthHandler) SetupRoutes(r *gin.Engine) {
	r.GET("/health", ServerUpHandler)
}

// HomeHandler handles requests to the home route
func ServerUpHandler(c *gin.Context) {
	c.String(http.StatusOK, "Server Up!")
}
