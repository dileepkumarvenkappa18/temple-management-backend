package superadmin

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/internal/entity"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ================== TENANT ==================

func (s *Service) ApproveTenant(ctx context.Context, userID uint, adminID uint) error {
	// Check existence and current status
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return errors.New("tenant not found")
	}

	if user.Status == "active" {
		return errors.New("tenant already approved")
	}
	if user.Status == "rejected" {
		return errors.New("tenant already rejected")
	}

	if err := s.repo.ApproveTenant(ctx, userID, adminID); err != nil {
		return err
	}
	return s.repo.MarkTenantApprovalApproved(ctx, userID, adminID)
}

func (s *Service) RejectTenant(ctx context.Context, userID uint, adminID uint, reason string) error {
	if reason == "" {
		return errors.New("rejection reason is required")
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return errors.New("tenant not found")
	}

	if user.Status == "rejected" {
		return errors.New("tenant already rejected")
	}
	if user.Status == "active" {
		return errors.New("tenant already approved")
	}

	return s.repo.RejectTenant(ctx, userID, adminID, reason)
}

func (s *Service) GetPendingTenants(ctx context.Context) ([]auth.User, error) {
	return s.repo.GetPendingTenants(ctx)
}

func (s *Service) GetTenantsWithFilters(ctx context.Context, status string, limit, page int) ([]TenantWithDetails, int64, error) {
	return s.repo.GetTenantsWithFilters(ctx, status, limit, page)
}

func (s *Service) UpdateTenantApprovalStatus(ctx context.Context, userID, adminID uint, action string, reason string) error {
	switch action {
	case "approve":
		return s.ApproveTenant(ctx, userID, adminID)
	case "reject":
		return s.RejectTenant(ctx, userID, adminID, reason)
	default:
		return errors.New("invalid action: must be approve or reject")
	}
}

// ================== ENTITY ==================

func (s *Service) GetPendingEntities(ctx context.Context) ([]entity.Entity, error) {
	return s.repo.GetPendingEntities(ctx)
}

func (s *Service) GetEntitiesWithFilters(ctx context.Context, status string, limit, page int) ([]entity.Entity, int64, error) {
	return s.repo.GetEntitiesWithFilters(ctx, status, limit, page)
}

func (s *Service) ApproveEntity(ctx context.Context, entityID uint, adminID uint) error {
	ent, err := s.repo.GetEntityByID(ctx, entityID)
	if err != nil {
		return errors.New("entity not found")
	}

	if ent.Status == "approved" {
		return errors.New("entity already approved")
	}
	if ent.Status == "rejected" {
		return errors.New("entity already rejected")
	}

	if err := s.repo.ApproveEntity(ctx, entityID, adminID); err != nil {
		return err
	}

	return s.repo.LinkEntityToUser(ctx, ent.CreatedBy, ent.ID)
}

func (s *Service) RejectEntity(ctx context.Context, entityID uint, adminID uint, reason string) error {
	if reason == "" {
		return errors.New("rejection reason is required")
	}

	ent, err := s.repo.GetEntityByID(ctx, entityID)
	if err != nil {
		return errors.New("entity not found")
	}

	if ent.Status == "rejected" {
		return errors.New("entity already rejected")
	}
	if ent.Status == "approved" {
		return errors.New("entity already approved")
	}

	rejectedAt := time.Now()
	return s.repo.RejectEntity(ctx, entityID, adminID, reason, rejectedAt)
}

func (s *Service) UpdateEntityApprovalStatus(ctx context.Context, entityID, adminID uint, action string, reason string) error {
	switch action {
	case "approve":
		return s.ApproveEntity(ctx, entityID, adminID)
	case "reject":
		return s.RejectEntity(ctx, entityID, adminID, reason)
	default:
		return errors.New("invalid action: must be approve or reject")
	}
}

// ================== METRIC COUNTS ==================

// Tenant approval counts for SuperAdmin dashboard
func (s *Service) GetTenantApprovalCounts(ctx context.Context) (*TenantApprovalCount, error) {
	approved, err := s.repo.CountTenantsByStatus(ctx, "active") // assuming "active" means approved
	if err != nil {
		return nil, err
	}

	pending, err := s.repo.CountTenantsByStatus(ctx, "pending")
	if err != nil {
		return nil, err
	}

	rejected, err := s.repo.CountTenantsByStatus(ctx, "rejected")
	if err != nil {
		return nil, err
	}

	return &TenantApprovalCount{
		Approved: approved,
		Pending:  pending,
		Rejected: rejected,
	}, nil
}

