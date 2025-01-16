package main

import (
	"Goal-Storage/config"
	"Goal-Storage/controllers"
	"Goal-Storage/factories"
	"Goal-Storage/initializers"
	"Goal-Storage/middleware"
	"Goal-Storage/repositories"
	"Goal-Storage/service"
	"context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
		_ = MongoClient.Disconnect(ctx)
	}(initializers.MongoClient, nil)

	// Initialize repository and factory
	repo := repositories.NewMongoGoalRepository("ptrainer_goals")
	factory := factories.NewConcreteGoalFactory(repo)

	// Create GraphQL schema
	schema := service.CreateGraphQLSchema(factory)

	// Set up router
	r := mux.NewRouter()

	// Register routes, referencing controller functions
	r.HandleFunc("/", controllers.HandleGraphQL(schema)).Methods(http.MethodGet)
	r.HandleFunc("/register/goal", controllers.RegisterGoal(factory, schema)).Methods(http.MethodPost)
	r.HandleFunc("/modify/goal", controllers.ModifyGoal(factory, schema)).Methods(http.MethodPost)
	r.HandleFunc("/get/goal", controllers.GetGoal(factory)).Methods(http.MethodGet)

	// Apply middleware
	r.Use(middleware.AuthMiddleware)

	// Start the HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("Server is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, r)))
}
