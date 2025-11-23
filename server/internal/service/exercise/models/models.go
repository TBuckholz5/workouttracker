package models

type Exercise struct {
	ID           int64
	Name         string
	Description  string
	TargetMuscle string
	PictureURL   string
}

type GetExerciseForUserParams struct {
	UserID int64
	Offset int
	Limit  int
}

type CreateExerciseForUserParams struct {
	Name         string
	Description  string
	TargetMuscle string
	PictureURL   string
	UserID       int64
}
