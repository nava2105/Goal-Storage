package main

import (
	"Goal-Storage/config"
	"Goal-Storage/dtos"
	"Goal-Storage/factories"
	"Goal-Storage/initializers"
	"Goal-Storage/middleware"
	"Goal-Storage/repositories"
	"Goal-Storage/service"
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/joho/godotenv"
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
	defer initializers.MongoClient.Disconnect(nil)
	// Initialize repository
	repo := repositories.NewMongoGoalRepository("ptrainer_goals")

	// Initialize factory
	factory := factories.NewConcreteGoalFactory(repo)

	// Example: Create a New Goal
	input := dtos.CreateGoalInput{
		User_id:        1,
		Weight:         75.5,
		Body_structure: "Mesomorph",
	}

	goal, err := factory.CreateGoal(input)
	if err != nil {
		log.Fatalf("Failed to create goal: %v", err)
	}
	log.Printf("Created Goal: %+v\n", goal)
	// Create a new GraphQL schema
	schema := service.CreateGraphQLSchema()
	// Set up the HTTP server
	r := mux.NewRouter()
	// GraphQL handler
	r.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		// Parse the GraphQL query
		var requestBody struct {
			Query string `json:"query"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}

		// Execute the query
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: requestBody.Query,
			Context:       r.Context(), // Pass context for token extraction
		})

		// Check for errors
		if len(result.Errors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": result.Errors,
			})
			return
		}

		// Custom response format: extract and return only the required field
		// Extract "userId" from the data object
		userId, ok := result.Data.(map[string]interface{})["userId"]
		if !ok {
			http.Error(w, "Failed to process data", http.StatusInternalServerError)
			return
		}

		// Send the custom response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"userId": userId,
		})
	}).Methods(http.MethodGet)
	r.Use(middleware.AuthMiddleware)
	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("Server is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, r)))
}
