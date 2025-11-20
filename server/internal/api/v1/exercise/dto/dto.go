package dto

type CreateExerciseRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	TargetMuscle string `json:"targetMuscle"`
}

type GetExerciseForUserRequest struct {
	UserID int64 `json:"userId" binding:"required"`
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}
