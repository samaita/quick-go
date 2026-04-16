package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionStore struct {
	client *redis.Client
}

func NewSessionStore(addr, password string, db int) *SessionStore {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &SessionStore{client: client}
}

func (s *SessionStore) Ping(ctx context.Context) error {
	return s.client.Ping(ctx).Err()
}

func sessionKey(jti string) string {
	return "session:" + jti
}

// Set stores the encrypted JWT string for the given JWT ID (jti) with a TTL.
func (s *SessionStore) Set(ctx context.Context, jti, encryptedJWT string, ttl time.Duration) error {
	if err := s.client.Set(ctx, sessionKey(jti), encryptedJWT, ttl).Err(); err != nil {
		return fmt.Errorf("session store: set %q: %w", jti, err)
	}
	return nil
}

// Get retrieves the encrypted JWT string for the given JWT ID.
// Returns ("", nil) if the key does not exist (session expired/revoked).
func (s *SessionStore) Get(ctx context.Context, jti string) (string, error) {
	val, err := s.client.Get(ctx, sessionKey(jti)).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("session store: get %q: %w", jti, err)
	}
	return val, nil
}

// Delete removes the session, effectively revoking it.
func (s *SessionStore) Delete(ctx context.Context, jti string) error {
	if err := s.client.Del(ctx, sessionKey(jti)).Err(); err != nil {
		return fmt.Errorf("session store: delete %q: %w", jti, err)
	}
	return nil
}
