package utils

import (
	"account/domain/model"
	. "account/proto/account"
	"encoding/json"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SwapTo 通过json tag 进行结构体赋值
func SwapTo(request, category interface{}) (err error) {
	dataByte, err := json.Marshal(request)
	if err != nil {
		return
	}
	return json.Unmarshal(dataByte, category)
}

// UserForResponse 类型转化
func UserForResponse(response *UserInfoResponse, userModel *model.User) *UserInfoResponse {
	response.UserId = userModel.ID
	response.Username = userModel.UserName
	response.FirstName = userModel.FirstName
	response.LastName = userModel.LastName
	response.Email = userModel.Email
	response.IsActive = userModel.IsActive
	response.Permission = userModel.Permission
	response.CreateDate = timestamppb.New(userModel.CreateDate)
	response.UpdateDate = timestamppb.New(userModel.UpdateDate)
	return response
}
