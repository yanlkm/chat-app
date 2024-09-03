package code

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// createCodeHandler creates a new code.
func CreateCodeHandler(codeService CodeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newCode Code
		if err := c.ShouldBindJSON(&newCode); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		// Check if the code string is empty
		if newCode.Code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Code cannot be empty"})
			return
		}
		// Create the code
		if err := codeService.CreateCode(c.Request.Context(), &newCode); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create code"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Code created successfully"})
	}
}
