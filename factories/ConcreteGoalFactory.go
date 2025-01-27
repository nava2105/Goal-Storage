package factories

import (
	"Goal-Storage/dtos"
	"Goal-Storage/models"
	"Goal-Storage/repositories"
	"errors"
)

type ConcreteGoalFactory struct {
	Repository repositories.GoalRepository
}

func NewConcreteGoalFactory(repo repositories.GoalRepository) *ConcreteGoalFactory {
	return &ConcreteGoalFactory{Repository: repo}
}

func (f *ConcreteGoalFactory) CreateGoal(input dtos.CreateGoalInput) (*models.GoalsModel, error) {
	// Validation checks
	if input.UserId == 0 || input.Weight <= 0 {
		return nil, errors.New("invalid goal input data")
	}

	goal := &models.GoalsModel{
		UserId:        input.UserId,
		Weight:        input.Weight,
		BodyStructure: input.BodyStructure,
	}

	// Delegates creation to repository
	return f.Repository.Create(goal)
}

func (f *ConcreteGoalFactory) UpdateGoal(userID int64, input dtos.CreateGoalInput) (*models.GoalsModel, error) {
	existingGoal, err := f.Repository.GetByID(userID)
	if err != nil {
		return nil, err
	}

	existingGoal.Weight = input.Weight
	existingGoal.BodyStructure = input.BodyStructure

	// Delegates update to repository
	return f.Repository.Update(userID, existingGoal)
}

func (f *ConcreteGoalFactory) GetGoalByID(userID int64) (*models.GoalsModel, error) {
	// Delegates retrieval to repository
	return f.Repository.GetByID(userID)
}
