package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"todo_with_auth/internals/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(config *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header found"})
			ctx.Abort()
			return
		}

		trimmedToken := strings.TrimPrefix(authHeader, "Bearer ")
		if trimmedToken == "" || trimmedToken == authHeader {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Malformed Token"})
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(trimmedToken, func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				ctx.Abort()
				return nil, fmt.Errorf("Unexpected Signing Method: %v", token.Header["alg"])
			}
			return []byte(config.JWTSecret), nil
		})
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}

		userId, ok := claims["user_id"].(string)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}

		expiration, ok := claims["exp"].(float64)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		} else {
			expirationTime := time.Unix(int64(expiration), 0)
			if time.Now().After(expirationTime) {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired!"})
				ctx.Abort()
				return
			}
		}

		ctx.Set("user_id", userId)
		ctx.Next()
	}
}
