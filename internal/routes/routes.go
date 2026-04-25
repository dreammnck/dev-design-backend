package routes

import (
	"backend/internal/auth"
	authHandler "backend/internal/auth/handler"
	authRepo "backend/internal/auth/repository"
	authSvc "backend/internal/auth/service"
	eHandler "backend/internal/events/handler"
	eRepo "backend/internal/events/repository"
	eSvc "backend/internal/events/service"
	pAdapter "backend/internal/payment/adapter"
	pHandler "backend/internal/payment/handler"
	pRepo "backend/internal/payment/repository"
	pSvc "backend/internal/payment/service"
	"backend/pkg/middleware"
	"backend/pkg/notification"
	"backend/pkg/ping"
	"fmt"
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

	// ── Auth ──────────────────────────────────────────────────────────────────
	userRepo := authRepo.NewUserRepository(db)
	authService := authSvc.NewAuthService(userRepo)
	loginH := authHandler.NewLoginHandler(authService)
	userH := authHandler.NewUserHandler(authService)

	r.POST("/auth/login", loginH.Handle)
	r.POST("/auth/register", userH.Register)

	// Gmail OAuth
	gmailAuthH, err := authHandler.NewGmailAuthHandler()
	if err == nil {
		r.GET("/api/auth/gmail", gmailAuthH.Login)
		r.GET("/api/auth/gmail/callback", gmailAuthH.Callback)
	} else {
		fmt.Printf("Warning: Gmail OAuth handlers not registered: %v\n", err)
	}

	// ── Dependency Injection ──────────────────────────────────────────────────

	// Events
	eventRepo := eRepo.NewEventRepository(db)
	eventSvc := eSvc.NewEventService(eventRepo)
	eventH := eHandler.NewEventHandler(eventSvc)

	// Seats
	seatRepo := sRepo.NewSeatRepository(db)
	seatSvc := sSvc.NewSeatService(seatRepo)
	getSeatsH := sHandler.NewGetSeatsHandler(seatSvc)
	seatAdminH := sHandler.NewSeatAdminHandler(eventSvc)

	// Payment
	paymentRepo := pRepo.NewPaymentRepository(db)
	gatewayURL := os.Getenv("PAYMENT_GATEWAY_URL")
	if gatewayURL == "" {
		gatewayURL = "http://localhost:8082"
	}
	paymentGatewayAda := pAdapter.NewPaymentGatewayAdapter(gatewayURL)
	notificationSvc := notification.NewNotificationService()
	paymentSvc := pSvc.NewPaymentService(paymentRepo, seatSvc, paymentGatewayAda, authService, notificationSvc)
	confirmPaymentH := pHandler.NewConfirmPaymentHandler(paymentSvc)
	payoutH := pHandler.NewPayoutHandler(eventSvc)
	bookingAdminH := pHandler.NewBookingAdminHandler(eventSvc)
	globalAdminH := pHandler.NewGlobalAdminHandler(paymentSvc)
	myBookingH := pHandler.NewMyBookingHandler(paymentSvc)

	api := r.Group("/api")

	// ── Public routes (no auth required) ─────────────────────────────────────
	{
		eventsGroup := api.Group("/events")
		eventsGroup.Use(middleware.OptionalAuth())
		{
			eventsGroup.GET("/banner", eventH.Banner)
			eventsGroup.GET("/list", eventH.ListAll)
			eventsGroup.GET("/list/all", eventH.ListAll)
			eventsGroup.GET("/list/recommend", eventH.ListRecommend)
			eventsGroup.GET("/list/coming-soon", eventH.ListComingSoon)
		}
		api.GET("/events-detail/:id", middleware.OptionalAuth(), eventH.Detail)

		// Seat listing is public; seat submission requires auth (see below)
		seatsPublic := api.Group("/events-seat")
		{
			seatsPublic.GET("/list", getSeatsH.Handle) // ?event_id=...
		}
		api.GET("/getLocation", eventH.GetLocations)
		api.GET("/seats/:id", getSeatsH.Handle)
	}

	// ── Protected routes (JWT required) ──────────────────────────────────────
	protected := api.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		// Auth profile
		protected.GET("/auth/me", userH.GetMe)
		protected.GET("/my-tickets", myBookingH.GetMyBookings)
		protected.GET("/tickets/:id", myBookingH.GetTicketByID)
		protected.GET("/favorite/list", eventH.ListFavorist)

		// Seat reservation — any authenticated user
		submitSeatH := sHandler.NewSubmitSeatHandler(seatSvc)
		protected.POST("/events-seat/submit/:id", submitSeatH.Handle)
		protected.POST("/events/favorite", eventH.ToggleFavorite)

		// Payment — customer or admin
		processPaymentH := pHandler.NewProcessPaymentHandler(paymentSvc)
		cancelPaymentH := pHandler.NewCancelPaymentHandler(paymentSvc)
		paymentGroup := protected.Group("/payment")
		paymentGroup.Use(middleware.RolesAllowed(auth.RoleCustomer, auth.RoleAdmin))
		{
			paymentGroup.POST("/process", processPaymentH.Handle)
			paymentGroup.POST("/cancel", cancelPaymentH.Handle)
		}
		protected.POST("/confirm-payment",
			middleware.RolesAllowed(auth.RoleAdmin, auth.RoleOrganization),
			confirmPaymentH.Handle,
		)
	}

	// ── Event Management (Organization) ──────────────────────────────────────
	eventsOrg := api.Group("/events")
	eventsOrg.Use(middleware.AuthRequired(), middleware.RolesAllowed(auth.RoleOrganization))
	{
		eventsOrg.GET("/my", eventH.GetMyEvents)
		eventsOrg.POST("", eventH.CreateEvent)
		eventsOrg.PUT("/:id", eventH.UpdateEvent)
		eventsOrg.POST("/:id/submit", eventH.SubmitForReview)
		eventsOrg.GET("/:id/bookings", bookingAdminH.GetBookingSummary)
		eventsOrg.POST("/:id/seats", seatAdminH.AddSeats)
		eventsOrg.POST("/scanQr", globalAdminH.ScanQR)
	}

	organizerGroup := api.Group("/organizer")
	organizerGroup.Use(middleware.AuthRequired(), middleware.RolesAllowed(auth.RoleOrganization))
	{
		organizerGroup.GET("/event", eventH.OrganizerListEvents)
		organizerGroup.GET("/summary", eventH.GetOrganizerOverallSummary)
		organizerGroup.GET("/events/booking/:eventId", eventH.GetOrganizerEventBookings)
		organizerGroup.GET("/events/summary/:eventId", eventH.GetOrganizerEventSummary)
		organizerGroup.GET("/events/:id/sales", eventH.GetOrganizerEventSales)
		organizerGroup.GET("/events/:id/seats/summary", eventH.GetOrganizerEventSeatsDashboard)
		organizerGroup.GET("/events/compare", eventH.GetOrganizerEventsCompare)
	}

	// ── Admin Dashboard ───────────────────────────────────────────────────
	adminGroup := api.Group("/admin")
	adminGroup.Use(middleware.AuthRequired())
	{
		adminGroup.GET("/summary", middleware.RolesAllowed(auth.RoleAdmin), eventH.GetAdminOverallSummary)
		adminGroup.GET("/orders/stats", middleware.RolesAllowed(auth.RoleAdmin), eventH.GetAdminOrderStats)
		adminGroup.GET("/orders/summary", middleware.RolesAllowed(auth.RoleAdmin), eventH.GetAdminOrdersSummaryStatus)
		adminGroup.GET("/events/top-selling", middleware.RolesAllowed(auth.RoleAdmin), eventH.GetAdminTopSellingEvents)

		adminGroup.GET("/users", middleware.RolesAllowed(auth.RoleAdmin), userH.GetAll)
		adminGroup.PATCH("/users/:id", middleware.RolesAllowed(auth.RoleAdmin), userH.AdminUpdateUser)
		adminGroup.PUT("/users/:id/role", middleware.RolesAllowed(auth.RoleAdmin), userH.UpdateRole)
		adminGroup.DELETE("/users/:id", middleware.RolesAllowed(auth.RoleAdmin), userH.DeleteUser)
		adminGroup.GET("/bookings", middleware.RolesAllowed(auth.RoleAdmin), globalAdminH.GetAllBookings)
		adminGroup.GET("/payments", middleware.RolesAllowed(auth.RoleAdmin), globalAdminH.GetAllPayments)
		adminGroup.GET("/event", middleware.RolesAllowed(auth.RoleAdmin), eventH.AdminListEvents)
		adminGroup.GET("/payouts", middleware.RolesAllowed(auth.RoleAdmin), payoutH.GetAllPayouts)
		adminGroup.POST("/payouts/:id", middleware.RolesAllowed(auth.RoleAdmin), payoutH.ProcessPayout)

		// Shared: Admin or Org
		adminGroup.PATCH("/editEvent/:id",
			middleware.RolesAllowed(auth.RoleAdmin, auth.RoleOrganization),
			eventH.AdminEditEvent,
		)
	}

	// ── Event Review (Admin) ─────────────────────────────────────────────────
	eventsAdmin := api.Group("/events")
	eventsAdmin.Use(middleware.AuthRequired(), middleware.RolesAllowed(auth.RoleAdmin))
	{
		eventsAdmin.GET("/pending", eventH.GetPendingEvents)
		eventsAdmin.POST("/:id/review", eventH.ReviewEvent)
	}

	// ── Payouts (Organization & Admin) ───────────────────────────────────────
	payoutsGroup := api.Group("/payouts")
	payoutsGroup.Use(middleware.AuthRequired())
	{
		// Org: Request and View
		orgPayouts := payoutsGroup.Group("")
		orgPayouts.Use(middleware.RolesAllowed(auth.RoleOrganization))
		{
			orgPayouts.GET("", payoutH.GetPayouts)
			orgPayouts.POST("", payoutH.RequestPayout)
		}
	}

	// Aliases as per design images
	api.GET("/getPayout", middleware.AuthRequired(), middleware.RolesAllowed(auth.RoleOrganization), payoutH.GetPayouts)
	api.POST("/createPayout", middleware.AuthRequired(), middleware.RolesAllowed(auth.RoleOrganization), payoutH.RequestPayout)
}
