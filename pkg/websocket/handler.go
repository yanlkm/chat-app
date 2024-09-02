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
	// get the room from the rooms map
	room, ok := rooms[roomID]
	if !ok {
		fmt.Println("Room not found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room not found"})
		roomsMu.Unlock()
		return
	}
	roomsMu.Unlock()

	// upgrade the HTTP connection to a WebSocket connection
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
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
		// create the message
		_, err = messageService.CreateMessage(c.Request.Context(), &messageDB)
		if err != nil {
			fmt.Println(err)
			roomsMu.Lock()
			delete(room.Members, ws)
			roomsMu.Unlock()
			break
		}

		// broadcast the message to all members in the room
		room.broadcast <- msg
	}
}
