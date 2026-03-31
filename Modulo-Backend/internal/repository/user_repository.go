package repository

import (
	"backend/internal/model"

	"github.com/google/uuid"
)

// UserRepository define las operaciones disponibles para usuarios
type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id uuid.UUID) (*model.User, error)
}
