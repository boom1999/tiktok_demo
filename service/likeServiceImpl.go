package service

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"tiktok_demo/config"
	"tiktok_demo/middleware/rabbitmq"
	"tiktok_demo/middleware/redis"
	"tiktok_demo/repository"
	"tiktok_demo/util"
	"time"
)

type LikeServiceImpl struct {
	VideoService
	UserService
}

/*
FavouriteAction 当前用户对视频进行点赞或取消点赞操作。1：点赞，2：取消点赞
保持数据一致性策略：先维护缓存，再更新数据库
*/
func (like *LikeServiceImpl) FavouriteAction(userId int64, videoId int64, actionType int32) error {
	strUserId := strconv.FormatInt(userId, 10)
	strVideoId := strconv.FormatInt(videoId, 10)
	// 拼接打入消息队列的信息
	sb := strings.Builder{}
	sb.WriteString(strUserId)
	sb.WriteString(" ")
	sb.WriteString(strVideoId)

	// 执行点赞操作维护
	if actionType == config.LikeAction {
		// 查询 RdbLikeUserId 是否已经加载过 strUserId
		if n, err := redis.RdbLikeUserId.Exists(redis.Ctx, strUserId).Result(); n > 0 {
			if err != nil {
				util.Log.Error("method:FavouriteAction RedisLikeUserId query key failed" + err.Error())
				return err
			}
			// 为了保证 数据库 与 redis 中的数据是一致的，只有 redis 操作成功才执行数据库 likes 表操作
			if _, err1 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, videoId).Result(); err1 != nil {
				util.Log.Error("method:FavouriteAction RedisLikeUserId add value failed" + err1.Error())
				return err1
			} else {
				// 更新数据库
				rabbitmq.RmqLikeAdd.Publish(sb.String())
			}
		} else {
			// 加入 DefaultRedisValue 目的：防止删最后一个数据的时候数据库还没更新完出现脏读，或者数据库操作失败造成的脏读
			if _, err := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, config.DefaultRedisValue).Result(); err != nil {
				util.Log.Error("method:FavouriteAction RedisLikeUserId add value failed" + err.Error())
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return err
			}
			// 给 strUserId 设置过期时间
			_, err := redis.RdbLikeUserId.Expire(redis.Ctx, strUserId, time.Duration(config.Config.OneDayOfHours.OneMonth)*time.Second).Result()
			if err != nil {
				util.Log.Error("method:FavouriteAction RedisLikeUserId set expiration time failed" + err.Error())
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return err
			}
			// 根据 userId 查询所有点赞的 videoId
			videoIdList, err1 := repository.GetLikeVideoIdList(userId)
			if err1 != nil {
				return err1
			}
			// 将 videoIdList 里的数据添加到 strUserId 集合中，若失败，删除 strUserId，防止脏读
			for _, likeVideoId := range videoIdList {
				if _, err1 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, likeVideoId).Result(); err1 != nil {
					util.Log.Error("method:FavouriteAction RedisLikeUserId add value failed" + err1.Error())
					redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
					return err1
				}
			}
			// 将当前的 videoId 添加到 strUserId 集合中
			if _, err2 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, videoId).Result(); err2 != nil {
				util.Log.Error("method:FavouriteAction RedisLikeUserId add value failed" + err2.Error())
				return err2
			} else {
				rabbitmq.RmqLikeAdd.Publish(sb.String())
			}
		}
		// 查询 RdbLikeVideoId 是否已经加载过 strVideoId
		if n, err := redis.RdbLikeVideoId.Exists(redis.Ctx, strVideoId).Result(); n > 0 {
			if err != nil {
				util.Log.Error("method:FavouriteAction RedisLikeVideoId query key failed" + err.Error())
				return err
			}
			if _, err1 := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, userId).Result(); err1 != nil {
				util.Log.Error("method:FavouriteAction RedisLikeVideoId add value failed" + err1.Error())
				return err1
			}
		} else {
			//如果不存在，则维护 RdbLikeVideoId 新建 strVideoId，设置过期时间，加入DefaultRedisValue
			if _, err := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, config.DefaultRedisValue).Result(); err != nil {
				util.Log.Error("method:FavouriteAction RedisLikeVideoId add value failed" + err.Error())
				redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
				return err
			}
			//给 strVideoId 设置过期时间
			_, err := redis.RdbLikeVideoId.Expire(redis.Ctx, strVideoId, time.Duration(config.Config.OneDayOfHours.OneMonth)*time.Second).Result()
			if err != nil {
				util.Log.Error("method:FavouriteAction RedisLikeVideoId set expiration time failed" + err.Error())
				redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
				return err
			}

			// 根据 videoId 获取点赞用户的 userId
			userIdList, err1 := repository.GetLikeUserIdList(videoId)
			if err1 != nil {
				return err1
			}
			// 将 userIdList 里的数据添加到 strVideoId 集合中，若失败，删除 strVideoId，防止脏读
			for _, likeUserId := range userIdList {
				if _, err1 := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, likeUserId).Result(); err1 != nil {
					util.Log.Error("method:FavouriteAction RedisLikeVideoId add value failed" + err1.Error())
					redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
					return err1
				}
			}
			// 将当前的 userId 添加到 strVideoId 集合中
			if _, err2 := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, userId).Result(); err2 != nil {
				util.Log.Error("method:FavouriteAction RedisLikeVideoId add value failed" + err2.Error())
				return err2
			}
		}
	} else {
		// 执行取消赞操作维护
		// 查询 RdbLikeUserId 是否已经加载过 strUserId
		if n, err := redis.RdbLikeUserId.Exists(redis.Ctx, strUserId).Result(); n > 0 {
			if err != nil {
				util.Log.Error("method:FavouriteAction RedisLikeUserId query key failed" + err.Error())
				return err
			}
			// 防止出现 redis 与 数据库 数据不一致情况，当 redis 删除操作成功，才执行数据库更新操作
			if _, err1 := redis.RdbLikeUserId.SRem(redis.Ctx, strUserId, videoId).Result(); err1 != nil {
				util.Log.Error("method:FavouriteAction RedisLikeUserId del value failed" + err1.Error())
				return err1
			} else {
				// 更新数据库
				rabbitmq.RmqLikeDel.Publish(sb.String())
			}
		} else {
			// 加入 DefaultRedisValue 目的：防止删最后一个数据的时候数据库还没更新完出现脏读，或者数据库操作失败造成的脏读
			if _, err := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, config.DefaultRedisValue).Result(); err != nil {
				util.Log.Error("method:FavouriteAction RedisLikeUserId add value failed" + err.Error())
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return err
			}
			//给 strUserId 设置过期时间
			_, err := redis.RdbLikeUserId.Expire(redis.Ctx, strUserId, time.Duration(config.Config.OneDayOfHours.OneMonth)*time.Second).Result()
			if err != nil {
				util.Log.Error("method:FavouriteAction RedisLikeUserId set expiration time failed" + err.Error())
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return err
			}
			// 根据 userId 查询所有点赞的 videoId
			videoIdList, err1 := repository.GetLikeVideoIdList(userId)
			if err1 != nil {
				return err1
			}
			// 将 videoIdList 里的数据添加到 strUserId 集合中，若失败，删除 strUserId，防止脏读
			for _, likeVideoId := range videoIdList {
				if _, err1 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, likeVideoId).Result(); err1 != nil {
					util.Log.Error("method:FavouriteAction RedisLikeUserId add value failed" + err1.Error())
					redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
					return err1
				}
			}
			// 将当前的 videoId 从 strUserId 集合中删除
			if _, err2 := redis.RdbLikeUserId.SRem(redis.Ctx, strUserId, videoId).Result(); err2 != nil {
				util.Log.Error("method:FavouriteAction RedisLikeUserId del value failed" + err2.Error())
				return err2
			} else {
				// 更新数据库
				rabbitmq.RmqLikeDel.Publish(sb.String())
			}
		}

		// 查询 RdbLikeVideoId 是否已经加载过 strVideoId
		if n, err := redis.RdbLikeVideoId.Exists(redis.Ctx, strVideoId).Result(); n > 0 {
			if err != nil {
				util.Log.Error("method:FavouriteAction RedisLikeVideoId query key failed" + err.Error())
				return err
			}
			if _, err1 := redis.RdbLikeVideoId.SRem(redis.Ctx, strVideoId, userId).Result(); err1 != nil {
				util.Log.Error("method:FavouriteAction RedisLikeVideoId del value failed" + err1.Error())
				return err1
			}
		} else {
			// 加入 DefaultRedisValue 目的：防止删最后一个数据的时候数据库还没更新完出现脏读，或者数据库操作失败造成的脏读
			if _, err := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, config.DefaultRedisValue).Result(); err != nil {
				util.Log.Error("method:FavouriteAction RedisLikeVideoId add value" + err.Error())
				redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
				return err
			}
			//给 strVideoId 设置过期时间
			_, err := redis.RdbLikeVideoId.Expire(redis.Ctx, strVideoId, time.Duration(config.Config.OneDayOfHours.OneMonth)*time.Second).Result()
			if err != nil {
				util.Log.Error("method:FavouriteAction RedisLikeVideoId set expiration time failed" + err.Error())
				redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
				return err
			}

			// 根据 videoId 获取点赞用户的 userId
			userIdList, err1 := repository.GetLikeUserIdList(videoId)
			if err1 != nil {
				redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
				return err1
			}
			// 将 userIdList 里的数据添加到 strVideoId 集合中，若失败，删除 strVideoId，防止脏读
			for _, likeUserId := range userIdList {
				if _, err1 := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, likeUserId).Result(); err1 != nil {
					util.Log.Error("method:FavouriteAction RedisLikeVideoId add value failed" + err1.Error())
					redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
					return err1
				}
			}
			// 将当前的 userId 从 strVideoId 集合中删除
			if _, err2 := redis.RdbLikeVideoId.SRem(redis.Ctx, strVideoId, userId).Result(); err2 != nil {
				util.Log.Error("method:FavouriteAction RedisLikeVideoId del value failed" + err2.Error())
				return err2
			}
		}
	}
	return nil
}

