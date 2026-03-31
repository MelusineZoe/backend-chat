package impl

import (
	"backend/internal/repository"

	"gorm.io/gorm"
)

// Repositories agrupa todas las interfaces de repositorios
type Repositories struct {
	User repository.UserRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User: NewUserRepository(db),
	}
}
