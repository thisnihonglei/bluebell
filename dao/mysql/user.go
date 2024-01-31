package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

const secret = "nihonglei"

// CheckUserExist 根据用户名查询用户是否重复存在
func CheckUserExist(username string) (err error) {
	SqlStr := `select count(id) from user where username= ?`
	var count int
	err = db.Get(&count, SqlStr, username)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) error {
	//对密码进行加密
	user.Password = encryptPassword(user.Password)
	//执行SQL语句入库
	SqlStr := `insert into user(user_id,username,password) values(?,?,?)`
	_, err := db.Exec(SqlStr, user.UserID, user.Username, user.Password)
	return err
}

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(user *models.User) (err error) {
	oPassword := user.Password
	SqlStr := `select user_id,username,password from user where username=?`
	err = db.Get(user, SqlStr, user.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrorNoUserExist
	}
	if err != nil {
		return
	}
	//判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}

func GetUserById(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id,username from user where user_id=?`
	err = db.Get(user, sqlStr, uid)
	return
}
