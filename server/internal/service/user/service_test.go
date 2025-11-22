package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/TBuckholz5/workouttracker/internal/api/v1/user/dto"
	userRepo "github.com/TBuckholz5/workouttracker/internal/repository/user"
	"github.com/TBuckholz5/workouttracker/internal/service/user/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) CreateUser(ctx context.Context, params *userRepo.CreateUserParams) (models.User, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *mockUserRepo) GetUserForUsername(ctx context.Context, username string) (models.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(models.User), args.Error(1)
}

type mockHasher struct {
	mock.Mock
}

func (m *mockHasher) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *mockHasher) VerifyPassword(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

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

func TestCreateUser_Success(t *testing.T) {
	password := "password123"
	hashedPassword := "hashedpassword123"

	repo := &mockUserRepo{}
	repo.On("CreateUser", mock.Anything, mock.Anything).Return(models.User{ID: 1}, nil)

	hasher := &mockHasher{}
	hasher.On("HashPassword", password).Return(hashedPassword, nil)

	s := NewService(repo, hasher, nil)
	req := &dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: password,
	}
	err := s.CreateUser(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	repo.AssertNumberOfCalls(t, "CreateUser", 1)
	repo.AssertCalled(t, "CreateUser", mock.Anything, mock.MatchedBy(func(arg *userRepo.CreateUserParams) bool {
		return arg.Username == req.Username && arg.Email == req.Email && arg.PwHash == hashedPassword
	}))
}

func TestCreateUser_HashError(t *testing.T) {
	password := "password123"

	repo := &mockUserRepo{}
	repo.On("CreateUser", mock.Anything, mock.Anything).Return(models.User{ID: 1}, nil)

	hasher := &mockHasher{}
	hasher.On("HashPassword", password).Return("", fmt.Errorf("hash error"))

	s := NewService(repo, hasher, nil)
	req := &dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: password,
	}
	err := s.CreateUser(context.Background(), req)
	assert.NotEqual(t, err, nil)

	repo.AssertNumberOfCalls(t, "CreateUser", 0)
}

func TestAuthenticateUser_Success(t *testing.T) {
	password := "password123"
	hashedPassword := []byte("test")
	tokenString := "validtoken"

	repo := &mockUserRepo{}
	repo.On("GetUserForUsername", mock.Anything, mock.MatchedBy(func(arg string) bool {
		return arg == "testuser"
	})).Return(models.User{
		ID:     1,
		PwHash: string(hashedPassword),
	}, nil)

	hasher := &mockHasher{}
	hasher.On("VerifyPassword", string(hashedPassword), password).Return(nil)

	jwtService := &mockJwtService{}
	jwtService.On("GenerateJwt", int64(1)).Return(tokenString, nil)

	s := NewService(repo, hasher, jwtService)
	req := &dto.LoginRequest{
		Username: "testuser",
		Password: password,
	}
	token, err := s.AuthenticateUser(context.Background(), req)

	assert.Nil(t, err)
	assert.Equal(t, tokenString, token)
	repo.AssertNumberOfCalls(t, "GetUserForUsername", 1)
	repo.AssertCalled(t, "GetUserForUsername", mock.Anything, mock.MatchedBy(func(arg string) bool {
		return arg == "testuser"
	}))
}

func TestAuthenticateUser_UserNotFound(t *testing.T) {
	password := "password123"

	repo := &mockUserRepo{}
	repo.On("GetUserForUsername", mock.Anything, mock.MatchedBy(func(arg string) bool {
		return arg == "testuser"
	})).Return(models.User{}, fmt.Errorf("user not found"))

	s := NewService(repo, nil, nil)
	req := &dto.LoginRequest{
		Username: "testuser",
		Password: password,
	}
	_, err := s.AuthenticateUser(context.Background(), req)

	assert.NotNil(t, err)
	repo.AssertNumberOfCalls(t, "GetUserForUsername", 1)
}

func TestAuthenticateUser_PasswordMismatchError(t *testing.T) {
	password := "password123"
	hashedPassword := []byte("test")

	repo := &mockUserRepo{}
	repo.On("GetUserForUsername", mock.Anything, mock.MatchedBy(func(arg string) bool {
		return arg == "testuser"
	})).Return(models.User{
		ID:     1,
		PwHash: string(hashedPassword),
	}, nil)

	hasher := &mockHasher{}
	hasher.On("VerifyPassword", string(hashedPassword), password).Return(fmt.Errorf("passwords do not match"))

	s := NewService(repo, hasher, nil)
	req := &dto.LoginRequest{
		Username: "testuser",
		Password: password,
	}
	_, err := s.AuthenticateUser(context.Background(), req)

	assert.NotNil(t, err)
	repo.AssertNumberOfCalls(t, "GetUserForUsername", 1)
	repo.AssertCalled(t, "GetUserForUsername", mock.Anything, mock.MatchedBy(func(arg string) bool {
		return arg == "testuser"
	}))
	hasher.AssertNumberOfCalls(t, "VerifyPassword", 1)
	hasher.AssertCalled(t, "VerifyPassword", string(hashedPassword), password)
}

func TestAuthenticateUser_JwtError(t *testing.T) {
	password := "password123"
	hashedPassword := []byte("test")

	repo := &mockUserRepo{}
	repo.On("GetUserForUsername", mock.Anything, mock.MatchedBy(func(arg string) bool {
		return arg == "testuser"
	})).Return(models.User{
		ID:     1,
		PwHash: string(hashedPassword),
	}, nil)

	hasher := &mockHasher{}
	hasher.On("VerifyPassword", string(hashedPassword), password).Return(nil)

	jwtService := &mockJwtService{}
	jwtService.On("GenerateJwt", int64(1)).Return("", fmt.Errorf("jwt generation error"))

	s := NewService(repo, hasher, jwtService)
	req := &dto.LoginRequest{
		Username: "testuser",
		Password: password,
	}
	_, err := s.AuthenticateUser(context.Background(), req)

	assert.NotNil(t, err)
	repo.AssertNumberOfCalls(t, "GetUserForUsername", 1)
	repo.AssertCalled(t, "GetUserForUsername", mock.Anything, mock.MatchedBy(func(arg string) bool {
		return arg == "testuser"
	}))
	hasher.AssertNumberOfCalls(t, "VerifyPassword", 1)
	hasher.AssertCalled(t, "VerifyPassword", string(hashedPassword), password)
	jwtService.AssertNumberOfCalls(t, "GenerateJwt", 1)
	jwtService.AssertCalled(t, "GenerateJwt", int64(1))
}
