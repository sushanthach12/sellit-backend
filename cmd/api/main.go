package main

import (
	"log"
	"net/http"

	"github.com/sushanthach12/sellit-backend/config"
	"github.com/sushanthach12/sellit-backend/routes"
)

func main() {
	mux := http.NewServeMux() // Create a new ServeMux for routing

	routes.RegisterRoutes(mux)

	// server := &http.Server{
	// 	Addr:         ":8000",
	// 	Handler:      mux,
	// 	ReadTimeout:  time.Second * 10,
	// 	WriteTimeout: time.Second * 30,
	// 	IdleTimeout:  time.Second * 60,
	// }

	server := config.GetServerConfig(mux)

	if serverErr := http.ListenAndServe(server.Addr, server.Handler); serverErr != nil {
		log.Fatalf("Server failed to start: %v", serverErr)
	}

}
