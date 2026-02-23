package entity

import (
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
	UploadDir string
	BaseURL   string
	MaxSize   int64
}

func NewHandler(s *Service, uploadDir, baseURL string) *Handler {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("Failed to create upload directory: %v", err)
	}
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "/uploads"
	}
	return &Handler{
		Service:   s,
		UploadDir: uploadDir,
		BaseURL:   baseURL,
		MaxSize:   1000 * 1024 * 1024,
	}
}

// canAccessEntityID is a helper that checks if the access context allows access to a given entity ID.
// For standarduser/monitoringuser, it also checks entity ownership via TenantID at DB level.
func canAccessEntityID(accessContext middleware.AccessContext, entityID uint, h *Handler) bool {
	switch accessContext.RoleName {
	case "superadmin":
		return true

	case "templeadmin":
		return accessContext.DirectEntityID != nil && *accessContext.DirectEntityID == entityID

	case "standarduser", "monitoringuser":
		// Check direct entity match first
		if accessContext.AssignedEntityID != nil && *accessContext.AssignedEntityID == entityID {
			return true
		}
		// Check via TenantID: does this entity belong to the user's tenant?
		if accessContext.TenantID > 0 && h != nil {
			tenantID, err := h.Service.GetTenantIDByEntityID(entityID)
			if err == nil && tenantID == accessContext.TenantID {
				return true
			}
		}
		return false

	case "devotee", "volunteer":
		if accessContext.AssignedEntityID != nil && *accessContext.AssignedEntityID == entityID {
			return true
		}
		if accessContext.DirectEntityID != nil && *accessContext.DirectEntityID == entityID {
			return true
		}
		return false
	}
	return false
}

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

	userVal, exists := c.Get("user")
	if !exists {
		h.cleanupTempFiles(tempFiles)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := userVal.(auth.User)

	userID := user.ID
	userRole := user.Role.RoleName
	userRoleID := user.Role.ID

	switch userRole {
	case "superadmin":
		tenantID, err := h.Service.Repo.GetTenantIDForUser(userID)
		if err != nil || tenantID == 0 {
			h.cleanupTempFiles(tempFiles)
			c.JSON(http.StatusForbidden, gin.H{"error": "User is not assigned to any tenant"})
			return
		}
		input.CreatedBy = tenantID

	case "templeadmin":
		input.CreatedBy = userID

	case "standarduser", "monitoringuser":
		tenantID, err := h.Service.Repo.GetTenantIDForUser(userID)
		if err != nil || tenantID == 0 {
			h.cleanupTempFiles(tempFiles)
			c.JSON(http.StatusForbidden, gin.H{"error": "User is not assigned to any tenant"})
			return
		}
		input.CreatedBy = tenantID

	default:
		h.cleanupTempFiles(tempFiles)
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid user role"})
		return
	}

	if input.Status == "" {
		if userRoleID == 1 {
			input.Status = "approved"
		} else {
			input.Status = "pending"
		}
	}

	ip := middleware.GetIPFromContext(c)

	if err := h.Service.CreateEntity(&input, userID, userRoleID, ip); err != nil {
		log.Printf("CreateEntity Error: %v", err)
		h.cleanupTempFiles(tempFiles)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entity"})
		return
	}

	log.Printf("✅ Entity created: ID=%d Status=%s", input.ID, input.Status)

	finalFileInfos := make(map[string]FileInfo)

	if len(tempFiles) > 0 {
		if err := h.moveFilesToFinalLocation(c, &input, tempFiles, &finalFileInfos); err != nil {
			log.Printf("File move error for entity %d: %v", input.ID, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":     "Temple created but file processing failed",
				"temple_id": input.ID,
			})
			return
		}

		if err := h.updateEntityWithFileInfo(&input); err != nil {
			log.Printf("❌ MEDIA SAVE FAILED for entity %d: %v", input.ID, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Temple created but media could not be saved",
			})
			return
		}

		log.Printf("✅ Files + media saved for entity %d", input.ID)
	}

	message := "Temple registration request submitted successfully"
	if input.Status == "approved" {
		message = "Temple created and approved successfully"
	}

	response := gin.H{
		"message":        message,
		"temple_id":      input.ID,
		"status":         input.Status,
		"auto_approved":  input.Status == "approved",
		"uploaded_files": finalFileInfos,
	}

	if input.Media != "" {
		response["media"] = input.Media
	}

	c.JSON(http.StatusCreated, response)
}

