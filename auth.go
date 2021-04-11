package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func handleLogin(c *gin.Context) {
	var (
		newCredential  Credential
		credentialType string
		err            error
		errMsg         string
	)

	newCredential.Key = c.PostForm("credential_key")
	newCredential.Access = c.PostForm("credential_access")
	credentialType = c.DefaultPostForm("credential_type", fmt.Sprint(credentialTypeEmail))

	if newCredential.Type, err = strconv.ParseInt(credentialType, 10, 64); err != nil {
		errMsg = fmt.Sprintf("Credential Type '%s' is invalid", credentialType)
		log.Printf("[handleLogin][Parse Credential] Input: %v, Output %v", credentialType, err)
		APIResponseBadRequest(c, nil, errMsg)
		return
	}

	if err = newCredential.login(); err != nil {
		log.Printf("[handleLogin][login] Input: %v, Output %v", credentialType, err.Error())
		APIResponseInternalServerError(c, nil, err.Error())
		return
	}

	response := map[string]interface{}{
		"access_token": newCredential.Token,
		"uid":          newCredential.UID,
	}

	APIResponseOK(c, response)
}

func handleRegister(c *gin.Context) {
	var (
		newCredential  Credential
		credentialType string
		err            error
		errMsg         string
	)

	newCredential.Key = c.PostForm("credential_key")
	newCredential.Access = c.PostForm("credential_access")
	credentialType = c.DefaultPostForm("credential_type", fmt.Sprint(credentialTypeEmail))

	if newCredential.Type, err = strconv.ParseInt(credentialType, 10, 64); err != nil {
		errMsg = fmt.Sprintf("Credential Type '%s' is invalid", credentialType)
		log.Printf("[handleRegister][Parse Credential] Input: %v, Output %v", credentialType, err)
		APIResponseBadRequest(c, nil, errMsg)
		return
	}

	if err = newCredential.register(); err != nil {
		log.Printf("[handleRegister][register] Input: %v, Output %v", credentialType, err.Error())
		APIResponseInternalServerError(c, nil, err.Error())
		return
	}

	response := map[string]interface{}{
		"access_token": newCredential.Token,
		"uid":          newCredential.UID,
	}

	APIResponseOK(c, response)
}
