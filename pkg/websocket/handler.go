package websocket

import (
	"chat-app/pkg/message"
	"chat-app/pkg/room"
	"chat-app/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// WebSocketHandler handles WebSocket connections for a specific room
func WebSocketHandler(c *gin.Context, messageService message.MessageService, roomService room.RoomService) {
	roomID := c.Query("id")
	// update the room from the database
	UpdateRoomsFromDatabase(c, roomService)
	roomsMu.Lock()
	// get the room from the rooms map
	room, ok := rooms[roomID]
	if !ok {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Room not found"})
		roomsMu.Unlock()
		return
	}
	roomsMu.Unlock()

	// upgrade the HTTP connection to a WebSocket connection
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
			return
	}
	defer ws.Close()

	// add the WebSocket connection to the room
	room.Members[ws] = true

	// read messages from the WebSocket connection
	for {
		var msg MessageSocket
		err := ws.ReadJSON(&msg)
		if err != nil {
			roomsMu.Lock()
			delete(room.Members, ws)
			roomsMu.Unlock()
			break
		}

		// check the token on the message
		if msg.Token == "" {
			roomsMu.Lock()
			roomsMu.Unlock()
			return
		}
		if msg.Token != "" {
			_, err := utils.VerifyToken(&msg.Token)
			if err != nil {
				roomsMu.Lock()
				roomsMu.Unlock()
				return
			}
		}
		// save the message to the database
		messageDB := message.MessageEntity{
			RoomID:   roomID,
			UserID:   msg.UserID,
			Username: msg.Username,
			Content:  msg.Message,
		}
		// create the message
		_, err = messageService.CreateMessage(c.Request.Context(), &messageDB)
		if err != nil {
			roomsMu.Lock()
			delete(room.Members, ws)
			roomsMu.Unlock()
			break
		}

		// broadcast the message to all members in the room
		room.broadcast <- msg
	}
}
