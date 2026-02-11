package main

import (
	"ebpf-dashboard/database"
	"ebpf-dashboard/handlers"
	"ebpf-dashboard/repository"
	"ebpf-dashboard/services"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting eBPF Dashboard Backend...")

	// Initialize database
	db, err := database.InitDB("./metrics.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	processRepo := repository.NewProcessRepository(db)
	networkRepo := repository.NewNetworkRepository(db)
	diskRepo := repository.NewDiskRepository(db)

	// Initialize services
	processService := services.NewProcessService(processRepo)
	networkService := services.NewNetworkService(networkRepo)
	diskService := services.NewDiskService(diskRepo)

	// Start background collectors
	processService.StartCollecting()
	networkService.StartCollecting()
	diskService.StartCollecting()

	// Initialize handlers
	processHandler := handlers.NewProcessHandler(processService)
	networkHandler := handlers.NewNetworkHandler(networkService)
	diskHandler := handlers.NewDiskHandler(diskService)
	healthHandler := handlers.NewHealthHandler()

	// Setup Gin router
	router := gin.Default()

	// Register routes
	api := router.Group("/api/metrics")
	{
		api.GET("/processes", processHandler.GetRecentProcesses)
		api.GET("/network", networkHandler.GetRecentConnections)
		api.GET("/disk", diskHandler.GetLatestLatency)
	}
	router.GET("/health", healthHandler.GetHealth)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down gracefully...")
		processService.StopCollecting()
		networkService.StopCollecting()
		diskService.StopCollecting()
		os.Exit(0)
	}()

	// Start server
	log.Println("Server is running on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
