package main

import (
	"log"
	"net/http"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handler"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db := database.Connect(cfg)
	database.AutoMigrate(db)

	repos := database.NewRepositories(db)
	authHandler := handler.NewAuthHandler(repos.User, cfg)
	userHandler := handler.NewUserHandler(repos.User)

	r := gin.Default()

	// Middlewares
	r.Use(middleware.CORS())
	r.Use(gin.Recovery()) // Esto ayuda a ver errores mejor

	// Rutas
	r.POST("/api/auth/register", authHandler.Register)
	r.POST("/api/auth/login", authHandler.Login)
	r.GET("/api/users", userHandler.GetAllUsers)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Backend funcionando correctamente - " + cfg.ServerPort,
		})
	})

	log.Printf("🚀 Servidor corriendo en http://localhost:%s", cfg.ServerPort)

	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Error al iniciar servidor:", err)
	}
}
