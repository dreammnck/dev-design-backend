package routes

import (
	eHandler "backend/internal/events/handler"
	eRepo "backend/internal/events/repository"
	eSvc "backend/internal/events/service"
	pAdapter "backend/internal/payment/adapter"
	pHandler "backend/internal/payment/handler"
	pRepo "backend/internal/payment/repository"
	pSvc "backend/internal/payment/service"
	"backend/pkg/ping"
	"os"

	sHandler "backend/internal/seats/handler"
	sRepo "backend/internal/seats/repository"
	sSvc "backend/internal/seats/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes registers all the routes for the application
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Ping
	r.GET("/ping", ping.Handler)

	// Dependency Injection for Events
	eventRepo := eRepo.NewEventRepository(db)
	eventSvc := eSvc.NewEventService(eventRepo)
	eventH := eHandler.NewEventHandler(eventSvc)

	// Dependency Injection for Seats
	seatRepo := sRepo.NewSeatRepository(db)
	seatSvc := sSvc.NewSeatService(seatRepo)
	getSeatsH := sHandler.NewGetSeatsHandler(seatSvc)

	// Dependency Injection for Payment
	paymentRepo := pRepo.NewPaymentRepository(db)
	gatewayURL := os.Getenv("PAYMENT_GATEWAY_URL")
	if gatewayURL == "" {
		gatewayURL = "http://localhost:8082" // default for local development
	}
	paymentGatewayAda := pAdapter.NewPaymentGatewayAdapter(gatewayURL)
	paymentSvc := pSvc.NewPaymentService(paymentRepo, seatSvc, paymentGatewayAda)
	confirmPaymentH := pHandler.NewConfirmPaymentHandler(paymentSvc)

	api := r.Group("/api")
	{
		// Events
		eventsGroup := api.Group("/events")
		{
			eventsGroup.GET("/banner", eventH.Banner)
			eventsGroup.GET("/list", eventH.ListAll)
			eventsGroup.GET("/list/all", eventH.ListAll)
			eventsGroup.GET("/list/recommend", eventH.ListRecommend)
			eventsGroup.GET("/list/coming-soon", eventH.ListComingSoon)
		}
		api.GET("/events-detail/:id", eventH.Detail)

		// Seats
		submitSeatH := sHandler.NewSubmitSeatHandler(seatSvc)
		seatsGroup := api.Group("/events-seat")
		{
			seatsGroup.GET("/list", getSeatsH.Handle) // Use query param ?event_id=...
			seatsGroup.POST("/submit/:id", submitSeatH.Handle)
		}
		api.GET("/seats/:id", getSeatsH.Handle)

		// Payment
		processPaymentH := pHandler.NewProcessPaymentHandler(paymentSvc)
		cancelPaymentH := pHandler.NewCancelPaymentHandler(paymentSvc)
		paymentGroup := api.Group("/payment")
		{
			paymentGroup.POST("/process", processPaymentH.Handle)
			paymentGroup.POST("/cancel", cancelPaymentH.Handle)
		}
		api.POST("/confirm-payment", confirmPaymentH.Handle)
	}
}
