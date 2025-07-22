package userprofile

import (
	"errors"
	"time"

	"github.com/sharath018/temple-management-backend/internal/auth"
	"gorm.io/gorm"
)

// ========== INTERFACES ==========

type Service interface {
	CreateOrUpdateProfile(userID uint, entityID uint, input DevoteeProfileInput) (*DevoteeProfile, error)
	Get(userID uint) (*DevoteeProfile, error)
	JoinTemple(userID uint, entityID uint) (*UserEntityMembership, error)
	ListMemberships(userID uint) ([]UserEntityMembership, error)
}

// ========== SERVICE INIT ==========

type service struct {
	repo     Repository
	authRepo auth.Repository
}

func NewService(repo Repository, authRepo auth.Repository) Service {
	return &service{repo: repo, authRepo: authRepo}
}

// ========== PROFILE DTO ==========

type DevoteeProfileInput struct {
	// Section 1
	FullName      *string    `json:"full_name"`
	DOB           *time.Time `json:"dob"`
	Gender        *string    `json:"gender"`
	StreetAddress *string    `json:"street_address"`
	City          *string    `json:"city"`
	State         *string    `json:"state"`
	Pincode       *string    `json:"pincode"`
	Country       *string    `json:"country"`

	// Section 2
	Gotra     *string `json:"gotra"`
	Nakshatra *string `json:"nakshatra"`
	Rashi     *string `json:"rashi"`
	Lagna     *string `json:"lagna"`
	VedaShaka *string `json:"veda_shaka"`

	// Section 3
	FatherName              *string `json:"father_name"`
	FatherGotra             *string `json:"father_gotra"`
	FatherNativePlace       *string `json:"father_native_place"`
	FatherVedaShaka         *string `json:"father_veda_shaka"`
	MotherName              *string `json:"mother_name"`
	MaidenGotra             *string `json:"maiden_gotra"`
	MotherNativePlace       *string `json:"mother_native_place"`
	MaternalGrandfatherName *string `json:"maternal_grandfather_name"`
	PaternalGrandfatherName *string `json:"paternal_grandfather_name"`
	PaternalGrandmotherName *string `json:"paternal_grandmother_name"`
	MaternalGrandmotherName *string `json:"maternal_grandmother_name"`

	// Section 4
	SevaAbhisheka              *bool   `json:"seva_abhisheka"`
	SevaArti                   *bool   `json:"seva_arti"`
	SevaAnnadana               *bool   `json:"seva_annadana"`
	SevaArchana                *bool   `json:"seva_archana"`
	SevaKalyanam               *bool   `json:"seva_kalyanam"`
	SevaHomam                  *bool   `json:"seva_homam"`
	DonateTempleMaintenance    *bool   `json:"donate_temple_maintenance"`
	DonateAnnadanaProgram      *bool   `json:"donate_annadana_program"`
	DonateFestivalCelebrations *bool   `json:"donate_festival_celebrations"`
	DonateReligiousEducation   *bool   `json:"donate_religious_education"`
	DonateTempleConstruction   *bool   `json:"donate_temple_construction"`
	DonateGeneral              *bool   `json:"donate_general"`
	SpecialInterestsOrNotes    *string `json:"special_interests_or_notes"`

	// Section 5
	SpouseName      *string     `json:"spouse_name"`
	SpouseEmail     *string     `json:"spouse_email"`
	SpousePhone     *string     `json:"spouse_phone"`
	SpouseDOB       *time.Time  `json:"spouse_dob"`
	SpouseGotra     *string     `json:"spouse_gotra"`
	SpouseNakshatra *string     `json:"spouse_nakshatra"`
	Children        []*Child    `json:"children"`
	EmergencyContacts []*EmergencyContact `json:"emergency_contacts"`

	// Section 6
	HealthNotes           *string `json:"health_notes"`
	AllergiesOrConditions *string `json:"allergies_or_conditions"`
	DietaryRestrictions   *string `json:"dietary_restrictions"`
	PersonalSankalpa      *string `json:"personal_sankalpa"`
	AdditionalNotes       *string `json:"additional_notes"`
}

// ========== PROFILE LOGIC ==========

func (s *service) Get(userID uint) (*DevoteeProfile, error) {
	return s.repo.GetByUserID(userID)
}

