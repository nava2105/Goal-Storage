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

func (r *MongoGoalRepository) Update(userID int64, goal *models.GoalsModel) (*models.GoalsModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := initializers.GetMongoCollection("ptrainer_goals", r.Collection)
	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": goal}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return goal, nil
}

func (r *MongoGoalRepository) GetByID(userID int64) (*models.GoalsModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := initializers.GetMongoCollection("ptrainer_goals", r.Collection)
	filter := bson.M{"user_id": userID}

	var goal models.GoalsModel
	err := collection.FindOne(ctx, filter).Decode(&goal)
	if err != nil {
		return nil, err
	}
	return &goal, nil
}
