package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User representa un usuario del sistema
type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Username  string         `gorm:"size:50;unique;not null" json:"username" validate:"required,min=3,max=50"`
	Email     string         `gorm:"size:100;unique;not null" json:"email" validate:"required,email"`
	Password  string         `gorm:"column:password_hash;size:255;not null" json:"-"` // Nunca devolver la contraseña
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName especifica el nombre de la tabla en la base de datos
func (User) TableName() string {
	return "users"
}

// Models lista todos los modelos para AutoMigrate
var Models = []interface{}{
	&User{},
	&Room{},
	&Message{},
}

// Room represents a chat room or similar entity
type Room struct {
	ID   uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Type string    `json:"type" gorm:"not null"`
	// Add other fields as needed, e.g., Name, CreatedAt, etc.
}

// Message represents a chat message
type Message struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	RoomID    uuid.UUID `json:"room_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Content   string    `json:"content" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
