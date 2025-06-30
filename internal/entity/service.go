package entity

import (
	"errors"
	"strings"
)

type Service struct {
	Repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{Repo: r}
}

// ========== ENTITY CORE ==========

func (s *Service) CreateEntity(e *Entity) error {
	if strings.TrimSpace(e.TempleName) == "" ||
		strings.TrimSpace(e.EntityCode) == "" ||
		strings.TrimSpace(e.Email) == "" {
		return errors.New("temple name, code, and email are required")
	}
	return s.Repo.CreateEntity(e)
}

func (s *Service) GetAllEntities() ([]Entity, error) {
	return s.Repo.GetAllEntities()
}

func (s *Service) GetEntityByID(id int) (Entity, error) {
	return s.Repo.GetEntityByID(id)
}

func (s *Service) UpdateEntity(e Entity) error {
	return s.Repo.UpdateEntity(e)
}

func (s *Service) ToggleEntityStatus(id int, isActive bool) error {
	return s.Repo.ToggleEntityStatus(id, isActive)
}

func (s *Service) DeleteEntity(id int) error {
	return s.Repo.DeleteEntity(id)
}

// ========== ADDRESS ==========

func (s *Service) AddEntityAddress(addr EntityAddress) error {
	return s.Repo.AddEntityAddress(addr)
}

func (s *Service) GetEntityAddress(entityID int) (EntityAddress, error) {
	return s.Repo.GetEntityAddress(entityID)
}

// ========== DOCUMENTS ==========

func (s *Service) AddEntityDocument(doc EntityDocument) error {
	return s.Repo.AddEntityDocument(doc)
}

func (s *Service) GetEntityDocuments(entityID int) ([]EntityDocument, error) {
	return s.Repo.GetEntityDocuments(entityID)
}
