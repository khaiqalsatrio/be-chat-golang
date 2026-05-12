package main

import (
	"chat-golang/src/config"
	"chat-golang/src/internal/infrastructure/database"
	"chat-golang/src/internal/infrastructure/repositories"
	"chat-golang/src/internal/infrastructure/services"
	"chat-golang/src/internal/infrastructure/worker"
	"chat-golang/src/internal/interfaces/http/handlers"
	"chat-golang/src/internal/interfaces/http/routes"
	"chat-golang/src/internal/interfaces/websocket"
	"chat-golang/src/internal/usecases/agenda"
	"chat-golang/src/internal/usecases/auth"
	"chat-golang/src/internal/usecases/chat"
	"chat-golang/src/internal/usecases/room"
	"chat-golang/src/internal/usecases/user"
	"log"

	"github.com/gin-gonic/gin"
)

// @title Chat Golang API
// @version 1.0
// @description This is a sample chat server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token JWT Anda. Contoh: "eyJhbGci..." atau "Bearer eyJhbGci..."
func main() {
	// Load Configuration
	cfg := config.LoadConfig()

	// Initialize Database
	database.InitDB(cfg)

	// Start Background Workers
	worker.StartReminderWorker()

	// Initialize Services
	jwtService := services.NewJWTService(cfg.JWTSecret)
	googleAuthService := services.NewGoogleAuthService(cfg.GoogleID)

	// Initialize Repositories
	userRepo := repositories.NewPostgresUserRepository(database.DB)
	agendaRepo := repositories.NewPostgresAgendaRepository(database.DB)
	chatRepo := repositories.NewPostgresChatRepository(database.DB)
	roomRepo := repositories.NewPostgresRoomRepository(database.DB)

	// Initialize WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize Usecases
	registerUsecase := auth.NewRegisterUsecase(userRepo)
	loginUsecase := auth.NewLoginUsecase(userRepo, jwtService)
	getMeUsecase := auth.NewGetMeUsecase(userRepo)
	googleLoginUsecase := auth.NewGoogleLoginUsecase(userRepo, googleAuthService, jwtService)

	createAgendaUsecase := agenda.NewCreateAgendaUsecase(agendaRepo)
	getAgendasUsecase := agenda.NewGetAgendasUsecase(agendaRepo)

	sendMessageUsecase := chat.NewSendMessageUsecase(chatRepo, hub)
	getMessagesUsecase := chat.NewGetMessagesUsecase(chatRepo)

	getAllUsersUsecase := user.NewGetAllUsersUsecase(userRepo)
	getUserByIDUsecase := user.NewGetUserByIDUsecase(userRepo)

	createRoomUsecase := room.NewCreateRoomUsecase(roomRepo)
	getRoomsUsecase := room.NewGetRoomsUsecase(roomRepo)

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(registerUsecase, loginUsecase, getMeUsecase, googleLoginUsecase, jwtService)
	agendaHandler := handlers.NewAgendaHandler(createAgendaUsecase, getAgendasUsecase)
	chatHandler := handlers.NewChatHandler(sendMessageUsecase, getMessagesUsecase)
	userHandler := handlers.NewUserHandler(getAllUsersUsecase, getUserByIDUsecase)
	roomHandler := handlers.NewRoomHandler(createRoomUsecase, getRoomsUsecase)
	wsHandler := websocket.NewWSHandler(hub, jwtService)

	// Setup Router
	r := gin.Default()

	// Setup Routes
	routes.SetupRoutes(r, authHandler, agendaHandler, chatHandler, userHandler, roomHandler, wsHandler, jwtService)

	// Start Server
	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
