package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TBuckholz5/workouttracker/internal/workoutsession/models"
	"github.com/TBuckholz5/workouttracker/internal/workoutsession/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWorkoutSessionService struct {
	mock.Mock
}

func (m *MockWorkoutSessionService) Create(ctx context.Context, session *models.WorkoutSession) (*models.WorkoutSession, error) {
	args := m.Called(ctx, session)
	return args.Get(0).(*models.WorkoutSession), args.Error(1)
}

func setupRouterWithUserID(svc service.WorkoutSessionService, userID int64) *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := NewHandler(svc)
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.POST("/workout-session", h.Create)
	return r
}

func setupRouter(svc service.WorkoutSessionService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := NewHandler(svc)
	r := gin.Default()
	r.POST("/workout-session", h.Create)
	return r
}

func TestCreate_Success(t *testing.T) {
	mockSvc := new(MockWorkoutSessionService)
	userID := int64(42)
	reqBody := map[string]any{
		"name":        "Morning Workout",
		"description": "Upper body workout",
		"duration":    60,
		"workouts": []map[string]any{
			{
				"exerciseID":  1,
				"description": "Bench press",
				"sets": []map[string]any{
					{
						"reps":     10,
						"weight":   135.5,
						"setType":  "working",
						"setOrder": 1,
					},
				},
			},
		},
	}

	expectedPayload := &models.WorkoutSession{
		UserID:      userID,
		Name:        "Morning Workout",
		Description: "Upper body workout",
		Duration:    60,
	}

	expectedResponse := &models.WorkoutSession{
		ID:          1,
		UserID:      userID,
		Name:        "Morning Workout",
		Description: "Upper body workout",
		Duration:    60,
	}

	mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(session *models.WorkoutSession) bool {
		return session.UserID == expectedPayload.UserID &&
			session.Name == expectedPayload.Name &&
			session.Description == expectedPayload.Description &&
			session.Duration == expectedPayload.Duration
	})).Return(expectedResponse, nil)

	router := setupRouterWithUserID(mockSvc, userID)
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/workout-session", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"Morning Workout"`)
	mockSvc.AssertExpectations(t)
}

func TestCreate_BadRequest(t *testing.T) {
	mockSvc := new(MockWorkoutSessionService)
	router := setupRouter(mockSvc)

	req, _ := http.NewRequest("POST", "/workout-session", bytes.NewBuffer([]byte("bad json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreate_MissingUserID(t *testing.T) {
	mockSvc := new(MockWorkoutSessionService)
	reqBody := map[string]any{
		"name":        "Morning Workout",
		"description": "Upper body workout",
		"duration":    60,
	}
	body, _ := json.Marshal(reqBody)

	gin.SetMode(gin.TestMode)
	h := NewHandler(mockSvc)
	r := gin.Default()
	r.POST("/workout-session", func(c *gin.Context) {
		h.Create(c)
	})

	req, _ := http.NewRequest("POST", "/workout-session", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "userID not found in context")
}

func TestCreate_ServiceError(t *testing.T) {
	mockSvc := new(MockWorkoutSessionService)
	userID := int64(42)
	reqBody := map[string]any{
		"name":        "Morning Workout",
		"description": "Upper body workout",
		"duration":    60,
	}

	mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(session *models.WorkoutSession) bool {
		return session.UserID == userID
	})).Return((*models.WorkoutSession)(nil), errors.New("database error"))

	router := setupRouterWithUserID(mockSvc, userID)
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/workout-session", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "database error")
	mockSvc.AssertExpectations(t)
}

func TestCreate_EmptyPayload(t *testing.T) {
	mockSvc := new(MockWorkoutSessionService)
	userID := int64(42)
	reqBody := map[string]any{}

	expectedResponse := &models.WorkoutSession{
		ID:     1,
		UserID: userID,
	}

	mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(session *models.WorkoutSession) bool {
		return session.UserID == userID
	})).Return(expectedResponse, nil)

	router := setupRouterWithUserID(mockSvc, userID)
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/workout-session", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}