// GetFavouriteList 返回 当前用户 的点赞列表
func (like *LikeServiceImpl) GetFavouriteList(userId int64, curId int64) ([]Video, error) {
	strUserId := strconv.FormatInt(userId, 10)
	// 查询 RdbLikeUserId，如果 strUserId 存在,则获取集合中全部 videoId
	if n, err := redis.RdbLikeUserId.Exists(redis.Ctx, strUserId).Result(); n > 0 {
		if err != nil {
			util.Log.Error("method:GetFavouriteList RedisLikeVideoId query key failed" + err.Error())
			return nil, err
		}
		// 获取集合中全部 videoId
		videoIdList, err1 := redis.RdbLikeUserId.SMembers(redis.Ctx, strUserId).Result()
		if err1 != nil {
			util.Log.Error("method:GetFavouriteList RedisLikeVideoId get values failed" + err1.Error())
			return nil, err1
		}

		favoriteVideoList := new([]Video)
		// 采用协程并发将 Video 类型对象添加到集合中去
		i := len(videoIdList) - 1 // 去掉 DefaultRedisValue
		if i == 0 {
			return *favoriteVideoList, nil
		}
		var wg sync.WaitGroup
		wg.Add(i)
		for j := 0; j <= i; j++ {
			videoId, _ := strconv.ParseInt(videoIdList[j], 10, 64)
			if videoId == config.DefaultRedisValue {
				continue
			}
			go like.addFavouriteVideoList(videoId, curId, favoriteVideoList, &wg)
		}
		wg.Wait()
		return *favoriteVideoList, nil
	} else {
		if _, err := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, config.DefaultRedisValue).Result(); err != nil {
			util.Log.Error("method:GetFavouriteList RedisLikeUserId add value failed" + err.Error())
			redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
			return nil, err
		}
		// 给 strUserId 设置过期时间
		_, err := redis.RdbLikeUserId.Expire(redis.Ctx, strUserId, time.Duration(config.Config.OneDayOfHours.OneMonth)*time.Second).Result()
		if err != nil {
			util.Log.Error("method:GetFavouriteList RedisLikeUserId set expiration time failed" + err.Error())
			redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
			return nil, err
		}
		videoIdList, err1 := repository.GetLikeVideoIdList(userId)
		if err1 != nil {
			util.Log.Error("method:GetFavouriteList get likeVideoIdList failed" + err1.Error())
			redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
			return nil, err1
		}

		for _, likeVideoId := range videoIdList {
			if _, err2 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, likeVideoId).Result(); err2 != nil {
				util.Log.Error("method:GetFavouriteList RedisLikeUserId add value failed" + err2.Error())
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return nil, err2
			}
		}

		favoriteVideoList := new([]Video)
		// 采用协程并发将 Video 类型对象添加到集合中去
		i := len(videoIdList) - 1 // 去掉 DefaultRedisValue
		if i == 0 {
			return *favoriteVideoList, nil
		}
		var wg sync.WaitGroup
		wg.Add(i)
		for j := 0; j <= i; j++ {
			if videoIdList[j] == config.DefaultRedisValue {
				continue
			}
			go like.addFavouriteVideoList(videoIdList[j], curId, favoriteVideoList, &wg)
		}
		wg.Wait()
		return *favoriteVideoList, nil
	}
}

