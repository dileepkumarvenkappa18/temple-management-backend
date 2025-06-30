package entity

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
}

//
// ========== ENTITY CORE ==========
//

func (h *Handler) CreateEntity(c *gin.Context) {
	var e Entity
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}
	if err := h.Service.CreateEntity(&e); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Entity created successfully",
		"entity":  e,
	})
}

func (h *Handler) GetAllEntities(c *gin.Context) {
	entities, err := h.Service.GetAllEntities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch entities"})
		return
	}
	c.JSON(http.StatusOK, entities)
}

func (h *Handler) GetEntityByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}
	entity, err := h.Service.GetEntityByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Entity not found"})
		return
	}
	c.JSON(http.StatusOK, entity)
}

func (h *Handler) UpdateEntity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	var e Entity
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}
	e.ID = uint(id)

	if err := h.Service.UpdateEntity(e); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update entity: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Entity updated successfully"})
}

func (h *Handler) ToggleStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}
	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}
	if err := h.Service.ToggleEntityStatus(id, req.IsActive); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Entity status updated"})
}

func (h *Handler) DeleteEntity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}
	if err := h.Service.DeleteEntity(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete entity: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Entity deleted successfully"})
}

//
// ========== ADDRESS ==========
//

func (h *Handler) AddEntityAddress(c *gin.Context) {
	var addr EntityAddress
	if err := c.ShouldBindJSON(&addr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}
	if err := h.Service.AddEntityAddress(addr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add address: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Address added successfully"})
}

func (h *Handler) GetEntityAddress(c *gin.Context) {
	entityID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}
	addr, err := h.Service.GetEntityAddress(entityID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}
	c.JSON(http.StatusOK, addr)
}

//
// ========== DOCUMENTS ==========
//

func (h *Handler) AddEntityDocument(c *gin.Context) {
	var doc EntityDocument
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}
	if err := h.Service.AddEntityDocument(doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload document: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Document uploaded successfully"})
}

func (h *Handler) GetEntityDocuments(c *gin.Context) {
	entityID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}
	docs, err := h.Service.GetEntityDocuments(entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch documents"})
		return
	}
	c.JSON(http.StatusOK, docs)
}

//
// ========== FILE UPLOAD ==========
//

func (h *Handler) UploadEntityDocument(c *gin.Context) {
	entityIDStr := c.PostForm("entity_id")
	docType := c.PostForm("document_type")

	entityID, err := strconv.Atoi(entityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity_id"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed: " + err.Error()})
		return
	}

	safeFileName := filepath.Base(file.Filename)
	filePath := fmt.Sprintf("uploads/entity_%d_%s_%s", entityID, docType, safeFileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file: " + err.Error()})
		return
	}

	doc := EntityDocument{
		EntityID:     uint(entityID),
		DocumentType: docType,
		DocumentURL:  filePath,
	}

	if err := h.Service.AddEntityDocument(doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record document: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded", "path": filePath})
}
