package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

//定义请求参数结构体

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	UserName   string `json:"username" binding:"required"`
	PassWord   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=PassWord"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	UserName string `json:"username" binding:"required"`
	PassWord string `json:"password" binding:"required"`
}

// ParamVoteData 投票参数
type ParamVoteData struct {
	//UserID 从请求中获取当前用户
	PostID    string `json:"post_id" binding:"required"`              //帖子ID
	Direction int8   `json:"direction,string" binding:"oneof=1 -1 0"` //赞成票(1)反对票(-1)取消投票(0)
}

// ParamPostList 帖子列表参数
type ParamPostList struct {
	CommunityId int64  `json:"community_id" form:"community_id"` //可以为空
	Page        int64  `json:"page" form:"page"`
	Size        int64  `json:"size" form:"size"`
	Order       string `json:"order" form:"order"`
}

// ParamCommunityPostList 按社区获取帖子列表参数
type ParamCommunityPostList struct {
	*ParamPostList
	CommunityId int64 `json:"community_id" form:"community_id"`
}