func (s *service) CreateOrUpdateProfile(userID, entityID uint, input DevoteeProfileInput) (*DevoteeProfile, error) {
	existing, err := s.repo.GetByUserID(userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	profile := &DevoteeProfile{
		UserID:   userID,
		EntityID: entityID,

		FullName:      input.FullName,
		DOB:           input.DOB,
		Gender:        input.Gender,
		StreetAddress: input.StreetAddress,
		City:          input.City,
		State:         input.State,
		Pincode:       input.Pincode,
		Country:       input.Country,

		Gotra:     input.Gotra,
		Nakshatra: input.Nakshatra,
		Rashi:     input.Rashi,
		Lagna:     input.Lagna,
		VedaShaka: input.VedaShaka,

		FatherName:               input.FatherName,
		FatherGotra:              input.FatherGotra,
		FatherNativePlace:        input.FatherNativePlace,
		FatherVedaShaka:          input.FatherVedaShaka,
		MotherName:               input.MotherName,
		MaidenGotra:              input.MaidenGotra,
		MotherNativePlace:        input.MotherNativePlace,
		MaternalGrandfatherName:  input.MaternalGrandfatherName,
		PaternalGrandfatherName:  input.PaternalGrandfatherName,
		PaternalGrandmotherName:  input.PaternalGrandmotherName,
		MaternalGrandmotherName:  input.MaternalGrandmotherName,

		SevaAbhisheka:              input.SevaAbhisheka,
		SevaArti:                   input.SevaArti,
		SevaAnnadana:               input.SevaAnnadana,
		SevaArchana:                input.SevaArchana,
		SevaKalyanam:               input.SevaKalyanam,
		SevaHomam:                  input.SevaHomam,
		DonateTempleMaintenance:    input.DonateTempleMaintenance,
		DonateAnnadanaProgram:      input.DonateAnnadanaProgram,
		DonateFestivalCelebrations: input.DonateFestivalCelebrations,
		DonateReligiousEducation:   input.DonateReligiousEducation,
		DonateTempleConstruction:   input.DonateTempleConstruction,
		DonateGeneral:              input.DonateGeneral,
		SpecialInterestsOrNotes:    input.SpecialInterestsOrNotes,

		SpouseName:      input.SpouseName,
		SpouseEmail:     input.SpouseEmail,
		SpousePhone:     input.SpousePhone,
		SpouseDOB:       input.SpouseDOB,
		SpouseGotra:     input.SpouseGotra,
		SpouseNakshatra: input.SpouseNakshatra,
		Children:        input.Children,
		EmergencyContacts: input.EmergencyContacts,

		HealthNotes:           input.HealthNotes,
		AllergiesOrConditions: input.AllergiesOrConditions,
		DietaryRestrictions:   input.DietaryRestrictions,
		PersonalSankalpa:      input.PersonalSankalpa,
		AdditionalNotes:       input.AdditionalNotes,

		ProfileCompletionPercentage: calculateCompletionPercentage(input),
		UpdatedAt:                   time.Now(),
	}

	if existing != nil && existing.ID > 0 {
		profile.ID = existing.ID
		err = s.repo.Update(profile)
		return profile, err
	}

	err = s.repo.Create(profile)
	return profile, err
}

// ========== MEMBERSHIP LOGIC ==========

func (s *service) JoinTemple(userID uint, entityID uint) (*UserEntityMembership, error) {
	existing, err := s.repo.GetMembership(userID, entityID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	membership := &UserEntityMembership{
		UserID:   userID,
		EntityID: entityID,
		JoinedAt: time.Now(),
	}
	if err := s.repo.CreateMembership(membership); err != nil {
		return nil, err
	}

	if err := s.authRepo.UpdateEntityID(userID, entityID); err != nil {
		return nil, err
	}

	return membership, nil
}

func (s *service) ListMemberships(userID uint) ([]UserEntityMembership, error) {
	return s.repo.ListMembershipsByUser(userID)
}

// ========== PROFILE COMPLETION LOGIC ==========

func calculateCompletionPercentage(p DevoteeProfileInput) int {
	filled := 0
	total := 12

	if p.FullName != nil && *p.FullName != "" {
		filled++
	}
	if p.DOB != nil {
		filled++
	}
	if p.Gender != nil && *p.Gender != "" {
		filled++
	}
	if p.StreetAddress != nil && *p.StreetAddress != "" {
		filled++
	}
	if p.Gotra != nil && *p.Gotra != "" {
		filled++
	}
	if p.FatherName != nil && *p.FatherName != "" {
		filled++
	}
	if p.MotherName != nil && *p.MotherName != "" {
		filled++
	}
	if p.HealthNotes != nil && *p.HealthNotes != "" {
		filled++
	}
	if p.PersonalSankalpa != nil && *p.PersonalSankalpa != "" {
		filled++
	}
	if len(p.Children) > 0 {
		filled++
	}
	if len(p.EmergencyContacts) > 0 {
		filled++
	}
	if p.SpecialInterestsOrNotes != nil && *p.SpecialInterestsOrNotes != "" {
		filled++
	}

	return int(float64(filled) / float64(total) * 100)
}

