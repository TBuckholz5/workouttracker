package user

import (
	"context"
	"fmt"

	serviceModels "github.com/TBuckholz5/workouttracker/internal/service/user/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(ctx context.Context, params *CreateUserParams) (serviceModels.User, error)
	GetUserForUsername(ctx context.Context, username string) (serviceModels.User, error)
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

type CreateUserParams struct {
	Username string
	Email    string
	PwHash   string
}

func (r *Repository) CreateUser(ctx context.Context, params *CreateUserParams) (serviceModels.User, error) {
	row := r.pool.QueryRow(ctx, createUser, params.Username, params.Email, params.PwHash)
	var user user
	err := row.Scan(
		&user.id,
		&user.username,
		&user.email,
		&user.pwHash,
		&user.createdAt,
		&user.updatedAt,
	)
	if err != nil {
		return serviceModels.User{}, fmt.Errorf("could not create user: %w", err)
	}
	return serviceModels.User{
		ID:       user.id,
		Username: user.username.String,
		Email:    user.email.String,
		PwHash:   user.pwHash.String,
	}, nil
}

func (r *Repository) GetUserForUsername(ctx context.Context, username string) (serviceModels.User, error) {
	row := r.pool.QueryRow(ctx, getUserByUsername, username)
	var user user
	err := row.Scan(
		&user.id,
		&user.username,
		&user.email,
		&user.pwHash,
		&user.createdAt,
		&user.updatedAt,
	)
	if err != nil {
		return serviceModels.User{}, fmt.Errorf("could not get user for username: %s", username)
	}
	return serviceModels.User{
		ID:       user.id,
		Username: user.username.String,
		Email:    user.email.String,
		PwHash:   user.pwHash.String,
	}, nil
}