type TempFileInfo struct {
	TempPath     string
	FileType     string
	FileName     string
	OriginalName string
	FileSize     int64
	ContentType  string
	UploadedAt   time.Time
}

func (h *Handler) handleMultipartFormData(c *gin.Context, input *Entity, tempFiles *[]TempFileInfo) error {
	form, err := c.MultipartForm()
	if err != nil {
		return fmt.Errorf("failed to parse multipart form: %v", err)
	}

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
	tempSessionDir := filepath.Join(h.UploadDir, "temp_uploads", uuid.New().String())
	if err := os.MkdirAll(tempSessionDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	log.Printf("Created temp directory: %s", tempSessionDir)

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

	if additionalDocs := form.File["additional_docs"]; len(additionalDocs) > 0 {
		log.Printf("📎 Found %d additional docs (non-indexed format)", len(additionalDocs))
		for idx, file := range additionalDocs {
			info, err := h.uploadFileToTemp(file, tempSessionDir, "additional_docs")
			if err != nil {
				log.Printf("Warning: Failed to upload additional document %d: %v", idx, err)
				continue
			}
			*tempFiles = append(*tempFiles, info)
			log.Printf("✅ Additional doc %d uploaded: %s", idx, file.Filename)
		}
	} else {
		log.Printf("📎 Checking for indexed additional docs format...")
		for i := 0; i < 10; i++ {
			field := fmt.Sprintf("additional_docs_%d", i)
			if add := form.File[field]; len(add) > 0 {
				info, err := h.uploadFileToTemp(add[0], tempSessionDir, "additional_docs")
				if err != nil {
					log.Printf("Warning: Failed to upload additional document %d: %v", i, err)
					continue
				}
				*tempFiles = append(*tempFiles, info)
				log.Printf("✅ Additional doc %d uploaded (indexed): %s", i, add[0].Filename)
			}
		}
	}

	if logoFiles := form.File["temple_logo"]; len(logoFiles) > 0 {
		info, err := h.uploadFileToTemp(logoFiles[0], tempSessionDir, "temple_logo")
		if err != nil {
			log.Printf("Warning: Failed to upload temple logo: %v", err)
		} else {
			*tempFiles = append(*tempFiles, info)
			log.Printf("✅ Temple logo uploaded to temp: %s", info.FileName)
		}
	}

	if videoFiles := form.File["temple_video"]; len(videoFiles) > 0 {
		info, err := h.uploadFileToTemp(videoFiles[0], tempSessionDir, "temple_video")
		if err != nil {
			log.Printf("Warning: Failed to upload temple video: %v", err)
		} else {
			*tempFiles = append(*tempFiles, info)
			log.Printf("✅ Temple video uploaded to temp: %s", info.FileName)
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

func (h *Handler) moveFilesToFinalLocation(c *gin.Context, entity *Entity, tempFiles []TempFileInfo, finalFileInfos *map[string]FileInfo) error {
	entityDir := filepath.Join(h.UploadDir, strconv.FormatUint(uint64(entity.ID), 10))
	if err := os.MkdirAll(entityDir, 0755); err != nil {
		return fmt.Errorf("failed to create entity directory: %v", err)
	}
	tempEntityDir := entityDir
	log.Printf("Created entity directory: %s for entity %d", entityDir, entity.ID)

	*finalFileInfos = make(map[string]FileInfo)
	var additionalFiles []FileInfo

	mediaInfo := MediaInfo{}
	if oldMediaVal, exists := c.Get("old_media"); exists {
		if oldMedia, ok := oldMediaVal.(MediaInfo); ok {
			mediaInfo = oldMedia
			log.Printf("🧩 Loaded old media for merge: %+v", mediaInfo)
		}
	}

	var tenantID uint
	tenantIDFromContext := middleware.GetTenantIDFromAccessContext(c)
	if tenantIDFromContext > 0 {
		tenantID = tenantIDFromContext
		log.Printf("📁 Using tenant ID from access context: %d", tenantID)
	}

	if tenantID == 0 && entity.CreatedBy > 0 {
		tenantID = entity.CreatedBy
		log.Printf("📁 Using tenant ID from entity.CreatedBy: %d", tenantID)
	}

	for _, tf := range tempFiles {
		finalFileName := tf.FileName

		if tf.FileType == "registration_cert" ||
			tf.FileType == "trust_deed" ||
			tf.FileType == "property_docs" ||
			tf.FileType == "additional_docs" {
			entityDir = filepath.Join(h.UploadDir, strconv.Itoa(int(tenantID)), strconv.FormatUint(uint64(entity.ID), 10))
			if err := os.MkdirAll(entityDir, 0755); err != nil {
				return fmt.Errorf("failed to create entity directory: %v", err)
			}
		} else {
			entityDir = tempEntityDir
		}

		finalPath := filepath.Join(entityDir, finalFileName)

		log.Println("finalPath:", finalPath)
		if err := os.Rename(tf.TempPath, finalPath); err != nil {
			if err := copyFile(tf.TempPath, finalPath); err != nil {
				log.Printf("Failed to move/copy file %s to %s: %v", tf.TempPath, finalPath, err)
				return fmt.Errorf("failed to persist file %s: %v", tf.FileName, err)
			}
			_ = os.Remove(tf.TempPath)
		}

		rel := filepath.ToSlash(filepath.Join(strconv.FormatUint(uint64(entity.ID), 10), finalFileName))
		fileURL := h.buildFileURL(rel)
		log.Println("fileURL:", fileURL)

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

		case "temple_logo":
			(*finalFileInfos)["temple_logo"] = fi
			mediaInfo.Logo = fileURL
			log.Printf("✅ Temple logo URL set: %s", fileURL)

		case "temple_video":
			(*finalFileInfos)["temple_video"] = fi
			mediaInfo.Video = fileURL
			log.Printf("✅ Temple video URL set: %s", fileURL)
		}
	}

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

	if mediaJSON, err := json.Marshal(mediaInfo); err == nil {
		entity.Media = string(mediaJSON)
		log.Printf("💾 Final merged media saved: %s", entity.Media)
	} else {
		log.Printf("⚠️ Failed to marshal media info: %v", err)
	}

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
		".pdf":  true,
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".doc":  true,
		".docx": true,
		".mp4":  true,
		".mov":  true,
		".avi":  true,
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

func (h *Handler) buildFileURL(rel string) string {
	rel = strings.TrimLeft(rel, "/")
	return fmt.Sprintf("%s/%s", strings.TrimRight(h.BaseURL, "/"), rel)
}

type EntityDirectory struct {
	EntityID   string   `json:"entity_id"`
	FilesCount int      `json:"files_count"`
	Files      []string `json:"files"`
}

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

type FileDetails struct {
	FileName string `json:"file_name"`
	FileURL  string `json:"file_url"`
	Size     int64  `json:"size"`
}

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

	if !canAccessEntityID(accessCtx, entityIDUint, h) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to files for this entity"})
		return
	}

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
		rel := filepath.ToSlash(filepath.Join(entityID, e.Name()))
		url := h.buildFileURL(rel)
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

// GetAllEntities retrieves entities based on user role and permissions
// FIXED: standarduser/monitoringuser now use TenantID to find their entities
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
		// FIXED: Use TenantID (from DB assignment) to find entities
		// TenantID is the tenant/creator who owns the entities this user can access
		if accessContext.TenantID > 0 {
			entities, err = h.Service.GetEntitiesByCreator(accessContext.TenantID)
			if err != nil || len(entities) == 0 {
				log.Printf("No entities found for tenant %d", accessContext.TenantID)
				entities = []Entity{}
				err = nil
			}
		} else if accessContext.AssignedEntityID != nil {
			// Fallback: if AssignedEntityID is set, get entities by that creator
			entities, err = h.Service.GetEntitiesByCreator(*accessContext.AssignedEntityID)
			if err != nil || len(entities) == 0 {
				entities = []Entity{}
				err = nil
			}
		} else {
			log.Printf("⚠️ standarduser %d has no TenantID or AssignedEntityID", user.ID)
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
// FIXED: standarduser/monitoringuser now check via TenantID as well
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Temple not found"})
		return
	}

	mediaObj := map[string]interface{}{
		"logo":  "",
		"video": "",
	}
	if entity.Media != "" {
		log.Printf("📦 Raw Media from DB: %s", entity.Media)
		if err := json.Unmarshal([]byte(entity.Media), &mediaObj); err != nil {
			log.Printf("⚠️ Failed to parse media JSON for entity %d: %v", entity.ID, err)
		}
	}

	var additionalDocs []string
	if entity.AdditionalDocsURLs != "" {
		if err := json.Unmarshal([]byte(entity.AdditionalDocsURLs), &additionalDocs); err != nil {
			additionalDocs = []string{}
		}
	}

	var creatorDetails *CreatorDetails
	if entity.CreatedBy > 0 {
		cd, err := h.Service.GetCreatorDetails(entity.CreatedBy)
		if err == nil && cd != nil {
			creatorDetails = cd
		}
	}

	// Devotee safe read-only access
	if user.Role.RoleName == "devotee" && strings.HasSuffix(c.FullPath(), "/details") {
		response := gin.H{
			"id":          entity.ID,
			"name":        entity.Name,
			"main_deity":  entity.MainDeity,
			"temple_type": entity.TempleType,
			"city":        entity.City,
			"district":    entity.District,
			"state":       entity.State,
			"map_link":    entity.MapLink,
			"status":      entity.Status,
			"isactive":    entity.IsActive,
			"media":       mediaObj,
		}
		if creatorDetails != nil {
			response["creator"] = creatorDetails
		}
		c.JSON(http.StatusOK, response)
		return
	}

	// FIXED: Use canAccessEntityID helper which checks TenantID for standarduser
	hasAccess := false
	switch user.Role.RoleName {
	case "superadmin":
		hasAccess = true
	case "templeadmin":
		hasAccess = (accessContext.DirectEntityID != nil && *accessContext.DirectEntityID == uint(id)) ||
			entity.CreatedBy == user.ID
	case "standarduser", "monitoringuser":
		hasAccess = canAccessEntityID(accessContext, uint(id), h)
	case "devotee":
		hasAccess = false
	}

	if !hasAccess {
		log.Printf("🔒 Access denied: UserID=%d Role=%s EntityID=%d TenantID=%d EntityCreatedBy=%d",
			user.ID, user.Role.RoleName, id, accessContext.TenantID, entity.CreatedBy)
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this entity"})
		return
	}

	response := gin.H{
		"id":               entity.ID,
		"name":             entity.Name,
		"main_deity":       entity.MainDeity,
		"temple_type":      entity.TempleType,
		"established_year": entity.EstablishedYear,
		"phone":            entity.Phone,
		"email":            entity.Email,
		"description":      entity.Description,
		"street_address":   entity.StreetAddress,
		"city":             entity.City,
		"district":         entity.District,
		"state":            entity.State,
		"pincode":          entity.Pincode,
		"landmark":         entity.Landmark,
		"map_link":         entity.MapLink,
		"status":           entity.Status,
		"isactive":         entity.IsActive,
		"registration_cert_url": entity.RegistrationCertURL,
		"trust_deed_url":        entity.TrustDeedURL,
		"property_docs_url":     entity.PropertyDocsURL,
		"additional_docs_urls":  additionalDocs,
		"media":                 mediaObj,
		"created_by":            entity.CreatedBy,
		"creator_role_id":       entity.CreatorRoleID,
	}

	if creatorDetails != nil {
		response["creator"] = creatorDetails
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateEntity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := userVal.(auth.User)

	accessVal, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access context"})
		return
	}
	accessContext := accessVal.(middleware.AccessContext)

	existingEntity, err := h.Service.GetEntityByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Temple not found"})
		return
	}

	var oldMedia MediaInfo
	if existingEntity.Media != "" {
		_ = json.Unmarshal([]byte(existingEntity.Media), &oldMedia)
	}
	c.Set("old_media", oldMedia)

	// FIXED: Use canAccessEntityID helper
	hasAccess := false
	switch user.Role.RoleName {
	case "superadmin":
		hasAccess = true
	case "templeadmin":
		hasAccess = existingEntity.CreatedBy == user.ID ||
			(accessContext.DirectEntityID != nil && *accessContext.DirectEntityID == uint(id))
	case "standarduser", "monitoringuser":
		hasAccess = canAccessEntityID(accessContext, uint(id), h)
	}

	if !hasAccess || !accessContext.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	contentType := c.GetHeader("Content-Type")
	isMultipart := strings.Contains(contentType, "multipart/form-data")

	var input Entity
	var tempFiles []TempFileInfo

	if isMultipart {
		if err := h.handleMultipartFormData(c, &input, &tempFiles); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
	}

	input.ID = uint(id)
	input.CreatedBy = existingEntity.CreatedBy
	input.CreatedAt = existingEntity.CreatedAt
	input.UpdatedAt = time.Now()

	wasRejected := existingEntity.Status == "rejected"

	if input.RegistrationCertURL == "" {
		input.RegistrationCertURL = existingEntity.RegistrationCertURL
		input.RegistrationCertInfo = existingEntity.RegistrationCertInfo
	}
	if input.TrustDeedURL == "" {
		input.TrustDeedURL = existingEntity.TrustDeedURL
		input.TrustDeedInfo = existingEntity.TrustDeedInfo
	}
	if input.PropertyDocsURL == "" {
		input.PropertyDocsURL = existingEntity.PropertyDocsURL
		input.PropertyDocsInfo = existingEntity.PropertyDocsInfo
	}
	if input.AdditionalDocsURLs == "" {
		input.AdditionalDocsURLs = existingEntity.AdditionalDocsURLs
		input.AdditionalDocsInfo = existingEntity.AdditionalDocsInfo
	}

	if wasRejected && user.Role.RoleName != "superadmin" {
		input.Status = "pending"
	} else if user.Role.RoleName != "superadmin" {
		input.Status = existingEntity.Status
	}

	if !input.IsActive {
		input.IsActive = existingEntity.IsActive
	}

	finalFileInfos := make(map[string]FileInfo)

	if len(tempFiles) > 0 {
		_ = h.deleteOldEntityFiles(&existingEntity, tempFiles)

		if err := h.moveFilesToFinalLocation(c, &input, tempFiles, &finalFileInfos); err != nil {
			h.cleanupTempFiles(tempFiles)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process uploaded files"})
			return
		}
	}

	ip := middleware.GetIPFromContext(c)

	if err := h.Service.UpdateEntity(input, user.ID, user.Role.ID, ip, wasRejected); err != nil {
		h.cleanupTempFiles(tempFiles)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := gin.H{
		"message":   "Temple updated successfully",
		"temple_id": id,
	}

	if len(finalFileInfos) > 0 {
		resp["uploaded_files"] = finalFileInfos
		resp["files_updated"] = len(finalFileInfos)
	}

	if input.Media != "" {
		resp["media"] = input.Media
	}

	if wasRejected && input.Status == "pending" {
		resp["status_changed"] = true
		resp["new_status"] = "pending"
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) deleteOldEntityFiles(entity *Entity, newFiles []TempFileInfo) error {
	entityDir := filepath.Join(h.UploadDir, strconv.FormatUint(uint64(entity.ID), 10))

	fileTypesBeingReplaced := make(map[string]bool)
	for _, tf := range newFiles {
		fileTypesBeingReplaced[tf.FileType] = true
	}

	if fileTypesBeingReplaced["registration_cert"] && entity.RegistrationCertInfo != "" {
		var oldFileInfo FileInfo
		if err := json.Unmarshal([]byte(entity.RegistrationCertInfo), &oldFileInfo); err == nil {
			oldPath := filepath.Join(entityDir, oldFileInfo.FileName)
			if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
				log.Printf("⚠️ Failed to delete old registration cert: %v", err)
			}
		}
	}

	if fileTypesBeingReplaced["trust_deed"] && entity.TrustDeedInfo != "" {
		var oldFileInfo FileInfo
		if err := json.Unmarshal([]byte(entity.TrustDeedInfo), &oldFileInfo); err == nil {
			oldPath := filepath.Join(entityDir, oldFileInfo.FileName)
			if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
				log.Printf("⚠️ Failed to delete old trust deed: %v", err)
			}
		}
	}

	if fileTypesBeingReplaced["property_docs"] && entity.PropertyDocsInfo != "" {
		var oldFileInfo FileInfo
		if err := json.Unmarshal([]byte(entity.PropertyDocsInfo), &oldFileInfo); err == nil {
			oldPath := filepath.Join(entityDir, oldFileInfo.FileName)
			if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
				log.Printf("⚠️ Failed to delete old property docs: %v", err)
			}
		}
	}

	if fileTypesBeingReplaced["additional_docs"] && entity.AdditionalDocsInfo != "" {
		var oldFiles []FileInfo
		if err := json.Unmarshal([]byte(entity.AdditionalDocsInfo), &oldFiles); err == nil {
			for _, oldFile := range oldFiles {
				oldPath := filepath.Join(entityDir, oldFile.FileName)
				if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
					log.Printf("⚠️ Failed to delete old additional doc: %v", err)
				}
			}
		}
	}

	if (fileTypesBeingReplaced["temple_logo"] || fileTypesBeingReplaced["temple_video"]) && entity.Media != "" {
		var oldMediaInfo MediaInfo
		if err := json.Unmarshal([]byte(entity.Media), &oldMediaInfo); err == nil {
			if fileTypesBeingReplaced["temple_logo"] && oldMediaInfo.Logo != "" {
				logoFileName := filepath.Base(oldMediaInfo.Logo)
				oldLogoPath := filepath.Join(entityDir, logoFileName)
				if err := os.Remove(oldLogoPath); err != nil && !os.IsNotExist(err) {
					log.Printf("⚠️ Failed to delete old logo: %v", err)
				}
			}
			if fileTypesBeingReplaced["temple_video"] && oldMediaInfo.Video != "" {
				videoFileName := filepath.Base(oldMediaInfo.Video)
				oldVideoPath := filepath.Join(entityDir, videoFileName)
				if err := os.Remove(oldVideoPath); err != nil && !os.IsNotExist(err) {
					log.Printf("⚠️ Failed to delete old video: %v", err)
				}
			}
		}
	}

	return nil
}

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

