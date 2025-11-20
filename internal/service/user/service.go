package user

import (
	"context"
	"fmt"

	"github.com/TBuckholz5/workouttracker/internal/api/v1/user/dto"
	repo "github.com/TBuckholz5/workouttracker/internal/repository/user"
	"github.com/TBuckholz5/workouttracker/internal/util/hash"
	"github.com/TBuckholz5/workouttracker/internal/util/jwt"
)

type UserService interface {
	CreateUser(reqContext context.Context, userDto *dto.RegisterRequest) error
	AuthenticateUser(reqContext context.Context, loginDto *dto.LoginRequest) (string, error)
}

type Service struct {
	repo       repo.UserRepository
	jwtService jwt.JwtService
	hasher     hash.Hasher
}

func NewService(r repo.UserRepository, hasher hash.Hasher, jwtService jwt.JwtService) *Service {
	return &Service{
		repo:       r,
		hasher:     hasher,
		jwtService: jwtService,
	}
}

func (s *Service) CreateUser(reqContext context.Context, userDto *dto.RegisterRequest) error {
	hashedPassword, err := s.hasher.HashPassword(userDto.Password)
	if err != nil {
		return err
	}
	_, err = s.repo.CreateUser(reqContext, &repo.CreateUserParams{
		Username: userDto.Username,
		Email:    userDto.Email,
		PwHash:   hashedPassword,
	})
	return err
}

func (s *Service) AuthenticateUser(reqContext context.Context, loginDto *dto.LoginRequest) (string, error) {
	user, err := s.repo.GetUserForUsername(reqContext, loginDto.Username)
	if err != nil {
		return "", err
	}

	err = s.hasher.VerifyPassword(user.PwHash, loginDto.Password)
	if err != nil {
		return "", fmt.Errorf("passwords do not match")
	}

	token, err := s.jwtService.GenerateJwt(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
