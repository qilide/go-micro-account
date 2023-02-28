package service

import (
	"account/domain/model"
	"account/domain/repository"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	// AddUser 插入用户
	AddUser(user *model.User) (int64, error)
	// DeleteUser 删除用户
	DeleteUser(int64) error
	// UpdateUser 更新用户
	UpdateUser(user *model.User, isChangePwd bool) (err error)
	// FindUserByName 根据用户名称查找用户信息
	FindUserByName(string) (*model.User, error)
	// FindUserByID 根据用户ID查找用户信息
	FindUserByID(int64) (*model.User, error)
	// CheckPwd 比对账号密码是否正确
	CheckPwd(userName string, pwd string) (isOk bool, err error)
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

// NewUserService 创建实例
func NewUserService(userRepository repository.IUserRepository) IUserService {
	return &UserService{UserRepository: userRepository}
}

type UserService struct {
	UserRepository repository.IUserRepository
}

// GeneratePassword 加密用户密码
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

// ValidatePassword 验证用户密码
func ValidatePassword(userPassword string, hashed string) (isOk bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, errors.New("密码比对错误")
	}
	return true, nil
}

// AddUser 插入用户
func (u *UserService) AddUser(user *model.User) (userID int64, err error) {
	pwdByte, err := GeneratePassword(user.PassWord)
	if err != nil {
		return user.ID, err
	}
	user.PassWord = string(pwdByte)
	return u.UserRepository.CreateUser(user)
}

// DeleteUser 删除用户
func (u *UserService) DeleteUser(userID int64) error {
	return u.UserRepository.DeleteUserByID(userID)
}

// UpdateUser 更新用户
func (u *UserService) UpdateUser(user *model.User, isChangePwd bool) (err error) {
	if isChangePwd {
		pwdByte, err := GeneratePassword(user.PassWord)
		if err != nil {
			return nil
		}
		user.PassWord = string(pwdByte)
	}
	return u.UserRepository.UpdateUser(user)
}

// FindUserByName 根据用户名称查找用户信息
func (u *UserService) FindUserByName(userName string) (user *model.User, err error) {
	return u.UserRepository.FindUserByName(userName)
}

// FindUserByID 根据用户名称查找用户信息
func (u *UserService) FindUserByID(userId int64) (user *model.User, err error) {
	return u.UserRepository.FindUserByID(userId)
}

// CheckPwd 比对账号密码是否正确
func (u *UserService) CheckPwd(userName string, pwd string) (isOk bool, err error) {
	user, err := u.UserRepository.FindUserByName(userName)
	if err != nil {
		return false, err
	}
	return ValidatePassword(pwd, user.PassWord)
}

// ResetPwd 重置密码
func (u *UserService) ResetPwd(userID int64, pwd string) error {
	return u.UserRepository.ResetPwd(userID, pwd)
}

// GetPermission 获取权限
func (u *UserService) GetPermission(userID int64) (permission int64, err error) {
	return u.UserRepository.GetPermission(userID)
}

// UpdatePermission 修改权限
func (u *UserService) UpdatePermission(userID int64, permission int64) error {
	return u.UserRepository.UpdatePermission(userID, permission)
}

// EnableUser 启用账号
func (u *UserService) EnableUser(userID int64) error {
	return u.UserRepository.EnableUser(userID)
}

// DisableUser 禁用账号
func (u *UserService) DisableUser(userID int64) error {
	return u.UserRepository.DisableUser(userID)
}
