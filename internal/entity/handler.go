package entity

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/middleware"
)

type Handler struct {
	Service   *Service
	UploadDir string // filesystem base, now using persistent volume at /data/uploads
	BaseURL   string // URL base for secure endpoints, e.g. "/api/v1/entities"
	MaxSize   int64  // 10MB default
}

type TempFileInfo struct {
	TempPath     string
	FileType     string // registration_cert, trust_deed, property_docs, additional_docs
	FileName     string
	OriginalName string
	FileSize     int64
	ContentType  string
	UploadedAt   time.Time
}

type EntityDirectory struct {
	EntityID   string   `json:"entity_id"`
	FilesCount int      `json:"files_count"`
	Files      []string `json:"files"`
}

type FileDetails struct {
	FileName string `json:"file_name"`
	FileURL  string `json:"file_url"`
	Size     int64  `json:"size"`
}

func NewHandler(s *Service, uploadDir, baseURL string) *Handler {
	// Use persistent volume path if not specified
	if uploadDir == "" || uploadDir == "./uploads" {
		uploadDir = "/data/uploads"
	}
	
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("Failed to create upload directory: %v", err)
	}
	// Use secure endpoint base URL
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "/api/v1/entities"
	}
	return &Handler{
		Service:   s,
		UploadDir: uploadDir,
		BaseURL:   baseURL,
		MaxSize:   10 * 1024 * 1024,
	}
}

// ========== SECURE FILE SERVING HANDLERS ==========

// SecureFileHandler handles authenticated individual file downloads
func (h *Handler) SecureFileHandler(c *gin.Context) {
	entityID := c.Param("entityID")
	filename := c.Param("filename")

	// Validate parameters
	if entityID == "" || filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid parameters",
			"message": "Entity ID and filename are required",
		})
		return
	}

	// Get authenticated user
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication required",
			"message": "You must be logged in to access files",
		})
		return
	}

	user, ok := userVal.(auth.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user context",
		})
		return
	}

	// Convert entityID to uint
	entityIDInt, err := strconv.Atoi(entityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid entity ID format",
		})
		return
	}
	entityIDUint := uint(entityIDInt)

	// Check if user has access to this entity
	hasAccess, err := h.validateEntityAccess(user, entityIDUint)
	if err != nil {
		log.Printf("Error validating entity access for user %d, entity %d: %v", user.ID, entityIDUint, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Access validation failed",
		})
		return
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Access denied",
			"message": "You don't have permission to access files for this entity",
		})
		return
	}

	// Construct and validate file path using persistent volume path
	filePath := filepath.Join(h.UploadDir, entityID, filename)
	cleanPath := filepath.Clean(filePath)
	expectedPrefix := filepath.Clean(filepath.Join(h.UploadDir, entityID))

	// Security check - ensure path is within expected directory
	if !strings.HasPrefix(cleanPath, expectedPrefix) {
		log.Printf("Security violation: Attempted path traversal by user %d for path %s", user.ID, cleanPath)
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Access denied",
			"message": "Invalid file path",
		})
		return
	}

	// Check if file exists
	fileInfo, err := os.Stat(cleanPath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "File not found",
			"message": "The requested file does not exist",
		})
		return
	}

	if err != nil {
		log.Printf("Error accessing file %s: %v", cleanPath, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "File access error",
		})
		return
	}

	// Set appropriate headers based on file type and request
	h.setFileHeaders(c, filename, fileInfo.Size())

	// Log file access for audit
	log.Printf("File accessed: %s by user %d (%s)", cleanPath, user.ID, user.Email)

	// Serve the file
	c.File(cleanPath)
}

