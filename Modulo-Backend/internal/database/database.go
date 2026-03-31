package database

import (
	"fmt"
	"log"

	"backend/internal/config"
	"backend/internal/model"
	"backend/internal/repository/impl"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Error al conectar a PostgreSQL:", err)
	}

	log.Println("✅ Conexión exitosa a PostgreSQL")
	return db
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(model.Models...)
	if err != nil {
		log.Fatal("❌ Error durante AutoMigrate:", err)
	}
	log.Println("✅ Migración de tablas completada (users)")
}

// NewRepositories devuelve todas las implementaciones de repositories
func NewRepositories(db *gorm.DB) *impl.Repositories {
	return impl.NewRepositories(db)
}
