package repository

import "time"

type exercise struct {
	id           int64
	name         string
	description  string
	targetMuscle string
	pictureUrl   string
	createdAt    time.Time
	updatedAt    time.Time
	userId       int64
}
