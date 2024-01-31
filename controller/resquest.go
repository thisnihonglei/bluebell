package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

const CtxUserIDKey = "UserId"

var ErrorUserNotLogin = errors.New("用户未登录")

// GetCurrenUserID 获取当前登录用户ID
func GetCurrenUserID(c *gin.Context) (UserId int64, err error) {
	uid, ok := c.Get(CtxUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	UserId, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

func getPageInfo(c *gin.Context) (int64, int64) {
	pageNumStr := c.Query("pageNum")
	pageSizeStr := c.Query("pageSize")

	var (
		pageNum  int64
		pageSize int64
		err      error
	)

	pageNum, err = strconv.ParseInt(pageNumStr, 10, 64)
	if err != nil {
		pageNum = 1
	}
	pageSize, err = strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		pageSize = 10
	}
	return pageNum, pageSize
}
