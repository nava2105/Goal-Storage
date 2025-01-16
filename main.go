package main

import (
	"Goal-Storage/config"
	"Goal-Storage/factories"
	"Goal-Storage/initializers"
	"Goal-Storage/middleware"
	"Goal-Storage/repositories"
	"Goal-Storage/service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os"
)

func init() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	// Initialize MongoDB
	initializers.InitializeMongoDB(cfg.MongoURI)
	defer func(MongoClient *mongo.Client, ctx context.Context) {
		err := MongoClient.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(initializers.MongoClient, nil)
	// Initialize repository
	repo := repositories.NewMongoGoalRepository("ptrainer_goals")
	// Initialize factory
	factory := factories.NewConcreteGoalFactory(repo)
	// Create a new GraphQL schema
	schema := service.CreateGraphQLSchema(factory)
	// Set up the HTTP server
	r := mux.NewRouter()
	// GraphQL handler
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Parse the GraphQL query
		var requestBody struct {
			Query string `json:"query"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}
		result := graphql.Do(graphql.Params{
			// Execute the query
			Schema:        schema,
			RequestString: requestBody.Query,
			Context:       r.Context(), // Pass context for token extraction
		})
		if len(result.Errors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": result.Errors,
			})
			if err != nil {
				return
			}
			return
		}
		userId, ok := result.Data.(map[string]interface{})["userId"]
		if !ok {
			http.Error(w, "Failed to process data", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"userId": userId,
		})
		if err != nil {
			return
		}
	}).Methods(http.MethodGet)
	r.HandleFunc("/register/goal", func(w http.ResponseWriter, r *http.Request) {
		userId, err := service.GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		// Parse the GraphQL query
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
			// Execute the query
			Schema:        schema,
			RequestString: mutation,
			Context:       context.WithValue(r.Context(), "factory", factory), // Pass context for token extraction
		})

		if len(result.Errors) > 0 {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": result.Errors,
			})
			if err != nil {
				return
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		er := json.NewEncoder(w).Encode(result.Data)
		if er != nil {
			return
		}
	}).Methods(http.MethodPost)
	r.HandleFunc("/modify/goal", func(w http.ResponseWriter, r *http.Request) {
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
		factory := context.WithValue(r.Context(), "factory", factory).Value("factory").(factories.GoalFactory)
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
		mutation := fmt.Sprintf(query, existingGoal.Goal_id, userId, requestBody.Weight, requestBody.BodyStructure)
		// Execute the mutation via the GraphQL schema
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: mutation,
			Context:       context.WithValue(r.Context(), "factory", factory),
		})

		// Handle GraphQL execution errors
		if len(result.Errors) > 0 {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": result.Errors,
			})
			if err != nil {
				return
			}
			return
		}

		// Return the updated goal to the client
		w.Header().Set("Content-Type", "application/json")
		er := json.NewEncoder(w).Encode(result.Data)
		if er != nil {
			return
		}
	}).Methods(http.MethodPost)
	r.Use(middleware.AuthMiddleware)
	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("Server is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, r)))
}
