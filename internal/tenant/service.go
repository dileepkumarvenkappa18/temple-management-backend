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

// GetTenantUsers fetches users assigned to a tenant
func (s *Service) GetTenantUsers(tenantID uint, role string) ([]UserResponse, error) {
    log.Printf("SERVICE: Getting users for tenant ID %d", tenantID)
    users, err := s.repo.GetTenantUsers(tenantID, role)
    if err != nil {
        log.Printf("Service: Error getting users: %v", err)
        return nil, err
    }
    
    // Ensure we always return an empty array instead of nil
    if users == nil {
        log.Printf("Service: No users found, returning empty array")
        return []UserResponse{}, nil
    }
    
    // Add role for frontend compatibility
    for i := range users {
        // Default role if not available from DB
        users[i].Role = "StandardUser"
    }
    
    log.Printf("Service: Returning %d users", len(users))
    return users, nil
}

// CreateOrUpdateTenantUser creates a new user or updates an existing user's tenant assignment
func (s *Service) CreateOrUpdateTenantUser(tenantID uint, input UserInput) (*UserResponse, error) {
    log.Printf("🔴 SERVICE: Creating/updating user for tenant %d: %s (%s)", tenantID, input.Name, input.Email)
    
    // Check if user exists
    existingUser, err := s.repo.GetUserByEmail(input.Email)
    if err != nil {
        log.Printf("Error checking for existing user: %v", err)
        return nil, err
    }
    
    var userID uint
    createdBy := uint(1) // Default to system user, ideally get from context
    
    if existingUser != nil {
        // User exists, use their ID
        log.Printf("User already exists (ID: %d), will update assignment", existingUser.ID)
        userID = existingUser.ID
    } else {
        // Create new user
        log.Printf("User does not exist, creating new user")
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
        if err != nil {
            log.Printf("Failed to hash password: %v", err)
            return nil, errors.New("failed to hash password")
        }
        
        // Get role ID from name
        roleID, err := s.repo.GetRoleIDByName(input.Role)
        if err != nil {
            log.Printf("Invalid role '%s': %v", input.Role, err)
            roleID = 5 // Default to standarduser if lookup fails
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
    
    // Create or update tenant user assignment - explicitly passing tenantID parameter
    log.Printf("🔴 SERVICE: Passing tenant ID %d to repository", tenantID)
    err = s.repo.UpdateTenantUserAssignment(userID, tenantID, createdBy)
    if err != nil {
        log.Printf("Failed to assign user to tenant: %v", err)
        return nil, errors.New("failed to assign user to tenant: " + err.Error())
    }
    
    log.Printf("User successfully assigned to tenant %d", tenantID)
    
    // Construct user response directly rather than re-fetching
    response := &UserResponse{
        ID:        userID,
        Name:      input.Name,
        Email:     input.Email,
        Phone:     input.Phone,
        Status:    "active",
        CreatedAt: time.Now(),
        Role:      input.Role, // Use input role for consistency
    }
    
    return response, nil
}