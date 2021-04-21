package main

import (
	"github.com/gin-gonic/gin"
)

func initHandler() {
	router := gin.Default()

	router.POST("/api/register", handleRegister)
	router.POST("/api/login", handleLogin)

	authorized := router.Group("/api")
	authorized.Use(AuthRequired())
	{
		authorized.GET("/user/info", handleGetUserInfo)
	}

	router.Run(":3000")
}
