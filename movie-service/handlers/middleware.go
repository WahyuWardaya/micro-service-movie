package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "jwt secret not configured"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if uidRaw, ok := claims["user_id"]; ok {
			if f, ok2 := uidRaw.(float64); ok2 {
				c.Set("user_id", uint(f))
			} else {
				c.Set("user_id", 0)
			}
		} else {
			c.Set("user_id", 0)
		}

		c.Next()
	}
}