// SecureBulkDownloadHandler handles authenticated bulk file downloads as ZIP
func (h *Handler) SecureBulkDownloadHandler(c *gin.Context) {
	entityID := c.Param("entityID")

	// Validate parameters
	if entityID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid parameters",
			"message": "Entity ID is required",
		})
		return
	}

	// Get authenticated user
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication required",
			"message": "You must be logged in to download files",
		})
		return
	}

	user, ok := userVal.(auth.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user context",
		})
		return
	}

	// Convert entityID to uint
	entityIDInt, err := strconv.Atoi(entityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid entity ID format",
		})
		return
	}
	entityIDUint := uint(entityIDInt)

	// Check if user has access to this entity
	hasAccess, err := h.validateEntityAccess(user, entityIDUint)
	if err != nil {
		log.Printf("Error validating entity access for user %d, entity %d: %v", user.ID, entityIDUint, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Access validation failed",
		})
		return
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Access denied",
			"message": "You don't have permission to download files for this entity",
		})
		return
	}

	// Check if entity directory exists using persistent volume path
	entityDir := filepath.Join(h.UploadDir, entityID)
	if _, err := os.Stat(entityDir); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "No files found",
			"message": "No files found for this entity",
		})
		return
	}

	// Set ZIP download headers
	zipFileName := fmt.Sprintf("Entity_%s_Documents_%s.zip", entityID, time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", zipFileName))
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	// Create ZIP and stream to response
	err = h.createZipStream(c.Writer, entityDir)
	if err != nil {
		log.Printf("Error creating ZIP for entity %s: %v", entityID, err)
		// Can't send JSON error here as headers are already sent
		return
	}

	// Log bulk download for audit
	log.Printf("Bulk download: Entity %s by user %d (%s)", entityID, user.ID, user.Email)
}

// ========== ACCESS VALIDATION ==========

// validateEntityAccess checks if user has access to the specified entity
func (h *Handler) validateEntityAccess(user auth.User, entityID uint) (bool, error) {
	switch user.Role.RoleName {
	case "superadmin":
		// Super admins can access all entities
		return true, nil

	case "templeadmin":
		// Temple admins can access entities they created
		entities, err := h.Service.GetEntitiesByCreator(user.ID)
		if err != nil {
			return false, err
		}
		for _, entity := range entities {
			if entity.ID == entityID {
				return true, nil
			}
		}
		return false, nil

	case "standarduser", "monitoringuser":
		// Check if user is assigned to this entity via tenant
		tenantID, err := h.Service.Repo.GetTenantIDForUser(user.ID)
		if err != nil {
			return false, err
		}

		// Check if the entity belongs to the user's tenant
		entities, err := h.Service.GetEntitiesByCreator(tenantID)
		if err != nil {
			return false, err
		}

		for _, entity := range entities {
			if entity.ID == entityID {
				return true, nil
			}
		}

		// Also check if entity ID matches tenant ID (for mock entities)
		if entityID == tenantID {
			return true, nil
		}

		return false, nil

	case "devotee", "volunteer":
		// Devotees and volunteers can only access entities they're members of
		// Use the repository to check membership directly
		hasActiveMembership, err := h.Service.Repo.CheckUserEntityMembership(user.ID, entityID)
		if err != nil {
			return false, err
		}
		return hasActiveMembership, nil

	default:
		return false, nil
	}
}

// ========== HELPER FUNCTIONS ==========

// setFileHeaders sets appropriate HTTP headers for file serving
func (h *Handler) setFileHeaders(c *gin.Context, filename string, fileSize int64) {
	ext := strings.ToLower(filepath.Ext(filename))

	// Set content type based on file extension
	switch ext {
	case ".pdf":
		c.Header("Content-Type", "application/pdf")
	case ".jpg", ".jpeg":
		c.Header("Content-Type", "image/jpeg")
	case ".png":
		c.Header("Content-Type", "image/png")
	case ".gif":
		c.Header("Content-Type", "image/gif")
	case ".doc":
		c.Header("Content-Type", "application/msword")
	case ".docx":
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	case ".xls":
		c.Header("Content-Type", "application/vnd.ms-excel")
	case ".xlsx":
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	case ".txt":
		c.Header("Content-Type", "text/plain")
	default:
		c.Header("Content-Type", "application/octet-stream")
	}

	// Set content length
	c.Header("Content-Length", fmt.Sprintf("%d", fileSize))

	// Check if download is requested
	download := c.Query("download")
	if download == "true" || download == "1" {
		// Force download
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
	} else {
		// Inline viewing (for PDFs, images, etc.)
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
		// Set cache headers for better performance
		c.Header("Cache-Control", "public, max-age=3600")
		c.Header("Last-Modified", time.Now().Format(http.TimeFormat))
	}

	// Security headers
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-Frame-Options", "SAMEORIGIN")
}

