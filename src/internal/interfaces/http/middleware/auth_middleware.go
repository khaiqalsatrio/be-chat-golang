package middleware

import (
	"chat-golang/src/internal/infrastructure/services"
	"chat-golang/src/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtService *services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		tokenString := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Allow raw token if Bearer prefix is missing (for easier testing)
			tokenString = authHeader
		}

		if tokenString == "" {
			response.Error(c, http.StatusUnauthorized, "Token is empty")
			c.Abort()
			return
		}

		token, err := jwtService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}

		userID := claims["user_id"].(string)
		c.Set("user_id", userID)
		c.Next()
	}
}
