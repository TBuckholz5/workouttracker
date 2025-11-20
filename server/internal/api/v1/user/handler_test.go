package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TBuckholz5/workouttracker/internal/api/v1/user/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	createUserErr     error
	authenticateToken string
	authenticateErr   error
}

func (m *mockService) CreateUser(ctx context.Context, req *dto.RegisterRequest) error {
	return m.createUserErr
}

func (m *mockService) AuthenticateUser(ctx context.Context, req *dto.LoginRequest) (string, error) {
	return m.authenticateToken, m.authenticateErr
}

func setupRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	return r
}

func TestRegister_Success(t *testing.T) {
	svc := &mockService{}
	h := NewHandler(svc)
	router := setupRouter(h)

	payload := dto.RegisterRequest{Username: "test", Email: "test@gmail.com", Password: "passwordtest"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestRegister_BadRequest(t *testing.T) {
	svc := &mockService{}
	h := NewHandler(svc)
	router := setupRouter(h)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(`bad json`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestRegister_ServiceError(t *testing.T) {
	svc := &mockService{createUserErr: errors.New("fail")}
	h := NewHandler(svc)
	router := setupRouter(h)

	payload := dto.RegisterRequest{Username: "test", Email: "test@gmail.com", Password: "passwordtest"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, 500, w.Code)
}

func TestLogin_Success(t *testing.T) {
	svc := &mockService{authenticateToken: "token123"}
	h := NewHandler(svc)
	router := setupRouter(h)

	payload := dto.LoginRequest{Username: "test", Password: "pass"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "token123")
}

func TestLogin_BadRequest(t *testing.T) {
	svc := &mockService{}
	h := NewHandler(svc)
	router := setupRouter(h)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(`bad json`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestLogin_AuthError(t *testing.T) {
	svc := &mockService{authenticateErr: errors.New("unauthorized")}
	h := NewHandler(svc)
	router := setupRouter(h)

	payload := dto.LoginRequest{Username: "test", Password: "pass"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
}
