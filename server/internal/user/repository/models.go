package repository

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type user struct {
	id        int64
	email     pgtype.Text
	username  pgtype.Text
	pwHash    pgtype.Text
	createdAt pgtype.Timestamp
	updatedAt pgtype.Timestamp
}
