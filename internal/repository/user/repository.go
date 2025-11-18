package user

import (
	"context"
	"fmt"
	"time"

	db "github.com/TBuckholz5/workouttracker/internal/db/user"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg *CreateUserParams) (User, error)
	GetUserForUsername(ctx context.Context, username string) (User, error)
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

type User struct {
	ID        int64
	Username  string
	Email     string
	PwHash    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *Repository) CreateUser(ctx context.Context, params *CreateUserParams) (User, error) {
	arg := &db.CreateUserParams{
		Username: pgtype.Text{
			String: params.Username,
			Valid:  true,
		},
		Email: pgtype.Text{
			String: params.Email,
			Valid:  true,
		},
		PwHash: pgtype.Text{
			String: params.PwHash,
			Valid:  true,
		},
	}
	row, err := db.New(r.pool).CreateUser(ctx, *arg)
	if err != nil {
		return User{}, err
	}
	return User{
		ID:        row.ID,
		Username:  row.Username.String,
		Email:     row.Email.String,
		PwHash:    row.PwHash.String,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

func (r *Repository) GetUserForUsername(ctx context.Context, username string) (User, error) {
	row, err := db.New(r.pool).GetUserByUsername(ctx, pgtype.Text{
		String: username,
		Valid:  true,
	})
	if err != nil {
		return User{}, fmt.Errorf("could not get user for username: %s", username)
	}
	return User{
		ID:        row.ID,
		Username:  row.Username.String,
		Email:     row.Email.String,
		PwHash:    row.PwHash.String,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}
