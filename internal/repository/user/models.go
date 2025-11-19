package user

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        int64
	Email     pgtype.Text
	Username  pgtype.Text
	PwHash    pgtype.Text
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}
