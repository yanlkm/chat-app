package websocket

import (
	"chat-app/pkg/message"
	"chat-app/pkg/room"
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
	room, ok := rooms[roomID]
	if !ok {
		fmt.Println("Room not found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room not found"})
		roomsMu.Unlock()
		return
	}
	roomsMu.Unlock()

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	room.Members[ws] = true

	for {
		var msg MessageSocket
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			roomsMu.Lock()
			delete(room.Members, ws)
			roomsMu.Unlock()
			break
		}
		// save the message to the database
		messageDB := message.MessageEntity{
			RoomID:   roomID,
			UserID:   msg.UserID,
			Username: msg.Username,
			Content:  msg.Message,
		}
		_, err = messageService.CreateMessage(c.Request.Context(), &messageDB)
		if err != nil {
			fmt.Println(err)
			roomsMu.Lock()
			delete(room.Members, ws)
			roomsMu.Unlock()
			break
		}

		room.broadcast <- msg
	}
}
