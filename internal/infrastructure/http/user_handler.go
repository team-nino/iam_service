package http

import (
	"net/http"
	"strconv"

	"github.com/team-nino/iam_service/internal/usecase"
)

// UserHandler handles user management HTTP requests
type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// GetUser handles retrieving a user by ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, Response{
			Success: false,
			Error:   "method not allowed",
		})
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		respondJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "missing user id",
		})
		return
	}

	user, err := h.userUseCase.GetUser(r.Context(), id)
	if err != nil {
		respondJSON(w, http.StatusNotFound, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    user,
	})
}

// ListUsers handles retrieving a list of users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, Response{
			Success: false,
			Error:   "method not allowed",
		})
		return
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if limit == 0 {
		limit = 10
	}

	users, err := h.userUseCase.ListUsers(r.Context(), offset, limit)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    users,
	})
}
