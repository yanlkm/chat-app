package websocket

import (
	"chat-app/pkg/room"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

// RoomSocket struct from the websocket package
type RoomSocket struct {
	ID        string
	Members   map[*websocket.Conn]bool
	broadcast chan MessageSocket
}

// MessageSocket struct from the websocket package
type MessageSocket struct {
	RoomID    string    `json:"roomId,omitempty"`
	Username  string    `json:"username,omitempty"`
	UserID    string    `json:"userId,omitempty"`
	Message   string    `json:"message,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	Token     string    `json:"token,omitempty"`
}

// upgrader variable from the websocket package
var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	rooms   = make(map[string]*RoomSocket)
	roomsMu sync.Mutex
)

// GetRoomsFromDatabase retrieves rooms from the database
func GetRoomsFromDatabase(c *gin.Context, roomService room.RoomService) []*RoomSocket {
	// Get all rooms from the database
	roomsFromDB, err := roomService.GetAllRooms(c)
	if err != nil {
		fmt.Printf("Error getting rooms from database: %v\n", err)
		return nil
	}

	// Create a slice of RoomSocket structs
	var rooms []*RoomSocket
	// Append each room to the slice
	for _, r := range roomsFromDB {
		rooms = append(rooms, &RoomSocket{
			ID:        r.ID,
			Members:   make(map[*websocket.Conn]bool),
			broadcast: make(chan MessageSocket),
		})
	}
	// Return the slice of rooms
	return rooms
}

// UpdateRoomsFromDatabase updates the rooms map with the latest rooms from the database
func UpdateRoomsFromDatabase(c *gin.Context, roomService room.RoomService) {
	// Get all rooms from the database
	roomsFromDB, err := roomService.GetAllRooms(c)
	if err != nil {
		fmt.Printf("Error updating rooms from database: %v\n", err)
		return
	}

	// Lock the rooms map
	roomsMu.Lock()
	defer roomsMu.Unlock()

	// Update the global rooms map
	for _, r := range roomsFromDB {
		roomID := r.ID
		if _, ok := rooms[roomID]; !ok {
			rooms[roomID] = &RoomSocket{
				ID:        roomID,
				Members:   make(map[*websocket.Conn]bool),
				broadcast: make(chan MessageSocket),
			}
			// Start broadcasting messages to all members in the room
			go handleRoomBroadcast(rooms[roomID])
		}
	}

	// Remove any rooms that no longer exist in the database
	for id := range rooms {
		found := false
		for _, r := range roomsFromDB {
			if r.ID == id {
				found = true
				break
			}
		}
		// If the room is not found in the database, delete it from the rooms map
		if !found {
			delete(rooms, id)
		}
	}
}

// handleRoomBroadcast handles broadcasting messages to all members in the room
func handleRoomBroadcast(room *RoomSocket) {
	for {
		// Get the message from the broadcast channel
		msg := <-room.broadcast

		// Broadcast the message to all members in the room
		roomsMu.Lock()
		// Iterate over all members in the room
		for client := range room.Members {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				client.Close()
				delete(room.Members, client)
			}
		}
		// Unlock the rooms mutex
		roomsMu.Unlock()
	}
}

// HandleRooms creates goroutines to manage existing rooms
func HandleRooms(c *gin.Context, roomService room.RoomService) {
	// Get all rooms from the database
	roomsFromDB := GetRoomsFromDatabase(c, roomService)

	// Lock the rooms map
	for _, roomDB := range roomsFromDB {
		roomSocket := &RoomSocket{
			ID:        roomDB.ID,
			Members:   make(map[*websocket.Conn]bool),
			broadcast: make(chan MessageSocket),
		}
		// Start broadcasting messages to all members in the room
		roomsMu.Lock()
		rooms[roomSocket.ID] = roomSocket
		roomsMu.Unlock()
		// Start broadcasting messages to all members in the room
		go func(r *RoomSocket) {
			for {
				// Get the message from the broadcast channel
				msg := <-r.broadcast

				// Broadcast the message to all members in the room
				for client := range r.Members {
					err := client.WriteJSON(msg)
					if err != nil {
						fmt.Printf("error: %v\n", err)
						// Close the client connection
						client.Close()
						roomsMu.Lock()
						// Delete the client from the room
						delete(r.Members, client)
						roomsMu.Unlock()
					}
				}
			}
		}(roomSocket)
	}
}
