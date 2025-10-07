package superadmin

import (
	"time"
)

type TenantApprovalCount struct {
	Approved int64 `json:"approved"`
	Pending  int64 `json:"pending"`
	Rejected int64 `json:"rejected"`
}

// ================ TEMPLE APPROVAL COUNTS ================

type TempleApprovalCount struct {
	PendingApproval int64 `json:"pending_approval"`
	ActiveTemples   int64 `json:"active_temples"`
	Rejected        int64 `json:"rejected"`
	TotalDevotees   int64 `json:"total_devotees"`
}

// ================ TENANT WITH TEMPLE DETAILS ================

type TenantWithDetails struct {
	// User details
<<<<<<< HEAD
	ID        uint      `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	RoleID    uint      `json:"role_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

=======
	ID           uint      `json:"id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	RoleID       uint      `json:"role_id"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	// Temple details
	TempleDetails *TenantTempleDetails `json:"temple_details,omitempty"`
}

type TenantTempleDetails struct {
	ID                uint      `json:"id"`
	TempleName        string    `json:"temple_name"`
	TemplePlace       string    `json:"temple_place"`
	TempleAddress     string    `json:"temple_address"`
	TemplePhoneNo     string    `json:"temple_phone_no"`
	TempleDescription string    `json:"temple_description"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// ================ USER MANAGEMENT ================

type CreateUserRequest struct {
	FullName          string `json:"fullName" binding:"required"`
	Email             string `json:"email" binding:"required,email"`
	Password          string `json:"password" binding:"required,min=6"`
	Phone             string `json:"phone" binding:"required"`
	Role              string `json:"role" binding:"required"`
<<<<<<< HEAD

=======
	
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	// Temple admin specific fields (required only for templeadmin role)
	TempleName        string `json:"templeName"`
	TemplePlace       string `json:"templePlace"`
	TempleAddress     string `json:"templeAddress"`
	TemplePhoneNo     string `json:"templePhoneNo"`
	TempleDescription string `json:"templeDescription"`
}

type UpdateUserRequest struct {
	FullName          string `json:"fullName"`
	Email             string `json:"email" binding:"email"`
	Phone             string `json:"phone"`
<<<<<<< HEAD

=======
	
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	// Temple admin specific fields (only for templeadmin role)
	TempleName        string `json:"templeName"`
	TemplePlace       string `json:"templePlace"`
	TempleAddress     string `json:"templeAddress"`
	TemplePhoneNo     string `json:"templePhoneNo"`
	TempleDescription string `json:"templeDescription"`
}

// NEW: Struct for tenant selection (different from assignment)
type TenantSelectionResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
<<<<<<< HEAD
	TempleName   string `json:"templeName"`
=======
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	Location     string `json:"location"`
	Status       string `json:"status"`
	TemplesCount int    `json:"templesCount"`
	ImageUrl     string `json:"imageUrl,omitempty"`
}

type UserResponse struct {
<<<<<<< HEAD
	ID        uint      `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Role      UserRole  `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Assignment details exposed directly for Vue
	TenantAssigned     string     `json:"tenant_assigned,omitempty"` // now stores tenant name
	AssignedDate       *time.Time `json:"assignedDate,omitempty"`
	ReassignmentDate   *time.Time `json:"reassignmentDate,omitempty"`

	// Temple details for templeadmin users
	TempleDetails           *TenantTempleDetails     `json:"temple_details,omitempty"`
	TenantAssignmentDetails *TenantAssignmentDetails `json:"tenant_assignment_details,omitempty"`
=======
    ID        uint      `json:"id"`
    FullName  string    `json:"full_name"`
    Email     string    `json:"email"`
    Phone     string    `json:"phone"`
    Role      UserRole  `json:"role"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`

    // Temple details for templeadmin users
    TempleDetails           *TenantTempleDetails    `json:"temple_details,omitempty"`
    // Add backticks here
    TenantAssignmentDetails *TenantAssignmentDetails `json:"tenant_assignment_details,omitempty"`
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
}

type UserRole struct {
	ID                  uint   `json:"id"`
	RoleName            string `json:"role_name"`
	Description         string `json:"description"`
	CanRegisterPublicly bool   `json:"can_register_publicly"`
}

// AssignableTenant represents a single tenant for the assignment list.
type AssignableTenant struct {
	UserID        uint   `json:"userID"`
	TenantName    string `json:"tenantName"`
	Email         string `json:"email"`
	TempleAddress string `json:"templeAddress"`
	TempleName    string `json:"templeName"`
}

// AssignRequest holds the user IDs and the single tenant ID for assignment.
<<<<<<< HEAD
type AssignRequest struct {
	UserID   uint `json:"userId" binding:"required"`
	TenantID uint `json:"tenantId" binding:"required"`
=======
// Corrected backend struct definition
type AssignRequest struct {
    UserID   uint `json:"userId" binding:"required"`
    TenantID uint `json:"tenantId" binding:"required"`
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
}

// AssignTenantsRequest holds the tenant IDs to be assigned to a user.
type AssignTenantsRequest struct {
	TenantIDs []uint `json:"tenant_ids" binding:"required"`
}

<<<<<<< HEAD
=======

>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
// TenantListResponse holds the details of a single tenant for the list view.
type TenantListResponse struct {
	ID            uint   `json:"id"`
	UserID        uint   `json:"user_id"`
	TempleName    string `json:"temple_name"`
	TempleAddress string `json:"temple_address"`
}

type TenantAssignmentDetails struct {
<<<<<<< HEAD
	TenantName string    `json:"tenant_name"`
	AssignedOn time.Time `json:"assigned_on"`
	UpdatedOn  time.Time `json:"updated_on"`
=======
	TenantName  string    `json:"tenant_name"`
	AssignedOn  time.Time `json:"assigned_on"`
	UpdatedOn   time.Time `json:"updated_on"`
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
}

// ================ TEMPLE DETAILS FOR API RESPONSE ================

<<<<<<< HEAD
=======
// TempleDetails represents the temple details for API response
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
type TempleDetails struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	City  string `json:"city"`
	State string `json:"state"`
}

// TenantResponse represents a tenant with their temple details for API response
type TenantResponse struct {
	ID       uint           `json:"id"`
	FullName string         `json:"fullName"`
	Email    string         `json:"email"`
	Role     string         `json:"role"`
	Status   string         `json:"status"`
	Temple   *TempleDetails `json:"temple,omitempty"`
}

// BulkUserCSV maps each row of the uploaded CSV file.
type BulkUserCSV struct {
	FullName string `csv:"Full Name" json:"full_name"`
	Email    string `csv:"Email" json:"email"`
	Phone    string `csv:"Phone" json:"phone"`
	Password string `csv:"Password" json:"password"`
	Role     string `csv:"Role" json:"role"`
	Status   string `csv:"Status" json:"status"`
}

// BulkUploadResult represents response after processing bulk upload.
type BulkUploadResult struct {
<<<<<<< HEAD
	TotalRows    int      `json:"total_rows"`
	SuccessCount int      `json:"success_count"`
	FailedCount  int      `json:"failed_count"`
	Errors       []string `json:"errors,omitempty"`
}
=======
	
	TotalRows    int `json:"total_rows"`
	SuccessCount int `json:"success_count"`
	FailedCount  int `json:"failed_count"`
	Errors       []string `json:"errors,omitempty"` // ✅ add this
}
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
