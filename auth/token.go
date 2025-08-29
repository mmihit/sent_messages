package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"my_app/helper"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("mihit")

type JWTClaims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type contextKey string

const UserContextKey = contextKey("user")

type UserInfo struct {
	UserID int
	Role   string
}

func CreateToken(profile helper.Profile) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":        profile.Id,
			"user_name": profile.UserName,
			"role":      profile.Role,
			"exp":       time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		fmt.Println("invalid token: ", err)
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		fmt.Println("invalid token claims: ", err)
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				sendError(w, http.StatusUnauthorized, "missing or invalid authorization header")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := VerifyToken(tokenString)
			if err != nil {
				sendError(w, http.StatusUnauthorized, err.Error())
				return
			}

			hasPermission := slices.Contains(allowedRoles, claims.Role)

			if !hasPermission {
				sendError(w, http.StatusForbidden, "insufficient permissions")
				return
			}

			userInfo := UserInfo{UserID: claims.UserID, Role: claims.Role}
			ctx := context.WithValue(r.Context(), UserContextKey, userInfo)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func sendError(w http.ResponseWriter, code int, message string) {
	response := helper.ApiResponse{
		Code: code,
		Data: message,
	}
	response.Sent(w)
}
