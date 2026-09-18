package main

import (
	"log"
	"net/http"

	"github.com/sushanthach12/sellit-backend/internal/config"
	"github.com/sushanthach12/sellit-backend/internal/database"
	"github.com/sushanthach12/sellit-backend/internal/routes"
)

func main() {
	// Load environment variables and configuration
	cfg := config.MustLoad()

	// Initialize the database connection
	db, err := database.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	mux := http.NewServeMux() // Create a new ServeMux for routing

	routes.RegisterRoutes(mux, db)

	server := config.GetServerConfig(mux, cfg)

	log.Printf("Server listening at port %s and running in %s mode", server.Addr, cfg.Env)
	if serverErr := http.ListenAndServe(server.Addr, server.Handler); serverErr != nil {
		log.Fatalf("Server failed to start: %v", serverErr)
	}

}
