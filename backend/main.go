package main

import (
	"context"
	"ebpf-dashboard/config"
	"ebpf-dashboard/database"
	"ebpf-dashboard/handlers"
	"ebpf-dashboard/logger"
	"ebpf-dashboard/repository"
	"ebpf-dashboard/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Initialize logger
	if err := logger.Init(cfg.LogPath); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	logger.Info("Starting eBPF Dashboard Backend...")

	// Initialize database
	db, err := database.InitDB(cfg.DBPath)
	if err != nil {
		logger.Error("Failed to initialize database: %v", err)
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	processRepo := repository.NewProcessRepository(db)
	networkRepo := repository.NewNetworkRepository(db)
	diskRepo := repository.NewDiskRepository(db)
	cpuProfileRepo := repository.NewCPUProfileRepository(db)
	tcpLifeRepo := repository.NewTCPLifeRepository(db)
	syscallRepo := repository.NewSyscallRepository(db)

	// Initialize services
	processService := services.NewProcessService(processRepo)
	networkService := services.NewNetworkService(networkRepo)
	diskService := services.NewDiskService(diskRepo)
	cpuProfileService := services.NewCPUProfileService(cpuProfileRepo)
	tcpLifeService := services.NewTCPLifeService(tcpLifeRepo)
	syscallService := services.NewSyscallService(syscallRepo)

	// Start background collectors
	processService.StartCollecting()
	networkService.StartCollecting()
	diskService.StartCollecting()
	cpuProfileService.Start()
	if err := tcpLifeService.StartCollecting(); err != nil {
		logger.Error("Failed to start tcplife collector: %v", err)
	}
	syscallService.Start()

	// Initialize handlers
	processHandler := handlers.NewProcessHandler(processService)
	networkHandler := handlers.NewNetworkHandler(networkService)
	diskHandler := handlers.NewDiskHandler(diskService)
	cpuProfileHandler := handlers.NewCPUProfileHandler(cpuProfileService)
	tcpLifeHandler := handlers.NewTCPLifeHandler(tcpLifeService)
	syscallHandler := handlers.NewSyscallHandler(syscallService)
	healthHandler := handlers.NewHealthHandler()

	// Setup Gin router
	router := gin.Default()

	// CORS middleware
	if cfg.CORSEnabled {
		corsConfig := cors.DefaultConfig()
		origins := strings.Split(cfg.CORSOrigins, ",")
		corsConfig.AllowOrigins = origins
		corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
		corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
		router.Use(cors.New(corsConfig))
		logger.Info("CORS enabled for origins: %s", cfg.CORSOrigins)
	}

	// Register routes
	api := router.Group("/api/metrics")
	{
		api.GET("/processes", processHandler.GetRecentProcesses)
		api.GET("/network", networkHandler.GetRecentConnections)
		api.GET("/disk", diskHandler.GetLatestLatency)
		api.GET("/cpuprofile", cpuProfileHandler.GetCPUProfiles)
		api.GET("/tcplife", tcpLifeHandler.GetTCPLifeEvents)
		api.GET("/syscalls", syscallHandler.GetSyscallStats)
	}
	router.GET("/health", healthHandler.GetHealth)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		logger.Info("Shutting down gracefully...")

		// Stop collectors
		processService.StopCollecting()
		networkService.StopCollecting()
		diskService.StopCollecting()
		cpuProfileService.Stop()
		tcpLifeService.StopCollecting()
		syscallService.Stop()

		// Shutdown HTTP server with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("Server forced to shutdown: %v", err)
		}

		logger.Info("Server stopped")
	}()

	// Start server
	logger.Info("Server is running on http://localhost:%s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Failed to start server: %v", err)
		log.Fatalf("Failed to start server: %v", err)
	}
}
