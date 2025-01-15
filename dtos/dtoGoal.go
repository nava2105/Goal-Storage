package dtos

type CreateGoalInput struct {
	User_id        int64   `json:"user_id" validate:"required"`
	Weight         float64 `json:"weight" validate:"required"`
	Body_structure string  `json:"body_structure" validate:"required,email"`
}