// createZipStream creates a ZIP file and streams it directly to the response writer
func (h *Handler) createZipStream(w io.Writer, sourceDir string) error {
	zipWriter := zip.NewWriter(w)
	defer func() {
		if err := zipWriter.Close(); err != nil {
			log.Printf("Error closing zip writer: %v", err)
		}
	}()

	return filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error walking file %s: %v", filePath, err)
			return nil // Continue with other files
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(sourceDir, filePath)
		if err != nil {
			log.Printf("Error getting relative path for %s: %v", filePath, err)
			return nil // Continue with other files
		}

		zipFile, err := zipWriter.Create(relPath)
		if err != nil {
			log.Printf("Error creating zip entry for %s: %v", relPath, err)
			return nil // Continue with other files
		}

		srcFile, err := os.Open(filePath)
		if err != nil {
			log.Printf("Error opening file %s: %v", filePath, err)
			return nil // Continue with other files
		}
		defer srcFile.Close()

		_, err = io.Copy(zipFile, srcFile)
		if err != nil {
			log.Printf("Error copying file %s to zip: %v", filePath, err)
			return nil // Continue with other files
		}

		return nil
	})
}

// buildSecureFileURL creates a URL that points to the secure file endpoint
func (h *Handler) buildSecureFileURL(entityID uint, filename string) string {
	return fmt.Sprintf("%s/%d/files/%s", h.BaseURL, entityID, filename)
}

// ========== ENTITY CREATION WITH FILE UPLOADS ==========

// CreateEntity handles temple creation with optional uploads
func (h *Handler) CreateEntity(c *gin.Context) {
	contentType := c.GetHeader("Content-Type")
	isMultipart := strings.Contains(contentType, "multipart/form-data")

	var input Entity
	var tempFiles []TempFileInfo

	if isMultipart {
		if err := h.handleMultipartFormData(c, &input, &tempFiles); err != nil {
			log.Printf("Multipart Form Error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data", "details": err.Error()})
			return
		}
	} else {
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Printf("JSON Bind Error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
			return
		}
	}

	// Required fields validation
	if input.TempleType == "" || input.State == "" || input.EstablishedYear == nil {
		h.cleanupTempFiles(tempFiles)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Temple Type, State, and Established Year are required"})
		return
	}
	if strings.TrimSpace(input.StreetAddress) == "" {
		h.cleanupTempFiles(tempFiles)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Street address is required"})
		return
	}

	// Auth and access validation
	userVal, exists := c.Get("user")
	if !exists {
		h.cleanupTempFiles(tempFiles)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userObj := userVal.(auth.User)
	userID := userObj.ID
	userRole := userObj.Role.RoleName

	var accessContext middleware.AccessContext
	if v, ok := c.Get("access_context"); ok {
		accessContext, _ = v.(middleware.AccessContext)
	}

	// Determine CreatedBy based on user role
	switch userRole {
	case "superadmin":
		if accessContext.AssignedEntityID != nil {
			input.CreatedBy = *accessContext.AssignedEntityID
		} else {
			tenantID, err := h.Service.Repo.GetTenantIDForUser(userID)
			if err != nil || tenantID == 0 {
				h.cleanupTempFiles(tempFiles)
				c.JSON(http.StatusForbidden, gin.H{"error": "User is not assigned to any tenant"})
				return
			}
			input.CreatedBy = tenantID
		}
	case "templeadmin":
		input.CreatedBy = userID
	case "standarduser", "monitoringuser":
		if accessContext.AssignedEntityID != nil {
			input.CreatedBy = *accessContext.AssignedEntityID
		} else {
			tenantID, err := h.Service.Repo.GetTenantIDForUser(userID)
			if err != nil || tenantID == 0 {
				h.cleanupTempFiles(tempFiles)
				c.JSON(http.StatusForbidden, gin.H{"error": "User is not assigned to any tenant"})
				return
			}
			input.CreatedBy = tenantID
		}
	default:
		h.cleanupTempFiles(tempFiles)
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid user role for temple creation"})
		return
	}

	if input.Status == "" {
		input.Status = "pending"
	}
	ip := middleware.GetIPFromContext(c)

	// Create entity without file URLs first
	if err := h.Service.CreateEntity(&input, userID, ip); err != nil {
		log.Printf("Service Error: %v", err)
		h.cleanupTempFiles(tempFiles)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entity", "details": err.Error()})
		return
	}
	log.Printf("Entity created successfully with ID: %d", input.ID)

	// Move files and update URLs
	finalFileInfos := make(map[string]FileInfo)
	if len(tempFiles) > 0 {
		if err := h.moveFilesToFinalLocation(&input, tempFiles, &finalFileInfos); err != nil {
			log.Printf("Error moving files for entity %d: %v", input.ID, err)
			c.JSON(http.StatusCreated, gin.H{
				"message":    "Temple created but some files could not be processed",
				"temple_id":  input.ID,
				"file_error": err.Error(),
			})
			return
		}
		if err := h.updateEntityWithFileInfo(&input); err != nil {
			log.Printf("Error updating entity %d with file info: %v", input.ID, err)
		}
		log.Printf("Files processed successfully for entity %d", input.ID)
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":        "Temple registration request submitted successfully",
		"temple_id":      input.ID,
		"uploaded_files": finalFileInfos,
	})
}

