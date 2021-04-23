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

	if err = newCredential.grantAccess(); err != nil {
		log.Printf("[handleLogin][grantAccess] Input: %v, Output %v", credentialType, err.Error())
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
		newUser        User
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
		log.Printf("[handleRegister][register] Input: %v, Output %v", newCredential.Access, err.Error())
		APIResponseInternalServerError(c, nil, err.Error())
		return
	}

	newUser.FirstName = c.PostForm("first_name")
	newUser.LastName = c.PostForm("last_name")
	newUser.UID = newCredential.UID

	if newUser.FirstName == "" || newUser.LastName == "" {
		errMsg = fmt.Sprintf("FirstName or LastName empty")
		log.Printf("[handleRegister][Check First & LastName] Input: %v, Output %v", newCredential.UID, err)
		APIResponseBadRequest(c, nil, errMsg)
		return
	}

	if err = newUser.addUserInfo(userStatusRegistered); err != nil {
		log.Printf("[handleLogin][addUserInfo] Input: %v, Output %v", newCredential.Access, err.Error())
		APIResponseInternalServerError(c, nil, err.Error())
		return
	}

	if err = newCredential.grantAccess(); err != nil {
		log.Printf("[handleLogin][grantAccess] Input: %v, Output %v", newCredential.Access, err.Error())
		APIResponseInternalServerError(c, nil, err.Error())
		return
	}

	response := map[string]interface{}{
		"access_token": newCredential.Token,
		"uid":          newCredential.UID,
	}

	APIResponseOK(c, response)
}

func handleGetUserInfo(c *gin.Context) {
	var (
		err     error
		newUser User
	)

	newUser.UID = c.GetString(CtxUID)

	if err = newUser.getUserInfo(); err != nil {
		log.Printf("[handleGetUserInfo][getUserInfo] Input: %v, Output %v", newUser.UID, err.Error())
		APIResponseInternalServerError(c, nil, err.Error())
		return
	}

	response := map[string]interface{}{
		"user_info": newUser,
		"uid":       newUser.UID,
	}

	APIResponseOK(c, response)
}
