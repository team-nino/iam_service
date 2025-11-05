package http

import (
	"encoding/json"
	"net/http"

	"github.com/team-nino/iam_service/internal/usecase"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username"`
	Password        string `json:"password"`
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, Response{
			Success: false,
			Error:   "method not allowed",
		})
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid request body",
		})
		return
	}

	user, err := h.authUseCase.Register(r.Context(), req.Email, req.Username, req.Password)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusCreated, Response{
		Success: true,
		Data:    user,
	})
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, Response{
			Success: false,
			Error:   "method not allowed",
		})
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid request body",
		})
		return
	}

	session, err := h.authUseCase.Login(r.Context(), req.EmailOrUsername, req.Password)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    session,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, Response{
			Success: false,
			Error:   "method not allowed",
		})
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		respondJSON(w, http.StatusUnauthorized, Response{
			Success: false,
			Error:   "missing authorization token",
		})
		return
	}

	if err := h.authUseCase.Logout(r.Context(), token); err != nil {
		respondJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    "logged out successfully",
	})
}