func (h *Handler) handleMultipartFormData(c *gin.Context, input *Entity, tempFiles *[]TempFileInfo) error {
	form, err := c.MultipartForm()
	if err != nil {
		return fmt.Errorf("failed to parse multipart form: %v", err)
	}

	// Parse text fields
	input.Name = h.getFormValue(form, "name")
	if v := h.getFormValue(form, "main_deity"); v != "" {
		input.MainDeity = &v
	}
	input.TempleType = h.getFormValue(form, "temple_type")
	if yearStr := h.getFormValue(form, "established_year"); yearStr != "" {
		if year, err := strconv.ParseUint(yearStr, 10, 32); err == nil {
			yy := uint(year)
			input.EstablishedYear = &yy
		}
	}
	input.Phone = h.getFormValue(form, "phone")
	input.Email = h.getFormValue(form, "email")
	input.Description = h.getFormValue(form, "description")
	input.StreetAddress = h.getFormValue(form, "street_address")
	input.City = h.getFormValue(form, "city")
	input.District = h.getFormValue(form, "district")
	input.State = h.getFormValue(form, "state")
	input.Pincode = h.getFormValue(form, "pincode")
	input.Landmark = h.getFormValue(form, "landmark")
	input.MapLink = h.getFormValue(form, "map_link")

	if err := h.processFileUploadsToTemp(form, tempFiles); err != nil {
		return fmt.Errorf("failed to process file uploads: %v", err)
	}
	return nil
}

