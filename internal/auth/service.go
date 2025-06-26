package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/sharath018/temple-management-backend/config"
)

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type Service interface {
	Register(input RegisterInput) error
	Login(input LoginInput) (*TokenPair, *User, error) // ✅ Updated to return user
	Refresh(refreshToken string) (string, error)
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
	role, err := s.repo.FindRoleByName(in.Role)
	if err != nil {
		return errors.New("invalid role")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &User{
		FullName:     in.FullName,
		Email:        in.Email,
		PasswordHash: string(hash),
		RoleID:       role.ID,
	}
	return s.repo.Create(user)
}

func (s *service) Login(in LoginInput) (*TokenPair, *User, error) {
	user, err := s.repo.FindByEmail(in.Email)
	if err != nil {
		return nil, nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)) != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	// Access token
	accessClaims := jwt.MapClaims{
		"user_id": user.ID,
		"role_id": user.RoleID,
		"exp":     time.Now().Add(s.accessTTL).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	at, err := accessToken.SignedString([]byte(s.accessSecret))
	if err != nil {
		return nil, nil, err
	}

	// Refresh token
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
	}, user, nil // ✅ Return the user along with token pair
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
