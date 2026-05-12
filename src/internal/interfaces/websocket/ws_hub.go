package websocket

import (
	"chat-golang/src/internal/domain/entities"
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID    string
	Conn  *websocket.Conn
	Rooms map[string]bool
	Send  chan []byte
}

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *entities.Message
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *entities.Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.ID] = client
			h.mu.Unlock()
			log.Printf("Client registered: %s", client.ID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("Client unregistered: %s", client.ID)

		case message := <-h.Broadcast:
			h.mu.Lock()
			payload, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				h.mu.Unlock()
				continue
			}

			// Broadcast to participants of the room
			// For now, we'll iterate through all clients and check if they should receive it
			// In a more advanced implementation, we'd have a map[roomID][]Clients
			for _, client := range h.Clients {
				// Simple broadcast to everyone for now, but with full JSON
				// TODO: Add logic to check if client is a participant in message.RoomID
				select {
				case client.Send <- payload:
				default:
					close(client.Send)
					delete(h.Clients, client.ID)
				}
			}
			h.mu.Unlock()
		}
	}
}
