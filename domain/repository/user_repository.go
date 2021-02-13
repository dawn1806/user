package repository

import (
	"github.com/dawn1806/user/domain/model"
	"github.com/jinzhu/gorm"
)

type IUserRepository interface {
	// 初始化用户表
	InitTable() error
	// 创建用户
	CreateUser(*model.User) (int64, error)
	// 根据ID删除用户
	DeleteUserByID(int64) error
	// 更新用户
	UpdateUser(*model.User) error
	// 通过ID查询用户信息
	GetUserInfoByID(int64) (*model.User, error)
	// 通过用户名查询用户信息
	GetUserInfoByName(string) (*model.User, error)
	// 查询所有用户信息
	FindAll() ([]*model.User, error)
}

type UserRepository struct {
	mysqlDB *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{
		mysqlDB: db,
	}
}

func (ur *UserRepository) InitTable() error {
	return ur.mysqlDB.CreateTable(&model.User{}).Error
}

func (ur *UserRepository) CreateUser(user *model.User) (int64, error) {
	return user.ID, ur.mysqlDB.Create(user).Error
}

func (ur *UserRepository) DeleteUserByID(ID int64) error {
	return ur.mysqlDB.Where("id = ?", ID).Delete(&model.User{}).Error
}

func (ur *UserRepository) UpdateUser(user *model.User) error {
	return ur.mysqlDB.Model(user).Update(&user).Error
}

func (ur *UserRepository) GetUserInfoByID(ID int64) (user *model.User, err error) {
	user = &model.User{}
	return user, ur.mysqlDB.First(user, ID).Error
}

func (ur *UserRepository) GetUserInfoByName(name string) (user *model.User, err error) {
	user = &model.User{}
	return user, ur.mysqlDB.Where("user_name = ?", name).Find(user).Error
}

func (ur *UserRepository) FindAll() (userAll []*model.User, err error) {
	return userAll, ur.mysqlDB.Find(userAll).Error
}