func (h *Handler) processFileUploadsToTemp(form *multipart.Form, tempFiles *[]TempFileInfo) error {
	// Create temp directory in persistent volume
	tempSessionDir := filepath.Join(h.UploadDir, "temp_uploads", uuid.New().String())
	if err := os.MkdirAll(tempSessionDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	log.Printf("Created temp directory: %s", tempSessionDir)

	// Single-file fields
	if reg := form.File["registration_cert"]; len(reg) > 0 {
		info, err := h.uploadFileToTemp(reg[0], tempSessionDir, "registration_cert")
		if err != nil {
			return fmt.Errorf("failed to upload registration certificate: %v", err)
		}
		*tempFiles = append(*tempFiles, info)
	}
	if trust := form.File["trust_deed"]; len(trust) > 0 {
		info, err := h.uploadFileToTemp(trust[0], tempSessionDir, "trust_deed")
		if err != nil {
			return fmt.Errorf("failed to upload trust deed: %v", err)
		}
		*tempFiles = append(*tempFiles, info)
	}
	if prop := form.File["property_docs"]; len(prop) > 0 {
		info, err := h.uploadFileToTemp(prop[0], tempSessionDir, "property_docs")
		if err != nil {
			return fmt.Errorf("failed to upload property documents: %v", err)
		}
		*tempFiles = append(*tempFiles, info)
	}

	// Multiple additional docs
	for i := 0; i < 10; i++ {
		field := fmt.Sprintf("additional_docs_%d", i)
		if add := form.File[field]; len(add) > 0 {
			info, err := h.uploadFileToTemp(add[0], tempSessionDir, "additional_docs")
			if err != nil {
				log.Printf("Warning: Failed to upload additional document %d: %v", i, err)
				continue
			}
			*tempFiles = append(*tempFiles, info)
		}
	}
	log.Printf("Total temp files processed: %d", len(*tempFiles))
	return nil
}

func (h *Handler) uploadFileToTemp(file *multipart.FileHeader, tempDir, fileType string) (TempFileInfo, error) {
	var out TempFileInfo

	if err := h.validateFile(file); err != nil {
		return out, err
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	fileName := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	tempPath := filepath.Join(tempDir, fileName)

	src, err := file.Open()
	if err != nil {
		return out, fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(tempPath)
	if err != nil {
		return out, fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return out, fmt.Errorf("failed to copy file: %v", err)
	}

	out = TempFileInfo{
		TempPath:     tempPath,
		FileType:     fileType,
		FileName:     fileName,
		OriginalName: file.Filename,
		FileSize:     file.Size,
		ContentType:  sniffOrByExt(ext),
		UploadedAt:   time.Now(),
	}
	return out, nil
}

func (h *Handler) moveFilesToFinalLocation(entity *Entity, tempFiles []TempFileInfo, finalFileInfos *map[string]FileInfo) error {
	// Use persistent volume path for entity directory
	entityDir := filepath.Join(h.UploadDir, strconv.FormatUint(uint64(entity.ID), 10))
	if err := os.MkdirAll(entityDir, 0755); err != nil {
		return fmt.Errorf("failed to create entity directory: %v", err)
	}
	log.Printf("Created entity directory: %s for entity %d", entityDir, entity.ID)

	*finalFileInfos = make(map[string]FileInfo)
	var additionalFiles []FileInfo

	for _, tf := range tempFiles {
		finalFileName := tf.FileName
		finalPath := filepath.Join(entityDir, finalFileName)

		// Move file
		if err := os.Rename(tf.TempPath, finalPath); err != nil {
			if err := copyFile(tf.TempPath, finalPath); err != nil {
				log.Printf("Failed to move/copy file %s to %s: %v", tf.TempPath, finalPath, err)
				return fmt.Errorf("failed to persist file %s: %v", tf.FileName, err)
			}
			_ = os.Remove(tf.TempPath)
		}

		// Build secure URL using the new endpoint
		fileURL := h.buildSecureFileURL(entity.ID, finalFileName)

		fi := FileInfo{
			FileName:     finalFileName,
			FileURL:      fileURL,
			FileSize:     tf.FileSize,
			FileType:     tf.ContentType,
			UploadedAt:   tf.UploadedAt,
			OriginalName: tf.OriginalName,
		}

		switch tf.FileType {
		case "registration_cert":
			(*finalFileInfos)["registration_cert"] = fi
			entity.RegistrationCertURL = fileURL
			if b, err := json.Marshal(fi); err == nil {
				entity.RegistrationCertInfo = string(b)
			}
		case "trust_deed":
			(*finalFileInfos)["trust_deed"] = fi
			entity.TrustDeedURL = fileURL
			if b, err := json.Marshal(fi); err == nil {
				entity.TrustDeedInfo = string(b)
			}
		case "property_docs":
			(*finalFileInfos)["property_docs"] = fi
			entity.PropertyDocsURL = fileURL
			if b, err := json.Marshal(fi); err == nil {
				entity.PropertyDocsInfo = string(b)
			}
		case "additional_docs":
			additionalFiles = append(additionalFiles, fi)
		}
	}

	// Handle additional docs
	if len(additionalFiles) > 0 {
		var urlList []string
		for _, x := range additionalFiles {
			urlList = append(urlList, x.FileURL)
		}
		if b, err := json.Marshal(urlList); err == nil {
			entity.AdditionalDocsURLs = string(b)
		}
		if b, err := json.Marshal(additionalFiles); err == nil {
			entity.AdditionalDocsInfo = string(b)
		}
		(*finalFileInfos)["additional_docs"] = FileInfo{
			FileName: fmt.Sprintf("%d_additional_files", len(additionalFiles)),
			FileURL:  "",
		}
	}

	// Clean temp files
	h.cleanupTempFiles(tempFiles)

	log.Printf("Successfully processed %d files for entity %d", len(tempFiles), entity.ID)
	return nil
}

func (h *Handler) updateEntityWithFileInfo(entity *Entity) error {
	return h.Service.Repo.UpdateEntity(*entity)
}

func (h *Handler) validateFile(file *multipart.FileHeader) error {
	if file.Size > h.MaxSize {
		return fmt.Errorf("file size exceeds %dMB limit", h.MaxSize/(1024*1024))
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{
		".pdf": true, ".jpg": true, ".jpeg": true, ".png": true, ".doc": true, ".docx": true,
	}
	if !allowed[ext] {
		return fmt.Errorf("file type %s not allowed", ext)
	}
	return nil
}

func (h *Handler) getFormValue(form *multipart.Form, key string) string {
	if v, ok := form.Value[key]; ok && len(v) > 0 {
		return strings.TrimSpace(v[0])
	}
	return ""
}

func (h *Handler) cleanupTempFiles(tempFiles []TempFileInfo) {
	for _, tf := range tempFiles {
		_ = os.Remove(tf.TempPath)
		_ = os.Remove(filepath.Dir(tf.TempPath))
	}
}

// ========== DIRECTORY INTROSPECTION ==========

// GetAllEntityDirectories - Superadmin only
func (h *Handler) GetAllEntityDirectories(c *gin.Context) {
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userObj := userVal.(auth.User)
	if userObj.Role.RoleName != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only superadmins can view all entity directories"})
		return
	}

	var directories []EntityDirectory
	entries, err := os.ReadDir(h.UploadDir)
	if err != nil {
		log.Printf("Error reading upload directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read upload directory"})
		return
	}
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "temp_uploads" {
			entityDir := filepath.Join(h.UploadDir, entry.Name())
			files, err := os.ReadDir(entityDir)
			if err != nil {
				log.Printf("Error reading entity directory %s: %v", entry.Name(), err)
				continue
			}
			var names []string
			for _, f := range files {
				if !f.IsDir() {
					names = append(names, f.Name())
				}
			}
			if len(names) > 0 {
				directories = append(directories, EntityDirectory{
					EntityID:   entry.Name(),
					FilesCount: len(names),
					Files:      names,
				})
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"total_entities_with_files": len(directories),
		"directories":               directories,
	})
}

// GetEntityFiles - Requires temple access
func (h *Handler) GetEntityFiles(c *gin.Context) {
	entityID := c.Param("id")

	accessVal, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access context"})
		return
	}
	accessCtx, ok := accessVal.(middleware.AccessContext)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access context"})
		return
	}

	idInt, err := strconv.Atoi(entityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}
	entityIDUint := uint(idInt)

	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userObj := userVal.(auth.User)

	hasAccess := false
	switch userObj.Role.RoleName {
	case "superadmin":
		hasAccess = true
	case "templeadmin":
		hasAccess = (accessCtx.DirectEntityID != nil && *accessCtx.DirectEntityID == entityIDUint)
	case "standarduser", "monitoringuser":
		hasAccess = (accessCtx.AssignedEntityID != nil && *accessCtx.AssignedEntityID == entityIDUint)
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to files for this entity"})
		return
	}

	// Use persistent volume path
	entityDir := filepath.Join(h.UploadDir, entityID)
	if _, err := os.Stat(entityDir); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "No files found for this entity"})
		return
	}

	entries, err := os.ReadDir(entityDir)
	if err != nil {
		log.Printf("Error reading entity directory %s: %v", entityID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read entity files"})
		return
	}

	var out []FileDetails
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		// Use secure URL for files
		url := h.buildSecureFileURL(entityIDUint, e.Name())
		out = append(out, FileDetails{
			FileName: e.Name(),
			FileURL:  url,
			Size:     info.Size(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"entity_id":   entityID,
		"files_count": len(out),
		"files":       out,
	})
}

// ========== EXISTING ENTITY CRUD METHODS ==========

// GetAllEntities retrieves entities based on user role and permissions
func (h *Handler) GetAllEntities(c *gin.Context) {
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	user, ok := userVal.(auth.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user object"})
		return
	}

	accessContextVal, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access context"})
		return
	}
	accessContext, ok := accessContextVal.(middleware.AccessContext)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access context"})
		return
	}

	var entities []Entity
	var err error

	switch user.Role.RoleName {
	case "superadmin":
		entities, err = h.Service.GetAllEntities()
		
	case "templeadmin":
		entities, err = h.Service.GetEntitiesByCreator(user.ID)
		if err != nil || len(entities) == 0 {
			log.Printf("No entities found for templeadmin %d, returning empty list", user.ID)
			entities = []Entity{}
		}
		
	case "standarduser", "monitoringuser":
		if accessContext.AssignedEntityID != nil {
			tenantID := *accessContext.AssignedEntityID
			
			entities, err = h.Service.GetEntitiesByCreator(tenantID)
			
			if err != nil || len(entities) == 0 {
				log.Printf("No entities found for tenant %d, creating mock entity", tenantID)
				mockEntity := Entity{
					ID:          tenantID,
					Name:        "Temple " + strconv.FormatUint(uint64(tenantID), 10),
					Description: "Temple associated with your account",
					Status:      "active",
					CreatedBy:   tenantID,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				entities = []Entity{mockEntity}
				err = nil
			}
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "No entity assigned to this user"})
			return
		}
		
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid user role"})
		return
	}

	if err != nil {
		log.Printf("Error fetching entities: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch temples", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, entities)
}

