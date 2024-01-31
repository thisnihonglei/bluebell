package mysql

import (
	"errors"
)

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorNoUserExist     = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
	ErrorInvalidId       = errors.New("无效的ID")
)
