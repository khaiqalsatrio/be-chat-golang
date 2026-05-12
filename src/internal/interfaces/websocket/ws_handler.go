package websocket

import (
	"chat-golang/src/internal/infrastructure/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

type WSHandler struct {
	Hub        *Hub
	JWTService *services.JWTService
}

func NewWSHandler(hub *Hub, jwtService *services.JWTService) *WSHandler {
	return &WSHandler{
		Hub:        hub,
		JWTService: jwtService,
	}
}

func (h *WSHandler) HandleWS(c *gin.Context) {
	// 1. Authenticate (optional, can be done via query param or handshake)
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
		return
	}

	token, err := h.JWTService.ValidateToken(tokenString)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// 2. Upgrade connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}

	// 3. Register client
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userID := claims["user_id"].(string)
	client := &Client{
		ID:    userID,
		Conn:  conn,
		Rooms: make(map[string]bool),
		Send:  make(chan []byte, 256),
	}
	h.Hub.Register <- client

	// 4. Start read/write loops
	go client.writePump()
	go client.readPump(h.Hub)
}

func (c *Client) readPump(h *Hub) {
	defer func() {
		h.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		log.Printf("Received: %s", message)
		// Process message (e.g. broadcast)
	}
}

func (c *Client) writePump() {
	defer func() {
		c.Conn.Close()
	}()

	for message := range c.Send {
		c.Conn.WriteMessage(websocket.TextMessage, message)
	}
	c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}
