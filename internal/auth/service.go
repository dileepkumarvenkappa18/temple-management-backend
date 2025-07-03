package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sharath018/temple-management-backend/config"
	"golang.org/x/crypto/bcrypt"
)

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type Service interface {
	Register(input RegisterInput) error
	Login(input LoginInput) (*TokenPair, *User, error)
	Refresh(refreshToken string) (string, error)
	GetUserByID(userID uint) (User, error)
	
}

type service struct {
	repo           Repository
	accessSecret   string
	refreshSecret  string
	accessTTL      time.Duration
	refreshTTL     time.Duration
}

func NewService(r Repository, cfg *config.Config) Service {
	return &service{
		repo:          r,
		accessSecret:  cfg.JWTAccessSecret,
		refreshSecret: cfg.JWTRefreshSecret,
		accessTTL:     time.Duration(cfg.JWTAccessTTLHours) * time.Hour,
		refreshTTL:    time.Duration(cfg.JWTRefreshTTLHours) * time.Hour,
	}
}

type RegisterInput struct {
	FullName string
	Email    string
	Password string
	Role     string
}

type LoginInput struct {
	Email    string
	Password string
}



func (s *service) Register(in RegisterInput) error {
	roleName := strings.ToLower(in.Role)
	role, err := s.repo.FindRoleByName(roleName)
	if err != nil {
		return errors.New("invalid role")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	status := "active"
	if roleName == "templeadmin" {
		status = "pending"
	}

	user := &User{
		FullName:     in.FullName,
		Email:        in.Email,
		PasswordHash: string(hash),
		RoleID:       role.ID,
		Status:       status,
	}

	if err := s.repo.Create(user); err != nil {
		return err
	}

	if roleName == "templeadmin" {
		// ✅ FIXED: Correct approval request type
		if err := s.repo.CreateApprovalRequest(user.ID, "tenant_approval"); err != nil {
			return errors.New("failed to create approval request")
		}
	}

	return nil
}


func (s *service) Login(in LoginInput) (*TokenPair, *User, error) {
	user, err := s.repo.FindByEmail(in.Email)
	if err != nil {
		return nil, nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	switch user.Status {
	case "pending":
		return nil, nil, errors.New("your account is pending approval")
	case "rejected":
		return nil, nil, errors.New("your account was rejected by admin")
	case "inactive":
		return nil, nil, errors.New("your account is inactive")
	}

	if user.Role.RoleName == "templeadmin" && user.Status != "active" {
		return nil, nil, errors.New("your account is pending approval by Super Admin")
	}

	// ✅ Inject entity_id from membership/approval for devotee/volunteer/templeadmin
	if user.EntityID == nil && (user.Role.RoleName == "templeadmin" || user.Role.RoleName == "devotee" || user.Role.RoleName == "volunteer") {
		entityID, err := s.repo.FindEntityIDByUserID(user.ID)
		if err == nil && entityID != nil {
			user.EntityID = entityID
		}
	}

	// ✅ Build access token with entity_id if available
	accessClaims := jwt.MapClaims{
		"user_id": user.ID,
		"role_id": user.RoleID,
		"exp":     time.Now().Add(s.accessTTL).Unix(),
	}
	if user.EntityID != nil {
		accessClaims["entity_id"] = *user.EntityID
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	at, err := accessToken.SignedString([]byte(s.accessSecret))
	if err != nil {
		return nil, nil, err
	}

	// ✅ Build refresh token
	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(s.refreshTTL).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	rt, err := refreshToken.SignedString([]byte(s.refreshSecret))
	if err != nil {
		return nil, nil, err
	}

	return &TokenPair{
		AccessToken:  at,
		RefreshToken: rt,
	}, user, nil
}



func (s *service) Refresh(refreshToken string) (string, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.refreshSecret), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid refresh token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		return "", errors.New("invalid token claims")
	}
	newAccessClaims := jwt.MapClaims{
		"user_id": claims["user_id"],
		"exp":     time.Now().Add(s.accessTTL).Unix(),
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessClaims)
	return newToken.SignedString([]byte(s.accessSecret))
}

func (s *service) GetUserByID(userID uint) (User, error) {
	return s.repo.FindByID(userID)
}
