package routes

import (
	_ "chat-golang/docs" // Import generated docs
	"chat-golang/src/internal/infrastructure/services"
	"chat-golang/src/internal/interfaces/http/handlers"
	"chat-golang/src/internal/interfaces/http/middleware"
	"chat-golang/src/internal/interfaces/websocket"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(
	r *gin.Engine,
	authHandler *handlers.AuthHandler,
	agendaHandler *handlers.AgendaHandler,
	chatHandler *handlers.ChatHandler,
	userHandler *handlers.UserHandler,
	roomHandler *handlers.RoomHandler,
	wsHandler *websocket.WSHandler,
	jwtService *services.JWTService,
) {
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/google", authHandler.GoogleLogin)

			// Protected routes
			auth.GET("/me", middleware.AuthMiddleware(jwtService), authHandler.GetMe)
		}

		// User routes
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(jwtService))
		{
			users.GET("", userHandler.GetAll)
			users.GET("/:id", userHandler.GetByID)
		}

		// Room routes
		roomsBase := api.Group("/rooms")
		roomsBase.Use(middleware.AuthMiddleware(jwtService))
		{
			roomsBase.POST("", roomHandler.Create)
			roomsBase.GET("", roomHandler.GetAll)
			roomsBase.POST("/:roomId/agendas", agendaHandler.Create)
			roomsBase.GET("/:roomId/agendas", agendaHandler.GetByRoom)
			roomsBase.GET("/:roomId/messages", chatHandler.GetMessages)
			roomsBase.POST("/:roomId/messages", chatHandler.SendMessage)
		}

		// Chat routes (Aligned with Mobile FE)
		chatGroup := api.Group("/chat")
		chatGroup.Use(middleware.AuthMiddleware(jwtService))
		{
			chatGroup.GET("/conversations", roomHandler.GetAll)
			chatGroup.GET("/conversations/:roomId", roomHandler.GetAll) // Fallback for single detail
			chatGroup.GET("/conversations/:roomId/messages", chatHandler.GetMessages)
			chatGroup.POST("/messages", chatHandler.SendMessage)
			chatGroup.POST("/conversations", roomHandler.Create)
		}

		// WebSocket
		api.GET("/ws", wsHandler.HandleWS)
	}

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
