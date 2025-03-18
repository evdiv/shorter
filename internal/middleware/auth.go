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

const UserIDKey contextKey = "userID"
const TokenExp = time.Hour * 3
const SecretKey = "supersecretkey"

func WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authRequired := false
		if r.RequestURI == "/api/user/urls" && r.Method == "GET" {
			authRequired = true
		}

		token, err := getJWTFromCookie(r)

		if authRequired && err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userID, err := getUserID(token)
		if authRequired && err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if userID == "" {
			userID, err = generateUserID()
			if err != nil {
				http.Error(w, "Failed to generate UserID", http.StatusInternalServerError)
				return
			}
			token, err = buildJWTString(userID)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:  "jwt",
				Value: token,
				Path:  "/",
			})
		}

		log.Println("WithAuth. userID: ", userID)

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getJWTFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return "", err
	}

	if cookie == nil || cookie.Value == "" {
		return "", errors.New("cookie is empty")
	}

	return cookie.Value, nil
}

func getUserID(tokenString string) (string, error) {
	if tokenString == "" {
		log.Println("getUserID: Token is empty")
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

	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("invalid token")
	}
	return claims.UserID, nil
}

func generateUserID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// BuildJWTString - creates JWT token
func buildJWTString(userID string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
