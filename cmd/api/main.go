package main

import (
	"log"
	"net/http"

	"github.com/sushanthach12/sellit-backend/internal/config"
	"github.com/sushanthach12/sellit-backend/internal/routes"
)

func main() {
	// Load environment variables and configuration
	cfg := config.MustLoad()

	mux := http.NewServeMux() // Create a new ServeMux for routing

	routes.RegisterRoutes(mux)

	// server := &http.Server{
	// 	Addr:         ":8000",
	// 	Handler:      mux,
	// 	ReadTimeout:  time.Second * 10,
	// 	WriteTimeout: time.Second * 30,
	// 	IdleTimeout:  time.Second * 60,
	// }

	server := config.GetServerConfig(mux, cfg)

	log.Printf("Server listening at port %s and running in %s mode", server.Addr, cfg.Env)
	if serverErr := http.ListenAndServe(server.Addr, server.Handler); serverErr != nil {
		log.Fatalf("Server failed to start: %v", serverErr)
	}

}
