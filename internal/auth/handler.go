package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	cookieName  = "token"
	stateCookie = "oauth_state"
)

type Handler struct {
	oauthCfg     *oauth2.Config
	jwtSecret    string
	aesKey       string
	sessions     *SessionStore
	db           *sql.DB
	appURL       string
	cookieDomain string
	isProd       bool
}

func NewHandler(clientID, clientSecret, redirectURL, jwtSecret, aesKey, appURL, cookieDomain string, isProd bool, sessions *SessionStore, db *sql.DB) *Handler {
	return &Handler{
		oauthCfg: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
		jwtSecret:    jwtSecret,
		aesKey:       aesKey,
		sessions:     sessions,
		db:           db,
		appURL:       appURL,
		cookieDomain: cookieDomain,
		isProd:       isProd,
	}
}

// GoogleLogin redirects the user to the Google OAuth2 consent screen.
func (h *Handler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state, err := generateState()
	if err != nil {
		http.Error(w, "failed to generate state", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     stateCookie,
		Value:    state,
		Path:     "/",
		MaxAge:   300, // 5 minutes
		HttpOnly: true,
		Secure:   h.isProd,
		SameSite: http.SameSiteLaxMode, // Lax required for OAuth redirect flow
	})

	url := h.oauthCfg.AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallback handles the OAuth2 callback from Google.
func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Validate state parameter
	stateCookieVal, err := r.Cookie(stateCookie)
	if err != nil || stateCookieVal.Value != r.URL.Query().Get("state") {
		http.Redirect(w, r, h.appURL+"/?error=oauth_state_mismatch", http.StatusTemporaryRedirect)
		return
	}

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     stateCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.isProd,
	})

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, h.appURL+"/?error=oauth_no_code", http.StatusTemporaryRedirect)
		return
	}

	token, err := h.oauthCfg.Exchange(r.Context(), code)
	if err != nil {
		log.Printf("auth: oauth exchange error: %v", err)
		http.Redirect(w, r, h.appURL+"/?error=oauth_exchange_failed", http.StatusTemporaryRedirect)
		return
	}

	userInfo, err := fetchGoogleUserInfo(r.Context(), h.oauthCfg, token)
	if err != nil {
		log.Printf("auth: fetch user info error: %v", err)
		http.Redirect(w, r, h.appURL+"/?error=oauth_userinfo_failed", http.StatusTemporaryRedirect)
		return
	}

	if err := upsertUser(h.db, userInfo); err != nil {
		log.Printf("auth: upsert user error: %v", err)
		http.Redirect(w, r, h.appURL+"/?error=db_error", http.StatusTemporaryRedirect)
		return
	}

	tokenStr, jti, err := Issue(h.jwtSecret, userInfo.Email, userInfo.Name, userInfo.Picture)
	if err != nil {
		log.Printf("auth: issue jwt error: %v", err)
		http.Redirect(w, r, h.appURL+"/?error=jwt_error", http.StatusTemporaryRedirect)
		return
	}

	encrypted, err := Encrypt(h.aesKey, tokenStr)
	if err != nil {
		log.Printf("auth: encrypt jwt error: %v", err)
		http.Redirect(w, r, h.appURL+"/?error=encrypt_error", http.StatusTemporaryRedirect)
		return
	}

	if err := h.sessions.Set(r.Context(), jti, encrypted, TokenTTL); err != nil {
		// Non-fatal: fall back to JWT-only mode, log the warning
		log.Printf("auth: WARNING: redis unavailable, session not stored: %v", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    tokenStr,
		Path:     "/",
		MaxAge:   int(TokenTTL.Seconds()),
		HttpOnly: true,
		Secure:   h.isProd,
		SameSite: http.SameSiteStrictMode,
		Domain:   h.cookieDomain,
	})

	http.Redirect(w, r, h.appURL+"/", http.StatusTemporaryRedirect)
}

// Logout clears the session from Redis and removes the cookie.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieName)
	if err == nil {
		claims, err := Parse(h.jwtSecret, cookie.Value)
		if err == nil {
			if delErr := h.sessions.Delete(r.Context(), claims.ID); delErr != nil {
				log.Printf("auth: logout: delete session: %v", delErr)
			}
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.isProd,
		SameSite: http.SameSiteStrictMode,
		Domain:   h.cookieDomain,
	})

	http.Redirect(w, r, h.appURL+"/", http.StatusTemporaryRedirect)
}

// Me returns the current user's claims as JSON (used by frontend Alpine.js).
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(ClaimsKey{}).(*Claims)
	if !ok || claims == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"email":      claims.Email,
		"name":       claims.Name,
		"avatar_url": claims.AvatarURL,
	})
}

type googleUserInfo struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func fetchGoogleUserInfo(ctx context.Context, cfg *oauth2.Config, token *oauth2.Token) (*googleUserInfo, error) {
	client := cfg.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("fetch userinfo: %w", err)
	}
	defer resp.Body.Close()

	var info googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("decode userinfo: %w", err)
	}
	if info.Email == "" {
		return nil, fmt.Errorf("userinfo: missing email")
	}
	return &info, nil
}

func upsertUser(db *sql.DB, info *googleUserInfo) error {
	_, err := db.Exec(`
		INSERT INTO users (email, name, avatar_url, provider, updated_at)
		VALUES (?, ?, ?, 'google', datetime('now'))
		ON CONFLICT(email) DO UPDATE SET
			name       = excluded.name,
			avatar_url = excluded.avatar_url,
			updated_at = excluded.updated_at
	`, info.Email, info.Name, info.Picture)
	if err != nil {
		return fmt.Errorf("upsert user: %w", err)
	}
	return nil
}

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// ClaimsKey is the context key for storing JWT claims.
type ClaimsKey struct{}

// SessionTTLRemaining returns how long until the session expires.
func SessionTTLRemaining(claims *Claims) time.Duration {
	return time.Until(claims.ExpiresAt.Time)
}
