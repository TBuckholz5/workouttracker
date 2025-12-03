package service

type CreateExerciseRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	TargetMuscle string `json:"targetMuscle"`
}

type GetExerciseForUserParams struct {
	UserID int64 `json:"userID"`
	Offset int   `json:"offset"`
	Limit  int   `json:"limit"`
}

type CreateExerciseForUserParams struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	TargetMuscle string `json:"targetMuscle"`
	PictureURL   string `json:"pictureURL"`
	UserID       int64  `json:"userID"`
}
