package services

import (
	"gin-user-app/models"
	"gin-user-app/repositories"
)

type UserService interface {
	GetUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{userRepo: repo}
}

func (s *userService) GetUsers() ([]models.User, error) {
	return s.userRepo.GetUsers()
}

func (s *userService) GetUserByID(id int) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *userService) CreateUser(user *models.User) error {
	return s.userRepo.CreateUser(user)
}

func (s *userService) UpdateUser(user *models.User) error {
	return s.userRepo.UpdateUser(user)
}

func (s *userService) DeleteUser(id int) error {
	return s.userRepo.DeleteUser(id)
}
