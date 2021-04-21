package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	missingAuthorization = "Missing Authorization"
	invalidSession       = "Invalid Session"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err error
		)

		if err = verifyLogin(c); err != nil {
			APIResponseInternalServerError(c, nil, err.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}

func verifyLogin(c *gin.Context) error {
	var (
		err           error
		BearerToken   string
		FirebaseToken string
		valid         bool
	)

	BearerToken = strings.Replace(c.Request.Header.Get("Authorization"), "Bearer ", "", -1)
	FirebaseToken = c.Request.Header.Get("Firebase-Token")
	if BearerToken == "" && FirebaseToken == "" {
		return fmt.Errorf("%s", missingAuthorization)
	}

	if BearerToken != "" {
		if valid, err = verifyBearerToken(BearerToken); err != nil || !valid {
			log.Printf("[verifyLogin][verifyCustomToken] Input: %s Output: %v", BearerToken, err)
			return fmt.Errorf("%s", invalidSession)
		}
	}

	if FirebaseToken != "" {
		if valid, err = verifyFirebaseToken(FirebaseToken); err != nil || !valid {
			log.Printf("[verifyLogin][verifyFirebaseToken] Input: %s Output: %v", FirebaseToken, err)
			return fmt.Errorf("%s", invalidSession)
		}
	}

	return err
}

func verifyBearerToken(token string) (bool, error) {
	var (
		err error
	)

	// CHECK TO REDIS

	return true, err
}
