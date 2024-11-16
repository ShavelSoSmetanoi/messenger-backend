package app

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/config"
	services2 "github.com/ShavelSoSmetanoi/messenger-backend/internal/services"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/transport/rest"
	"log"
)

func Run() {
	// Setup configuration (e.g., environment variables, settings)
	config.SetupConfig()

	// Initialize services by calling the InitServices function
	services := services2.InitServices()

	// Create a new handler for the REST API and initialize the route handlers
	handler := rest.NewHandler(services)
	router := handler.Init()

	// Run the server on port 8080. If there is an error, log it and exit.
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server startup error: %v", err)
	}
}
