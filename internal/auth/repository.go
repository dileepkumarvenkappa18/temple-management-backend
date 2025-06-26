package auth

import "gorm.io/gorm"

type Repository interface {
  Create(user *User) error
  FindByEmail(email string) (*User, error)
  FindRoleByName(name string) (*UserRole, error) // ✅ Added
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository {
  return &repository{db}
}

func (r *repository) Create(user *User) error {
  return r.db.Create(user).Error
}

func (r *repository) FindByEmail(email string) (*User, error) {
  var u User
  err := r.db.Preload("Role").Where("email = ?", email).First(&u).Error
  return &u, err
}

func (r *repository) FindRoleByName(name string) (*UserRole, error) {
  var role UserRole
  err := r.db.Where("role_name = ?", name).First(&role).Error
  return &role, err
}
