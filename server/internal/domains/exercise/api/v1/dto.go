package v1

type Exercise struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	TargetMuscle string `json:"targetMuscle"`
	PictureURL   string `json:"pictureURL"`
}

type CreateExerciseRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	TargetMuscle string `json:"targetMuscle" binding:"required"`
}

type CreateExerciseResponse struct {
	Exercise Exercise `json:"exercise"`
}

type GetExerciseListResponse struct {
	Exercises []Exercise `json:"exercises"`
}
