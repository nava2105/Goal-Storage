package dtos

type CreateGoalInput struct {
	UserId        int64   `json:"user_id" validate:"required"`
	Weight        float64 `json:"weight" validate:"required"`
	BodyStructure string  `json:"body_structure" validate:"required,email"`
}
