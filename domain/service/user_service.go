package service

import (
	"errors"
	"github.com/dawn1806/user/domain/model"
	"github.com/dawn1806/user/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	AddUser(*model.User) (int64, error)
	DeleteUser(int64) error
	UpdateUser(user *model.User, isChangedPwd bool) error
	FindUserByName(string) (*model.User, error)
	CheckPwd(userName string, pwd string) (isOK bool, err error)
}

type UserService struct {
	userRepository repository.IUserRepository
}

func NewUserService(userRepository repository.IUserRepository) IUserService {
	return &UserService{
		userRepository: userRepository,
	}
}

// 加密用户密码
func GenPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// 验证用户密码
func ValidatePassword(password string, hashed string) (isOK bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return false, errors.New("密码比对错误")
	}
	return true, nil
}

func (us *UserService) AddUser(user *model.User) (int64, error) {
	bts, err := GenPassword(user.HashPassword)
	if err != nil {
		return user.ID, errors.New("密码加密失败")
	}
	user.HashPassword = string(bts)
	return us.userRepository.CreateUser(user)
}

func (us *UserService) DeleteUser(id int64) error {
	return us.userRepository.DeleteUserByID(id)
}

func (us *UserService) UpdateUser(user *model.User, isOK bool) error {
	// 判断是否更新了密码
	if isOK {
		bts, err := GenPassword(user.HashPassword)
		if err != nil {
			return err
		}
		user.HashPassword = string(bts)
	}
	return us.userRepository.UpdateUser(user)
}

func (us *UserService) FindUserByName(name string) (user *model.User, err error) {
	return us.userRepository.GetUserInfoByName(name)
}

func (us *UserService) CheckPwd(name string, pwd string) (isOK bool, err error) {
	user, err := us.userRepository.GetUserInfoByName(name)
	if err != nil {
		return false, err
	}
	return ValidatePassword(pwd, user.HashPassword)
}
