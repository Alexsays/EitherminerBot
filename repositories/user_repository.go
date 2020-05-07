package repositories

import (
	"../models"
	"github.com/jinzhu/gorm"
)

// RepositoryResult ...
type RepositoryResult struct {
	Result interface{}
	Error  error
}

// UserRepository ...
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository ...
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Save ...
func (r *UserRepository) Save(user *models.User) RepositoryResult {
	err := r.db.Save(user).Error

	if err != nil {
		return RepositoryResult{Error: err}
	}

	return RepositoryResult{Result: user}
}

// FindAll ...
func (r *UserRepository) FindAll() RepositoryResult {
	var users models.Users

	err := r.db.Find(&users).Error

	if err != nil {
		return RepositoryResult{Error: err}
	}

	return RepositoryResult{Result: &users}
}

// FindOneByID ...
func (r *UserRepository) FindOneByID(id string) RepositoryResult {
	var user models.User

	err := r.db.Where(&models.User{ID: id}).Take(&user).Error

	if err != nil {
		return RepositoryResult{Error: err}
	}

	return RepositoryResult{Result: &user}
}

// FindOneByUsername ...
func (r *UserRepository) FindOneByUsername(username string) RepositoryResult {
	var user models.User

	err := r.db.Where(&models.User{Username: username}).Take(&user).Error

	if err != nil {
		return RepositoryResult{Error: err}
	}

	return RepositoryResult{Result: &user}
}
