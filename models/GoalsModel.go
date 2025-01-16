package models

type GoalsModel struct {
	GoalId        string  `bson:"_id,omitempty"`
	UserId        int64   `bson:"user_id"`
	Weight        float64 `bson:"weight"`
	BodyStructure string  `bson:"body_structure"`
}
