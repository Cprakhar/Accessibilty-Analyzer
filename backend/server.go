package main

import (
	"backend/api"
	"backend/jobs"
	"backend/middleware"
	"backend/services"
	"backend/utils"
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Initialize the logger
	utils.InitLogger()
	defer utils.CloseLogger()

	// Load environment variables or config here if needed
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017/accessibility_analyser"
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	db := client.Database("accessibility_analyser")
	services.InitUserService(db)
	services.InitReportService(db)

	r := gin.Default()

	// Add CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Register authentication routes
	api.RegisterAuthRoutes(r)
	api.RegisterAnalyzeRoutes(r)

	// TODO: Register other API routes here

	// Start background worker for analysis jobs
	jobs.StartAnalyzeWorker()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}
