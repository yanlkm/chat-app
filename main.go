package main

import (
	"chat-app/pkg/auth"
	"chat-app/pkg/code"
	"chat-app/pkg/database"
	"chat-app/pkg/message"
	"chat-app/pkg/room"
	"chat-app/pkg/router"
	"chat-app/pkg/user"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func main() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Connect to MongoDB
	mongoURI := os.Getenv("MONGO_URI")
	fmt.Println(mongoURI)
	client, err := mongo.NewMongoClient(mongoURI)
	if err != nil {
		fmt.Println("Failed to connect to MongoDB: %v", err)
	}

	db := client.Database("chat_app")
	userCollection := db.Collection("users")
	roomCollection := db.Collection("rooms")
	messageCollection := db.Collection("messages")
	codeCollection := db.Collection("codes")

	//debug
	fmt.Println(userCollection)

	// Initialize code repository and service
	codeRepo := code.NewCodeRepository(codeCollection)
	codeService := code.NewCodeService(codeRepo)
	// Initialize user repository and service
	userRepo := user.NewUserRepository(userCollection)
	userService := user.NewUserService(userRepo)
	// Initialize auth repository and service
	authRepo := auth.NewAuthRepository(userCollection)
	authService := auth.NewAuthService(authRepo)
	// Initialize room repository and service
	roomRepo := room.NewRoomRepository(roomCollection)
	roomService := room.NewRoomService(roomRepo)
	// Initialize message repository and service
	messageRepo := message.NewMessageRepository(messageCollection, roomCollection)
	messageService := message.NewMessageService(messageRepo)

	// Initialize router
	r := router.NewRouter(userService, codeService, authService, roomService, messageService)

	// Start HTTP server
	port := os.Getenv("PORT")
	fmt.Printf("Server started on %s", port)
	err = http.ListenAndServe(port, r)
	if err != nil {
		fmt.Println("Failed to start server: %v", err)
	}
}
