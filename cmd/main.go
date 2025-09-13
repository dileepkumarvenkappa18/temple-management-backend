package main

import (
	"fmt"

	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/sharath018/temple-management-backend/config"
	"github.com/sharath018/temple-management-backend/database"
	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/internal/auditlog"
	"github.com/sharath018/temple-management-backend/internal/entity"
	"github.com/sharath018/temple-management-backend/internal/event"
	"github.com/sharath018/temple-management-backend/internal/eventrsvp"
	"github.com/sharath018/temple-management-backend/internal/notification"
	"github.com/sharath018/temple-management-backend/routes"
	"github.com/sharath018/temple-management-backend/utils"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg)

	// ✅ Init Redis
	if err := utils.InitRedis(); err != nil {
		log.Fatalf("❌ Redis init failed: %v", err)
	}

	// ✅ Init Kafka
	utils.InitializeKafka()

	// ✅ Initialize repositories and services
	authRepo := auth.NewRepository(db)

	// ✅ Initialize audit log service
	auditRepo := auditlog.NewRepository(db)
	auditSvc := auditlog.NewService(auditRepo)

	// ✅ Inject authRepo and auditSvc into notification service
	notificationRepo := notification.NewRepository(db)
	notificationService := notification.NewService(notificationRepo, authRepo, cfg, auditSvc)
	notification.StartKafkaConsumer(notificationService)

	// ✅ Seed roles and super admin
	if err := auth.SeedUserRoles(db); err != nil {
		panic(fmt.Sprintf("❌ Failed to seed roles: %v", err))
	}
	if err := auth.SeedSuperAdminUser(db); err != nil {
		panic(fmt.Sprintf("❌ Failed to seed Super Admin: %v", err))
	}

	// ✅ Auto-migrate models
	if err := db.AutoMigrate(
		&auditlog.AuditLog{},
		&entity.Entity{},
		&event.Event{},
		&eventrsvp.RSVP{},
	); err != nil {
		panic(fmt.Sprintf("❌ DB AutoMigrate failed: %v", err))
	}

	// 🌐 Setup router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.LoadHTMLGlob("templates/*")

	// ✅ Optional request logger for CORS debugging
	router.Use(func(c *gin.Context) {
		log.Printf("👉 %s %s from origin %s", c.Request.Method, c.Request.URL.Path, c.Request.Header.Get("Origin"))
		c.Next()
	})

	// ✅ Global CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
		},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Tenant-ID", "Content-Length", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204)
	})

	// ✅ Register existing routes
	routes.Setup(router, cfg)

	// ✅ File Upload Route
	// Create uploads folder if it doesn't exist
	os.MkdirAll("./uploads", os.ModePerm)

	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"error": "File not found"})
			return
		}

		// Secure the filename
		dst := "./uploads/" + file.Filename

		// Save the uploaded file
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(500, gin.H{"error": "Failed to save file"})
			return
		}

		c.JSON(200, gin.H{"message": fmt.Sprintf("File '%s' uploaded successfully!", file.Filename)})
	})

	// 🚀 Run server
	fmt.Printf("🚀 Server starting on port %s\n", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
