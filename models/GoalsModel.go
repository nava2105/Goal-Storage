package models

type GoalsModel struct {
	Goal_id        int64   `bson:"_id,omitempty"`
	User_id        int64   `bson:"user_id"`
	Weight         float64 `bson:"weight"`
	Body_structure string  `bson:"body_structure"`
}
