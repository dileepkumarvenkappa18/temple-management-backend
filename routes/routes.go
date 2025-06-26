package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/config"
	"github.com/sharath018/temple-management-backend/database"
	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/sharath018/temple-management-backend/docs" // needed for Swagger UI
)

func Setup(r *gin.Engine, cfg *config.Config) {
	// Health check
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	// ✅ Swagger docs route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create service, repo, handler for auth
	repo := auth.NewRepository(database.DB)
	svc := auth.NewService(repo, cfg)
	h := auth.NewHandler(svc)

	// ✅ Grouped under /api/v1
	api := r.Group("/api/v1")

	// Public Auth routes
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
		authGroup.POST("/refresh", h.Refresh)
	}

	// ✅ Authenticated routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg))

	// Dashboard routes (role-based RBAC)
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

	// // Optional: Protected test route
	// protected.GET("/protected", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{"message": "authenticated"})
	// })
}
