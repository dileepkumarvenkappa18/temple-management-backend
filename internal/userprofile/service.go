package userprofile

import (
	"errors"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"github.com/sharath018/temple-management-backend/internal/auth"
)

type Service interface {
	CreateOrUpdateProfile(userID uint, input CreateOrUpdateProfileDTO) (*DevoteeProfile, error)
	Get(userID uint) (*DevoteeProfile, error)

	// Membership logic
	JoinTemple(userID uint, entityID uint) (*UserEntityMembership, error)
	ListMemberships(userID uint) ([]UserEntityMembership, error)
}

type service struct {
	repo Repository
	authRepo auth.Repository
}

func NewService(repo Repository, authRepo auth.Repository) Service {
	return &service{
		repo:     repo,
		authRepo: authRepo,
	}
}

// ========== DevoteeProfile ==========

func (s *service) Get(userID uint) (*DevoteeProfile, error) {
	return s.repo.GetByUserID(userID)
}

func (s *service) CreateOrUpdateProfile(userID uint, input CreateOrUpdateProfileDTO) (*DevoteeProfile, error) {
	existing, err := s.repo.GetByUserID(userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	profile := &DevoteeProfile{
		UserID:                      userID,
		DateOfBirth:                 input.DateOfBirth,
		Gender:                      input.Gender,
		Address:                     input.Address,
		Gotra:                       input.Gotra,
		Nakshatra:                   input.Nakshatra,
		Rashi:                       input.Rashi,
		Lagna:                       input.Lagna,
		VedaShaka:                   input.VedaShaka,
		FatherName:                  input.FatherName,
		MotherName:                  input.MotherName,
		GrandfatherName:             input.GrandfatherName,
		GrandmotherName:             input.GrandmotherName,
		FamilyTreeURL:               input.FamilyTreeURL,
		FavoriteSevas:               input.FavoriteSevas,
		DonationInterests:           input.DonationInterests,
		SpouseInfo:                  input.SpouseInfo,
		ChildrenInfo:                input.ChildrenInfo,
		EmergencyContacts:           input.EmergencyContacts,
		HealthInfo:                  input.HealthInfo,
		SankalpaNotes:               input.SankalpaNotes,
		PhotosURL:                   input.PhotosURL,
		ProfileCompletionPercentage: calculateCompletionPercentage(input),
		UpdatedAt:                   time.Now(),
	}

	if existing != nil && existing.ID > 0 {
		profile.ID = existing.ID
		err := s.repo.Update(profile)
		return profile, err
	}

	err = s.repo.Create(profile)
	return profile, err
}

// ========== Membership ==========

func (s *service) JoinTemple(userID uint, entityID uint) (*UserEntityMembership, error) {
	existing, err := s.repo.GetMembership(userID, entityID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil // already joined
	}

	membership := &UserEntityMembership{
		UserID:   userID,
		EntityID: entityID,
		JoinedAt: time.Now(),
	}
	if err := s.repo.CreateMembership(membership); err != nil {
		return nil, err
	}

	// ✅ Update auth.users table with entity_id
	if err := s.authRepo.UpdateEntityID(userID, entityID); err != nil {
		return nil, err
	}

	return membership, nil
}



func (s *service) ListMemberships(userID uint) ([]UserEntityMembership, error) {
	return s.repo.ListMembershipsByUser(userID)
}

// ========== DTO & Completion ==========

type CreateOrUpdateProfileDTO struct {
	DateOfBirth       *time.Time     `json:"date_of_birth"`
	Gender            *string        `json:"gender"`
	Address           *string        `json:"address"`
	Gotra             *string        `json:"gotra"`
	Nakshatra         *string        `json:"nakshatra"`
	Rashi             *string        `json:"rashi"`
	Lagna             *string        `json:"lagna"`
	VedaShaka         *string        `json:"veda_shaka"`
	FatherName        *string        `json:"father_name"`
	MotherName        *string        `json:"mother_name"`
	GrandfatherName   *string        `json:"grandfather_name"`
	GrandmotherName   *string        `json:"grandmother_name"`
	FamilyTreeURL     *string        `json:"family_tree_url"`
	FavoriteSevas     datatypes.JSON `json:"favorite_sevas"`
	DonationInterests datatypes.JSON `json:"donation_interests"`
	SpouseInfo        datatypes.JSON `json:"spouse_info"`
	ChildrenInfo      datatypes.JSON `json:"children_info"`
	EmergencyContacts datatypes.JSON `json:"emergency_contacts"`
	HealthInfo        *string        `json:"health_info"`
	SankalpaNotes     *string        `json:"sankalpa_notes"`
	PhotosURL         datatypes.JSON `json:"photos_url"`
}


func calculateCompletionPercentage(p CreateOrUpdateProfileDTO) int {
	filled := 0
	total := 10

	if p.DateOfBirth != nil {
		filled++
	}
	if p.Gender != nil && *p.Gender != "" {
		filled++
	}
	if p.Address != nil && *p.Address != "" {
		filled++
	}
	if p.Gotra != nil && *p.Gotra != "" {
		filled++
	}
	if p.Nakshatra != nil && *p.Nakshatra != "" {
		filled++
	}
	if p.FatherName != nil && *p.FatherName != "" {
		filled++
	}
	if p.MotherName != nil && *p.MotherName != "" {
		filled++
	}
	if p.HealthInfo != nil && *p.HealthInfo != "" {
		filled++
	}
	if len(p.PhotosURL) > 0 {
		filled++
	}
	if len(p.FavoriteSevas) > 0 {
		filled++
	}

	return int(float64(filled) / float64(total) * 100)
}

