package service

import (
	"strings"

	"server/internal/model"
	"server/internal/repository"
	"server/internal/utils"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(req model.CreateUserRequest) (model.User, error) {
	if strings.TrimSpace(req.Name) == "" || !utils.IsValidEmail(req.Email) {
		return model.User{}, ErrInvalidInput
	}

	return s.repo.Create(model.User{
		Name:  strings.TrimSpace(req.Name),
		Email: strings.TrimSpace(req.Email),
	}), nil
}

func (s *UserService) GetUser(id int) (model.User, error) {
	user, ok := s.repo.GetByID(id)
	if !ok {
		return model.User{}, ErrNotFound
	}
	return user, nil
}

func (s *UserService) ListUsers() []model.User {
	return s.repo.List()
}
