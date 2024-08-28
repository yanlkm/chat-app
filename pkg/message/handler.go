package message

import (
	"chat-app/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func CreateMessageHandler(messageService MessageService) gin.HandlerFunc {
	return func(c *gin.Context) {

		var message Message
		if err := c.BindJSON(&message); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// check if username and content are empty
		if message.RoomID == "" || message.Username == "" || message.Content == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Message error"})
			return
		}

		// check if userConnected is the one who sends a creation message request
		_, usernameConnected, errConnection := utils.GetUserIDAndUsernameFromContext(c)
		if errConnection != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not send a message"})
			return
		}
		if message.Username != usernameConnected {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not send a message"})
			return
		}

		// check if roomID is a valid objectID, and convert it to an objectID
		_, err := primitive.ObjectIDFromHex(message.RoomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The room does not exist"})
			return
		}
		// check if username is not too long
		if len(message.Username) > 20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username"})
			return

		}
		// check if content is not too long
		if len(message.Content) > 5000 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Message is too long"})
			return
		}
		// check if content is not too short
		if len(message.Content) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Message is too short"})
			return
		}
		// create message
		messageCreated, err := messageService.CreateMessage(c.Request.Context(), &message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, messageCreated)
	}
}

func GetMessagesHandler(messageService MessageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomIDString := c.Param("id")

		// check if roomID is a valid objectID, and convert it to an objectID
		roomID, err := primitive.ObjectIDFromHex(roomIDString)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The room does not exist"})
			return
		}
		// get messages
		messages, err := messageService.GetMessages(c.Request.Context(), roomID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, messages)
	}
}

func DeleteMessageHandler(messageService MessageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		messageIDString := c.Param("id")

		// check if messageID is a valid objectID, and convert it to an objectID
		messageID, err := primitive.ObjectIDFromHex(messageIDString)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No message found"})
			return
		}
		message, err := messageService.GetMessage(c.Request.Context(), messageID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete message"})
			return
		}

		// check if userConnected is the one who sends a deletion message request
		_, usernameConnected, errConnection := utils.GetUserIDAndUsernameFromContext(c)
		if errConnection != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete a message"})
			return
		}
		if message.Username != usernameConnected {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete a message"})
			return
		}

		// delete message
		if err := messageService.DeleteMessage(c.Request.Context(), messageID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete message"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": true, "message:": "your message has been successfully deleted !"})
	}

}
