package ws

import (
	"encoding/json"
	"log"
	"time"

	"backend/internal/model"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type Client struct {
	ID     uuid.UUID
	UserID uuid.UUID
	RoomID uuid.UUID
	Conn   *websocket.Conn
	Send   chan []byte
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 2048
)

// WritePump envía mensajes al cliente
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Agregar mensajes pendientes
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ReadPump lee mensajes del cliente y los guarda + broadcast
func (c *Client) ReadPump(hub *Hub, db *gorm.DB) {
	defer func() {
		hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Parsear mensaje simple (puedes mejorar con un struct DTO)
		var incoming struct {
			Content string `json:"content"`
		}
		if json.Unmarshal(msg, &incoming) != nil {
			continue
		}

		// Guardar en base de datos
		message := model.Message{
			ID:      uuid.New(),
			RoomID:  c.RoomID,
			UserID:  c.UserID,
			Content: incoming.Content,
		}

		if err := db.Create(&message).Error; err != nil {
			log.Println("Error guardando mensaje:", err)
			continue
		}

		// Cargar username para el broadcast
		var user model.User
		db.First(&user, "id = ?", c.UserID)

		// Crear mensaje para broadcast
		broadcastMsg := &Message{
			ID:        message.ID,
			RoomID:    message.RoomID,
			UserID:    message.UserID,
			Username:  user.Username,
			Content:   message.Content,
			CreatedAt: message.CreatedAt.Format(time.RFC3339),
		}

		// Enviar al hub para broadcast
		hub.broadcast <- broadcastMsg
	}
}
