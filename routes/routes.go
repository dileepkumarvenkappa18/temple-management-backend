package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/config"
	"github.com/sharath018/temple-management-backend/database"
	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/internal/entity"
	"github.com/sharath018/temple-management-backend/internal/event"
	"github.com/sharath018/temple-management-backend/internal/eventrsvp"
	"github.com/sharath018/temple-management-backend/internal/seva"
	"github.com/sharath018/temple-management-backend/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/sharath018/temple-management-backend/docs"
)

func Setup(r *gin.Engine, cfg *config.Config) {
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")

	// ========== Auth ==========
	authRepo := auth.NewRepository(database.DB)
	authSvc := auth.NewService(authRepo, cfg)
	authHandler := auth.NewHandler(authSvc)

	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.Refresh)
	}

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg))

	// Dashboards
	protected.GET("/superadmin/dashboard", middleware.RBACMiddleware("superadmin"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Super Admin dashboard access granted!"})
	})
	protected.GET("/tenant/dashboard", middleware.RBACMiddleware("templeadmin"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Temple Admin dashboard access granted!"})
	})
	protected.GET("/entity/:id/devotee/dashboard", middleware.RBACMiddleware("devotee"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Devotee dashboard access granted!"})
	})
	protected.GET("/entity/:id/volunteer/dashboard", middleware.RBACMiddleware("volunteer"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Volunteer dashboard access granted!"})
	})

	// ========== Seva ==========
	{
		sevaRepo := seva.NewRepository(database.DB)
		sevaService := seva.NewService(sevaRepo)
		sevaHandler := seva.NewHandler(sevaService)

		sevaRoutes := protected.Group("/sevas")
		sevaRoutes.POST("/", middleware.RBACMiddleware("templeadmin"), sevaHandler.CreateSeva)
		sevaRoutes.GET("/", sevaHandler.GetSevas)
		sevaRoutes.POST("/bookings", middleware.RBACMiddleware("devotee"), sevaHandler.BookSeva)
		sevaRoutes.GET("/my-bookings", middleware.RBACMiddleware("devotee"), sevaHandler.GetMyBookings)
		sevaRoutes.GET("/entity-bookings", middleware.RBACMiddleware("templeadmin"), sevaHandler.GetEntityBookings)
		sevaRoutes.PATCH("/bookings/:id/status", middleware.RBACMiddleware("templeadmin"), sevaHandler.UpdateBookingStatus)
		sevaRoutes.PATCH("/bookings/:id/cancel", middleware.RBACMiddleware("devotee"), sevaHandler.CancelBooking)
	}

	// ========== Entity ==========
	{
		entityRepo := entity.NewRepository(database.DB)
		entityService := entity.NewService(entityRepo)
		entityHandler := entity.NewHandler(entityService)

		entityRoutes := protected.Group("/entities")
		entityRoutes.POST("/", middleware.RBACMiddleware("superadmin", "templeadmin"), entityHandler.CreateEntity)
		entityRoutes.GET("/", middleware.RBACMiddleware("superadmin", "templeadmin"), entityHandler.GetAllEntities)
		entityRoutes.GET("/:id", middleware.RBACMiddleware("superadmin", "templeadmin"), entityHandler.GetEntityByID)
		entityRoutes.PUT("/:id", middleware.RBACMiddleware("superadmin", "templeadmin"), entityHandler.UpdateEntity)
		entityRoutes.PATCH("/:id/status", middleware.RBACMiddleware("superadmin", "templeadmin"), entityHandler.ToggleStatus)
		entityRoutes.DELETE("/:id", middleware.RBACMiddleware("superadmin", "templeadmin"), entityHandler.DeleteEntity)
		entityRoutes.POST("/address", middleware.RBACMiddleware("templeadmin"), entityHandler.AddEntityAddress)
		entityRoutes.GET("/:id/address", middleware.RBACMiddleware("templeadmin", "superadmin"), entityHandler.GetEntityAddress)
		entityRoutes.POST("/documents", middleware.RBACMiddleware("templeadmin"), entityHandler.AddEntityDocument)
		entityRoutes.GET("/:id/documents", middleware.RBACMiddleware("templeadmin", "superadmin"), entityHandler.GetEntityDocuments)
		entityRoutes.POST("/documents/upload", middleware.RBACMiddleware("templeadmin"), entityHandler.UploadEntityDocument)
	}

	// ========== Event & RSVP ==========
	eventRepo := event.NewRepository(database.DB)
	eventService := event.NewService(eventRepo)
	eventHandler := event.NewHandler(eventService)

	{
		eventRoutes := protected.Group("/events")
		eventRoutes.POST("/", middleware.RBACMiddleware("templeadmin"), eventHandler.CreateEvent)
		eventRoutes.GET("/", eventHandler.ListEvents)
		eventRoutes.GET("/upcoming", eventHandler.GetUpcomingEvents)
		eventRoutes.GET("/:id", eventHandler.GetEventByID)
		eventRoutes.PUT("/:id", middleware.RBACMiddleware("templeadmin"), eventHandler.UpdateEvent)
		eventRoutes.DELETE("/:id", middleware.RBACMiddleware("templeadmin"), eventHandler.DeleteEvent)
	}

	{
		rsvpRepo := eventrsvp.NewRepository(database.DB)
		rsvpService := eventrsvp.NewService(rsvpRepo, eventService)
		rsvpHandler := eventrsvp.NewHandler(rsvpService, eventService)

		rsvpRoutes := protected.Group("/event-rsvps")
		rsvpRoutes.POST("/:eventID", middleware.RBACMiddleware("devotee", "volunteer"), rsvpHandler.CreateRSVP)
		rsvpRoutes.GET("/:eventID", middleware.RBACMiddleware("templeadmin"), rsvpHandler.GetRSVPsByEvent)
		rsvpRoutes.GET("/my", middleware.RBACMiddleware("devotee", "volunteer"), rsvpHandler.GetMyRSVPs)
	}
}
