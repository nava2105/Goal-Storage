package repositories

import (
	"Goal-Storage/initializers"
	"Goal-Storage/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type GoalRepository interface {
	Create(goal *models.GoalsModel) (*models.GoalsModel, error)
	Update(goalID int64, goal *models.GoalsModel) (*models.GoalsModel, error)
	GetByID(goalID int64) (*models.GoalsModel, error)
	Delete(goalID int64) error
}

type MongoGoalRepository struct {
	Collection string
}

func NewMongoGoalRepository(collection string) *MongoGoalRepository {
	return &MongoGoalRepository{Collection: collection}
}

func (r *MongoGoalRepository) Create(goal *models.GoalsModel) (*models.GoalsModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := initializers.GetMongoCollection("ptrainer_goals", r.Collection)
	result, err := collection.InsertOne(ctx, goal)
	if err != nil {
		return nil, err
	}
	goal.Goal_id, _ = result.InsertedID.(int64)
	return goal, nil
}

func (r *MongoGoalRepository) Update(goalID int64, goal *models.GoalsModel) (*models.GoalsModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := initializers.GetMongoCollection("ptrainer_goals", r.Collection)
	filter := bson.M{"_id": goalID}
	update := bson.M{"$set": goal}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return goal, nil
}

func (r *MongoGoalRepository) GetByID(goalID int64) (*models.GoalsModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := initializers.GetMongoCollection("ptrainer_goals", r.Collection)
	filter := bson.M{"_id": goalID}

	var goal models.GoalsModel
	err := collection.FindOne(ctx, filter).Decode(&goal)
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

func (r *MongoGoalRepository) Delete(goalID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := initializers.GetMongoCollection("ptrainer_goals", r.Collection)
	filter := bson.M{"_id": goalID}

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
