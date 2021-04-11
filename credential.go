package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	credentialTypeEmail = 1

	credentialTypeAccessEmail = "Email"

	errNotRegistered    = "is not registered"
	errRegistered       = "is already registered"
	errInvalidPassword  = "password is not match"
	errPasswordTooShort = "password too short"
	errDefault          = "Please try again in a while"

	minimumPasswordLength = 8

	sizeHashSeed = 32
	costHash     = 8
)

type Credential struct {
	UID        string
	Type       int64
	Access     string
	Key        string
	StoredKey  string
	StoredSalt string
	Token      string

	typeAccess string
}

func (c *Credential) login() error {
	var (
		err                     error
		isRegistered, isCorrect bool
	)

	c.getTypeAccess()

	if isRegistered, err = c.isRegistered(); err != nil || !isRegistered {
		log.Printf("[login][isRegistered] Input: %s Output: %v", c.Access, err)
		return fmt.Errorf("%s %s", c.typeAccess, errNotRegistered)
	}

	if isCorrect, err = c.isCorrect(); err != nil || !isCorrect {
		log.Printf("[login][isCorrect] Input: %s Output: %v", c.Access, err)
		return fmt.Errorf("%s or %s", c.typeAccess, errInvalidPassword)
	}

	if err = c.getToken(); err != nil {
		log.Printf("[login][getToken] Input: %s Output: %v", c.Access, err)
		return err
	}

	return err
}

func (c *Credential) register() error {
	var (
		err          error
		isRegistered bool
	)

	c.getTypeAccess()

	if c.Type == credentialTypeEmail && len(c.Key) < minimumPasswordLength {
		return fmt.Errorf("%s", errPasswordTooShort)
	}

	if isRegistered, err = c.isRegistered(); (err != nil && err != sql.ErrNoRows) || isRegistered {
		log.Printf("[register][isRegistered] Input: %s Output: %v", c.Access, err)
		return fmt.Errorf("%s %s", c.typeAccess, errRegistered)
	}

	if err = c.addCredential(); err != nil {
		log.Printf("[register][addCredential] Input: %s Output: %v", c.Access, err)
		return fmt.Errorf("%s", errDefault)
	}

	if err = c.getToken(); err != nil {
		log.Printf("[register][getToken] Input: %s Output: %v", c.Access, err)
		return err
	}

	return err
}

func (c *Credential) isRegistered() (bool, error) {
	var (
		isExist bool
		err     error
		query   string
		isFound int64
	)

	query = `SELECT 1 FROM user WHERE credential_type = $1 AND credential_access = $2 LIMIT 1`
	if err = DB.QueryRowContext(context.Background(), query, c.Type, c.Access).Scan(&isFound); err != nil {
		log.Printf("[isRegistered][QueryRowContext] Input: %s Output: %v", c.Access, err)
		return isExist, err
	}

	isExist = isFound > 0

	return isExist, err
}

func (c *Credential) addCredential() error {
	var (
		err    error
		newKey []byte
		query  string
	)

	c.UID = generateUUID()
	c.StoredSalt = generateSalt()
	if newKey, err = bcrypt.GenerateFromPassword([]byte(c.Key+c.StoredSalt), costHash); err != nil {
		log.Printf("[addCredential][bcrypt.GenerateFromPassword] Input: %s Output: %v", c.Access, err)
		return err
	}
	c.StoredKey = string(newKey)

	query = `
		INSERT INTO user
		(uid, credential_type, credential_access, credential_key, credential_salt, create_time)
		VALUES
		($1, $2, $3, $4, $5, $6)
	`
	if _, err = DB.ExecContext(context.Background(), query, c.UID, c.Type, c.Access, c.StoredKey, c.StoredSalt, time.Now().Format(time.RFC3339)); err != nil {
		log.Printf("[addCredential][ExecContext] Input: %s Output: %v", c.Access, err)
		return err
	}

	return err
}

func (c *Credential) updateCredential() error {
	var (
		err error
	)

	return err
}

func (c *Credential) getToken() error {
	var (
		err error
	)

	if c.Token, err = generateCustomToken(c.Access); err != nil {
		log.Printf("[getToken][generateCustomToken] Input: %s Output: %v", c.Access, err)
		return err
	}

	return err
}

func (c *Credential) getTypeAccess() {
	switch c.Type {
	case credentialTypeEmail:
		c.typeAccess = credentialTypeAccessEmail
	default:
		c.typeAccess = credentialTypeAccessEmail
	}
}

func (c *Credential) getStoredCredential() error {
	var (
		err   error
		query string
	)

	query = `SELECT uid, credential_key, credential_salt FROM user WHERE credential_type = $1 AND credential_access = $2 LIMIT 1`
	if err = DB.QueryRowContext(context.Background(), query, c.Type, c.Access).Scan(&c.UID, &c.StoredKey, &c.StoredSalt); err != nil {
		log.Printf("[getStoredCredential][QueryRowContext] Input: %s Output: %v", c.Access, err)
		return err
	}

	return err
}

func generateSalt() string {
	var (
		b []byte
	)

	b = make([]byte, sizeHashSeed)
	rand.Read(b)

	return base64.URLEncoding.EncodeToString(b)
}

func generateUUID() string {
	return uuid.New().String()
}

func (c *Credential) isCorrect() (bool, error) {
	var (
		err error
	)

	if err = c.getStoredCredential(); err != nil {
		log.Printf("[isCorrect][getStoredCredential] Input: %s Output: %v", c.Access, err)
		return false, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(c.StoredKey), []byte(c.Key+c.StoredSalt)); err != nil {
		log.Printf("[isCorrect][bcrypt.CompareHashAndPassword] Input: %s Output: %v", c.Access, err)
		return false, err
	}

	return true, nil
}
