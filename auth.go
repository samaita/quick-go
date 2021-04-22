package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	missingAuthorization = "Missing Authorization"
	invalidSession       = "Invalid Session"

	headerAuthorization = "Authorization"
	headerAppToken      = "App-Token"
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
		UID           string
		valid         bool
	)

	BearerToken = strings.Replace(c.Request.Header.Get(headerAuthorization), "Bearer ", "", -1)
	FirebaseToken = c.Request.Header.Get(headerAppToken)
	if BearerToken == "" && FirebaseToken == "" {
		return fmt.Errorf("%s", missingAuthorization)
	}

	if BearerToken != "" {
		if UID, err = verifyBearerToken(BearerToken); err != nil {
			log.Printf("[verifyLogin][verifyCustomToken] Input: %s Output: %v", BearerToken, err)
		}
		if UID == "" {
			return fmt.Errorf("%s", invalidSession)
		}
		c.Set(CtxUID, UID)
	}

	if FirebaseToken != "" {
		if valid, err = verifyFirebaseToken(FirebaseToken); err != nil || !valid {
			log.Printf("[verifyLogin][verifyFirebaseToken] Input: %s Output: %v", FirebaseToken, err)
			return fmt.Errorf("%s", invalidSession)
		}
	}

	return err
}

func verifyBearerToken(token string) (string, error) {
	var (
		err error
	)

	val, err := RedisClient.Get(context.Background(), token).Result()
	if err != nil {
		return val, err
	}

	return val, err
}
