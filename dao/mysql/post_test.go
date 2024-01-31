package mysql

import (
	"bluebell/models"
	"bluebell/setting"
	"testing"
)

func init() {
	mysqlConfig := &setting.MySQLConfig{
		Host:         "127.0.0.1",
		User:         "root",
		Password:     "root",
		DB:           "bluebell",
		Port:         3306,
		MaxOpenConns: 20,
		MaxIdleConns: 10,
	}
	err := Init(mysqlConfig)
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	post := &models.Post{
		ID:          123,
		AuthorId:    123,
		CommunityId: 1,
		Title:       "test",
		Content:     "just test",
	}
	err := CreatePost(post)
	if err != nil {
		t.Fatalf("CreatePost insert record into mysql failed,err:%v\n", err)
	}
	t.Logf("CreatePost insert record into mysql success")
}
