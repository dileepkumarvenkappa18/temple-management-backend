package donation

import (
    "net/http"
    "strconv"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/sharath018/temple-management-backend/internal/auth"
)

type Handler struct {
    svc Service
}

func NewHandler(svc Service) *Handler {
    return &Handler{svc: svc}
}

// 1. Create Razorpay Order & Log Donation Intent
func (h *Handler) CreateDonation(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    currentUser, ok := user.(auth.User)
    if !ok || currentUser.EntityID == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing entity"})
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

// 2. Verify Razorpay Signature
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

// 3. Get My Donations
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

// 4. Get Donations for Admin - with Filters
func (h *Handler) GetDonationsByEntity(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    currentUser, ok := user.(auth.User)
    if !ok || currentUser.EntityID == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing entity"})
        return
    }

    pageStr := c.DefaultQuery("page", "1")
    limitStr := c.DefaultQuery("limit", "20")
    status := c.Query("status")
    from := c.Query("from")
    to := c.Query("to")
    dType := c.Query("type")
    method := c.Query("method")
    minAmount := c.Query("min")
    maxAmount := c.Query("max")
    search := c.Query("search")

    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1
    }
    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit < 1 || limit > 100 {
        limit = 20
    }

    donations, total, err := h.svc.GetDonationsByEntityWithFilters(
        *currentUser.EntityID,
        page, limit, status, from, to, dType, method, minAmount, maxAmount, search,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data":       donations,
        "total":      total,
        "page":       page,
        "limit":      limit,
        "page_count": (total + int64(limit) - 1) / int64(limit),
    })
}

// 5. Top 5 Donors
type TopDonor struct {
    Name   string  `json:"name"`
    Amount float64 `json:"amount"`
}
func (h *Handler) GetTopDonors(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    currentUser, ok := user.(auth.User)
    if !ok || currentUser.EntityID == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing entity"})
        return
    }

    topDonors, err := h.svc.GetTopDonors(*currentUser.EntityID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"top_donors": topDonors})
}


func (h *Handler) GetAllDonors(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    currentUser, ok := user.(auth.User)
    if !ok || currentUser.EntityID == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing entity"})
        return
    }

    donors, err := h.svc.GetAllDonors(*currentUser.EntityID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"donors": donors})
}

// 7. Dashboard Summary (30 days)
func (h *Handler) GetDonationSummary(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    currentUser, ok := user.(auth.User)
    if !ok || currentUser.EntityID == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing entity"})
        return
    }

    entityID := *currentUser.EntityID
    since := time.Now().AddDate(0, 0, -30)

    summary, err := h.svc.GetDonationSummary(entityID, since)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, summary)
}

// 8. Full Dashboard Overview
func (h *Handler) GetDonationDashboard(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    currentUser, ok := user.(auth.User)
    if !ok || currentUser.EntityID == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing entity"})
        return
    }

    dashboard, err := h.svc.GetDonationDashboard(*currentUser.EntityID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, dashboard)
}



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