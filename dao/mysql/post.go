package mysql

import (
	"bluebell/models"
	"github.com/jmoiron/sqlx"
	"strings"
)

func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post(post_id,title,content,author_id,community_id) values (?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorId, p.CommunityId)
	return
}

// GetPostById 根据id查询单个帖子详情
func GetPostById(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select post_id,title,content,author_id,community_id,create_time from post where post_id=?`
	err = db.Get(post, sqlStr, pid)
	return
}

// GetPostList 查询帖子列表函数
func GetPostList(pageNum, pageSize int64) (posts []*models.Post, err error) {
	posts = make([]*models.Post, 0, 2)
	sqlStr := `select post_id,title,content,author_id,community_id,create_time from post order by create_time desc limit ?,?`
	err = db.Select(&posts, sqlStr, (pageNum-1)*pageSize, pageSize)
	return
}

// GetPostListByIDs 根据ID列表查询帖子数据
func GetPostListByIDs(ids []string) (postsList []*models.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,create_time from post where post_id in (?) order by FIND_IN_SET(post_id,?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	err = db.Select(&postsList, query, args...)
	return
}
