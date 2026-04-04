package handler

import (
	"net/http"

	"backend/internal/dto"
	"backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// GetAllUsers retorna la lista de todos los usuarios registrados
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users"})
		return
	}

	// Convertir a DTO para no exponer contraseñas
	var userDTOs []dto.UserResponse
	for _, user := range users {
		userDTOs = append(userDTOs, dto.UserResponse{
			ID:       user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
		})
	}

	c.JSON(http.StatusOK, userDTOs)
}
