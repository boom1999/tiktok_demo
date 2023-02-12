package redis

import (
	"context"
	"tiktok_demo/config"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

var RdbFollowers *redis.Client
var RdbFollowing *redis.Client
var RdbFollowingPart *redis.Client

var RdbLikeUserId *redis.Client  //key:userId,value:VideoId
var RdbLikeVideoId *redis.Client //key:VideoId,value:userId

var RdbVCid *redis.Client //redis db11 -- video_id + comment_id
var RdbCVid *redis.Client //redis db12 -- comment_id + video_id

// InitRedis 初始化redis连接
func InitRedis() {
	Conf := config.GetConfig()
	host := Conf.Redis.Host
	port := Conf.Redis.Port
	password := Conf.Redis.Password
	addr := host + ":" + port
	RdbFollowers = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // 粉丝列表信息存入 DB0.
	})
	RdbFollowing = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       1, // 关注列表信息信息存入 DB1.
	})
	RdbFollowingPart = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       2, // 当前用户是否关注了自己粉丝信息存入 DB2.
	})

	RdbLikeUserId = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       4, //  选择将点赞视频id信息存入 DB4.
	})

	RdbLikeVideoId = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       6, //  选择将点赞用户id信息存入 DB6.
	})
	RdbVCid = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       7, // lsy 选择将video_id中的评论id s存入 DB7.
	})
	RdbCVid = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       9, // lsy 选择将comment_id对应video_id存入 DB8.
	})
}