// Temple (entity) approval counts for dashboard
func (s *Service) GetTempleApprovalCounts(ctx context.Context) (*TempleApprovalCount, error) {
	pending, err := s.repo.CountEntitiesByStatus(ctx, "PENDING")
	if err != nil {
		return nil, err
	}

	active, err := s.repo.CountEntitiesByStatus(ctx, "APPROVED")
	if err != nil {
		return nil, err
	}

	rejected, err := s.repo.CountEntitiesByStatus(ctx, "REJECTED")
	if err != nil {
		return nil, err
	}

	totalDevotees, err := s.repo.CountUsersByRole(ctx, "devotee")
	if err != nil {
		return nil, err
	}

	return &TempleApprovalCount{
		PendingApproval: pending,
		ActiveTemples:   active,
		Rejected:        rejected,
		TotalDevotees:   totalDevotees,
	}, nil
}

// ================== USER MANAGEMENT ==================

// Create user (admin-created users)
func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest, adminID uint) error {
	// Validate role exists
	role, err := s.repo.FindRoleByName(ctx, strings.ToLower(req.Role))
	if err != nil {
		return errors.New("invalid role")
	}

	// Check if email already exists
	exists, err := s.repo.UserExistsByEmail(ctx, req.Email)
	if err != nil {
		return errors.New("failed to check email availability")
	}
	if exists {
		return errors.New("email already exists")
	}

	// Clean and validate phone
	phone, err := cleanPhone(req.Phone)
	if err != nil {
		return err
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Create user - admin-created users are active by default
	user := &auth.User{
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: string(hash),
		Phone:        phone,
		RoleID:       role.ID,
		Status:       "active", // Admin-created users are active immediately
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return errors.New("failed to create user")
	}

	// If templeadmin role, create tenant details
	if strings.ToLower(req.Role) == "templeadmin" {
		tenantDetails := &auth.TenantDetails{
			UserID:            user.ID,
			TempleName:        req.TempleName,
			TemplePlace:       req.TemplePlace,
			TempleAddress:     req.TempleAddress,
			TemplePhoneNo:     req.TemplePhoneNo,
			TempleDescription: req.TempleDescription,
		}

		if err := s.repo.CreateTenantDetails(ctx, tenantDetails); err != nil {
			return errors.New("failed to create temple details")
		}
	}

	return nil
}

// Get users with pagination and filters
func (s *Service) GetUsers(ctx context.Context, limit, page int, search, roleFilter, statusFilter string) ([]UserResponse, int64, error) {
	return s.repo.GetUsers(ctx, limit, page, search, roleFilter, statusFilter)
}

// Get user by ID
func (s *Service) GetUserByID(ctx context.Context, userID uint) (*UserResponse, error) {
	return s.repo.GetUserWithDetails(ctx, userID)
}

// Update user - UPDATED: SuperAdmin restrictions removed
func (s *Service) UpdateUser(ctx context.Context, userID uint, req UpdateUserRequest, adminID uint) error {
	// Get existing user to check if it exists
	existingUser, err := s.repo.GetUserWithDetails(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	// ✅ REMOVED: SuperAdmin restriction
	// Now SuperAdmin users can be updated

	// Check if email is being changed and if new email already exists
	if req.Email != "" && req.Email != existingUser.Email {
		exists, err := s.repo.UserExistsByEmail(ctx, req.Email)
		if err != nil {
			return errors.New("failed to check email availability")
		}
		if exists {
			return errors.New("email already exists")
		}
	}

	// Prepare user updates
	userUpdates := &auth.User{}
	if req.FullName != "" {
		userUpdates.FullName = req.FullName
	}
	if req.Email != "" {
		userUpdates.Email = req.Email
	}
	if req.Phone != "" {
		phone, err := cleanPhone(req.Phone)
		if err != nil {
			return err
		}
		userUpdates.Phone = phone
	}

	// Update user
	if err := s.repo.UpdateUser(ctx, userID, userUpdates); err != nil {
		return errors.New("failed to update user")
	}

	// Update temple details if user is templeadmin and temple details provided
	if existingUser.Role.RoleName == "templeadmin" && 
		(req.TempleName != "" || req.TemplePlace != "" || req.TempleAddress != "" || 
		 req.TemplePhoneNo != "" || req.TempleDescription != "") {
		
		tenantDetails := &auth.TenantDetails{}
		if req.TempleName != "" {
			tenantDetails.TempleName = req.TempleName
		}
		if req.TemplePlace != "" {
			tenantDetails.TemplePlace = req.TemplePlace
		}
		if req.TempleAddress != "" {
			tenantDetails.TempleAddress = req.TempleAddress
		}
		if req.TemplePhoneNo != "" {
			tenantDetails.TemplePhoneNo = req.TemplePhoneNo
		}
		if req.TempleDescription != "" {
			tenantDetails.TempleDescription = req.TempleDescription
		}

		if err := s.repo.UpdateTenantDetails(ctx, userID, tenantDetails); err != nil {
			return errors.New("failed to update temple details")
		}
	}

	return nil
}

// Delete user - KEPT: Still prevent SuperAdmin deletion for safety
func (s *Service) DeleteUser(ctx context.Context, userID uint, adminID uint) error {
	// Get existing user to check role
	existingUser, err := s.repo.GetUserWithDetails(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Keep this restriction for safety - prevent deleting superadmin users
	if existingUser.Role.RoleName == "superadmin" {
		return errors.New("cannot delete superadmin user")
	}

	// Prevent self-deletion
	if userID == adminID {
		return errors.New("cannot delete your own account")
	}

	return s.repo.DeleteUser(ctx, userID)
}

// Update user status - UPDATED: SuperAdmin restriction removed
// Update user status - UPDATED: Simplified after removing SuperAdmin restrictions
func (s *Service) UpdateUserStatus(ctx context.Context, userID uint, status string, adminID uint) error {
	// Only keep the self-deactivation check - no need to fetch user details
	if userID == adminID && status == "inactive" {
		return errors.New("cannot deactivate your own account")
	}

	return s.repo.UpdateUserStatus(ctx, userID, status)
}

// Get all user roles
func (s *Service) GetUserRoles(ctx context.Context) ([]UserRole, error) {
	return s.repo.GetUserRoles(ctx)
}

// ================== HELPERS ==================

func cleanPhone(raw string) (string, error) {
	re := regexp.MustCompile(`\D`)
	cleaned := re.ReplaceAllString(raw, "")

	// Strip leading 91 if present and length is 12
	if len(cleaned) == 12 && strings.HasPrefix(cleaned, "91") {
		cleaned = cleaned[2:]
	}

	if len(cleaned) != 10 {
		return "", errors.New("invalid phone number format")
	}

	return cleaned, nil
}

// ================== USER ROLES ==================

// CreateRole handles the creation of a new user role.
func (s *Service) CreateRole(ctx context.Context, req *auth.CreateRoleRequest) error {
	// 1. Basic validation from the DTO
	if req.RoleName == "" || req.Description == "" {
		return errors.New("role name and description are required")
	}

	// 2. Check for uniqueness using the new repository method
	exists, err := s.repo.CheckIfRoleNameExists(ctx, req.RoleName)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("role name already exists")
	}

	// 3. Create the UserRole model instance
	newRole := &auth.UserRole{
		RoleName:            req.RoleName,
		Description:         req.Description,
		CanRegisterPublicly: false, // Defaulting to false as per UI analysis
		Status:              "active",
	}

	// 4. Save to the database via the repository
	return s.repo.CreateUserRole(ctx, newRole)
}

// GetRoles fetches all active roles for the UI.
func (s *Service) GetRoles(ctx context.Context) ([]auth.RoleResponse, error) {
	// 1. Fetch all active roles from the repository
	roles, err := s.repo.GetAllUserRoles(ctx)
	if err != nil {
		return nil, err
	}

	// 2. Convert the database models to the response DTOs
	var roleResponses []auth.RoleResponse
	for _, role := range roles {
		roleResponses = append(roleResponses, auth.RoleResponse{
			ID:                  role.ID,
			RoleName:            role.RoleName,
			Description:         role.Description,
			Status:              role.Status,
			CanRegisterPublicly: role.CanRegisterPublicly,
		})
	}

	return roleResponses, nil
}

// UpdateRole updates an existing user role's details.
// 🎯 FIX: Changed 'req' type from *auth.CreateRoleRequest to *auth.UpdateRoleRequest
func (s *Service) UpdateRole(ctx context.Context, roleID uint, req *auth.UpdateRoleRequest) error {
	role, err := s.repo.GetUserRoleByID(ctx, roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Update only if provided
	if req.RoleName != "" && req.RoleName != role.RoleName {
		exists, err := s.repo.CheckIfRoleNameExists(ctx, req.RoleName)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("role name already exists")
		}
		role.RoleName = req.RoleName
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.CanRegisterPublicly != nil {
		role.CanRegisterPublicly = *req.CanRegisterPublicly
	}

	role.UpdatedAt = time.Now()

	return s.repo.UpdateUserRole(ctx, role)
}

// 🎯 NEW: ToggleRoleStatus specifically handles updating only the status.
func (s *Service) ToggleRoleStatus(ctx context.Context, roleID uint, status string) error {
	role, err := s.repo.GetUserRoleByID(ctx, roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Check if the status is a valid value
	if status != "active" && status != "inactive" {
		return errors.New("invalid status provided")
	}

	role.Status = status
	role.UpdatedAt = time.Now()

	return s.repo.UpdateUserRole(ctx, role)
}
