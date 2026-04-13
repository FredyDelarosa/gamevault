package main

import (
	"log"

	"github.com/gin-gonic/gin"

	_ "gv/docs"

	"gv/core/config"
	"gv/core/database"
	"gv/core/logger"
	"gv/domain/models"
	domainServices "gv/domain/services"
	"gv/infrastructure/adapters/handlers"
	infraRepositories "gv/infrastructure/adapters/repositories"
	"gv/infrastructure/middleware"
	"gv/infrastructure/router"
)

// @title GameVault API
// @version 1.0
// @description API para gestionar autenticacion y catalogo de juegos.
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Ingresa el token con formato: Bearer <token>

func main() {
	// Inicializar logger
	logger.Init()

	// Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Conectar a base de datos
	db, err := database.NewMySQLConnection(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	logger.Info("Database connected successfully")

	// Auto-migrar tablas (crear tabla device_tokens si no existe)
	if err := db.AutoMigrate(&models.DeviceToken{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	logger.Info("Database migration completed")

	// Configurar modo de Gin
	if gin.Mode() == gin.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	// Inicializar repositorios (infraestructura)
	userRepo := infraRepositories.NewUserRepository(db)
	gameRepo := infraRepositories.NewGameRepository(db)
	deviceTokenRepo := infraRepositories.NewDeviceTokenRepository(db)

	// Inicializar servicios (dominio)
	notificationService := domainServices.NewNotificationService(deviceTokenRepo, cfg)
	authService := domainServices.NewAuthService(userRepo, cfg, notificationService)
	gameService := domainServices.NewGameService(gameRepo, notificationService)

	// Inicializar handlers
	authHandler := handlers.NewAuthHandler(authService)
	gameHandler := handlers.NewGameHandler(gameService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// Inicializar middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Configurar router
	r := router.NewRouter(authHandler, gameHandler, notificationHandler, authMiddleware)
	engine := r.Setup()

	// Iniciar servidor
	logger.Info("Server starting on port %s", cfg.ServerPort)
	if err := engine.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
