package api

import (
	"github.com/gin-gonic/gin"
	"github.com/elaurentium/exilium-blog-backend/internal/infra/api/handlers"
	"github.com/elaurentium/exilium-blog-backend/internal/infra/api/middleware"
	"github.com/elaurentium/exilium-blog-backend/internal/infra/persistence/redis"
	"github.com/elaurentium/exilium-blog-backend/pkg/logger"
)

func NewRouter(
	userHandler *handlers.UserHandler,
	postHandler *handlers.PostHandler,
	commentHandler *handlers.CommentHandler,
	subHandler *handlers.SubHandler,
	authMiddleware *middleware.AuthMiddleware,
	redisClient *redis.RedisClient,
) *gin.Engine {
	logger := logger.NewLogger()

	router := gin.Default()

	// Middleware
	router.Use(middleware.CorsMiddleware())
	router.Use(middleware.LoggerMiddleware(logger))
	router.Use(middleware.SecurityMiddleware())
	router.Use(middleware.NewRateLimiterMiddleware(redisClient.GetClient()))

	// Public routes
	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)

	// Protected routes
	authGroup := router.Group("/")
	authGroup.Use(authMiddleware.Authenticate(), authMiddleware.RequireRole("user"))
	{
		authGroup.GET("/profile", userHandler.GetProfile)
		authGroup.PUT("/profile", userHandler.UpdateProfile)
		authGroup.POST("/posts", postHandler.CreatePost)
		authGroup.PUT("/posts/:id", postHandler.UpdatePost)
		authGroup.DELETE("/posts/:id", postHandler.DeletePost)
		authGroup.POST("/comments", commentHandler.CreateComment)
		authGroup.PUT("/comments/:id", commentHandler.UpdateComment)
		authGroup.DELETE("/comments/:id", commentHandler.DeleteComment)
		authGroup.POST("/sub", subHandler.CreateSub)
		authGroup.PUT("/sub/:id", subHandler.UpdateSub)
		authGroup.DELETE("/sub/:id", subHandler.DeleteSub)
	}

	return router
}