package room

import (
	"chat-app/pkg/user"
	"chat-app/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

// CreateRoomHandler create a room
func CreateRoomHandler(roomService RoomService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newRoom RoomEntity
		if err := c.ShouldBindJSON(&newRoom); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		// check if name, creator and description are empty
		if newRoom.Name == "" || newRoom.Description == "" || newRoom.Creator == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Could not create room"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"room": room})

	}
}

// GetRoomHandler get a room by its id
func GetRoomHandler(roomService RoomService) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("id")
		// get room by id
		room, err := roomService.GetRoom(c.Request.Context(), roomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get room"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"room": room})
	}
}

// GetUserRoomsHandler get all rooms of a user
func GetUserRoomsHandler(roomService RoomService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")

		// get all rooms where user is a member
		rooms, err := roomService.GetUserRooms(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get rooms"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"rooms": rooms})
	}

}

// GetRoomsCreatedByAdminHandler get all rooms created by an admin
func GetRoomsCreatedByAdminHandler(roomService RoomService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		// get all rooms created by an admin
		rooms, err := roomService.GetRoomsCreatedByAdmin(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get rooms"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"rooms": rooms})
	}

}

// GetRoomMembersHandler get all members of a room
func GetRoomMembersHandler(roomService RoomService, userService user.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("id")

		room, err := roomService.GetRoom(c.Request.Context(), roomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get room"})
			return
		}
		// get all members id of a room in an array
		membersID := room.Members
		// get all users by id
		var users []user.UserEntity
		for _, memberID := range membersID {
			userRetrieved, err := userService.GetUser(c.Request.Context(), memberID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get members"})
				return
			}
			// remove password from user and more
			userRetrieved.Password = ""
			users = append(users, *userRetrieved)
		}
		// return users
		c.JSON(http.StatusOK, gin.H{"users": users})

	}
}

// AddMemberToRoom add a member to a room
func AddMemberToRoom(roomService RoomService) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("id")
		var member MemberEntity
		if err := c.ShouldBindJSON(&member); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		room, err := roomService.AddMember(c.Request.Context(), roomID, member.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Could not add member"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"room": room})
	}
}

// RemoveMemberFromRoom remove a member from a room
func RemoveMemberFromRoom(roomService RoomService) gin.HandlerFunc {
	return func(c *gin.Context) {

		roomID := c.Param("id")

		var member MemberEntity
		if err := c.ShouldBindJSON(&member); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		room, err := roomService.RemoveMember(c.Request.Context(), roomID, member.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Could not remove member"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"room": room})
	}
}

// AddHashtagToRoomHandler add a hashtag to a room
func AddHashtagToRoomHandler(roomService RoomService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get room by id
		var room *RoomEntity
		roomID := c.Param("id")

		var hashtagToAdd HashtagEntity
		if err := c.ShouldBindJSON(&hashtagToAdd); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": " Invalid hashtag"})
			return
		}
		// check if hashtag is a 3-min letters word
		if len(hashtagToAdd.Hashtag) < 3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Hashtag too short"})
			return
		}
		// check if hashtag respects hashtag name convention
		hashtagConvention := "^#[a-zA-Z]*$"
		if re, _ := regexp.Compile(hashtagConvention); !re.Match([]byte(hashtagToAdd.Hashtag)) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hashtag"})
			return
		}
		room, err := roomService.AddHashtag(c.Request.Context(), roomID, hashtagToAdd.Hashtag)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error adding hashtag " + hashtagToAdd.Hashtag})
			return
		}
		c.JSON(http.StatusOK, gin.H{"room": room})
		return
	}

}

// RemoveHashtagFromRoomHandler remove a hashtag from a room
func RemoveHashtagFromRoomHandler(roomService RoomService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var room *RoomEntity
		roomID := c.Param("id")

		var hashtagToRemove HashtagEntity
		if err := c.ShouldBindJSON(&hashtagToRemove); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": " Invalid hashtag"})
			return
		}
		// check if hashtag is a 3-min letters word
		if len(hashtagToRemove.Hashtag) < 3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Hashtag too short"})
			return
		}
		// check if hashtag respects hashtag name convention
		hashtagConvention := "^#[a-zA-Z]*$"
		if re, _ := regexp.Compile(hashtagConvention); !re.Match([]byte(hashtagToRemove.Hashtag)) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hashtag"})
			return
		}
		// remove hashtag from room
		room, err := roomService.RemoveHashtag(c.Request.Context(), roomID, hashtagToRemove.Hashtag)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error removing hashtag " + hashtagToRemove.Hashtag})
			return
		}
		c.JSON(http.StatusOK, gin.H{"room": room})
		return
	}

}

// GetRoomsHandler get all rooms
func GetRoomsHandler(roomService RoomService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get all rooms
		rooms, err := roomService.GetAllRooms(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get rooms"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": true, "rooms": rooms})
	}
}

// DeleteRoomHandler delete a room
func DeleteRoomHandler(roomService RoomService) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get user id from token
		userConnectedId, _, err := utils.GetUserIDAndUsernameFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Could not delete room"})
			return
		}

		// get room by id
		roomID := c.Param("id")

		room, err := roomService.GetRoom(c.Request.Context(), roomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get room"})
			return
		}
		// check if stringClaimsObjectID is the room creator
		if userConnectedId != room.Creator {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You are not allowed to do this action"})
			return
		}
		// check if there still members in room
		if len(room.Members) > 1 && room.Members[0] != room.Creator {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please, delete all room members before"})
			return
		}
		err = roomService.DeleteRoom(c.Request.Context(), room.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Fail to delete room"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Room deleted successfully"})

	}
}
