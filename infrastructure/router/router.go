package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"gv/infrastructure/adapters/handlers"
	"gv/infrastructure/middleware"
)

type Router struct {
	authHandler         *handlers.AuthHandler
	gameHandler         *handlers.GameHandler
	notificationHandler *handlers.NotificationHandler
	authMiddleware      *middleware.AuthMiddleware
}

func NewRouter(
	authHandler *handlers.AuthHandler,
	gameHandler *handlers.GameHandler,
	notificationHandler *handlers.NotificationHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	return &Router{
		authHandler:         authHandler,
		gameHandler:         gameHandler,
		notificationHandler: notificationHandler,
		authMiddleware:      authMiddleware,
	}
}

func (r *Router) Setup() *gin.Engine {
	router := gin.Default()

	// Configurar CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Rutas públicas
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
		}
	}

	// Rutas protegidas (requieren autenticación)
	protected := api.Group("/")
	protected.Use(r.authMiddleware.Authenticate())
	{
		// Games CRUD
		protected.GET("/games", r.gameHandler.GetGames)
		protected.GET("/games/:id", r.gameHandler.GetGameByID)
		protected.POST("/games", r.gameHandler.CreateGame)
		protected.PUT("/games/:id", r.gameHandler.UpdateGame)
		protected.DELETE("/games/:id", r.gameHandler.DeleteGame)

		// Notifications
		protected.POST("/notifications/token", r.notificationHandler.SaveFcmToken)
	}

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}