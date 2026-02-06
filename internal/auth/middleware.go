package auth

import (
	"errors"
	"net/http"
	"strings"

	"post/internal/entity"
	"post/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Middleware(jwtService JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header is missing", nil)
			c.Abort()
			return
		}

		// Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "Invalid authorization format", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwtService.ValidateToken(tokenString)

		var newToken string
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				// Try refresh
				refreshToken := c.GetHeader("X-Refresh-Token")
				if refreshToken == "" {
					response.Error(c, http.StatusUnauthorized, "Token expired and no refresh token provided", nil)
					c.Abort()
					return
				}

				refreshTok, refreshErr := jwtService.ValidateToken(refreshToken)
				if refreshErr != nil || !refreshTok.Valid {
					response.Error(c, http.StatusUnauthorized, "Invalid or expired refresh token", nil)
					c.Abort()
					return
				}

				claims, ok := refreshTok.Claims.(jwt.MapClaims)
				if !ok || claims["type"] != "refresh" {
					response.Error(c, http.StatusUnauthorized, "Invalid refresh token", nil)
					c.Abort()
					return
				}

				// Generate new access token
				userID, _ := claims["user_id"].(float64)
				role, _ := claims["role"].(float64)

				user := &entity.User{
					ID:   uint(userID),
					Role: entity.Role(role),
				}

				newToken, err = jwtService.GenerateToken(user)
				if err != nil {
					response.Error(c, http.StatusInternalServerError, "Failed to generate new token", nil)
					c.Abort()
					return
				}

				c.Header("X-New-Token", newToken)

				// Use refresh token claims for context since they are valid and contain the same info
				token = refreshTok
			} else {
				response.Error(c, http.StatusUnauthorized, "Invalid token", nil)
				c.Abort()
				return
			}
		} else if !token.Valid {
			response.Error(c, http.StatusUnauthorized, "Invalid token", nil)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Invalid token claims", nil)
			c.Abort()
			return
		}

		// Set userID and role in context
		// Note: JWT stores numbers as float64 in JSON
		userID, ok := claims["user_id"].(float64)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Invalid token claims: missing user_id", nil)
			c.Abort()
			return
		}
		c.Set("userID", uint(userID))

		if role, ok := claims["role"].(float64); ok {
			c.Set("role", entity.Role(role))
		}

		c.Next()
	}
}
