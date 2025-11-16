package user

import (
	"context"

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