func (h *Handler) ToggleEntityStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

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

	existingEntity, err := h.Service.GetEntityByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Temple not found"})
		return
	}

	hasAccess := false
	switch user.Role.RoleName {
	case "superadmin":
		hasAccess = true
	case "templeadmin":
		hasAccess = (existingEntity.CreatedBy == user.ID)
	case "standarduser", "monitoringuser":
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Insufficient permissions to toggle temple status",
			"message": "Only temple creators and administrators can change temple status",
		})
		return
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid user role"})
		return
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to toggle status for this temple"})
		return
	}

	var req struct {
		IsActive bool `json:"isactive"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	ip := middleware.GetIPFromContext(c)

	if err := h.Service.ToggleEntityStatus(id, req.IsActive, user.ID, ip); err != nil {
		log.Printf("Toggle Status Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle temple status", "details": err.Error()})
		return
	}

	statusText := "inactive"
	if req.IsActive {
		statusText = "active"
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Temple status updated successfully",
		"temple_id": id,
		"isactive":  req.IsActive,
		"status":    statusText,
	})
}

// GetDevoteesByEntity retrieves devotees for a specific entity
// FIXED: Now checks TenantID for standarduser/monitoringuser
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

	log.Printf("=== Access Check Debug ===")
	log.Printf("Requested Entity ID: %d", entityID)
	log.Printf("DirectEntityID: %v", accessContext.DirectEntityID)
	log.Printf("AssignedEntityID: %v", accessContext.AssignedEntityID)
	log.Printf("TenantID: %d", accessContext.TenantID)

	// FIXED: Use canAccessEntityID which checks TenantID for standarduser
	hasAccess := canAccessEntityID(accessContext, entityID, h)

	log.Printf("hasAccess result: %v", hasAccess)
	log.Printf("=========================")

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
// FIXED: Now checks TenantID for standarduser/monitoringuser
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

	// FIXED: Use canAccessEntityID which checks TenantID for standarduser
	if !canAccessEntityID(accessContext, entityID, h) {
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

func (h *Handler) UpdateDevoteeMembershipStatus(c *gin.Context) {
	log.Println("PARAMS DEBUG:", c.Params)

	entityIDUint, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	userIDUint, err := strconv.ParseUint(c.Param("userID"), 10, 64)
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
	userID := uint(userIDUint)

	// FIXED: Use canAccessEntityID helper
	if !canAccessEntityID(accessContext, entityID, h) {
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

	if err := h.Service.MembershipService.UpdateMembershipStatus(userID, entityID, req.Status); err != nil {
		log.Printf("Error updating membership status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Membership status updated successfully"})
}

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

func (h *Handler) GetTenantByEntityID(c *gin.Context) {
	entityIDStr := c.Param("id")
	entityID, err := strconv.ParseUint(entityIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	tenantID, err := h.Service.GetTenantIDByEntityID(uint(entityID))
	if err != nil {
		log.Printf("Error fetching tenant ID for entity %d: %v", entityID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tenant information"})
		return
	}

	if tenantID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No tenant found for this entity"})
		return
	}

	creatorDetails, err := h.Service.GetCreatorDetails(tenantID)
	if err != nil {
		log.Printf("Error fetching creator details for tenant %d: %v", tenantID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch creator details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"entity_id":      entityID,
		"tenant_id":      tenantID,
		"tenant_details": creatorDetails,
	})
}

func GetTenantByEntityID(entityID string) uint {
	if repo == nil {
		log.Printf("⚠️ Repository not initialized for GetTenantByEntityID")
		return 0
	}

	entityIDUint, err := strconv.ParseUint(entityID, 10, 64)
	if err != nil {
		log.Printf("⚠️ Invalid entity ID: %s", entityID)
		return 0
	}

	tenantID, err := repo.GetTenantIDByEntityID(uint(entityIDUint))
	if err != nil {
		log.Printf("⚠️ Error fetching tenant ID for entity %s: %v", entityID, err)
		return 0
	}

	return tenantID
}

func SetRepository(r *Repository) {
	repo = r
}

var repo *Repository