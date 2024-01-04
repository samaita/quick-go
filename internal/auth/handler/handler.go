// pkg/modules/auth/auth_handler.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler is an implementation of the Module interface for authentication
type AuthHandler struct{}

// SetupRoutes sets up authentication-related routes
func (a *AuthHandler) SetupRoutes(r *gin.Engine) {
	r.GET("/", HomeHandler)
}

// HomeHandler handles requests to the home route
func HomeHandler(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to the Home Page!")
}
