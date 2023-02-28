package handler

import (
	"account/common/mail"
	"account/common/snow_flake"
	"account/common/token"
	"account/common/utils"
	"account/domain/model"
	"account/domain/service"
	. "account/proto/account"
	"context"
	"errors"
	"fmt"
	"time"
)

type Account struct {
	AccountService service.IUserService
}

// Login 登录
func (a *Account) Login(ctx context.Context, req *LoginRequest, rsp *LoginResponse) error {
	userInfo, err := a.AccountService.FindUserByName(req.Username)
	isOk, err := a.AccountService.CheckPwd(req.Username, req.Password)
	if err != nil {
		return err
	}
	token2, err := token.GenToken(req.Username)
	if err != nil {
		return err
	}
	token.SetToken(req.Username, token2)
	rsp.IsSuccess = isOk
	rsp.Token = token2
	rsp.UserId = userInfo.UserID
	return nil
}

// Register 注册
func (a *Account) Register(ctx context.Context, req *RegisterRequest, rsp *RegisterResponse) error {
	var sf snow_flake.Snowflake
	userId := sf.NextVal()
	_, err := mail.CheckMail(req.RegisterRequest.Email, req.Code)
	if err != nil {
		return err
	}
	nowTime := time.Now()
	userRegister := &model.User{
		UserID:     userId,
		UserName:   req.RegisterRequest.Username,
		FirstName:  req.RegisterRequest.FirstName,
		LastName:   req.RegisterRequest.LastName,
		PassWord:   req.RegisterRequest.Password,
		Permission: 0,
		CreateDate: nowTime,
		UpdateDate: nowTime,
		IsActive:   1,
		Email:      req.RegisterRequest.Email,
	}
	_, err = a.AccountService.AddUser(userRegister)
	if err != nil {
		return err
	}
	mail.DelMail(req.RegisterRequest.Email)
	rsp.IsSuccess = true
	rsp.UserId = userId
	return nil
}

// GetUserInfo 查询用户信息
func (a *Account) GetUserInfo(ctx context.Context, req *UserIdRequest, rsp *UserInfoResponse) error {
	userInfo, err := a.AccountService.FindUserByID(req.UserId)
	if err != nil {
		return err
	}
	rsp = utils.UserForResponse(rsp, userInfo)
	fmt.Println(rsp)
	return nil
}

// UpdateUserInfo 修改信息
func (a *Account) UpdateUserInfo(ctx context.Context, req *UserInfoRequest, rsp *Response) error {
	var user *model.User
	err := utils.SwapTo(req, user)
	if err != nil {
		return err
	}
	isChangePwd := false
	if req.UserInfo.Password != "" {
		isChangePwd = true
	}
	if err = a.AccountService.UpdateUser(user, isChangePwd); err != nil {
		return err
	}
	rsp.Message = "修改信息成功"
	return nil
}

// SendRegisterMail 发送注册邮件
func (a *Account) SendRegisterMail(ctx context.Context, req *SendMailRequest, rsp *SendMailResponse) error {
	code, err := mail.SendRegisterMail(req.Email)
	if err != nil {
		return err
	}
	mail.SetMail(req.Email, code)
	rsp.Msg = "邮件发送成功"
	rsp.Code = code
	return nil
}

// SendResetPwdMail 发送重置密码邮件
func (a *Account) SendResetPwdMail(ctx context.Context, req *SendMailRequest, rsp *SendMailResponse) error {
	email := req.Email
	code, err := mail.SendResetPwdMail(email)
	if err != nil {
		return err
	}
	mail.SetMail(email, code)
	rsp.Msg = "邮件发送成功"
	return nil
}

// ResetPwd 重置密码
func (a *Account) ResetPwd(ctx context.Context, req *ResetPwdRequest, rsp *Response) error {
	userInfo, err := a.AccountService.FindUserByID(req.UserId)
	if err != nil {
		return err
	}
	code, err := mail.GetMail(userInfo.Email)
	if err != nil {
		return err
	}
	if code != req.Code {
		return errors.New("验证码错误")
	}
	if err := a.AccountService.ResetPwd(req.UserId, req.Password); err != nil {
		return err
	}
	rsp.Message = "重置密码成功"
	return nil
}

// GetUserPermission 获取权限
func (a *Account) GetUserPermission(ctx context.Context, req *UserIdRequest, rsp *GetPermissionResponse) error {
	permission, err := a.AccountService.GetPermission(req.UserId)
	if err != nil {
		return err
	}
	rsp.Permission = permission
	return nil
}

// UpdateUserPermission 修改权限
func (a *Account) UpdateUserPermission(ctx context.Context, req *UpdatePermissionRequest, rsp *Response) error {
	if err := a.AccountService.UpdatePermission(req.UserId, req.Permission); err != nil {
		return err
	}
	rsp.Message = "修改权限成功"
	return nil
}

// Logout 退出账号
func (a *Account) Logout(ctx context.Context, req *UserIdRequest, rsp *Response) error {
	userInfo, err := a.AccountService.FindUserByID(req.UserId)
	if err != nil {
		return err
	}
	_, err = token.GenToken(userInfo.UserName)
	if err != nil {
		return errors.New("账号未登录！")
	}
	token.DelToken(userInfo.UserName)
	rsp.Message = "退出账号成功"
	return nil
}

// DelUser 删除账号
func (a *Account) DelUser(ctx context.Context, req *UserIdRequest, rsp *Response) error {
	if err := a.AccountService.DeleteUser(req.UserId); err != nil {
		return err
	}
	rsp.Message = "账号删除成功"
	return nil
}

// DisableUser 禁用账号
func (a *Account) DisableUser(ctx context.Context, req *UserIdRequest, rsp *Response) error {
	if err := a.AccountService.DisableUser(req.UserId); err != nil {
		return err
	}
	rsp.Message = "账号禁用成功"
	return nil
}

// EnableUser 启用账号
func (a *Account) EnableUser(ctx context.Context, req *UserIdRequest, rsp *Response) error {
	if err := a.AccountService.EnableUser(req.UserId); err != nil {
		return err
	}
	rsp.Message = "账号启用成功"
	return nil
}
