package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snawflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	//判断用户是否存在
	if err = mysql.CheckUserExist(p.UserName); err != nil {
		//数据库查询错误
		return err
	}
	//生成UID
	UserID := snawflake.GenID()

	user := models.User{
		UserID:   UserID,
		Username: p.UserName,
		Password: p.PassWord,
	}

	//保存进数据库
	return mysql.InsertUser(&user)
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.UserName,
		Password: p.PassWord,
	}

	if err = mysql.Login(user); err != nil {
		return nil, err
	}

	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}
