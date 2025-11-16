package user

import (
	"context"
	"database/sql"

	"github.com/TBuckholz5/workouttracker/internal/api/v1/user/dto"
	db "github.com/TBuckholz5/workouttracker/internal/db/user"
	repo "github.com/TBuckholz5/workouttracker/internal/repository/user"
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
		Username: sql.NullString{String: userDto.Username, Valid: true},
		Email:    sql.NullString{String: userDto.Email, Valid: true},
		PwHash:   sql.NullString{String: hashedPassword, Valid: true},
	})
	return err
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
