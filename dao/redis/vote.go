package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"math"
	"strconv"
	"time"
)

const oneWeekInSecond = 7 * 24 * 3600

var ErrVoteTimeExpire = errors.New("投票时间已过")
var ErrVoteRepeated = errors.New("不允许重复投票")

const scorePerVote = 432

func CreatePost(postID, communityID int64) error {
	pipeline := client.TxPipeline()
	//帖子时间
	pipeline.ZAdd(context.Background(), getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	//帖子分数
	pipeline.ZAdd(context.Background(), getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	ckey := getRedisKey(KeyCommunitySet + strconv.Itoa(int(communityID)))
	pipeline.SAdd(context.Background(), ckey, postID)
	_, err := pipeline.Exec(context.Background())
	return err
}

func VoteForPost(userID string, postID string, value float64) error {
	//1.判断投票限制
	//去redis读取帖子发布时间
	postTime := client.ZScore(context.Background(), getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSecond {
		return ErrVoteTimeExpire
	}
	//2.更新帖子分数
	//先查当前用户给当前帖子的投票记录
	oValue := client.ZScore(context.Background(), getRedisKey(KeyPostVotedZSet+postID), userID).Val()
	//如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if value == oValue {
		return ErrVoteRepeated
	}

	var dir float64
	if value > oValue {
		dir = 1
	} else {
		dir = -1
	}
	//计算两次投票的差值
	diff := math.Abs(oValue - value)
	//更新分数
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(context.Background(), getRedisKey(KeyPostScoreZSet), dir*diff*scorePerVote, postID)
	//3.记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(context.Background(), getRedisKey(KeyPostVotedZSet+postID), userID)
	} else {
		pipeline.ZAdd(context.Background(), getRedisKey(KeyPostVotedZSet+postID), redis.Z{
			Score:  value, //赞成票还是反对票
			Member: userID,
		})
	}
	_, err := pipeline.Exec(context.Background())
	return err
}
