package routes

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/redis/go-redis/v9"
	"github.com/ryuudan/golang-rest-api/ent/generated"
	"github.com/ryuudan/golang-rest-api/src/database"
	"github.com/ryuudan/golang-rest-api/src/internal/handlers"
	"github.com/ryuudan/golang-rest-api/src/internal/repositories"
	"github.com/ryuudan/golang-rest-api/src/internal/services"
)

func PublicRouter(client *generated.Client, redis_client *redis.Client) http.Handler {
	public := chi.NewRouter()

	// 100 requests per minute
	public.Use(httprate.LimitByIP(100, 1*time.Minute))

	public.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the public API"))
	})

	return public
}

func PrivateRouter(client *generated.Client, redis_client *redis.Client) http.Handler {
	private := chi.NewRouter()

	// Add authentication middleware here
	// Add authorization middleware here
	// Add rate limiting middleware here

	// Initialize handlers
	cache := database.Cache(redis_client)

	// repositories
	userRepo := repositories.NewUserRepository(client.User)

	// services
	userService := services.NewUserService(userRepo)

	// handlers
	userHandler := handlers.NewUserHandler(userService, cache)

	private.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the private API"))
	})

	private.Route("/users", func(r chi.Router) {
		r.Get("/{id}", userHandler.GetOneByID)
		r.Post("/", userHandler.Create)
	})

	return private
}
