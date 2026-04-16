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
// @description API para gestionar autenticacion, catalogo de juegos y foro comunitario.
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	logger.Init()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err := database.NewMySQLConnection(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	logger.Info("Database connected successfully")

	// Auto-migrar tablas
	if err := db.AutoMigrate(
		&models.DeviceToken{},
		&models.Post{},
		&models.PostReaction{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	logger.Info("Database migration completed")

	if gin.Mode() == gin.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	// Repositorios
	userRepo := infraRepositories.NewUserRepository(db)
	gameRepo := infraRepositories.NewGameRepository(db)
	deviceTokenRepo := infraRepositories.NewDeviceTokenRepository(db)
	postRepo := infraRepositories.NewPostRepository(db)
	postReactionRepo := infraRepositories.NewPostReactionRepository(db)

	// Servicios (notificationService primero, los demás dependen de él)
	notificationService := domainServices.NewNotificationService(deviceTokenRepo, cfg)
	authService := domainServices.NewAuthService(userRepo, cfg, notificationService)
	gameService := domainServices.NewGameService(gameRepo, notificationService)
	postService := domainServices.NewPostService(postRepo, postReactionRepo, gameRepo, userRepo, notificationService)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	gameHandler := handlers.NewGameHandler(gameService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	postHandler := handlers.NewPostHandler(postService, userRepo)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	r := router.NewRouter(authHandler, gameHandler, notificationHandler, postHandler, authMiddleware)
	engine := r.Setup()

	logger.Info("Server starting on port %s", cfg.ServerPort)
	if err := engine.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
