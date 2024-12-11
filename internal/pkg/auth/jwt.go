package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

type Claims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}

func AuthMiddleware(JwtSecretKey string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString, err := getToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Проверяем токен
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(JwtSecretKey), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token "+err.Error(), http.StatusUnauthorized)
			return
		}

		claims, _ := token.Claims.(*Claims)
		ctx := r.Context()
		ctx = context.WithValue(ctx, "login", claims.Login)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func getToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		tokenString := strings.Split(authHeader, " ")
		if len(tokenString) == 2 {
			return tokenString[1], nil
		}
	}

	cookie, err := r.Cookie("auth_token")
	if err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}

	return "", errors.New("authorization token is missing")
}
