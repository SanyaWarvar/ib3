package service

import (
	"errors"

	"github.com/SanyaWarvar/ib3/pkg/models"
	"github.com/SanyaWarvar/ib3/pkg/repository"
)

type UserService struct {
	repo repository.IUserRepo
}

func NewUserService(repo repository.IUserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) HashPassword(password string) (string, error) {
	return s.repo.HashPassword(password)
}

func (s *UserService) CreateUser(user models.User) error {
	var err error
	user.Password, err = s.HashPassword(user.Password)
	if err != nil {
		return err
	}
	return s.repo.CreateUser(user)
}

func (s *UserService) GetUserByUP(user models.User) (models.User, error) {
	targetUser, err := s.repo.GetUserByU(user.Username)
	if err != nil {
		return user, err
	}

	if s.repo.ComparePassword(user.Password, targetUser.Password) {
		return targetUser, err
	}
	return user, errors.New("incorrect password")
}

func (s *UserService) GetUserByU(username string) (models.User, error) {

	return s.repo.GetUserByU(username)

}
