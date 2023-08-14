package user

import (
	"context"
	"transfer/api/types/errors"
	userModel "transfer/internal/app/model/user"
	"transfer/internal/app/pkg/logger"
	"transfer/internal/app/pkg/password"

	"github.com/google/wire"
)

var ServiceSet = wire.NewSet(wire.Struct(new(Service), "*"))

type Service struct {
	UserRepo *userModel.UserRepo
}

func (s *Service) Register(c context.Context, user *userModel.User) error {
	//先查询该用户是否存在， 如存在则直接返回错误
	if s.CheckExist(user) {
		logger.Infof("User is exist", user)
		return errors.ErrUserExist
	}

	password := password.EncryptPassword(user.Password)
	user.Password = password

	// 插入新用户
	if err := s.Create(c, user); err != nil {
		return err
	}

	return nil
}

func (s *Service) Login(user *userModel.User) error {
	pass := user.Password

	if err := s.UserRepo.Info(user); err != nil {
		return err
	}
	if !password.ValidatePassword(user.Password, pass) {
		return errors.ErrPassword
	}
	return nil

}

func (s *Service) Create(c context.Context, user *userModel.User) error {
	if err := s.UserRepo.Create(user); err != nil {
		return err
	}
	// if err := storage.CreateUser(c, user.Username, user.Password); err != nil {
	// 	return err
	// }
	return nil
}

func (s *Service) CheckExist(user *userModel.User) bool {
	return s.UserRepo.CheckExist(user)

}

func (s *Service) ByName(user *userModel.User) error {
	if err := s.UserRepo.Info(user); err != nil {
		return err
	}
	return nil
}

func (s *Service) Role(c context.Context, roleId int) string {
	return s.UserRepo.Role(roleId)
}
