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
	postHandler         *handlers.PostHandler
	authMiddleware      *middleware.AuthMiddleware
}

func NewRouter(
	authHandler *handlers.AuthHandler,
	gameHandler *handlers.GameHandler,
	notificationHandler *handlers.NotificationHandler,
	postHandler *handlers.PostHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	return &Router{
		authHandler:         authHandler,
		gameHandler:         gameHandler,
		notificationHandler: notificationHandler,
		postHandler:         postHandler,
		authMiddleware:      authMiddleware,
	}
}

func (r *Router) Setup() *gin.Engine {
	router := gin.Default()

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

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
		}
	}

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
		protected.POST("/notifications/send", r.notificationHandler.SendNotification)

		// Posts (Foro/Comunidad)
		protected.POST("/posts", r.postHandler.CreatePost)
		protected.GET("/posts", r.postHandler.GetAllPosts)
		protected.GET("/posts/my-games", r.postHandler.GetPostsForMyGames)
		protected.GET("/posts/:id", r.postHandler.GetPostByID)
		protected.DELETE("/posts/:id", r.postHandler.DeletePost)
		protected.POST("/posts/:id/react", r.postHandler.ToggleReaction)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
