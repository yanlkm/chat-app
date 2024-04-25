package room

import (
	"chat-app/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"regexp"
)

func CreateRoomHandler(roomService RoomService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newRoom Room
		if err := c.ShouldBindJSON(&newRoom); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		// check if name, creator and description are empty
		if newRoom.Name == "" || newRoom.Description == "" || newRoom.Creator == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
			return
		}
		// check if creator is a valid objectID, and convert it to an objectID
		_, err := primitive.ObjectIDFromHex(newRoom.Creator)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The creator does not exist"})
			return
		}

		// check if name is not too long
		if len(newRoom.Name) > 20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Name is too long"})
			return
		}
		// check if description is not too long
		if len(newRoom.Description) > 300 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Description is too long"})
			return
		}
		// check if description is not too short
		if len(newRoom.Description) < 10 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Description is too short"})
			return
		}
		// check if name respects convention with a regex check
		nameConvention := "^[a-zA-Z0-9_]*$"
		if re, _ := regexp.Compile(nameConvention); !re.Match([]byte(newRoom.Name)) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid name"})
			return
		}
		// check if description respects convention with a regex check
		descriptionConvention := "^[a-zA-Z0-9_ ]*$"
		if re, _ := regexp.Compile(descriptionConvention); !re.Match([]byte(newRoom.Description)) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid description"})
			return
		}
		// check if name is unique
		if err := roomService.CheckName(c.Request.Context(), newRoom.Name); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Name already exists"})
			return
		}
		// create room
		room, err := roomService.CreateRoom(c.Request.Context(), &newRoom)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create room"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": true, "room": room})

	}
}

func GetRoomHandler(roomService RoomService) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(roomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The room does not exist"})
			return
		}
		room, err := roomService.GetRoom(c.Request.Context(), objectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get room"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": true, "room": room})
	}
}

// get all members of a room
func GetMembersOfRoom(roomService RoomService) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(roomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The room does not exist"})
			return
		}
		room, err := roomService.GetRoom(c.Request.Context(), objectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get room"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": true, "members": room.Members})
	}

}

// add a member to a room
func AddMemberToRoom(roomService RoomService) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(roomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The room does not exist"})
			return
		}
		var member Member
		if err := c.ShouldBindJSON(&member); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		memberObjectID, err := primitive.ObjectIDFromHex(member.ID)
		if err != nil {
			// debug
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "The member does not exist"})
			return
		}
		room, err := roomService.AddMember(c.Request.Context(), objectID, memberObjectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add member"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": true, "room": room})
	}
}

// remove a member from a room
func RemoveMemberFromRoom(roomService RoomService) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(roomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The room does not exist"})
			return
		}
		var member Member
		if err := c.ShouldBindJSON(&member); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		// convert member ID string to string hex
		memberObjectID, err := primitive.ObjectIDFromHex(member.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The member does not exist"})
			return
		}
		room, err := roomService.RemoveMember(c.Request.Context(), objectID, memberObjectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not remove member"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": true, "room": room})
	}
}
func GetRoomsHandler(roomService RoomService) gin.HandlerFunc {
	return func(c *gin.Context) {
		rooms, err := roomService.GetAllRooms(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get rooms"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": true, "rooms": rooms})
	}
}

// TODO : Add a handler to delete a room
func DeleteRoomHandler(roomService RoomService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get userID (logged in) from token cookie/headers
		// Get token from cookie
		token, err := c.Cookie("token")
		if err != nil {
			token := c.GetHeader("Authorization")
			if token == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "1 - Unauthorized"})
				c.Abort()
				return
			}
		}

		// Verify and decode token
		claims, err := utils.VerifyToken(&token)
		if err != nil {
			fmt.Println("Error: ", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "2 - Unauthorized"})
			c.Abort()
			return
		}

		//check if user is logged in with id passed in params
		// convert claims ObjectId to string
		stringClaimsObjectId := claims.UserID.Hex()

		roomID := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(roomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The room does not exist"})
			return
		}
		room, err := roomService.GetRoom(c.Request.Context(), objectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get room"})
			return
		}
		// check if stringClaimsObjectID is the room creator
		if stringClaimsObjectId != room.Creator {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "You are not allowed to do this action"})
			return
		}
		// check if there still members in room
		if len(room.Members) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Please, delete all room members before"})
			return
		}
		err = roomService.DeleteRoom(c.Request.Context(), room.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to delete room"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Room deleted successfully"})

	}
}
