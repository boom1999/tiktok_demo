package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"tiktok_demo/config"
)

var Ctx = context.Background()
var RdbFollowers *redis.Client
var RdbFollowing *redis.Client
var RdbFollowingPart *redis.Client

var RdbLikeUserId *redis.Client  //key:userId,value:VideoId
var RdbLikeVideoId *redis.Client //key:VideoId,value:userId

var RdbVCid *redis.Client //redis db11 -- video_id + comment_id
var RdbCVid *redis.Client //redis db12 -- comment_id + video_id

// InitRedis 初始化Redis连接。
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
	_, err := RdbFollowers.Ping(Ctx).Result()
	if err != nil {
		panic("failed to connect redis, err:" + err.Error())
	}

	RdbFollowing = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       1, // 关注列表信息信息存入 DB1.
	})
	_, err = RdbFollowing.Ping(Ctx).Result()
	if err != nil {
		panic("failed to connect redis, err:" + err.Error())
	}

	RdbFollowingPart = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       3, // 当前用户是否关注了自己粉丝信息存入 DB1.
	})
	_, err = RdbFollowingPart.Ping(Ctx).Result()
	if err != nil {
		panic("failed to connect redis, err:" + err.Error())
	}

	RdbLikeUserId = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       5, //  选择将点赞视频id信息存入 DB5.
	})
	_, err = RdbLikeUserId.Ping(Ctx).Result()
	if err != nil {
		panic("failed to connect redis, err:" + err.Error())
	}

	RdbLikeVideoId = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       6, //  选择将点赞用户id信息存入 DB6.
	})
	_, err = RdbLikeVideoId.Ping(Ctx).Result()
	if err != nil {
		panic("failed to connect redis, err:" + err.Error())
	}

	RdbVCid = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       11, // lsy 选择将video_id中的评论id s存入 DB11.
	})
	_, err = RdbVCid.Ping(Ctx).Result()
	if err != nil {
		panic("failed to connect redis, err:" + err.Error())
	}

	RdbCVid = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       12, // lsy 选择将comment_id对应video_id存入 DB12.
	})
	_, err = RdbCVid.Ping(Ctx).Result()
	if err != nil {
		panic("failed to connect redis, err:" + err.Error())
	}

}
