package repository

import (
	"backend/internal/auth"
	"errors"

	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	FindByUsername(username string) (*auth.User, error)
	FindByEmail(email string) (*auth.User, error)
	FindByID(id string) (*auth.User, error)
	FindAll() ([]auth.User, error)
	Create(user *auth.User) error
	Update(id string, updates map[string]interface{}) error
	UpdateRole(id string, role auth.UserRole) error
	DeleteUser(id string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByUsername(username string) (*auth.User, error) {
	var user auth.User
	if err := r.db.Where("username = ? AND deleted_at IS NULL AND is_active = true", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*auth.User, error) {
	var user auth.User
	if err := r.db.Where("email = ? AND deleted_at IS NULL AND is_active = true", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id string) (*auth.User, error) {
	var user auth.User
	if err := r.db.Where("id = ? AND deleted_at IS NULL AND is_active = true", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll() ([]auth.User, error) {
	var users []auth.User
	if err := r.db.Where("deleted_at IS NULL AND is_active = true").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Create(user *auth.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(id string, updates map[string]interface{}) error {
	result := r.db.Model(&auth.User{}).Where("id = ? AND deleted_at IS NULL", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepository) UpdateRole(id string, role auth.UserRole) error {
	result := r.db.Model(&auth.User{}).Where("id = ? AND deleted_at IS NULL AND is_active = true", id).Update("role", role)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepository) DeleteUser(id string) error {
	// Soft delete and deactivate
	result := r.db.Model(&auth.User{}).Where("id = ? AND deleted_at IS NULL", id).Updates(map[string]interface{}{
		"is_active":  false,
		"deleted_at": gorm.Expr("NOW()"),
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}
