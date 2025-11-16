package user

import (
	"context"
	"fmt"

	"github.com/TBuckholz5/workouttracker/internal/api/v1/user/dto"
	db "github.com/TBuckholz5/workouttracker/internal/db/user"
	"github.com/TBuckholz5/workouttracker/internal/jwt"
	repo "github.com/TBuckholz5/workouttracker/internal/repository/user"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *repo.Repository
}

func NewService(r *repo.Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) CreateUser(reqContext context.Context, userDto *dto.RegisterRequest) error {
	hashedPassword, err := hashPassword(userDto.Password)
	if err != nil {
		return err
	}
	_, err = s.repo.CreateUser(reqContext, &db.CreateUserParams{
		Username: pgtype.Text{String: userDto.Username, Valid: true},
		Email:    pgtype.Text{String: userDto.Email, Valid: true},
		PwHash:   pgtype.Text{String: hashedPassword, Valid: true},
	})
	return err
}

func (s *Service) AuthenticateUser(reqContext context.Context, loginDto *dto.LoginRequest) (string, error) {
	user, err := s.repo.GetUserForUsername(reqContext, pgtype.Text{String: loginDto.Username, Valid: true})
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PwHash.String), []byte(loginDto.Password))
	if err != nil {
		return "", fmt.Errorf("passwords do not match")
	}

	token, err := jwt.GenerateJwt(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
