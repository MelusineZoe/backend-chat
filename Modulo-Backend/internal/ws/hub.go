package ws // ← Cambiado

import (
	"encoding/json"
	"sync"

	"github.com/google/uuid"
)

type Hub struct {
	clients    map[uuid.UUID]*Client
	rooms      map[uuid.UUID]map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

type Message struct {
	ID        uuid.UUID `json:"id"`
	RoomID    uuid.UUID `json:"room_id"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt string    `json:"created_at"`
}

var GlobalHub = NewHub()

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]*Client),
		rooms:      make(map[uuid.UUID]map[*Client]bool),
		broadcast:  make(chan *Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Run inicia el hub en un goroutine (se llama desde main)
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			if h.rooms[client.RoomID] == nil {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.rooms[client.RoomID][client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				if roomClients, ok := h.rooms[client.RoomID]; ok {
					delete(roomClients, client)
					if len(roomClients) == 0 {
						delete(h.rooms, client.RoomID)
					}
				}
				close(client.Send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			if clients, ok := h.rooms[message.RoomID]; ok {
				for client := range clients {
					select {
					case client.Send <- toJSON(message):
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func toJSON(msg *Message) []byte {
	data, _ := json.Marshal(msg)
	return data
}
