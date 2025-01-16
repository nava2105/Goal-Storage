package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"

	"Goal-Storage/factories"
	"Goal-Storage/service"
)

// RegisterGoal handles goal registration.
func RegisterGoal(factory factories.GoalFactory, schema graphql.Schema) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := service.GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var requestBody struct {
			Weight        float64 `json:"weight"`
			BodyStructure string  `json:"body_structure"`
		}

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}

		query := `
			mutation {
				createGoal(userId: %d, weight: %f, body_structure: "%s") {
					goalId
					userId
					weight
					body_structure
				}
			}
		`

		mutation := fmt.Sprintf(query, userId, requestBody.Weight, requestBody.BodyStructure)
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: mutation,
			Context:       context.WithValue(r.Context(), "factory", factory),
		})

		if len(result.Errors) > 0 {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": result.Errors,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(result.Data)
	}
}

// ModifyGoal handles goal modification.
func ModifyGoal(factory factories.GoalFactory, schema graphql.Schema) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the userId from the Authorization token
		userId, err := service.GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Parse the incoming request body for modifications
		var requestBody struct {
			Weight        float64 `json:"weight"`
			BodyStructure string  `json:"body_structure"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}

		// Fetch the active goal for the user
		existingGoal, err := factory.GetGoalByID(int64(userId))
		if err != nil || existingGoal == nil {
			http.Error(w, "No active goal found for the user", http.StatusNotFound)
			return
		}

		// Dynamically create the GraphQL mutation query for goal modification
		query := `
			mutation {
				updateGoal(goalId: "%s", userId: %d, weight: %f, body_structure: "%s") {
					goalId
					userId
					weight
					body_structure
				}
			}
		`
		mutation := fmt.Sprintf(query, existingGoal.GoalId, userId, requestBody.Weight, requestBody.BodyStructure)
		// Execute the mutation via the GraphQL schema
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: mutation,
			Context:       context.WithValue(r.Context(), "factory", factory),
		})

		// Handle GraphQL execution errors
		if len(result.Errors) > 0 {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": result.Errors,
			})
			return
		}

		// Return the updated goal to the client
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(result.Data)
	}
}

// GetGoal handles retrieving a user's goal.
func GetGoal(factory factories.GoalFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the user ID from the Authorization token
		userId, err := service.GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Fetch the goal for the authenticated user
		goal, err := factory.GetGoalByID(int64(userId))
		if err != nil {
			http.Error(w, "Failed to retrieve goal", http.StatusInternalServerError)
			return
		}

		if goal == nil {
			http.Error(w, "No active goal found for the user", http.StatusNotFound)
			return
		}

		// Return the goal as a JSON response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(goal); err != nil {
			http.Error(w, "Failed to encode goal response", http.StatusInternalServerError)
		}
	}
}
