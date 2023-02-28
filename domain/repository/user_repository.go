package repository

import (
	"account/domain/model"
	"github.com/jinzhu/gorm"
)

type IUserRepository interface {
	// InitTable 初始化数据表
	InitTable() error
	// FindUserByName 根据用户名称查找用户信息
	FindUserByName(string) (*model.User, error)
	// FindUserByID 根据用户ID查找用户信息
	FindUserByID(int64) (*model.User, error)
	// CreateUser 创建用户
	CreateUser(*model.User) (int64, error)
	// DeleteUserByID 根据用户ID删除用户
	DeleteUserByID(int64) error
	// UpdateUser 更新用户信息
	UpdateUser(*model.User) error
	// FindAll 查找所有用户
	FindAll() ([]model.User, error)
	// ResetPwd 重置密码
	ResetPwd(int64, string) error
	// GetPermission 获取权限
	GetPermission(int64) (int64, error)
	// UpdatePermission 修改权限
	UpdatePermission(int64, int64) error
	// EnableUser 启用账号
	EnableUser(int64) error
	// DisableUser 禁用账号
	DisableUser(int64) error
}

// NewUserRepository 创建UserRepository
func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{mysqlDb: db}
}

type UserRepository struct {
	mysqlDb *gorm.DB
}

// InitTable 初始化表
func (u *UserRepository) InitTable() error {
	return u.mysqlDb.CreateTable(&model.User{}).Error
}

// FindUserByName 根据用户名称查找用户信息
func (u *UserRepository) FindUserByName(name string) (user *model.User, err error) {
	user = &model.User{}
	return user, u.mysqlDb.Where("username=?", name).Find(user).Error
}

// FindUserByID 根据用户ID查找用户信息
func (u *UserRepository) FindUserByID(userID int64) (user *model.User, err error) {
	user = &model.User{}
	return user, u.mysqlDb.Where("user_id=?", userID).Find(user).Error
}

// CreateUser 创建用户
func (u *UserRepository) CreateUser(user *model.User) (userID int64, err error) {
	return user.ID, u.mysqlDb.Create(user).Error
}

// DeleteUserByID 删除用户
func (u *UserRepository) DeleteUserByID(userID int64) error {
	return u.mysqlDb.Where("user_id=?", userID).Delete(&model.User{}).Error
}

// UpdateUser 更新用户信息
func (u *UserRepository) UpdateUser(user *model.User) error {
	return u.mysqlDb.Model(user).Update(&user).Error
}

// FindAll 查找所有用户
func (u *UserRepository) FindAll() (userAll []model.User, err error) {
	return userAll, u.mysqlDb.Find(&userAll).Error
}

// ResetPwd 重置密码
func (u *UserRepository) ResetPwd(userID int64, Pwd string) error {
	return u.mysqlDb.Where("user_id=?",userID).Update(model.User{PassWord:Pwd}).Error
}

// GetPermission 获取权限
func (u *UserRepository) GetPermission(userID int64) (Permission int64, err error) {
	var user model.User
	return user.Permission,u.mysqlDb.Where("user_id=?",userID).Find(user).Error
}

// UpdatePermission 修改权限
func (u *UserRepository) UpdatePermission(userID int64, Permission int64) error {
	return u.mysqlDb.Where("user_id=?",userID).Update(model.User{Permission:Permission}).Error
}

// EnableUser 启用账号
func (u *UserRepository) EnableUser(userID int64) error {
	return u.mysqlDb.Where("user_id=?",userID).Update(&model.User{IsActive: 1}).Error
}

// DisableUser 禁用账号
func (u *UserRepository) DisableUser(userID int64) error {
	return u.mysqlDb.Where("user_id=?",userID).Update(&model.User{IsActive: 0}).Error
}