// GetEntityByID retrieves a specific entity by ID with permission checks
func (h *Handler) GetEntityByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	accessContextVal, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access context"})
		return
	}
	accessContext, ok := accessContextVal.(middleware.AccessContext)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access context"})
		return
	}

	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	user, ok := userVal.(auth.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user object"})
		return
	}
	
	entity, err := h.Service.GetEntityByID(id)
	if err != nil {
		if (user.Role.RoleName == "standarduser" || user.Role.RoleName == "monitoringuser") && 
		   accessContext.AssignedEntityID != nil && 
		   *accessContext.AssignedEntityID == uint(id) {
			
			entity = Entity{
				ID:          uint(id),
				Name:        "Temple " + strconv.Itoa(id),
				Description: "Temple associated with your account",
				Status:      "active",
				CreatedBy:   uint(id),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Temple not found"})
			return
		}
	}
	
	hasAccess := false
	
	switch user.Role.RoleName {
	case "superadmin":
		hasAccess = true
		
	case "templeadmin":
		hasAccess = (accessContext.DirectEntityID != nil && *accessContext.DirectEntityID == uint(id)) || 
			entity.CreatedBy == user.ID
			
	case "standarduser", "monitoringuser":
		if accessContext.AssignedEntityID != nil {
			hasAccess = (*accessContext.AssignedEntityID == uint(id)) || 
				entity.CreatedBy == *accessContext.AssignedEntityID
		}
		
	default:
		hasAccess = false
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this entity"})
		return
	}

	c.JSON(http.StatusOK, entity)
}