// addFavouriteVideoList 添加视频对象到点赞列表空间
func (like *LikeServiceImpl) addFavouriteVideoList(videoId int64, curId int64, favoriteVideoList *[]Video, wg *sync.WaitGroup) {
	defer wg.Done()
	// 调用 videoService 接口的 GetVideo() 函数
	videoService := new(VideoServiceImpl)
	video, err := videoService.GetVideo(videoId, curId)
	if err != nil {
		util.Log.Debug(errors.New("this favourite video is miss").Error())
		return
	}
	// 将 Video 类型对象添加到集合中去
	*favoriteVideoList = append(*favoriteVideoList, video)
}

// addVideoLikeCount 根据 videoId，获取该视频的点赞数
func (like *LikeServiceImpl) addVideoLikeCount(videoId int64, videoLikeCountList *[]int64, wg *sync.WaitGroup) {
	defer wg.Done()
	// 调用 FavouriteCount：根据 videoId 获取点赞数
	count, err := like.FavouriteCount(videoId)
	if err != nil {
		util.Log.Debug(err.Error())
		return
	}
	*videoLikeCountList = append(*videoLikeCountList, count)
}

// IsFavourite 查询点赞状态
func (like *LikeServiceImpl) IsFavourite(videoId int64, userId int64) (bool, error) {
	strUserId := strconv.FormatInt(userId, 10)
	strVideoId := strconv.FormatInt(videoId, 10)

	// 查询 RdbLikeUserId 的 strUserId 中是否存在 videoId
	if n, err := redis.RdbLikeUserId.Exists(redis.Ctx, strUserId).Result(); n > 0 {
		if err != nil {
			util.Log.Error("method:IsFavourite RedisLikeUserId query key failed" + err.Error())
			return false, err
		}
		exist, err1 := redis.RdbLikeUserId.SIsMember(redis.Ctx, strUserId, videoId).Result()
		if err1 != nil {
			util.Log.Error("method:IsFavourite RedisLikeUserId query value failed" + err1.Error())
			return false, err1
		}
		util.Log.Debug("method:IsFavourite RedisLikeUserId query value success")
		return exist, nil
	} else {
		// 如果 RdbLikeUserId 不存在 strUserId，查询 RdbLikeVideoId 的 strVideoId 中是否存在 userId
		if n, err := redis.RdbLikeVideoId.Exists(redis.Ctx, strVideoId).Result(); n > 0 {
			if err != nil {
				util.Log.Error("method:IsFavourite RedisLikeVideoId query key failed" + err.Error())
				return false, err
			}
			exist, err1 := redis.RdbLikeVideoId.SIsMember(redis.Ctx, strVideoId, userId).Result()
			if err1 != nil {
				util.Log.Error("method:IsFavourite RedisLikeVideoId query value failed" + err1.Error())
				return false, err1
			}
			util.Log.Debug("method:IsFavourite RedisLikeVideoId query value success")
			return exist, nil
		} else {
			// 如果 RdbLikeUserId 和 RdbLikeVideoId 都不存在对应数据。在 RdbLikeUserId 中创建 strUserId 集合
			if _, err := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, config.DefaultRedisValue).Result(); err != nil {
				util.Log.Error("method:IsFavourite RedisLikeUserId add value failed" + err.Error())
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return false, err
			}
			//给 strUserId 设置过期时间
			_, err := redis.RdbLikeUserId.Expire(redis.Ctx, strUserId, time.Duration(config.Config.OneDayOfHours.OneMonth)*time.Second).Result()
			if err != nil {
				util.Log.Error("method:IsFavourite RedisLikeUserId set expiration time" + err.Error())
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return false, err
			}
			// 通过 userId 查询 likes 表，返回所有点赞 videoId，并添加到 RdbLikeUserId 的 strUserId 集合中
			videoIdList, err1 := repository.GetLikeVideoIdList(userId)
			if err1 != nil {
				util.Log.Debug("method:IsFavourite get likeVideoIdList failed" + err1.Error())
				return false, err1
			}
			for _, likeVideoId := range videoIdList {
				redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, likeVideoId)
			}
			// 查询 RdbLikeUserId 的 strUserId 集合中是否存在 videoId
			exist, err2 := redis.RdbLikeUserId.SIsMember(redis.Ctx, strUserId, videoId).Result()
			if err2 != nil {
				util.Log.Error("method:IsFavourite RedisLikeUserId query value failed" + err2.Error())
				return false, err2
			}
			util.Log.Debug("method:IsFavourite RedisLikeUserId query value success")
			return exist, nil
		}
	}
}

