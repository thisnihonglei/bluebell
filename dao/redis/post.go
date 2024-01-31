package redis

import (
	"bluebell/models"
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

func getIDsFromKey(key string, pageSize, pageNum int64) ([]string, error) {
	start := (pageNum - 1) * pageSize
	end := start + pageSize - 1
	//按分数从大到小指定数量的查询
	return client.ZRevRange(context.Background(), key, start, end).Result()
}

func GetPostIDsInorder(p *models.ParamPostList) ([]string, error) {
	//从redis获取id
	//根据用户请求中携带的order参数确定要查询的redis的key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	return getIDsFromKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSet + id)
	//	//查找Key中分数为1的元素的数量
	//	v := client.ZCount(context.Background(), key, "1", "1").Val()
	//	data = append(data, v)
	//}
	//使用pipeline一次发送多条命令，减少RTT
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSet + id)
		pipeline.ZCount(context.Background(), key, "1", "1")
	}
	cmders, err := pipeline.Exec(context.Background())
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDsInorder 按社区查询ids
func GetCommunityPostIDsInorder(p *models.ParamPostList) ([]string, error) {
	//使用zinterstore 把分区的帖子set与帖子的分数zset 生成一个新的zset
	//针对新的zset按照之前的 按照之前的逻辑取数据
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	//社区的key
	cKey := getRedisKey(KeyCommunitySet + strconv.Itoa(int(p.CommunityId)))
	//利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityId))
	if client.Exists(context.Background(), orderKey).Val() < 1 {
		//不存在需要计算
		pipeline := client.Pipeline()
		pipeline.ZInterStore(context.Background(), key, &redis.ZStore{Keys: []string{cKey, orderKey}, Aggregate: "MAX"})
		pipeline.Expire(context.Background(), key, 60*time.Second) //设置超时时间
		_, err := pipeline.Exec(context.Background())
		if err != nil {
			return nil, err
		}
	}
	//存在的话，直接根据Key查找
	return getIDsFromKey(key, p.Page, p.Size)
}
