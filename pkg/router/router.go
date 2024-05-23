package router

import (
	"chat-app/pkg/auth"
	"chat-app/pkg/code"
	"chat-app/pkg/message"
	"chat-app/pkg/middlewares"
	"chat-app/pkg/room"
	"chat-app/pkg/user"
	"chat-app/pkg/websocket"
	"github.com/gin-gonic/gin"
)

func NewRouter(userService user.UserService, codeService code.CodeService, authService auth.AuthService, roomService room.RoomService, messageService message.MessageService) *gin.Engine {

	// Set Gin to default(debug) mode
	r := gin.Default()

	// User routes
	r.GET("users",
		user.GetUsersHandler(userService))
	r.GET("users/:id", middlewares.AuthMiddleware(),
		user.GetUserHandler(userService))
	r.POST("users",
		user.CreateUserHandler(userService, codeService))
	r.PUT("users/:id",
		middlewares.AuthMiddleware(),
		user.UpdateUserHandler(userService))
	r.PUT("users/:id/password",
		middlewares.AuthMiddleware(),
		user.UpdatePasswordHandler(userService))
	r.GET("users/ban/:id/:idBanned",
		middlewares.AuthMiddleware(),
		user.BanUserHandler(userService))
	r.GET("users/unban/:id/:idBanned",
		middlewares.AuthMiddleware(),
		user.UnBanUserHandler(userService))
	r.DELETE("users/:id",
		middlewares.AuthMiddleware(),
		user.DeleteUserHandler(userService))

	//code route
	r.POST("codes",
		code.CreateCodeHandler(codeService))

	// Room routes
	r.GET("rooms",
		room.GetRoomsHandler(roomService))
	r.GET("rooms/:id",
		room.GetRoomHandler(roomService))
	r.GET("rooms/user/:id",
		room.GetUserRoomsHandler(roomService))
	r.POST("rooms",
		room.CreateRoomHandler(roomService))
	r.PUT("rooms/add/:id",
		room.AddMemberToRoom(roomService))
	r.PUT("rooms/remove/:id",
		room.RemoveMemberFromRoom(roomService))
	// get all members of a room
	r.GET("rooms/members/:id",
		room.GetRoomMembersHandler(roomService, userService))
	r.DELETE("rooms/delete/:id",
		room.DeleteRoomHandler(roomService))

	// Message routes
	r.POST("messages",
		message.CreateMessageHandler(messageService))
	r.GET("messages/:id",
		message.GetMessagesHandler(messageService))
	r.DELETE("messages/:id",
		message.DeleteMessageHandler(messageService))

	// auth routes
	r.POST("auth/login",
		auth.LoginUserHandler(authService))
	r.GET("auth/logout",
		auth.LogoutUserHandler(authService))

	// websocket routes
	r.GET("/ws", func(c *gin.Context) {
		websocket.WebSocketHandler(c, messageService)
	})
	// starting handling rooms
	c := gin.Context{}
	func(c *gin.Context) {
		go websocket.HandleRooms(c, roomService)
	}(&c)

	return r
}
