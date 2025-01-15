package factories

import (
	"Goal-Storage/dtos"
	"Goal-Storage/models"
)

type GoalFactory interface {
	CreateGoal(input dtos.CreateGoalInput) (*models.GoalsModel, error)
	UpdateGoal(goalID int64, input dtos.CreateGoalInput) (*models.GoalsModel, error)
	GetGoalByID(goalID int64) (*models.GoalsModel, error)
	DeleteGoal(goalID int64) error
}