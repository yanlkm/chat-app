package user

import (
	"github.com/dchest/validator"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"regexp"
)

// CreateUserHandler creates a new user.
func CreateUserHandler(userService UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		//check if user has bad input
		if newUser.Username == "" || newUser.Password == "" || newUser.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
			return
		}

		//check if username is not too long
		if len(newUser.Username) > 10 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username is too long"})
			return
		}

		//check if username is not too short
		if len(newUser.Username) < 3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username is too short"})
			return
		}

		// check if username respects convention with a regex check
		usernameConvention := "^[a-zA-Z0-9_]*$"
		if re, _ := regexp.Compile(usernameConvention); !re.Match([]byte(newUser.Username)) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username"})
			return
		}

		// check if password has at least 6 characters
		if len(newUser.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must have at least 6 characters"})
			return
		}

		// check if email is valid
		if !validator.IsValidEmail(newUser.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
			return
		}
		// check if email is unique
		if err := userService.CheckEmail(c.Request.Context(), newUser.Email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}

		// check if username is unique
		if err := userService.CheckUsername(c.Request.Context(), newUser.Username); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}
		// create user
		if err := userService.CreateUser(c.Request.Context(), &newUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		c.JSON(http.StatusOK, newUser)
	}
}

// GetUserHandler retrieves a user by ID.
func GetUserHandler(userService UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting user informations"})
			return
		}

		user, err := userService.GetUser(c.Request.Context(), objectID)
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
		users, err := userService.GetAllUsers(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed displaying users"})
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

		userID := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(userID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting user informations"})
			return
		}

		var updatedUser UserUpdate
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

		if err := userService.UpdateUser(c.Request.Context(), objectID, updatedUser.Username); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
		c.JSON(http.StatusOK, updatedUser)
	}
}

// UpdatePasswordHandler updates the password for a user.
func UpdatePasswordHandler(userService UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting user informations"})
			return
		}
		var passwordUpdate PasswordUpdate
		if err := c.ShouldBindJSON(&passwordUpdate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// check if password has at least 6 characters
		if len(passwordUpdate.NewPassword) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must have at least 6 characters"})
			return
		}

		// update user password
		if err := userService.UpdatePassword(c.Request.Context(), objectID, passwordUpdate.NewPassword); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
	}
}

// DeleteUserHandler deletes a user by ID.
func DeleteUserHandler(userService UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting user informations"})
			return
		}
		if err := userService.DeleteUser(c.Request.Context(), objectID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user account"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
	}
}
