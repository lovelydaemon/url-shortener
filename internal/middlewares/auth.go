package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lovelydaemon/url-shortener/internal/pkg/logger"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

const TOKEN_EXP = time.Minute * 30

func Authorization(key string, log logger.Interface) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("jwt")
			if err != nil {
				log.Info("Unauthorized user")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userID, err := getUserID(cookie.Value, key)
			if err != nil {
				log.Info("Unauthorized user")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", userID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func Authentication(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("jwt")
			if err == nil {
				userID, err := getUserID(cookie.Value, key)
				if err == nil {
					ctx := context.WithValue(r.Context(), "userID", userID)
					r = r.WithContext(ctx)
					next.ServeHTTP(w, r)
					return
				}
			}

			cookie, userID, err := newJWTCookie(key)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			http.SetCookie(w, cookie)
			ctx := context.WithValue(r.Context(), "userID", userID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func newJWTCookie(key string) (*http.Cookie, uuid.UUID, error) {
	userID := uuid.New()
	tokenString, err := buildJWTString(userID, key)
	if err != nil {
		return nil, uuid.UUID{}, err
	}

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		HttpOnly: true,
	}
	return cookie, userID, nil
}

func getUserID(tokenString string, key string) (uuid.UUID, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return []byte(key), nil
	})
	if err != nil || !token.Valid {
		return uuid.UUID{}, err
	}

	return claims.UserID, nil
}

func buildJWTString(userID uuid.UUID, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: userID,
	})
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
