package v1

import (
	"encoding/json"
	"net/http"

	"github.com/TBuckholz5/workouttracker/internal/domains/user/service"
)

type Handler struct {
	service service.UserService
}

func NewHandler(s service.UserService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var payload RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.service.CreateUser(r.Context(), &service.RegisterParams{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var payload LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token, err := h.service.AuthenticateUser(r.Context(), &service.LoginParams{
		Username: payload.Username,
		Password: payload.Password,
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(LoginResponse{Token: token}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
