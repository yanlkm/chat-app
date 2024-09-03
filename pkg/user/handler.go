package user

import (
	"chat-app/pkg/code"
	"chat-app/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

// CreateUserHandler creates a new user.
func CreateUserHandler(userService UserService, codeService code.CodeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser UserEntity
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Check if user has bad input
		if newUser.Username == "" || newUser.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
			return
		}

		// Check if user has a role
		if newUser.Role == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
			return
		}

		// check if user role is user or admin, if not put user as default
		if newUser.Role != "user" && newUser.Role != "admin" {
			newUser.Role = "user"
		}

		// Check if user validity is valid, if not put invalid as default
		if newUser.Validity != "valid" {
			newUser.Validity = "valid"
		}

		// Check if username is not too long
		if len(newUser.Username) > 10 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username is too long"})
			return
		}

		// Check if username is not too short
		if len(newUser.Username) < 3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username is too short"})
			return
		}

		// Check if username respects convention with a regex check
		usernameConvention := "^[a-zA-Z0-9_]*$"
		if re, _ := regexp.Compile(usernameConvention); !re.Match([]byte(newUser.Username)) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username"})
			return
		}

		// Check if password has at least 6 characters
		if len(newUser.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must have at least 6 characters"})
			return
		}

		// Hash the password
		hashedPassword, err := utils.HashPassword(newUser.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please, change your password"})
			return
		}
		newUser.Password = hashedPassword

		// Check if username is unique
		if err := userService.CheckUsername(c.Request.Context(), newUser.Username); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}
		// Check if the provided code is valid
		if newUser.Code != "" {
			isValid, err := codeService.CheckCode(c.Request.Context(), &newUser.Code)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to check code"})
				return
			}
			if !isValid {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
				return
			}
			// Mark the code as used
			if err := codeService.UpdateCode(c.Request.Context(), &newUser.Code); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update code"})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Code is required"})
			return
		}

		// Create user
		if err := userService.CreateUser(c.Request.Context(), &newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
			return
		}
		// Remove password from the response
		newUser.Password = ""
		c.JSON(http.StatusOK, newUser)
	}
}

// GetUserHandler retrieves a user by ID.
func GetUserHandler(userService UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")

		user, err := userService.GetUser(c.Request.Context(), userID)
		// user without password
		user.Password = ""
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Error getting user informations"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

// GetUsersHandler retrieves all users.
func GetUsersHandler(userService UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get all users
		users, err := userService.GetAllUsers(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed displaying users"})
			return
		}
		// remove password from all users
		for i := range users {
			users[i].Password = ""
		}
		c.JSON(http.StatusOK, users)
	}

}

// UpdateUserHandler updates an existing user.
func UpdateUserHandler(userService UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get user id
		userID := c.Param("id")
		// entity to store user update
		var updatedUser UserUpdateEntity
		if err := c.ShouldBindJSON(&updatedUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		} //check if username is not too long
		if len(updatedUser.Username) > 10 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username is too long"})
			return
		}

		//check if username is not too short
		if len(updatedUser.Username) < 3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username is too short"})
			return
		}

		// check if username respects convention with a regex check
		usernameConvention := "^[a-zA-Z0-9_]*$"
		if re, _ := regexp.Compile(usernameConvention); !re.Match([]byte(updatedUser.Username)) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username"})
			return
		}
		// update user
		if err := userService.UpdateUser(c.Request.Context(), userID, updatedUser.Username); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update user"})
			return
		}
		c.JSON(http.StatusOK, updatedUser)
	}
}

// UpdatePasswordHandler updates the password for a user.
func UpdatePasswordHandler(userService UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get user id
		userID := c.Param("id")

		// entities to store user and password update
		var user *UserEntity
		var passwordUpdate PasswordUpdateEntity

		// check if password update is valid
		if err := c.ShouldBindJSON(&passwordUpdate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// check if password has at least 6 characters
		if len(passwordUpdate.NewPassword) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must have at least 6 characters"})
			return
		}

		// Hash the password
		hashedPassword, err := utils.HashPassword(passwordUpdate.NewPassword)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please, change your new password"})
			return
		}
		// update password
		passwordUpdate.NewPassword = hashedPassword

		// check if old password is correct
		user, err = userService.GetUser(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update password"})
			return
		}
		if !utils.ComparePasswords(user.Password, passwordUpdate.OldPassword) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
			return
		}

		// update user password
		if err := userService.UpdatePassword(c.Request.Context(), userID, passwordUpdate.NewPassword); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update password"})
			return
		}
		// remove password from the response
		c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
	}
}

// BanUserHandler bans a user by ID.
func BanUserHandler(userService UserService) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get user id : bannerID and bannedID
		var bannerID string
		var bannedID string

		// get user id : bannerID and bannedID
		bannerID = c.Param("id")
		bannedID = c.Param("idBanned")

		// ban user
		if err := userService.BanUser(c.Request.Context(), bannerID, bannedID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to ban user"})
			return
		}
		// remove password from the response
		c.JSON(http.StatusOK, gin.H{"message": "User banned successfully"})
	}
}

// UnBanUserHandler unbans a user by ID.
func UnBanUserHandler(userService UserService) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get user id : bannerID and bannedID
		var bannerID string
		var bannedID string
		bannerID = c.Param("id")
		bannedID = c.Param("idBanned")
		// unban user
		if err := userService.UnBanUser(c.Request.Context(), bannerID, bannedID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to unban user"})
			return
		}
		// remove password from the response
		c.JSON(http.StatusOK, gin.H{"message": "User unbanned successfully"})

	}
}

// DeleteUserHandler deletes a user by ID.
func DeleteUserHandler(userService UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get user id
		userID := c.Param("id")
		// delete user
		if err := userService.DeleteUser(c.Request.Context(), userID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete user account"})
			return
		}
		// remove password from the response
		c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
	}
}
