package redis

//redis key
//redis 使用命名空间的方式，方便查询和拆分

const (
	KeyPrefix        = "bluebell"
	KeyPostTimeZSet  = "post:time"   //ZSet 帖子及发帖时间
	KeyPostScoreZSet = "post:score"  //ZSet 帖子及投票分数
	KeyPostVotedZSet = "post:voted:" //ZSet 记录用户及投票类型
	KeyCommunitySet  = "community:"  //Set 保存每个分区下帖子的id
)

// 给Redis Key加上前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}
