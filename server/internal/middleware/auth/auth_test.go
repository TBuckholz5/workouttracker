package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockJwtService struct {
	mock.Mock
}

func (m *mockJwtService) GenerateJwt(userID int64) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *mockJwtService) ValidateJwt(tokenString string) (int64, error) {
	args := m.Called(tokenString)
	return args.Get(0).(int64), args.Error(1)
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authHeader     string
		validateUserID int64
		validateErr    error
		expectedStatus int
		expectUserID   bool
	}{
		{
			name:           "No Authorization header",
			authHeader:     "",
			validateUserID: 0,
			validateErr:    nil,
			expectedStatus: 401,
			expectUserID:   false,
		},
		{
			name:           "Invalid prefix",
			authHeader:     "Token sometoken",
			validateUserID: 0,
			validateErr:    nil,
			expectedStatus: 401,
			expectUserID:   false,
		},
		{
			name:           "Invalid JWT",
			authHeader:     "Bearer invalidtoken",
			validateUserID: 0,
			validateErr:    fmt.Errorf("invalid token"),
			expectedStatus: 401,
			expectUserID:   false,
		},
		{
			name:           "Valid JWT",
			authHeader:     "Bearer validtoken",
			validateUserID: 42,
			validateErr:    nil,
			expectedStatus: 200,
			expectUserID:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			mockService := &mockJwtService{}
			mockService.On("ValidateJwt", mock.Anything).Return(tt.validateUserID, tt.validateErr)

			r.Use(AuthMiddleware(mockService))
			r.GET("/test", func(c *gin.Context) {
				// Check if userID is set for valid JWT
				if tt.expectUserID {
					val, exists := c.Get("userID")
					assert.True(t, exists)
					assert.Equal(t, tt.validateUserID, val)
				}
				c.Status(200)
			})

			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}
