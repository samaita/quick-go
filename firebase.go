package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

const (
	pathServiceAccount = "quick-go-example-service-account.json"
	projectID          = "quick-go-example"
)

var FirebaseApp *firebase.App

func initFirebase() {
	var (
		err error
	)

	opt := option.WithCredentialsFile(pathServiceAccount)
	FirebaseApp, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("[initFirebase][firebase.NewApp] %v", err)
	}
}

func generateCustomToken(userIdentification string) (string, error) {
	var (
		token string
		err   error
	)

	client, err := FirebaseApp.Auth(context.Background())
	if err != nil {
		return token, fmt.Errorf("[FirebaseApp.Auth] %v", err)
	}

	token, err = client.CustomToken(context.Background(), userIdentification)
	if err != nil {
		return token, fmt.Errorf("[client.CustomToken] %v", err)
	}

	return token, nil
}
