package handler

import (
	"log"
	"net/http"

	"backend/internal/model"
	"backend/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func WebSocketConnect(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	userID := userIDInterface.(uuid.UUID)

	roomIDStr := c.Query("room_id")
	if roomIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room_id es requerido"})
		return
	}

	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room_id inválido"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var room model.Room
	if err := db.First(&room, "id = ? AND type = ?", roomID, "public").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sala pública no encontrada"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrade WebSocket:", err)
		return
	}

	client := &ws.Client{ // ← ws.Client
		ID:     uuid.New(),
		UserID: userID,
		RoomID: roomID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	ws.GlobalHub.Register(client) // ← ws.GlobalHub

	go client.WritePump()
	go client.ReadPump(ws.GlobalHub, db)

	log.Printf("Cliente WebSocket conectado - User: %s | Sala: %s", userID, roomID)
}
