package http

import (
	"net/http"
)

// Router sets up the HTTP routes
type Router struct {
	authHandler *AuthHandler
	userHandler *UserHandler
}

// NewRouter creates a new router
func NewRouter(authHandler *AuthHandler, userHandler *UserHandler) *Router {
	return &Router{
		authHandler: authHandler,
		userHandler: userHandler,
	}
}

// SetupRoutes sets up all HTTP routes
func (router *Router) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, Response{
			Success: true,
			Data:    "service is healthy",
		})
	})

	// Authentication routes
	mux.HandleFunc("/api/auth/register", router.authHandler.Register)
	mux.HandleFunc("/api/auth/login", router.authHandler.Login)
	mux.HandleFunc("/api/auth/logout", router.authHandler.Logout)

	// User management routes
	mux.HandleFunc("/api/users", router.userHandler.ListUsers)
	mux.HandleFunc("/api/users/get", router.userHandler.GetUser)

	return mux
}
