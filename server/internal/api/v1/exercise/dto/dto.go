package dto

type CreateExerciseRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	TargetMuscle string `json:"targetMuscle"`
}

type GetExerciseForUserRequest struct {
	UserID int64
	Offset int
	Limit  int
}
