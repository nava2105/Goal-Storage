package service

import (
	"Goal-Storage/dtos"
	"Goal-Storage/factories"
	"Goal-Storage/utils"
	"errors"
	"github.com/graphql-go/graphql"
)

// CreateGraphQLSchema defines and returns the GraphQL schema.
func CreateGraphQLSchema(factory factories.GoalFactory) graphql.Schema {
	// Define the Goal GraphQL Object.
	goalType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Goal",
		Fields: graphql.Fields{
			"goalId":         &graphql.Field{Type: graphql.Int},
			"userId":         &graphql.Field{Type: graphql.Int},
			"weight":         &graphql.Field{Type: graphql.Float},
			"body_structure": &graphql.Field{Type: graphql.String},
		},
	})
	// Existing schema: Root Query for userId.
	rootQueryFields := graphql.Fields{
		"userId": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Extract token from context and resolve userId
				token, ok := p.Context.Value("Authorization").(string)
				if !ok || token == "" {
					return nil, errors.New("no authorization token found")
				}
				userId, err := utils.FetchUserIdFromAuthAPI(token)
				if err != nil {
					return nil, err
				}
				return userId, nil
			},
		},
	}
	// Query: Get Goal by ID.
	rootQueryFields["getGoalById"] = &graphql.Field{
		Type: goalType, // Defined below
		Args: graphql.FieldConfigArgument{
			"userId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			userId := p.Args["userId"].(int)
			return factory.GetGoalByID(int64(userId))
		},
	}
	// Mutation: Create Goal.
	rootMutationFields := graphql.Fields{
		"createGoal": &graphql.Field{
			Type: goalType, // Defined below
			Args: graphql.FieldConfigArgument{
				"userId":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"weight":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
				"body_structure": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				factory := p.Context.Value("factory").(factories.GoalFactory)
				userID := int64(p.Args["userId"].(int))
				existingGoal, err := factory.GetGoalByID(userID)
				if err == nil && existingGoal != nil {
					return nil, errors.New("user already has an active goal")
				}
				weight := p.Args["weight"].(float64)
				bodyStructure := p.Args["body_structure"].(string)
				createGoalInput := dtos.CreateGoalInput{
					User_id:        userID,
					Weight:         weight,
					Body_structure: bodyStructure,
				}
				return factory.CreateGoal(createGoalInput)
			},
		},
		// Mutation: Update Goal.
		"updateGoal": &graphql.Field{
			Type: goalType, // Defined below
			Args: graphql.FieldConfigArgument{
				"goalId":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"userId":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"weight":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
				"body_structure": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				updateGoalInput := dtos.CreateGoalInput{
					User_id:        int64(p.Args["userId"].(int)),
					Weight:         p.Args["weight"].(float64),
					Body_structure: p.Args["body_structure"].(string),
				}
				goalId := int64(p.Args["goalId"].(int))
				return factory.UpdateGoal(goalId, updateGoalInput)
			},
		},
	}
	// Combine queries and mutations into schema configuration.
	query := graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: rootQueryFields})
	mutation := graphql.NewObject(graphql.ObjectConfig{Name: "Mutation", Fields: rootMutationFields})
	// Create Schema.
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    query,
		Mutation: mutation,
	})
	if err != nil {
		panic("Failed to create schema: " + err.Error())
	}
	return schema
}
