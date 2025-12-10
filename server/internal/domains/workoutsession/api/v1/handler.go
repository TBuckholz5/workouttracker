package v1

import (
	"encoding/json"
	"net/http"

	"github.com/TBuckholz5/workouttracker/internal/domains/workoutsession/models"
	"github.com/TBuckholz5/workouttracker/internal/domains/workoutsession/service"
)

type Handler struct {
	service service.WorkoutSessionService
}

func NewHandler(s service.WorkoutSessionService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var payload models.WorkoutSession
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID := r.Context().Value("userID")
	if userID == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	payload.UserID = userID.(int64)
	session, err := h.service.Create(r.Context(), &payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(CreateWorkoutSessionResponse{
		Session: *session,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
