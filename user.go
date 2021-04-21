package main

import (
	"context"
	"log"
)

const (
	userStatusDeactivate = -1
	userStatusRegistered = 0
	userStatusActive     = 1
)

type User struct {
	UID       string `json:"uid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Thumbnail string `json:"thumbnail"`
}

func (u *User) getUserInfo() error {
	var (
		err   error
		query string
	)

	query = `SELECT first_name, last_name, thumbnail FROM user_info WHERE uid = $1 AND status = $2`
	if err = DB.QueryRowContext(context.Background(), query, u.UID, userStatusActive).Scan(&u.FirstName, u.LastName, u.Thumbnail); err != nil {
		log.Printf("[getUserInfo][QueryRowContext] Input: %s Output: %v", u.UID, err)
		return err
	}

	return err
}
