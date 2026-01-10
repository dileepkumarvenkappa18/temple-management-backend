package tenant

import (
    "errors"
    "time"
    "golang.org/x/crypto/bcrypt"
    "log"
)

// Service provides tenant user management functionality
type Service struct {
    repo *Repository
}

// NewService creates a new service instance
func NewService(repo *Repository) *Service {
    return &Service{repo: repo}
}

// =========================
// GET TENANT PROFILE
// =========================
func (s *Service) GetTenantProfile(userID uint) (*TenantProfileResponse, error) {
	log.Printf("🔍 SERVICE: Fetching tenant profile for user ID: %d", userID)

	profile, err := s.repo.GetTenantProfileByUserID(userID)
	if err != nil {
		log.Printf("❌ SERVICE: Failed to fetch tenant profile: %v", err)
		return nil, err
	}

	log.Printf("✅ SERVICE: Tenant profile fetched for user %d", userID)
	return profile, nil
}

// =========================
// UPDATE TENANT PROFILE
// =========================
func (s *Service) UpdateTenantProfile(userID uint, input UpdateTenantProfileRequest) (*TenantProfileResponse, error) {
	log.Printf("🔄 SERVICE: Updating tenant profile for user ID: %d", userID)

	// Fetch current profile (to get tenantID)
	currentProfile, err := s.repo.GetTenantProfileByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Update tenant + user info
	if err := s.repo.UpdateTenantProfile(currentProfile.TenantID, userID, input); err != nil {
		return nil, err
	}

	// Return updated profile
	updatedProfile, err := s.repo.GetTenantProfileByUserID(userID)
	if err != nil {
		return nil, err
	}

	log.Printf("✅ SERVICE: Tenant profile updated for user %d", userID)
	return updatedProfile, nil
}

// GetTenantUsers fetches users assigned to a tenant
func (s *Service) GetTenantUsers(tenantID uint, role string) ([]UserResponse, error) {
    log.Printf("SERVICE: Getting users for tenant ID %d", tenantID)
    users, err := s.repo.GetTenantUsers(tenantID, role)
    if err != nil {
        log.Printf("Service: Error getting users: %v", err)
        return nil, err
    }
    
    if users == nil {
        log.Printf("Service: No users found, returning empty array")
        return []UserResponse{}, nil
    }
    
    for i := range users {
        if users[i].Role == "" {
            users[i].Role = "StandardUser"
        } else {
            switch users[i].Role {
            case "monitoringuser":
                users[i].Role = "MonitoringUser"
            case "standarduser":
                users[i].Role = "StandardUser"
            }
        }
    }
    
    log.Printf("Service: Returning %d users", len(users))
    return users, nil
}

// UpdateUser updates a user's details and/or status
func (s *Service) UpdateUser(tenantID, userID uint, input UserInput, status string) (*UserResponse, error) {
    log.Printf("🔵 SERVICE: Updating user %d for tenant %d", userID, tenantID)
    
    exists, err := s.repo.CheckUserBelongsToTenant(userID, tenantID)
    if err != nil {
        return nil, err
    }
    if !exists {
        return nil, errors.New("user does not belong to this tenant")
    }
    
    currentUser, err := s.repo.GetUserByID(userID)
    if err != nil {
        return nil, err
    }
    
    err = s.repo.UpdateUserDetails(userID, input)
    if err != nil {
        return nil, err
    }
    
    userStatus := status
    if userStatus == "" {
        userStatus = currentUser.Status
    }
    
    if status != "" {
        err = s.repo.UpdateUserStatus(userID, tenantID, status)
        if err != nil {
            return nil, err
        }
    }
    
    user, err := s.repo.GetUserByID(userID)
    if err != nil {
        return nil, err
    }
    
    roleName := input.Role
    if roleName == "" {
        switch user.RoleID {
        case 5:
            roleName = "StandardUser"
        case 6:
            roleName = "MonitoringUser"
        default:
            roleName = "StandardUser"
        }
    }
    
    response := &UserResponse{
        ID:        userID,
        Name:      user.FullName,
        Email:     user.Email,
        Phone:     user.Phone,
        Status:    userStatus,
        CreatedAt: user.CreatedAt,
        Role:      roleName,
    }
    
    return response, nil
}

// CreateOrUpdateTenantUser creates a new user or updates an existing user's tenant assignment
func (s *Service) CreateOrUpdateTenantUser(tenantID uint, input UserInput, creatorID uint) (*UserResponse, error) {
    log.Printf("🔴 SERVICE: Creating/updating user for tenant %d: %s (%s) by creator %d", 
               tenantID, input.Name, input.Email, creatorID)
    
    existingUser, err := s.repo.GetUserByEmail(input.Email)
    if err != nil {
        log.Printf("Error checking for existing user: %v", err)
        return nil, err
    }
    
    var userID uint
    
    if existingUser != nil {
        log.Printf("User already exists (ID: %d), will update assignment", existingUser.ID)
        userID = existingUser.ID
    } else {
        log.Printf("User does not exist, creating new user")
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
        if err != nil {
            log.Printf("Failed to hash password: %v", err)
            return nil, errors.New("failed to hash password")
        }
        
        roleID, err := s.repo.GetRoleIDByName(input.Role)
        if err != nil {
            log.Printf("Invalid role '%s': %v", input.Role, err)
            roleID = 5
        }
        
        newUser := User{
            FullName:     input.Name,
            Email:        input.Email,
            Phone:        input.Phone,
            PasswordHash: string(hashedPassword),
            RoleID:       roleID,
            Status:       "active",
            CreatedAt:    time.Now(),
            UpdatedAt:    time.Now(),
            CreatedBy:    "system",
        }
        
        if err := s.repo.CreateUser(&newUser); err != nil {
            log.Printf("Failed to create user: %v", err)
            return nil, errors.New("failed to create user: " + err.Error())
        }
        
        userID = newUser.ID
        log.Printf("New user created with ID: %d", userID)
    }
    
    log.Printf("🔴 SERVICE: Passing tenant ID %d and creator ID %d to repository", tenantID, creatorID)
    err = s.repo.UpdateTenantUserAssignment(userID, tenantID, creatorID)
    if err != nil {
        log.Printf("Failed to assign user to tenant: %v", err)
        return nil, errors.New("failed to assign user to tenant: " + err.Error())
    }
    
    log.Printf("User successfully assigned to tenant %d by creator %d", tenantID, creatorID)
    
    response := &UserResponse{
        ID:        userID,
        Name:      input.Name,
        Email:     input.Email,
        Phone:     input.Phone,
        Status:    "active",
        CreatedAt: time.Now(),
        Role:      input.Role,
    }
    
    return response, nil
}