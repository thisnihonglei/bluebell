package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snawflake"
	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	//1. 生成post id
	p.ID = snawflake.GenID()
	//2. 保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityId)
	if err != nil {
		return err
	}
	return
}

// GetPostById 根据帖子id查询帖子详情
func GetPostById(id int64) (data *models.ApiPostDetail, err error) {
	//查询并且拼接想要的数据
	post, err := mysql.GetPostById(id)
	if err != nil {
		zap.L().Error("mysql.GetPostById(id)", zap.Error(err))
		return
	}
	//根据作者查名称
	user, err := mysql.GetUserById(post.AuthorId)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorId)", zap.Error(err))
		return
	}
	//根据社区ID查询社区详细信息
	community, err := mysql.GetCommunityDetailById(post.CommunityId)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorId)", zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		AuthName:        user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

func GetPostList(pageNum, pageSize int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		//根据作者查名称
		user, err := mysql.GetUserById(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorId)", zap.Error(err))
			continue
		}
		//根据社区ID查询社区详细信息
		community, err := mysql.GetCommunityDetailById(post.CommunityId)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorId)", zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthName:        user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostListV2 获取帖子列表
func GetPostListV2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 2.去redis查询id列表
	ids, err := redis.GetPostIDsInorder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInorder(p) return 0 data")
		return
	}
	// 3.根据id去数据库查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	data = make([]*models.ApiPostDetail, 0, len(posts))
	//提前查好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	for theID, post := range posts {
		//根据作者查名称
		user, err := mysql.GetUserById(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorId)", zap.Error(err))
			continue
		}
		//根据社区ID查询社区详细信息
		community, err := mysql.GetCommunityDetailById(post.CommunityId)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorId)", zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthName:        user.Username,
			VoteNum:         voteData[theID],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

func GetCommunityPostListV2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 2.去redis查询id列表
	ids, err := redis.GetCommunityPostIDsInorder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInorder(p) return 0 data")
		return
	}
	// 3.根据id去数据库查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	data = make([]*models.ApiPostDetail, 0, len(posts))
	//提前查好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	for theID, post := range posts {
		//根据作者查名称
		user, err := mysql.GetUserById(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorId)", zap.Error(err))
			continue
		}
		//根据社区ID查询社区详细信息
		community, err := mysql.GetCommunityDetailById(post.CommunityId)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorId)", zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthName:        user.Username,
			VoteNum:         voteData[theID],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostListV3 将两个查询逻辑合二为一的
func GetPostListV3(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	if p.CommunityId == 0 {
		data, err = GetPostListV2(p)
	} else {
		data, err = GetCommunityPostListV2(p)
	}
	if err != nil {
		zap.L().Error("GetPostListV3 Failed", zap.Error(err))
		return nil, err
	}
	return
}