func (h *Handler) UpdateEntity(c *gin.Context) {
	// Parse entity ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}
	entityIDUint := uint(id)

	// Get user info
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userInfo := userVal.(auth.User)
	userID := userInfo.ID

	// Access context
	acVal, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access context"})
		return
	}
	ac, ok := acVal.(middleware.AccessContext)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access context"})
		return
	}

	log.Printf("🔍 UpdateEntity - UserID:%d TargetEntity:%d UserEntityID:%v RoleID:%d",
		userID, id, userInfo.EntityID, userInfo.RoleID)

	// Must have write access
	if !ac.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient write permissions"})
		return
	}

	// ---- Permission Logic ----
	hasAccess := false

	// 1️⃣ Superadmin check (if your DB defines it by RoleID == 1 or similar)
	//    Replace with your actual superadmin logic.
	if userInfo.RoleID == 1 {
		hasAccess = true
	}

	// 2️⃣ Entity checks for other roles (e.g., temple admins, staff, etc.)
	if !hasAccess {
		// Direct entity from JWT
		if ac.DirectEntityID != nil && *ac.DirectEntityID == entityIDUint {
			hasAccess = true
		}
		// Assigned entity from middleware context
		if !hasAccess && ac.AssignedEntityID != nil && *ac.AssignedEntityID == entityIDUint {
			hasAccess = true
		}
		// User's own EntityID
		if !hasAccess && userInfo.EntityID != nil && *userInfo.EntityID == entityIDUint {
			hasAccess = true
		}
		// Creator ownership
		if !hasAccess {
			entity, err := h.Service.GetEntityByID(id) // returns (Entity, error)
			if err == nil && entity.CreatedBy == userID {
				hasAccess = true
			}
		}
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to update this entity"})
		return
	}

	// ---- Update Operation ----
	var input Entity
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}
	input.ID = entityIDUint
	input.UpdatedAt = time.Now()

	ip := middleware.GetIPFromContext(c)

	if err := h.Service.UpdateEntity(input, userID, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update entity",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Entity updated successfully"})
}

