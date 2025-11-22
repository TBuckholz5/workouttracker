package exercise

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/TBuckholz5/workouttracker/internal/api/v1/exercise/dto"
	service "github.com/TBuckholz5/workouttracker/internal/service/exercise"
	"github.com/TBuckholz5/workouttracker/internal/service/exercise/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockExerciseService struct {
	mock.Mock
}

func (m *MockExerciseService) CreateExercise(ctx context.Context, req *dto.CreateExerciseRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockExerciseService) GetExercisesForUser(ctx context.Context, req *dto.GetExerciseForUserRequest) ([]models.Exercise, error) {
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
	reqBody := &dto.CreateExerciseRequest{Name: "Bench", TargetMuscle: "Chest"}
	mockSvc.On("CreateExercise", mock.Anything, reqBody).Return(nil)

	router := setupRouter(mockSvc)
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/exercise", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
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

func TestCreateExercise_ServiceError(t *testing.T) {
	mockSvc := new(MockExerciseService)
	reqBody := &dto.CreateExerciseRequest{Name: "Bench"}
	mockSvc.On("CreateExercise", mock.Anything, reqBody).Return(errors.New("fail"))

	router := setupRouter(mockSvc)
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
	reqPayload := &dto.GetExerciseForUserRequest{UserID: userID, Offset: offset, Limit: limit}
	resp := []models.Exercise{{ID: 1, Name: "Bench"}}
	mockSvc.On("GetExercisesForUser", mock.Anything, reqPayload).Return(resp, nil)

	router := setupRouterWithUserID(mockSvc, userID)
	req, _ := http.NewRequest("GET", "/exercise/user?offset="+strconv.Itoa(offset)+"&limit="+strconv.Itoa(limit), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Bench")
	mockSvc.AssertExpectations(t)
}

func TestGetExerciseForUser_DefaultParams(t *testing.T) {
	mockSvc := new(MockExerciseService)
	userID := int64(42)
	reqPayload := &dto.GetExerciseForUserRequest{UserID: userID, Offset: 0, Limit: 10}
	resp := []models.Exercise{{ID: 2, Name: "Squat"}}
	mockSvc.On("GetExercisesForUser", mock.Anything, reqPayload).Return(resp, nil)

	router := setupRouterWithUserID(mockSvc, userID)
	req, _ := http.NewRequest("GET", "/exercise/user", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Squat")
	mockSvc.AssertExpectations(t)
}

func TestGetExerciseForUser_ServiceError(t *testing.T) {
	mockSvc := new(MockExerciseService)
	userID := int64(42)
	offset := 0
	limit := 10
	reqPayload := &dto.GetExerciseForUserRequest{UserID: userID, Offset: offset, Limit: limit}
	mockSvc.On("GetExercisesForUser", mock.Anything, reqPayload).Return([]models.Exercise{}, errors.New("unauthorized"))

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
