package tenant

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

// GetUsers retrieves users by tenantID with optional role and name filters
func (s *Service) GetUsers(tenantID uint, role string, name string) ([]TenantUser, error) {
	return s.repo.GetUsers(tenantID, role, name)
}

// CreateOrUpdateUser updates an existing user by email (tenant-scoped) or creates a new user
func (s *Service) CreateOrUpdateUser(user TenantUser) (*TenantUser, error) {
	if user.TenantID == 0 {
		return nil, errors.New("tenant ID is required")
	}

	// Check if user exists in the same tenant
	existing, err := s.repo.GetByEmail(user.Email, user.TenantID)
	if err != nil {
		return nil, err
	}

	// Hash password if provided
	if user.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		user.Password = string(hashed)
	}

	if existing != nil {
		// Update existing user
		existing.Name = user.Name
		existing.Phone = user.Phone
		existing.Role = user.Role
		if user.Password != "" {
			existing.Password = user.Password
		}
		existing.UpdatedAt = time.Now()

		err := s.repo.Update(existing)
		if err != nil {
			return nil, errors.New("failed to update user: " + err.Error())
		}
		return existing, nil
	}

	// Create new user
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err = s.repo.Create(&user)
	if err != nil {
		return nil, errors.New("failed to create user: " + err.Error())
	}
	return &user, nil
}

// DeleteUser removes a user within the tenant
func (s *Service) DeleteUser(userID, tenantID uint) error {
	if tenantID == 0 {
		return errors.New("tenant ID is required")
	}
	return s.repo.Delete(userID, tenantID)
}
