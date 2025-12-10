package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/TBuckholz5/workouttracker/internal/domains/exercise/service"
)

type Handler struct {
	service service.ExerciseService
}

func NewHandler(s service.ExerciseService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateExercise(w http.ResponseWriter, r *http.Request) {
	var payload CreateExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID := r.Context().Value("userID")
	if userID == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	params := service.CreateExerciseForUserParams{
		UserID:       userID.(int64),
		Name:         payload.Name,
		Description:  payload.Description,
		TargetMuscle: payload.TargetMuscle,
	}
	exercise, err := h.service.CreateExercise(r.Context(), &params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(CreateExerciseResponse{
		Exercise: Exercise{
			ID:           exercise.ID,
			Name:         exercise.Name,
			Description:  exercise.Description,
			TargetMuscle: exercise.TargetMuscle,
			PictureURL:   exercise.PictureURL,
		},
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetExerciseForUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	if userID == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	queryParams := r.URL.Query()
	offset := 0
	limit := 10
	if val := queryParams.Get("offset"); val != "" {
		parsedOffset, err := strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		offset = parsedOffset
	}
	if val := queryParams.Get("limit"); val != "" {
		parsedLimit, err := strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		limit = parsedLimit
	}
	payload := service.GetExerciseForUserParams{
		UserID: userID.(int64),
		Offset: offset,
		Limit:  limit,
	}
	exercises, err := h.service.GetExercisesForUser(r.Context(), &payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	exercisesDTO := []Exercise{}
	for _, ex := range exercises {
		exercisesDTO = append(exercisesDTO, Exercise{
			ID:           ex.ID,
			Name:         ex.Name,
			Description:  ex.Description,
			TargetMuscle: ex.TargetMuscle,
			PictureURL:   ex.PictureURL,
		})
	}
	if err := json.NewEncoder(w).Encode(GetExerciseListResponse{Exercises: exercisesDTO}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
