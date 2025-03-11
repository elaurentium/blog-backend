package main

import (
	"net/http"
	"time"

	"github.com/elaurentium/exilium-blog-backend/internal/domain/services"
	"github.com/elaurentium/exilium-blog-backend/internal/infra/api"
	"github.com/elaurentium/exilium-blog-backend/internal/infra/api/handlers"
	"github.com/elaurentium/exilium-blog-backend/internal/infra/api/middleware"
	"github.com/elaurentium/exilium-blog-backend/internal/infra/auth"
	"github.com/elaurentium/exilium-blog-backend/internal/infra/persistence/db"
	"github.com/elaurentium/exilium-blog-backend/internal/infra/persistence/redis"
	"github.com/elaurentium/exilium-blog-backend/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	logger := logger.NewLogger()
	logger.Info("Starting application")

	// Inicializa o Redis
	redisClient, err := redis.NewRedisClient()
	if err != nil {
		logger.Info("Failed to connect to Redis: %v", err)
		return
	}
	defer redisClient.Close()

	// Inicializa o PostgreSQL
	pool, err := db.NewPostgresPool()
	if err != nil {
		logger.Info("Failed to connect to PostgreSQL: %v", err)
		return
	}
	defer pool.Close()

	// Inicializa os repositórios e serviços
	userRepo := db.NewUserRepository(pool)
	authService := auth.NewAuthService()
	userService := services.NewUserService(userRepo, authService)
	userHandler := handlers.NewUserHandler(userService)

	postHandler := &handlers.PostHandler{}
	commentHandler := &handlers.CommentHandler{}
	subHandler := &handlers.SubHandler{}
	authMiddleware := &middleware.AuthMiddleware{}

	// Cria o roteador
	router := api.NewRouter(userHandler, postHandler, commentHandler, subHandler, authMiddleware, redisClient)

	// Inicia o servidor HTTP
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Starting server on %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Info("Failed to start server: %v", err)
	}
}