// FavouriteCount 根据 videoId 获取对应点赞数量
func (like *LikeServiceImpl) FavouriteCount(videoId int64) (int64, error) {
	strVideoId := strconv.FormatInt(videoId, 10)

	// 如果 RdbLikeVideoId 中存在 strVideoId 集合，则计算集合中 userId 个数
	if n, err := redis.RdbLikeVideoId.Exists(redis.Ctx, strVideoId).Result(); n > 0 {
		if err != nil {
			util.Log.Error("method:FavouriteCount RedisLikeVideoId query key failed" + err.Error())
			return 0, err
		}
		// 获取集合中 userId 个数
		count, err1 := redis.RdbLikeVideoId.SCard(redis.Ctx, strVideoId).Result()
		if err1 != nil {
			util.Log.Error("method:FavouriteCount RedisLikeVideoId query count failed" + err1.Error())
			return 0, err1
		}
		util.Log.Debug("method:FavouriteCount RedisLikeVideoId query count success")
		return count - 1, nil // 去掉 DefaultRedisValue
	} else {
		// 在 RdbLikeVideoId 中创建 strVideoId 集合
		if _, err := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, config.DefaultRedisValue).Result(); err != nil {
			util.Log.Error("method:FavouriteCount RedisLikeVideoId add value failed" + err.Error())
			redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
			return 0, err
		}
		//给 strVideoId 设置过期时间
		_, err := redis.RdbLikeVideoId.Expire(redis.Ctx, strVideoId, time.Duration(config.Config.OneDayOfHours.OneMonth)*time.Second).Result()
		if err != nil {
			util.Log.Error("method:FavouriteCount RedisLikeVideoId set expiration time failed" + err.Error())
			redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
			return 0, err
		}
		// 通过 videoId 查询 likes 表，返回所有点赞 userId，并添加到 RdbLikeVideoId 的 strVideoId 集合中
		// 再通过 set 集合中 userId 个数，获取点赞数量
		userIdList, err1 := repository.GetLikeUserIdList(videoId)
		if err1 != nil {
			util.Log.Debug("method:FavouriteCount get likeUserIdList failed" + err1.Error())
			return 0, err1
		}
		for _, likeUserId := range userIdList {
			redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, likeUserId)
		}
		// 通过 strVideoId 集合中 userId 个数，获取点赞数量
		count, err2 := redis.RdbLikeVideoId.SCard(redis.Ctx, strVideoId).Result()
		if err2 != nil {
			util.Log.Error("method:FavouriteCount RedisLikeVideoId query count failed" + err2.Error())
			return 0, err2
		}
		util.Log.Debug("method:FavouriteCount RedisLikeVideoId query count success")
		return count - 1, nil // 去掉 DefaultRedisValue
	}
}

