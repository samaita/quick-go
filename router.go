package main

import (
	"github.com/gin-gonic/gin"
)

func initHandler() {
	router := gin.Default()

	router.POST("/api/register", handleRegister)
	router.POST("/api/login", handleLogin)

	router.Run(":3000")
}
