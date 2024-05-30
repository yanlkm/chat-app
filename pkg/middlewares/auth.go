package middlewares

import (
	"chat-app/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthMiddleware function to validate the token and authorize the user.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from cookie
		var token string
		token, err := c.Cookie("token")
		if err != nil {
			token = c.GetHeader("Authorization")
			if token == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "1 - Unauthorized"})
				c.Abort()
				return
			}
		}

		// Verify and decode token
		claims, err := utils.VerifyToken(&token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "2 - Unauthorized"})
			c.Abort()
			return
		}

		//check if user is logged in with id passed in params
		// convert claims ObjectId to string
		stringClaimsObjectId := claims.UserID.Hex()

		if c.Param("id") != "" {
			if c.Param("id") != stringClaimsObjectId {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "2 - Unauthorized"})
				c.Abort()
				return
			}

		}

		// TODO: Set user ID in context
		c.Set("userID", claims.UserID)

		// Continue with the request
		c.Next()
	}
}

// Check If User Is Logged In and If Is Admin
func IsAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from cookie
		var token string
		token, err := c.Cookie("token")
		if err != nil {
			token = c.GetHeader("Authorization")
			if token == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "1 - Unauthorized"})
				c.Abort()
				return
			}
		}

		// Verify and decode token
		claims, err := utils.VerifyToken(&token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "2 - Unauthorized"})
			c.Abort()
			return
		}

		//check if user is admin
		if claims.Role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "3 - Unauthorized"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		// Continue with the request
		c.Next()

	}
}

// Check If User Is Logged In
func IsLoggedInMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from cookie
		var token string
		token, err := c.Cookie("token")
		if err != nil {
			token = c.GetHeader("Authorization")
			if token == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "1 - Unauthorized"})
				c.Abort()
				return
			}
		}

		// Verify and decode token
		claims, err := utils.VerifyToken(&token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "2 - Unauthorized"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)

		// Continue with the request
		c.Next()
	}
}
