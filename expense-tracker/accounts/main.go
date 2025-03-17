package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jo/choreo-tutorial/accounts/config"
	"github.com/jo/choreo-tutorial/accounts/db"
	"github.com/jo/choreo-tutorial/accounts/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/swaggo/http-swagger"
	_ "github.com/swaggo/swag/example/basic/docs" // for swagger
)

// @title Accounts API
// @version 1.0
// @description Personal finance tracking API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

func main() {
	// Load .env file if it exists
	godotenv.Load()

	// Initialize configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Create tables if they don't exist
	if err := database.CreateTables(); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Initialize router
	r := mux.NewRouter()
	
	// API routes
	api := r.PathPrefix("/").Subrouter()
	
	// Bill handlers
	billHandler := handlers.NewBillHandler(database)
	api.HandleFunc("/bills", billHandler.GetBills).Methods("GET")
	api.HandleFunc("/bills", billHandler.CreateBill).Methods("POST")
	api.HandleFunc("/bills/{id}", billHandler.GetBill).Methods("GET")
	api.HandleFunc("/bills/{id}", billHandler.UpdateBill).Methods("PUT")
	api.HandleFunc("/bills/{id}", billHandler.DeleteBill).Methods("DELETE")
	
	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, r))
}
