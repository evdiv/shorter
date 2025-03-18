package middleware

import (
	"context"
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

		// If no JWT in the Cookie return 401 / Unauthorized
		token := getJWTFromCookie(r)
		if authRequired && token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		//if no created links return 204 / No Content
		userID, err := getUserID(token)

		if err != nil {
			userID = generateUserID()
			token = buildJWTString(userID)

			log.Println("Trying to create a token")

			if token == "" {
				http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
				return
			}
		}

		log.Println("get userID from token. userID: ", userID)

		http.SetCookie(w, &http.Cookie{
			Name:  "jwt",
			Value: token,
			Path:  "/",
		})

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getJWTFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("jwt")
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

func generateUserID() string {
	b := make([]byte, 16)
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
