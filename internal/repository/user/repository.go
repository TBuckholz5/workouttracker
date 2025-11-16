package user

import (
	"context"
	"fmt"

	"github.com/TBuckholz5/workouttracker/internal/db/user"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) CreateUser(ctx context.Context, arg *user.CreateUserParams) (user.User, error) {
	row, err := user.New(r.pool).CreateUser(ctx, *arg)
	if err != nil {
		return user.User{}, err
	}
	return user.User{
		ID:        row.ID,
		Username:  row.Username,
		Email:     row.Email,
		PwHash:    row.PwHash,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func (r *Repository) GetUserForUsername(ctx context.Context, username pgtype.Text) (user.User, error) {
	row, err := user.New(r.pool).GetUserByUsername(ctx, username)
	if err != nil {
		return user.User{}, fmt.Errorf("could not get user for username: %s", username.String)
	}
	return user.User{
		ID:        row.ID,
		Username:  row.Username,
		Email:     row.Email,
		PwHash:    row.PwHash,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}