// DeleteEntity handles entity deletion (superadmin only)
func (h *Handler) DeleteEntity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userObj := user.(auth.User)
	userID := userObj.ID

	if userObj.Role.RoleName != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only superadmins can delete temples"})
		return
	}

	ip := middleware.GetIPFromContext(c)

	if err := h.Service.DeleteEntity(id, userID, ip); err != nil {
		log.Printf("Delete Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete temple", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Temple deleted successfully"})
}

// GetDevoteesByEntity retrieves devotees for a specific entity
func (h *Handler) GetDevoteesByEntity(c *gin.Context) {
	entityIDParam := c.Param("id")
	entityIDUint, err := strconv.ParseUint(entityIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	accessContextVal, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access context"})
		return
	}
	accessContext, ok := accessContextVal.(middleware.AccessContext)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access context"})
		return
	}

	entityID := uint(entityIDUint)
	hasAccess := (accessContext.DirectEntityID != nil && *accessContext.DirectEntityID == entityID) ||
		(accessContext.AssignedEntityID != nil && *accessContext.AssignedEntityID == entityID)

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to devotees for this entity"})
		return
	}

	devotees, err := h.Service.GetDevotees(entityID)
	if err != nil {
		log.Printf("Error fetching devotees for entity %d: %v", entityID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch devotees", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, devotees)
}

// GetDevoteeStats retrieves devotee statistics for an entity
func (h *Handler) GetDevoteeStats(c *gin.Context) {
	entityIDStr := c.Param("id")
	entityIDUint, err := strconv.ParseUint(entityIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	accessContextVal, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access context"})
		return
	}
	accessContext, ok := accessContextVal.(middleware.AccessContext)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access context"})
		return
	}

	entityID := uint(entityIDUint)
	hasAccess := (accessContext.DirectEntityID != nil && *accessContext.DirectEntityID == entityID) ||
		(accessContext.AssignedEntityID != nil && *accessContext.AssignedEntityID == entityID)

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to devotee stats for this entity"})
		return
	}

	stats, err := h.Service.GetDevoteeStats(entityID)
	if err != nil {
		log.Printf("Error fetching devotee stats for entity %d: %v", entityID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch devotee stats", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// UpdateDevoteeMembershipStatus updates devotee membership status
func (h *Handler) UpdateDevoteeMembershipStatus(c *gin.Context) {
	entityIDUint, err := strconv.ParseUint(c.Param("entityID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	accessContextVal, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access context"})
		return
	}
	accessContext, ok := accessContextVal.(middleware.AccessContext)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access context"})
		return
	}

	if !accessContext.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient write permissions"})
		return
	}

	entityID := uint(entityIDUint)
	hasAccess := (accessContext.DirectEntityID != nil && *accessContext.DirectEntityID == entityID) ||
		(accessContext.AssignedEntityID != nil && *accessContext.AssignedEntityID == entityID)

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to manage devotees for this entity"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=active inactive"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	err = h.Service.MembershipService.UpdateMembershipStatus(uint(userID), entityID, req.Status)
	if err != nil {
		log.Printf("Error updating membership status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Membership status updated successfully"})
}

// GetDashboardSummary retrieves dashboard summary for the accessible entity
func (h *Handler) GetDashboardSummary(c *gin.Context) {
	accessContextVal, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access context"})
		return
	}
	accessContext, ok := accessContextVal.(middleware.AccessContext)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access context"})
		return
	}

	entityID := accessContext.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No accessible entity found"})
		return
	}

	summary, err := h.Service.GetDashboardSummary(*entityID)
	if err != nil {
		log.Printf("Dashboard Summary Error for entity %d: %v", *entityID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dashboard summary", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// ========== HELPER FUNCTIONS ==========

func sniffOrByExt(ext string) string {
	if mt := mime.TypeByExtension(ext); mt != "" {
		return mt
	}
	return "application/octet-stream"
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}