// TotalFavourite 根据 userId 获取该用户总共被点赞数量
func (like *LikeServiceImpl) TotalFavourite(userId int64) (int64, error) {
	// 根据 userId 获取该用户的发布视频列表信息
	videoService := new(VideoServiceImpl)
	videoIdList, err := videoService.GetVideoIdList(userId)
	if err != nil {
		util.Log.Debug(err.Error())
		return 0, err
	}
	var sum int64 //该用户的总被点赞数
	videoLikeCountList := new([]int64)

	// 采用协程并发将对应 videoId 的点赞数添加到集合中去
	i := len(videoIdList)
	var wg sync.WaitGroup
	wg.Add(i)
	for j := 0; j < i; j++ {
		go like.addVideoLikeCount(videoIdList[j], videoLikeCountList, &wg)
	}
	wg.Wait()
	for _, count := range *videoLikeCountList {
		sum += count
	}
	return sum, nil
}

// FavouriteVideoCount 根据 userId 获取该用户点赞视频数量
func (like *LikeServiceImpl) FavouriteVideoCount(userId int64) (int64, error) {
	strUserId := strconv.FormatInt(userId, 10)
	// 查询 RdbLikeUserId，如果 strUserId 集合存在，则获取集合中元素个数
	if n, err := redis.RdbLikeUserId.Exists(redis.Ctx, strUserId).Result(); n > 0 {
		if err != nil {
			util.Log.Error("method:FavouriteVideoCount RdbLikeUserId query key failed" + err.Error())
			return 0, err
		} else {
			count, err1 := redis.RdbLikeUserId.SCard(redis.Ctx, strUserId).Result()
			if err1 != nil {
				util.Log.Error("method:FavouriteVideoCount RdbLikeUserId query count failed" + err1.Error())
				return 0, err1
			}
			util.Log.Debug("method:FavouriteVideoCount RdbLikeUserId query count success")
			return count - 1, nil // 去掉 DefaultRedisValue

		}
	} else {
		// 在 RdbLikeUserId 中创建 strUserId 集合
		if _, err := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, config.DefaultRedisValue).Result(); err != nil {
			util.Log.Error("method:FavouriteVideoCount RedisLikeUserId add value failed" + err.Error())
			redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
			return 0, err
		}
		//给 strUserId 设置过期时间
		_, err := redis.RdbLikeUserId.Expire(redis.Ctx, strUserId, time.Duration(config.Config.OneDayOfHours.OneMonth)*time.Second).Result()
		if err != nil {
			util.Log.Error("method:FavouriteVideoCount RedisLikeUserId set expiration time" + err.Error())
			redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
			return 0, err
		}
		// 通过 userId 查询 likes 表，返回所有点赞 videoId，并添加到 RdbLikeUserId 的 strUserId 集合中
		videoIdList, err1 := repository.GetLikeVideoIdList(userId)
		if err1 != nil {
			util.Log.Debug(err1.Error())
			return 0, err1
		}
		for _, likeVideoId := range videoIdList {
			if _, err1 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, likeVideoId).Result(); err1 != nil {
				util.Log.Error("method:FavouriteVideoCount RedisLikeUserId add value failed" + err1.Error())
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return 0, err1
			}
		}
		// 再通过 strUserId 集合中 videoId 个数，获取点赞数量
		count, err2 := redis.RdbLikeUserId.SCard(redis.Ctx, strUserId).Result()
		if err2 != nil {
			util.Log.Error("method:FavouriteVideoCount RdbLikeUserId query count failed" + err2.Error())
			return 0, err2
		}
		util.Log.Debug("method:FavouriteVideoCount RdbLikeUserId query count success")
		return count - 1, nil // 去掉 DefaultRedisValue
	}
}
