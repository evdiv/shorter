package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

type contextKey string

const (
	UserIDKey  contextKey = "userID"
	TokenExp              = time.Hour * 3
	SecretKey             = "supersecretkey"
	CookieName            = "jwt"
)

func WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authRequired := r.RequestURI == "/api/user/urls" && r.Method == "GET"

		token := getJWTFromCookie(r)
		userID, err := getUserID(token)

		if err != nil {
			userID = generateUserID()
			token = buildJWTString(userID)

			// Set new JWT in cookie
			http.SetCookie(w, &http.Cookie{
				Name:    CookieName,
				Value:   token,
				Path:    "/",
				Expires: time.Now().Add(TokenExp),
			})
		}

		// If authentication is required and there's still no valid user ID, return 401
		if authRequired && userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getJWTFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return ""
	}

	if cookie == nil || cookie.Value == "" {
		return ""
	}
	return cookie.Value
}

func getUserID(tokenString string) (string, error) {
	if tokenString == "" {
		return "", errors.New("token is empty")
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SecretKey), nil
		})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}
	return claims.UserID, nil
}

func generateUserID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Println(err)
		return ""
	}
	return hex.EncodeToString(b)
}

// BuildJWTString - creates JWT token
func buildJWTString(userID string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return ""
	}
	return tokenString
}
