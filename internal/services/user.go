package services

import (
	"errors"
	"go_base/internal/models"
	"go_base/internal/repositories"
)

type UserService interface {
	Login(request models.UserLoginRequest) (string, error)
	GetUserByID(models.UserGetRequest) (models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Login(request models.UserLoginRequest) (string, error) {
	user, err := s.userRepo.GetUserByUsername(request.Username)
	if err != nil {
		return "", err
	}

	isPasswordMatch, err := s.userRepo.CompareHash(request.Password, user.Password)
	if err != nil {
		return "", err
	}

	if !isPasswordMatch {
		return "", errors.New("invalid password")
	}

	token, err := s.userRepo.GenereateJWTToken(user.Id)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) GetUserByID(request models.UserGetRequest) (models.User, error) {
	user, err := s.userRepo.GetUserByID(request.Id)
	if err != nil {
		return models.User{}, err
	}

	if user.Id == "" {
		return models.User{}, errors.New("user not found")
	}

	return user, nil
}
