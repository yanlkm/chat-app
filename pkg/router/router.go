package router

import (
	"chat-app/pkg/auth"
	"chat-app/pkg/middlewares"
	"chat-app/pkg/room"
	"chat-app/pkg/user"
	"github.com/gin-gonic/gin"
)

func NewRouter(userService user.UserService, authService auth.AuthService, roomService room.RoomService) *gin.Engine {

	// Set Gin to default(debug) mode
	r := gin.Default()

	// User routes
	r.GET("users",
		user.GetUsersHandler(userService))
	r.GET("users/:id", middlewares.AuthMiddleware(),
		user.GetUserHandler(userService))
	r.POST("users",
		user.CreateUserHandler(userService))
	r.PUT("users/:id",
		middlewares.AuthMiddleware(),
		user.UpdateUserHandler(userService))
	r.PUT("users/:id/password",
		middlewares.AuthMiddleware(),
		user.UpdatePasswordHandler(userService))
	r.DELETE("users/:id",
		middlewares.AuthMiddleware(),
		user.DeleteUserHandler(userService))

	// Room routes
	r.GET("rooms",
		room.GetRoomsHandler(roomService))
	r.GET("rooms/:id",
		room.GetRoomHandler(roomService))
	r.POST("rooms",
		room.CreateRoomHandler(roomService))
	r.PUT("rooms/add/:id",
		room.AddMemberToRoom(roomService))
	r.PUT("rooms/remove/:id",
		room.RemoveMemberFromRoom(roomService))
	// get all members of a room
	r.GET("rooms/members/:id",
		room.GetRoomMembersHandler(roomService))
	// TODO: Add a route to delete a room
	r.DELETE("rooms/delete/:id",
		room.DeleteRoomHandler(roomService))

	// auth routes
	r.POST("auth/login",
		auth.LoginUserHandler(authService))
	r.GET("auth/logout",
		auth.LogoutUserHandler(authService))

	return r
}
