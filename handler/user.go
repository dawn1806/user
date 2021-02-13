package handler

import (
	"context"
	"github.com/dawn1806/user/domain/model"
	"github.com/dawn1806/user/domain/service"

	user "github.com/dawn1806/user/proto/user"
)

type User struct {
	UserService service.IUserService
}

func (u *User) Login(ctx context.Context, req *user.LoginRequest, res *user.LoginResponse) error {
	isOK, err := u.UserService.CheckPwd(req.UserName, req.Password)
	if err != nil {
		return err
	}
	res.IsSuccess = isOK
	return nil
}

func (u *User) Register(ctx context.Context, req *user.RegisterRequest, res *user.RegisterResponse) error {
	userReg := &model.User{
		UserName:     req.UserName,
		FirstName:    req.FirstName,
		HashPassword: req.Password,
	}
	_, err := u.UserService.AddUser(userReg)
	if err != nil {
		return err
	}
	res.Message = "添加用户成功"
	return nil
}

func (u *User) GetUserInfo(ctx context.Context, req *user.UserInfoRequest, res *user.UserInfoResponse) error {
	userInfo, err := u.UserService.FindUserByName(req.UserName)
	if err != nil {
		return err
	}
	res = UserForResponse(userInfo)
	return nil
}

func UserForResponse(userInfo *model.User) *user.UserInfoResponse {
	return &user.UserInfoResponse{
		UserId:    userInfo.ID,
		UserName:  userInfo.UserName,
		FirstName: userInfo.FirstName,
	}
}
