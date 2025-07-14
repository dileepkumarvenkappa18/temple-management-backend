package donation

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auth"
)

// Handler represents the donation HTTP handler
type Handler struct {
	svc Service
}

// NewHandler creates a new donation handler
func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// ==============================
// 🌟 1. Create Razorpay Order & Log Donation Intent
// ==============================
// @Summary Create a donation order
// @Tags Donations
// @Accept json
// @Produce json
// @Param body body CreateDonationRequest true "Donation request"
// @Success 200 {object} RazorpayOrderResponse
// @Failure 400 {object} gin.H
// @Router /donations [post]
func (h *Handler) CreateDonation(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	currentUser, ok := user.(auth.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
		return
	}

	if currentUser.EntityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to any entity"})
		return
	}

	var req CreateDonationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UserID = currentUser.ID
	req.EntityID = *currentUser.EntityID

	order, err := h.svc.StartDonation(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// ==============================
// ✅ 2. Verify Razorpay Signature (Frontend calls this after payment success)
// ==============================
// @Summary Verify Razorpay donation payment
// @Tags Donations
// @Accept json
// @Produce json
// @Param body body VerifyPaymentRequest true "Payment verification"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /donations/verify [post]
func (h *Handler) VerifyDonation(c *gin.Context) {
	var req VerifyPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.VerifyAndUpdateDonation(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Donation verified successfully"})
}

// ==============================
// 🔍 3. Get My Donations (Devotee View)
// ==============================
// @Summary Get logged-in user's donations
// @Tags Donations
// @Produce json
// @Success 200 {array} Donation
// @Router /donations/me [get]
func (h *Handler) GetMyDonations(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	currentUser, ok := user.(auth.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
		return
	}

	donations, err := h.svc.GetDonationsByUser(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, donations)
}

// ==============================
// 🔍 4. Get All Donations for Temple (Temple Admin View)
// ==============================
// @Summary Get donations by temple entity ID
// @Tags Donations
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param status query string false "Payment status filter"
// @Success 200 {array} Donation
// @Failure 400 {object} gin.H
// @Router /donations/entity [get]
func (h *Handler) GetDonationsByEntity(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	currentUser, ok := user.(auth.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
		return
	}
	if currentUser.EntityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to any entity"})
		return
	}

	// ✅ Pagination
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")
	status := c.Query("status") // optional

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	donations, err := h.svc.GetDonationsByEntity(*currentUser.EntityID, page, limit, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, donations)
}

// ==============================
// 🧾 5. (Optional) Get Donation Receipt by ID (Future Enhancement)
// ==============================
// func (h *Handler) GetDonationReceipt(c *gin.Context) {
//     // TODO: Implement receipt endpoint that returns donation receipt
//     // as a signed URL or PDF metadata.
// }






// package donation

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/sharath018/temple-management-backend/internal/auth"
// )

// // Handler represents the donation HTTP handler
// type Handler struct {
// 	svc Service
// }

// // NewHandler creates a new donation handler
// func NewHandler(svc Service) *Handler {
// 	return &Handler{svc: svc}
// }

// // ==============================
// // 🌟 1. Create Razorpay Order & Log Donation
// // ==============================
// func (h *Handler) CreateDonation(c *gin.Context) {
// 	user, exists := c.Get("user")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}
// 	currentUser := user.(auth.User)

// 	var req CreateDonationRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	req.UserID = currentUser.ID

// 	// ✅ REMOVE this in test mode. Let Razorpay choose method (card, netbanking, etc.)
// 	// req.Method = "UPI"

// 	if currentUser.EntityID == nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to any entity"})
// 		return
// 	}
// 	req.EntityID = *currentUser.EntityID

// 	order, err := h.svc.StartDonation(req)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, order)
// }


// // ==============================
// // ✅ 2. Verify Payment Signature (Client-side call after success)
// // ==============================
// func (h *Handler) VerifyDonation(c *gin.Context) {
// 	var req VerifyPaymentRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if err := h.svc.VerifyAndUpdateDonation(req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Donation verified successfully"})
// }

// // ==============================
// // 🔍 3. Get My Donations
// // ==============================
// func (h *Handler) GetMyDonations(c *gin.Context) {
// 	user, exists := c.Get("user")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}
// 	currentUser := user.(auth.User)

// 	donations, err := h.svc.GetDonationsByUser(currentUser.ID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, donations)
// }

// // ==============================
// // 🔍 4. Get Donations for Temple Admin by entity_id
// // ==============================
// func (h *Handler) GetDonationsByEntity(c *gin.Context) {
// 	user, exists := c.Get("user")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}
// 	currentUser := user.(auth.User)

// 	if currentUser.EntityID == nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to any entity"})
// 		return
// 	}

// 	donations, err := h.svc.GetDonationsByEntity(*currentUser.EntityID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, donations)
// }