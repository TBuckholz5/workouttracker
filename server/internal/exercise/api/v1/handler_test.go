package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/TBuckholz5/workouttracker/internal/exercise/models"
	"github.com/TBuckholz5/workouttracker/internal/exercise/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockExerciseService struct {
	mock.Mock
}

func (m *MockExerciseService) CreateExercise(ctx context.Context, req *models.CreateExerciseForUserParams) (models.Exercise, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(models.Exercise), args.Error(1)
}

func (m *MockExerciseService) GetExercisesForUser(ctx context.Context, req *models.GetExerciseForUserParams) ([]models.Exercise, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]models.Exercise), args.Error(1)
}

func setupRouterWithUserID(svc service.ExerciseService, userID int64) *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := NewHandler(svc)
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.POST("/exercise", h.CreateExercise)
	r.GET("/exercise/user", h.GetExerciseForUser)
	return r
}

func setupRouter(svc service.ExerciseService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := NewHandler(svc)
	r := gin.Default()
	r.POST("/exercise", h.CreateExercise)
	r.GET("/exercise/user", h.GetExerciseForUser)
	return r
}

func TestCreateExercise_Success(t *testing.T) {
	mockSvc := new(MockExerciseService)
	userID := int64(42)
	reqBody := map[string]any{
		"name":         "Bench",
		"description":  "Chest exercise",
		"targetMuscle": "Chest",
	}
	params := &models.CreateExerciseForUserParams{
		UserID:       userID,
		Name:         "Bench",
		Description:  "Chest exercise",
		TargetMuscle: "Chest",
	}
	expectedExercise := models.Exercise{ID: 1, Name: "Bench", Description: "Chest exercise", TargetMuscle: "Chest"}
	mockSvc.On("CreateExercise", mock.Anything, params).Return(expectedExercise, nil)

	router := setupRouterWithUserID(mockSvc, userID)
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/exercise", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"Bench"`)
	mockSvc.AssertExpectations(t)
}

func TestCreateExercise_BadRequest(t *testing.T) {
	mockSvc := new(MockExerciseService)
	router := setupRouter(mockSvc)

	req, _ := http.NewRequest("POST", "/exercise", bytes.NewBuffer([]byte("bad json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateExercise_MissingUserID(t *testing.T) {
	mockSvc := new(MockExerciseService)
	reqBody := map[string]any{
		"name":         "Bench",
		"description":  "Chest exercise",
		"targetMuscle": "Chest",
	}
	body, _ := json.Marshal(reqBody)

	gin.SetMode(gin.TestMode)
	h := NewHandler(mockSvc)
	r := gin.Default()
	r.POST("/exercise", func(c *gin.Context) {
		h.CreateExercise(c)
	})

	req, _ := http.NewRequest("POST", "/exercise", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "userID not found in context")
}

func TestCreateExercise_ServiceError(t *testing.T) {
	mockSvc := new(MockExerciseService)
	userID := int64(42)
	reqBody := map[string]any{
		"name":         "Bench",
		"description":  "Chest exercise",
		"targetMuscle": "Chest",
	}
	params := &models.CreateExerciseForUserParams{
		UserID:       userID,
		Name:         "Bench",
		Description:  "Chest exercise",
		TargetMuscle: "Chest",
	}
	mockSvc.On("CreateExercise", mock.Anything, params).Return(models.Exercise{}, errors.New("fail"))

	router := setupRouterWithUserID(mockSvc, userID)
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/exercise", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetExerciseForUser_Success(t *testing.T) {
	mockSvc := new(MockExerciseService)
	userID := int64(42)
	offset := 0
	limit := 10
	params := &models.GetExerciseForUserParams{UserID: userID, Offset: offset, Limit: limit}
	resp := []models.Exercise{{ID: 1, Name: "Bench"}}
	mockSvc.On("GetExercisesForUser", mock.Anything, params).Return(resp, nil)

	router := setupRouterWithUserID(mockSvc, userID)
	req, _ := http.NewRequest("GET", "/exercise/user?offset="+strconv.Itoa(offset)+"&limit="+strconv.Itoa(limit), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Bench")
	mockSvc.AssertExpectations(t)
}

func TestGetExerciseForUser_EmptyResult(t *testing.T) {
	mockSvc := new(MockExerciseService)
	userID := int64(42)
	params := &models.GetExerciseForUserParams{UserID: userID, Offset: 0, Limit: 10}
	resp := []models.Exercise{}
	mockSvc.On("GetExercisesForUser", mock.Anything, params).Return(resp, nil)

	router := setupRouterWithUserID(mockSvc, userID)
	req, _ := http.NewRequest("GET", "/exercise/user", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"exercises":[]`)
	mockSvc.AssertExpectations(t)
}

func TestGetExerciseForUser_ServiceError(t *testing.T) {
	mockSvc := new(MockExerciseService)
	userID := int64(42)
	offset := 0
	limit := 10
	params := &models.GetExerciseForUserParams{UserID: userID, Offset: offset, Limit: limit}
	mockSvc.On("GetExercisesForUser", mock.Anything, params).Return([]models.Exercise{}, errors.New("unauthorized"))

	router := setupRouterWithUserID(mockSvc, userID)
	req, _ := http.NewRequest("GET", "/exercise/user?offset="+strconv.Itoa(offset)+"&limit="+strconv.Itoa(limit), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetExerciseForUser_MissingUserID(t *testing.T) {
	mockSvc := new(MockExerciseService)
	offset := 0
	limit := 10

	gin.SetMode(gin.TestMode)
	h := NewHandler(mockSvc)
	r := gin.Default()
	r.GET("/exercise/user", func(c *gin.Context) {
		h.GetExerciseForUser(c)
	})

	req, _ := http.NewRequest("GET", "/exercise/user?offset="+strconv.Itoa(offset)+"&limit="+strconv.Itoa(limit), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "userID not found in context")
}
