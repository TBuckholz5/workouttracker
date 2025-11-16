package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/TBuckholz5/workouttracker/internal/db/user"
)

type Repository struct {
	queries *user.Queries
}

func NewRepository(queries *user.Queries) *Repository {
	return &Repository{
		queries: queries,
	}
}

func (r *Repository) CreateUser(ctx context.Context, arg *user.CreateUserParams) (user.User, error) {
	row, err := r.queries.CreateUser(ctx, *arg)
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

func (r *Repository) GetUserForUsername(ctx context.Context, username sql.NullString) (user.User, error) {
	row, err := r.queries.GetUserByUsername(ctx, username)
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
