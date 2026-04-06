# Backend - Sistema de Chat en Tiempo Real

Repositorio del Backend del sistema de chat desarrollado en **Golang**.

**Responsable:** Gabriel Flores

## Tecnologías Utilizadas

- Go (Golang)
- Gin Gonic (Framework Web)
- GORM + PostgreSQL
- Gorilla WebSocket (Chat en tiempo real)
- JWT (Autenticación)
- Viper (Configuración)

## Funcionalidades Actuales

- Registro de usuarios
- Login con generación de JWT
- Chat en tiempo real usando WebSocket (solo salas públicas por ahora)
- Persistencia de mensajes en base de datos

## Cómo Montar y Ejecutar el Proyecto

### 1. Requisitos Previos

- Go 1.23 o superior instalado
- PostgreSQL instalado y corriendo
- Base de datos creada llamada `chatdb`

### 2. Configuración

1. Clona el repositorio:
   ```bash
   git clone https://github.com/MelusineZoe/backend-chat.git
   cd backend-chat

Instala las dependencias:bash

go mod tidy

Crea el archivo .env en la raíz del proyecto con el siguiente contenido:

env

SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=tu_contraseña_postgres
DB_NAME=chatdb
JWT_SECRET=supersecretkey2026cambialoencproduccion

Importante: Cambia tu_contraseña_postgres por la contraseña real de tu usuario postgres.
3. Ejecutar el Backendbash

go run cmd/app/main.go

Si todo está correcto verás:

✅ Conexión exitosa a PostgreSQL
✅ Migración de tablas completada
🚀 Servidor corriendo en http://localhost:8080

Cómo Probar el Sistema1. AutenticaciónRegistro:bash

POST http://localhost:8080/api/auth/register

Body (JSON):json

{
  "username": "gabriel",
  "email": "gabriel@test.com",
  "password": "123456"
}

Login:bash

POST http://localhost:8080/api/auth/login

Body (JSON):json

{
  "email": "gabriel@test.com",
  "password": "123456"
}

Guarda el token que te devuelve.2. Chat en Tiempo Real (WebSocket)Crea una sala pública (por ahora puedes crearla manualmente en la BD o agregar el endpoint después).
Conéctate al WebSocket usando la siguiente URL:

ws://localhost:8080/ws?room_id=UUID_DE_LA_SALA

Ejemplo de mensaje a enviar:json

{
  "content": "Hola, esto es un mensaje de prueba"
}

Los mensajes se enviarán en tiempo real a todos los usuarios conectados en la misma sala.Estructura del Proyecto

backend/
├── cmd/app/main.go                 → Punto de entrada
├── internal/
│   ├── config/                     → Configuración (.env)
│   ├── database/                   → Conexión y migraciones
│   ├── dto/                        → Request y Response
│   ├── handler/                    → Handlers (Auth + WebSocket)
│   ├── middleware/                 → JWT y CORS
│   ├── model/                      → Modelos de la BD
│   ├── repository/                 → Capa de datos
│   ├── service/                    → Lógica de negocio
│   └── ws/                         → WebSocket (Hub + Client)
└── .env

Cómo funciona el Chat en Tiempo RealEl usuario se autentica y obtiene un JWT.
Se conecta al endpoint /ws?room_id=xxx enviando el token en el header Authorization: Bearer <token>.
El Hub (un singleton en Go) maneja todas las conexiones activas.
Cuando un usuario envía un mensaje:Se guarda en la base de datos
Se transmite en tiempo real a todos los usuarios de la misma sala

