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
	r.GET("users", middlewares.IsAdminMiddleware(),
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
	r.POST("codes", middlewares.IsAdminMiddleware(),
		code.CreateCodeHandler(codeService))

	// Room routes
	r.GET("rooms",
		room.GetRoomsHandler(roomService))
	r.GET("rooms/:id", middlewares.IsLoggedInMiddleware(),
		room.GetRoomHandler(roomService))
	r.GET("rooms/user/:id", middlewares.IsLoggedInMiddleware(),
		room.GetUserRoomsHandler(roomService))
	r.POST("rooms", middlewares.IsAdminMiddleware(),
		room.CreateRoomHandler(roomService))
	r.PUT("rooms/add/:id", middlewares.IsLoggedInMiddleware(),
		room.AddMemberToRoom(roomService))
	r.PUT("rooms/remove/:id", middlewares.IsLoggedInMiddleware(),
		room.RemoveMemberFromRoom(roomService))
	r.PATCH("rooms/add/hashtag/:id", middlewares.IsAdminMiddleware(),
		room.AddHashtagToRoomHandler(roomService))
	r.PATCH("rooms/remove/hashtag/:id", middlewares.IsAdminMiddleware(),
		room.RemoveHashtagFromRoomHandler(roomService))
	// get all members of a room
	r.GET("rooms/members/:id", middlewares.IsLoggedInMiddleware(),
		room.GetRoomMembersHandler(roomService, userService))
	r.DELETE("rooms/delete/:id", middlewares.IsAdminMiddleware(),
		room.DeleteRoomHandler(roomService))

	// Message routes
	r.POST("messages", middlewares.IsLoggedInMiddleware(),
		message.CreateMessageHandler(messageService))
	r.GET("messages/:id", middlewares.IsLoggedInMiddleware(),
		message.GetMessagesHandler(messageService))
	r.DELETE("messages/:id", middlewares.IsLoggedInMiddleware(),
		message.DeleteMessageHandler(messageService))

	// auth routes
	r.POST("auth/login",
		auth.LoginUserHandler(authService))
	r.GET("auth/logout",
		auth.LogoutUserHandler(authService))

	// websocket routes
	r.GET("/ws", func(c *gin.Context) {
		websocket.WebSocketHandler(c, messageService, roomService)
	})
	// starting handling rooms
	c := gin.Context{}
	func(c *gin.Context) {
		go websocket.HandleRooms(c, roomService)
	}(&c)

	return r
}
