package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/corazawaf/coraza/v3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/config"
	"github.com/sharath018/temple-management-backend/database"
	"github.com/sharath018/temple-management-backend/internal/auditlog"
	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/internal/entity"
	"github.com/sharath018/temple-management-backend/internal/event"
	"github.com/sharath018/temple-management-backend/internal/eventrsvp"
	"github.com/sharath018/temple-management-backend/internal/notification"
	"github.com/sharath018/temple-management-backend/routes"
	"github.com/sharath018/temple-management-backend/utils"
)

func CorazaMiddleware(waf coraza.WAF, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("--> CorazaMiddleware Start")
		tx := waf.NewTransaction()
		defer tx.ProcessLogging()
		defer tx.Close()

		log.Printf("----------------------------------------------")
		log.Printf("c.Request.RemoteAddr %s", c.Request.RemoteAddr)
		log.Printf("c.Request.Host %s", c.Request.Host)
		log.Printf("----------------------------------------------")
		
		tx.ProcessConnection(c.Request.RemoteAddr, 0, c.Request.Host, 0)
		tx.ProcessURI(c.Request.RequestURI, c.Request.Method, "")

		if tx.Interruption() != nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		fmt.Printf("cfg.ProcessRequestHeaders: %s\n", cfg.ProcessRequestHeaders)
		if cfg.ProcessRequestHeaders == "true" {
			if it := tx.ProcessRequestHeaders(); it != nil {
				fmt.Printf("Request headers interrupted:\n")
				fmt.Printf("Reason: (Rule ID: %d, Status: %d)\n", it.RuleID, it.Status)
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		fmt.Printf("cfg.ProcessRequestBody: %s\n", cfg.ProcessRequestBody)
		if cfg.ProcessRequestBody == "true" {
			it, err := tx.ProcessRequestBody()
			if it != nil {
				fmt.Printf("Request body interrupted:\n")
				fmt.Printf("Reason: (Rule ID: %d, Status: %d)\n", it.RuleID, it.Status)
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
			if err != nil {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		if it := tx.Interruption(); it != nil {
			fmt.Printf("Request blocked by WAF\n")
			fmt.Printf("Reason: (Rule ID: %d, Status: %d)\n", it.RuleID, it.Status)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		log.Printf("--> CorazaMiddleware Done")
	}
}

func main() {
	cfg := config.Load()
	db := database.Connect(cfg)

	// Init Redis
	if err := utils.InitRedis(); err != nil {
		log.Fatalf("❌ Redis init failed: %v", err)
	}

	// Init Kafka
	utils.InitializeKafka()

	// Init repositories & services
	authRepo := auth.NewRepository(db)
	auditRepo := auditlog.NewRepository(db)
	auditSvc := auditlog.NewService(auditRepo)

	notificationRepo := notification.NewRepository(db)
	notificationService := notification.NewService(notificationRepo, authRepo, cfg, auditSvc)
	notification.StartKafkaConsumer(notificationService)

	// Seed roles & super admin
	if err := auth.SeedUserRoles(db); err != nil {
		panic(fmt.Sprintf("❌ Failed to seed roles: %v", err))
	}
	if err := auth.SeedSuperAdminUser(db); err != nil {
		panic(fmt.Sprintf("❌ Failed to seed Super Admin: %v", err))
	}

	// Auto-migrate models
	if err := db.AutoMigrate(
		&auditlog.AuditLog{},
		&entity.Entity{},
		&event.Event{},
		&eventrsvp.RSVP{},
	); err != nil {
		panic(fmt.Sprintf("❌ DB AutoMigrate failed: %v", err))
	}

	// Setup Gin router
	router := gin.New()

	// WAF setup
	log.Printf("cfg: %v", cfg)
	log.Printf("cfg.WafEnable: %s", cfg.WafEnable)
	if cfg.WafEnable == "true" {
		waf, err := coraza.NewWAF(coraza.NewWAFConfig().
			WithDirectivesFromFile("/data/wafdata/crs-setup.conf").
			WithDirectivesFromFile("/data/wafdata/rules/*.conf"))

		if err != nil {
			panic(fmt.Sprintf("Failed to create WAF: %s", err.Error()))
		}
		router.Use(CorazaMiddleware(waf, cfg))
	}

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.LoadHTMLGlob("templates/*")

	// Request logger
	router.Use(func(c *gin.Context) {
		log.Printf("👉 %s %s from origin %s", c.Request.Method, c.Request.URL.Path, c.Request.Header.Get("Origin"))
		c.Next()
	})

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://localhost:5173", "https://127.0.0.1:5173", "https://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Tenant-ID", "Content-Length", "X-Requested-With", "Cache-Control", "Pragma"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Content-Disposition", "Cache-Control", "Pragma", "Expires"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Handle preflight requests
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "https://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Tenant-ID, Content-Length, X-Requested-With, Cache-Control, Pragma")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, Content-Disposition, Cache-Control, Pragma, Expires")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Status(204)
	})

	// Create uploads directory
	uploadDir := "/data/uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		panic(fmt.Sprintf("❌ Failed to create upload directory: %v", err))
	}

	// Static file serving
	router.Static("/uploads", "/data/uploads")

	// File serving endpoint
	router.GET("/files/*filepath", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "https://localhost:5173")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, Content-Disposition, Cache-Control, Pragma, Expires")

		requestedPath := c.Param("filepath")
		fullPath := filepath.Join(uploadDir, requestedPath)
		cleanPath := filepath.Clean(fullPath)

		if !strings.HasPrefix(cleanPath, filepath.Clean(uploadDir)) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied", "message": "Invalid file path"})
			return
		}

		fileInfo, err := os.Stat(cleanPath)
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "File access error"})
			return
		}

		filename := filepath.Base(cleanPath)
		ext := strings.ToLower(filepath.Ext(filename))

		switch ext {
		case ".pdf":
			c.Header("Content-Type", "application/pdf")
		case ".jpg", ".jpeg":
			c.Header("Content-Type", "image/jpeg")
		case ".png":
			c.Header("Content-Type", "image/png")
		case ".doc":
			c.Header("Content-Type", "application/msword")
		case ".docx":
			c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
		default:
			c.Header("Content-Type", "application/octet-stream")
		}

		c.Header("Cache-Control", "public, max-age=3600")
		c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

		if c.Query("download") == "true" || c.GetHeader("Accept") == "application/octet-stream" {
			c.Header("Content-Description", "File Transfer")
			c.Header("Content-Transfer-Encoding", "binary")
			c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
		} else {
			c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
		}

		c.File(cleanPath)
	})

	// Entity file download endpoint
	router.GET("/api/v1/entities/:id/files/:filename", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "https://localhost:5173")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, Content-Disposition, Cache-Control, Pragma, Expires")

		entityID := c.Param("id")
		filename := c.Param("filename")

		if entityID == "" || filename == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
			return
		}

		filePath := filepath.Join(uploadDir, entityID, filename)
		cleanPath := filepath.Clean(filePath)
		expectedPrefix := filepath.Clean(filepath.Join(uploadDir, entityID))

		if !strings.HasPrefix(cleanPath, expectedPrefix) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		fileInfo, err := os.Stat(cleanPath)
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "File access error"})
			return
		}

		ext := strings.ToLower(filepath.Ext(filename))
		switch ext {
		case ".pdf":
			c.Header("Content-Type", "application/pdf")
		case ".jpg", ".jpeg":
			c.Header("Content-Type", "image/jpeg")
		case ".png":
			c.Header("Content-Type", "image/png")
		case ".doc":
			c.Header("Content-Type", "application/msword")
		case ".docx":
			c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
		default:
			c.Header("Content-Type", "application/octet-stream")
		}

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
		c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		c.File(cleanPath)
		log.Printf("File downloaded: %s/%s", entityID, filename)
	})

	// Bulk download endpoint
	router.GET("/api/v1/entities/:id/files-all", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "https://localhost:5173")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, Content-Disposition, Cache-Control, Pragma, Expires")

		entityID := c.Param("id")
		entityDir := filepath.Join(uploadDir, entityID)

		if _, err := os.Stat(entityDir); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No files found for this entity"})
			return
		}

		zipFileName := fmt.Sprintf("Entity_%s_Documents_%s.zip", entityID, time.Now().Format("20060102_150405"))

		c.Header("Content-Type", "application/zip")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", zipFileName))
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		zipWriter := zip.NewWriter(c.Writer)
		defer func() {
			if err := zipWriter.Close(); err != nil {
				log.Printf("Error closing zip writer: %v", err)
			}
		}()

		err := filepath.Walk(entityDir, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("Error walking file %s: %v", filePath, err)
				return nil
			}

			if info.IsDir() {
				return nil
			}

			relPath, err := filepath.Rel(entityDir, filePath)
			if err != nil {
				log.Printf("Error getting relative path for %s: %v", filePath, err)
				return nil
			}

			zipFile, err := zipWriter.Create(relPath)
			if err != nil {
				log.Printf("Error creating zip entry for %s: %v", relPath, err)
				return nil
			}

			srcFile, err := os.Open(filePath)
			if err != nil {
				log.Printf("Error opening file %s: %v", filePath, err)
				return nil
			}
			defer srcFile.Close()

			_, err = io.Copy(zipFile, srcFile)
			if err != nil {
				log.Printf("Error copying file %s to zip: %v", filePath, err)
				return nil
			}

			return nil
		})

		if err != nil {
			log.Printf("Error creating ZIP for entity %s: %v", entityID, err)
		}

		log.Printf("ZIP file created for entity %s", entityID)
	})

	// Upload endpoint
	router.POST("/upload", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "https://localhost:5173")
		c.Header("Access-Control-Allow-Credentials", "true")

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"error": "File not found in request"})
			return
		}

		filename := filepath.Base(file.Filename)
		dst := filepath.Join(uploadDir, filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to save file: %v", err)})
			return
		}

		c.JSON(200, gin.H{
			"message": fmt.Sprintf("File '%s' uploaded successfully!", filename),
			"path":    dst,
			"url":     fmt.Sprintf("/uploads/%s", filename),
		})
	})

	// Debug endpoint
	router.GET("/debug/entity-files", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "https://localhost:5173")
		c.Header("Access-Control-Allow-Credentials", "true")

		type EntityFileInfo struct {
			EntityID   string   `json:"entity_id"`
			FilesCount int      `json:"files_count"`
			Files      []string `json:"files"`
			TotalSize  int64    `json:"total_size"`
		}
		var entityFiles []EntityFileInfo

		entries, err := os.ReadDir(uploadDir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read upload directory"})
			return
		}

		for _, entry := range entries {
			if entry.IsDir() && entry.Name() != "temp_uploads" {
				entityDir := filepath.Join(uploadDir, entry.Name())
				files, _ := os.ReadDir(entityDir)
				var fileNames []string
				var totalSize int64
				for _, f := range files {
					if !f.IsDir() {
						fileNames = append(fileNames, f.Name())
						if info, err := f.Info(); err == nil {
							totalSize += info.Size()
						}
					}
				}
				if len(fileNames) > 0 {
					entityFiles = append(entityFiles, EntityFileInfo{
						EntityID:   entry.Name(),
						FilesCount: len(fileNames),
						Files:      fileNames,
						TotalSize:  totalSize,
					})
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"total_entities_with_files": len(entityFiles),
			"entity_files":              entityFiles,
		})
	})

	// File info endpoint
	router.GET("/api/v1/entities/:id/files/info", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "https://localhost:5173")
		c.Header("Access-Control-Allow-Credentials", "true")

		entityID := c.Param("id")
		entityDir := filepath.Join(uploadDir, entityID)
		if _, err := os.Stat(entityDir); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No files found for this entity"})
			return
		}

		files, err := os.ReadDir(entityDir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read entity files"})
			return
		}

		type FileInfo struct {
			FileName    string `json:"file_name"`
			Size        int64  `json:"size"`
			ModTime     string `json:"modified_time"`
			FileType    string `json:"file_type"`
			ViewURL     string `json:"view_url"`
			DownloadURL string `json:"download_url"`
		}

		var fileInfos []FileInfo
		var totalSize int64
		for _, file := range files {
			if !file.IsDir() {
				info, err := file.Info()
				if err != nil {
					continue
				}
				ext := strings.ToLower(filepath.Ext(file.Name()))
				fileType := strings.ToUpper(strings.TrimPrefix(ext, "."))
				if fileType == "" {
					fileType = "unknown"
				}
				fileInfos = append(fileInfos, FileInfo{
					FileName:    file.Name(),
					Size:        info.Size(),
					ModTime:     info.ModTime().Format("2006-01-02 15:04:05"),
					FileType:    fileType,
					ViewURL:     fmt.Sprintf("/files/%s/%s", entityID, file.Name()),
					DownloadURL: fmt.Sprintf("/api/v1/entities/%s/files/%s", entityID, file.Name()),
				})
				totalSize += info.Size()
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"entity_id":         entityID,
			"files_count":       len(fileInfos),
			"total_size":        totalSize,
			"files":             fileInfos,
			"bulk_download_url": fmt.Sprintf("/api/v1/entities/%s/files-all", entityID),
		})
	})

	// Register existing routes (this handles other API routes)
	routes.Setup(router, cfg)

	// Set trusted proxies for security
	router.SetTrustedProxies([]string{"127.0.0.1", "::1", "172.19.0.0/16"})

	// Server startup logs
	fmt.Printf("🚀 Server starting on port %s\n", cfg.Port)
	fmt.Printf("📁 Upload directory: %s\n", uploadDir)
	fmt.Printf("🌐 Static files: https://localhost:%s/files/{path}\n", cfg.Port)
	fmt.Printf("📥 Download file: https://localhost:%s/api/v1/entities/{id}/files/{filename}\n", cfg.Port)
	fmt.Printf("📦 Bulk download: https://localhost:%s/api/v1/entities/{id}/files-all\n", cfg.Port)

	// Start server (HTTPS or HTTP)
	if cfg.EnableHTTPS == "true" {
		fmt.Printf("🔒 Starting HTTPS server on port %s\n", cfg.Port)
		if err := router.RunTLS(":"+cfg.Port, cfg.TLSCertPath, cfg.TLSKeyPath); err != nil {
			panic(fmt.Sprintf("Failed to start HTTPS server: %v", err))
		}
	} else {
		fmt.Printf("🔓 Starting HTTP server on port %s\n", cfg.Port)
		if err := router.Run(":" + cfg.Port); err != nil {
			panic(fmt.Sprintf("Failed to start HTTP server: %v", err))
		}
	}
}