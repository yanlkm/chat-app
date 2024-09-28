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
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

func main() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Connect to MongoDB
	mongoURI := os.Getenv("MONGO_URI")
	client, err := mongo.NewMongoClient(mongoURI)
	if err != nil {
		fmt.Println("Failed to connect to MongoDB: %v", err)
	}

	db := client.Database("chat_app_test")
	userCollection := db.Collection("users")
	roomCollection := db.Collection("rooms")
	messageCollection := db.Collection("messages")
	codeCollection := db.Collection("codes")



	// Initialize code repository and service
	codeRepo := code.NewCodeRepository(codeCollection)
	codeService := code.NewCodeService(codeRepo)
	// Initialize user repository and service
	userRepo := user.NewUserRepository(userCollection, messageCollection)
	userService := user.NewUserService(userRepo)
	// Initialize auth repository and service
	authRepo := auth.NewAuthRepository(userCollection)
	authService := auth.NewAuthService(authRepo)
	// Initialize room repository and service
	roomRepo := room.NewRoomRepository(roomCollection, userCollection)
	roomService := room.NewRoomService(roomRepo)
	// Initialize message repository and service
	messageRepo := message.NewMessageRepository(messageCollection, roomCollection)
	messageService := message.NewMessageService(messageRepo)

	// Initialize router
	r := router.NewRouter(userService, codeService, authService, roomService, messageService)

	// Start HTTP server
	port := os.Getenv("PORT")
	fmt.Printf("Server started on %s", port)
	// CORS configuration
	err = http.ListenAndServe(port, handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}),
		handlers.AllowedOrigins([]string{os.Getenv("CORS_ORIGIN")}),
	)(r))
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		log.Fatal(err)
	}
}
