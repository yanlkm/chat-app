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

// Room struct from the websocket package
type RoomSocket struct {
	ID        string
	Members   map[*websocket.Conn]bool
	broadcast chan MessageSocket
}

// Message struct from the websocket package
type MessageSocket struct {
	RoomID    string    `json:"roomId,omitempty"`
	Username  string    `json:"username,omitempty"`
	UserID    string    `json:"userId,omitempty"`
	Message   string    `json:"message,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

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
	roomsFromDB, err := roomService.GetAllRooms(c)
	if err != nil {
		fmt.Printf("Error getting rooms from database: %v\n", err)
		return nil
	}

	var rooms []*RoomSocket
	for _, r := range roomsFromDB {
		rooms = append(rooms, &RoomSocket{
			ID:        r.ID.Hex(),
			Members:   make(map[*websocket.Conn]bool),
			broadcast: make(chan MessageSocket),
		})
	}
	return rooms
}

// HandleRooms creates goroutines to manage existing rooms
func HandleRooms(c *gin.Context, roomService room.RoomService) {
	roomsFromDB := GetRoomsFromDatabase(c, roomService)

	for _, roomDB := range roomsFromDB {
		room := &RoomSocket{
			ID:        roomDB.ID,
			Members:   make(map[*websocket.Conn]bool),
			broadcast: make(chan MessageSocket),
		}

		roomsMu.Lock()
		rooms[room.ID] = room
		roomsMu.Unlock()

		go func(r *RoomSocket) {
			for {
				msg := <-r.broadcast

				for client := range r.Members {
					err := client.WriteJSON(msg)
					if err != nil {
						fmt.Printf("error: %v\n", err)
						client.Close()
						roomsMu.Lock()
						delete(r.Members, client)
						roomsMu.Unlock()
					}
				}
			}
		}(room)
	}
}
