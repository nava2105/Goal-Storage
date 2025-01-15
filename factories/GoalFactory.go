package factories

import (
	"Goal-Storage/dtos"
	"Goal-Storage/models"
)

type GoalFactory interface {
	CreateGoal(input dtos.CreateGoalInput) (*models.GoalsModel, error)
	UpdateGoal(userID int64, input dtos.CreateGoalInput) (*models.GoalsModel, error)
	GetGoalByID(userID int64) (*models.GoalsModel, error)
}
