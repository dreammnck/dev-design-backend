package main

import (
	"backend/internal/routes"
	sRepo "backend/internal/seats/repository"
	"backend/pkg/database"
	"backend/pkg/middleware"
	"backend/pkg/storage"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db := database.InitDB()

	// Initialize Repository and Service
	seatRepo := sRepo.NewSeatRepository(db)

	// Initialize Gin engine
	r := gin.Default()
	r.Use(middleware.CORS())

	// Setup centralized routes
	routes.SetupRoutes(r, db)

	// Initialize Google Cloud Storage (requires env vars loaded by InitDB)
	storage.InitGCS()

	// Configuration for cleanup
	cleanupIntervalStr := os.Getenv("CLEANUP_INTERVAL")
	if cleanupIntervalStr == "" {
		cleanupIntervalStr = "5m"
	}
	cleanupInterval, err := time.ParseDuration(cleanupIntervalStr)
	if err != nil {
		log.Printf("Invalid CLEANUP_INTERVAL %s, using default 5m\n", cleanupIntervalStr)
		cleanupInterval = 5 * time.Minute
	}

	reservationTimeoutStr := os.Getenv("RESERVATION_TIMEOUT")
	if reservationTimeoutStr == "" {
		reservationTimeoutStr = "5m"
	}
	reservationTimeout, err := time.ParseDuration(reservationTimeoutStr)
	if err != nil {
		log.Printf("Invalid RESERVATION_TIMEOUT %s, using default 5m\n", reservationTimeoutStr)
		reservationTimeout = 5 * time.Minute
	}

	// Start background worker for reservation cleanup
	go func() {
		for {
			time.Sleep(cleanupInterval)
			log.Println("Running reservation cleanup...")
			if err := seatRepo.ClearExpiredReservations(reservationTimeout); err != nil {
				log.Printf("Failed to clear expired reservations: %v\n", err)
			}
		}
	}()

	// Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on :%s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Unable to start server: ", err)
	}
}
