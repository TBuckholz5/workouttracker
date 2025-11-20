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

func (m *mockJwtService) ValidateJwt(ctx *gin.Context, tokenString string) error {
	args := m.Called(ctx, tokenString)
	return args.Error(0)
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authHeader     string
		validateReturn error
		expectedStatus int
	}{
		{
			name:           "No Authorization header",
			authHeader:     "",
			validateReturn: nil,
			expectedStatus: 401,
		},
		{
			name:           "Invalid prefix",
			authHeader:     "Token sometoken",
			validateReturn: nil,
			expectedStatus: 401,
		},
		{
			name:           "Invalid JWT",
			authHeader:     "Bearer invalidtoken",
			validateReturn: fmt.Errorf("invalid token"),
			expectedStatus: 401,
		},
		{
			name:           "Valid JWT",
			authHeader:     "Bearer validtoken",
			validateReturn: nil,
			expectedStatus: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			mockService := &mockJwtService{}
			mockService.On("ValidateJwt", mock.Anything, mock.Anything).Return(tt.validateReturn)

			r.Use(AuthMiddleware(mockService))
			r.GET("/test", func(c *gin.Context) {
				c.Status(200)
			})

			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
