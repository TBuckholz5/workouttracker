package models

import "time"

type WorkoutSet struct {
	ID       int64   `json:"id,omitempty"`
	Reps     int     `json:"reps"`
	Weight   float64 `json:"weight"`
	SetType  string  `json:"set_type"`
	SetOrder int     `json:"set_order"`
}

type Workout struct {
	ID          int64        `json:"id,omitempty"`
	ExerciseID  int64        `json:"exerciseID"`
	Description string       `json:"description,omitempty"`
	Sets        []WorkoutSet `json:"sets"`
}

type WorkoutSession struct {
	ID          int64     `json:"id,omitempty"`
	UserID      int64     `json:"userID,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Duration    int       `json:"duration,omitempty"`
	Workouts    []Workout `json:"workouts"`
	CreatedAt   time.Time `json:"createdAt"`
}
