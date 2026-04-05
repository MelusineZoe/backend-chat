# 📋 Instrucciones para Aplicar Cambios del Backend

## Opción 1: Usar Patch File (RECOMENDADO)

### En tu máquina (MelusineZoe):

```bash
# 1. Ve al repositorio backend-chat
cd tu-ruta-backend-chat/Modulo-Backend

# 2. Asegúrate de estar en la rama main
git checkout main

# 3. Copia el archivo backend-changes.patch a este directorio

# 4. Aplica el patch
git apply backend-changes.patch

# 5. Verifica los cambios
git status

# 6. Commit y push
git add -A
git commit -m "feat: Implementar endpoint GET /api/users y correcciones de integración"
git push origin main
```

---

## Opción 2: Aplicar Cambios Manualmente

### Cambios a realizar:

#### 1. Crear archivo: `Modulo-Backend/.env`
```
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=admin
DB_NAME=chatdb
JWT_SECRET=tu-secreto-super-seguro-cambiar-en-produccion-minimo-32-caracteres
```

#### 2. Crear archivo: `Modulo-Backend/internal/handler/user_handler.go`
```go
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
```

#### 3. Modificar: `Modulo-Backend/cmd/app/main.go`

**Cambiar esta sección:**
```go
repos := database.NewRepositories(db)
authHandler := handler.NewAuthHandler(repos.User, cfg)
```

**Por esto:**
```go
repos := database.NewRepositories(db)
authHandler := handler.NewAuthHandler(repos.User, cfg)
userHandler := handler.NewUserHandler(repos.User)
```

**Agregar esta ruta después del login:**
```go
r.GET("/api/users", userHandler.GetAllUsers)
```

#### 4. Modificar: `Modulo-Backend/internal/dto/auth.go`

**Agregar al final del archivo:**
```go
type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
```

#### 5. Modificar: `Modulo-Backend/internal/repository/user_repository.go`

**Agregar este método a la interfaz:**
```go
GetAll() ([]model.User, error)
```

#### 6. Modificar: `Modulo-Backend/internal/repository/impl/user_respository.go`

**Agregar esta implementación:**
```go
func (r *userRepository) GetAll() ([]model.User, error) {
	var users []model.User
	err := r.db.Select("id", "username", "email", "created_at", "updated_at").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
```

#### 7. Modificar: `Modulo-Backend/internal/model/model.go`

**Cambiar esta línea (en el struct User):**
```go
Password  string         `gorm:"size:255;not null" json:"-"` // Nunca devolver la contraseña
```

**Por esto:**
```go
Password  string         `gorm:"column:password_hash;size:255;not null" json:"-"` // Nunca devolver la contraseña
```

---

## Verificación

Después de aplicar los cambios:

```bash
# Compilar y levantar
go run cmd/app/main.go

# Verificar endpoint (en otra terminal)
curl http://localhost:8080/api/users

# Debería retornar un JSON con lista de usuarios
```

---

## Resumen de Cambios

- ✅ Nuevo handler para usuarios
- ✅ Nuevo endpoint GET /api/users
- ✅ Nuevo DTO UserResponse
- ✅ Extensión de UserRepository con GetAll()
- ✅ Bugfix: password_hash column mapping

