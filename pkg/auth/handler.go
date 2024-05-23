package auth

import (
	"chat-app/pkg/user"
	"chat-app/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// LoginUserHandler handles user authentication and token generation upon successful login.
func LoginUserHandler(authService AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials UserCredentials
		var authenticatedUser *user.User

		// Bind request JSON to UserLogin struct
		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Check if required fields are empty
		if credentials.Password == "" || credentials.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
			return
		}

		// Authenticate user
		authenticatedUser, err := authService.LoginUser(c.Request.Context(), &credentials)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
			return
		}

		// compare if password is correct
		if !utils.ComparePasswords(authenticatedUser.Password, credentials.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
			return
		}

		// Check if user is valid
		if authenticatedUser.Validity != "valid" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You are not a valid user, " +
				"please contact the administrator"})
			return
		}

		// objectID from string to objectID

		objectID, err := primitive.ObjectIDFromHex(authenticatedUser.ID)

		// Generate JWT token
		token, err := utils.GenerateToken(&credentials.Username, &objectID, &authenticatedUser.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		// Set token in cookie
		c.SetCookie("token", token, 72*3600, "/", "localhost", false, true)
		// Set token in header
		c.Header("Authorization", token)

		// Send token as response
		c.JSON(http.StatusOK, gin.H{"status": true, "token": token, "message": "You're logged in!"})
	}
}

// LogoutUserHandler handles user logout by invalidating the token.
func LogoutUserHandler(authService AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve JWT token from header
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			return
		}

		// Verify and decode token
		claims, err := utils.VerifyToken(&token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Token", "message": err.Error()})
			return
		}

		// Call logout service
		err = authService.LogoutUser(c.Request.Context(), &claims.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to logout"})
			return
		}

		// Invalidate token by removing it from cookie and header
		c.SetCookie("token", "", -1, "/", "localhost", false, true)
		c.Header("Authorization", "")

		c.JSON(http.StatusOK, gin.H{"message": "You're logged out!"})
	}
}
