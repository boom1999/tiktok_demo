package service

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"tiktok_demo/middleware/rabbitmq"
	"tiktok_demo/middleware/redis"
	"tiktok_demo/repository"
	"tiktok_demo/util"
)

type FollowImpl struct {
	UserService
}

var (
	followImpl        *FollowImpl
	followServiceOnce sync.Once
)

var expireTime = time.Hour * 48

type FollowService interface {
	/*
	   一、其他同学需要调用的业务方法。
	*/
	// IsFollowing 根据当前用户id和目标用户id来判断当前用户是否关注了目标用户
	IsFollowing(userId int64, targetId int64) (bool, error)
	// GetFollowerCnt 根据用户id来查询用户被多少其他用户关注
	GetFollowerCnt(userId int64) (int64, error)
	// GetFollowingCnt 根据用户id来查询用户关注了多少其它用户
	GetFollowingCnt(userId int64) (int64, error)
	/*
	   二、直接request需要的业务方法
	*/
	// AddFollowRelation 当前用户关注目标用户
	AddFollowRelation(userId int64, targetId int64) (bool, error)
	// DeleteFollowRelation 当前用户取消对目标用户的关注
	DeleteFollowRelation(userId int64, targetId int64) (bool, error)
	// GetFollowing 获取当前用户的关注列表
	GetFollowing(userId int64) ([]User, error)
	// GetFollowers 获取当前用户的粉丝列表
	GetFollowers(userId int64) ([]User, error)
}

// NewFSIInstance 生成并返回FollowServiceImp结构体单例变量。
func NewFSIInstance() *FollowImpl {
	followServiceOnce.Do(
		func() {
			followImpl = &FollowImpl{
				UserService: &UserImpl{
					FollowService: &FollowImpl{},
				},
			}
		})
	return followImpl
}

// IsFollowing 给定当前用户和目标用户id，判断是否存在关注关系。
func (*FollowImpl) IsFollowing(userId int64, targetId int64) (bool, error) {
	// 先查Redis里面是否有此关系。
	if flag, err := redis.RdbFollowingPart.SIsMember(redis.Ctx, strconv.Itoa(int(userId)), targetId).Result(); flag {
		// 重现设置过期时间。
		redis.RdbFollowingPart.Expire(redis.Ctx, strconv.Itoa(int(userId)), expireTime)
		return true, err
	}
	// SQL 查询。
	relation, err := repository.NewFollowRepo().FindRelation(userId, targetId)

	if nil != err {
		return false, err
	}
	if nil == relation {
		return false, nil
	}
	// 存在此关系，将其注入Redis中。
	go addRelationToRedis(int(userId), int(targetId))

	return true, nil
}
func addRelationToRedis(userId int, targetId int) {
	// 第一次存入时，给该key添加一个-1为key，防止脏数据的写入。当然set可以去重，直接加，便于CPU。
	redis.RdbFollowingPart.SAdd(redis.Ctx, strconv.Itoa(int(userId)), -1)
	// 将查询到的关注关系注入Redis.
	redis.RdbFollowingPart.SAdd(redis.Ctx, strconv.Itoa(int(userId)), targetId)
	// 更新过期时间。
	redis.RdbFollowingPart.Expire(redis.Ctx, strconv.Itoa(int(userId)), expireTime)
}

// GetFollowerCnt 给定当前用户id，查询其粉丝数量。
func (*FollowImpl) GetFollowerCnt(userId int64) (int64, error) {
	// 查Redis中是否已经存在。
	if cnt, err := redis.RdbFollowers.SCard(redis.Ctx, strconv.Itoa(int(userId))).Result(); cnt > 0 {
		// 更新过期时间。
		redis.RdbFollowers.Expire(redis.Ctx, strconv.Itoa(int(userId)), expireTime)
		return cnt - 1, err
	}
	// SQL中查询。
	ids, err := repository.NewFollowRepo().GetFollowersList(userId)
	if nil != err {
		return 0, err
	}
	// 将数据存入Redis.
	// 更新followers 和 followingPart
	go addFollowersToRedis(int(userId), ids)

	return int64(len(ids)), err
}
func addFollowersToRedis(userId int, ids []int64) {
	redis.RdbFollowers.SAdd(redis.Ctx, strconv.Itoa(userId), -1)
	for i, id := range ids {
		redis.RdbFollowers.SAdd(redis.Ctx, strconv.Itoa(userId), id)
		redis.RdbFollowingPart.SAdd(redis.Ctx, strconv.Itoa(int(id)), userId)
		redis.RdbFollowingPart.SAdd(redis.Ctx, strconv.Itoa(int(id)), -1)
		// 更新部分关注者的时间
		redis.RdbFollowingPart.Expire(redis.Ctx, strconv.Itoa(int(id)),
			expireTime+time.Duration((i%10)<<8))
	}
	// 更新followers的过期时间。
	redis.RdbFollowers.Expire(redis.Ctx, strconv.Itoa(userId), expireTime)

}

// GetFollowingCnt 给定当前用户id，查询其关注者数量。
func (*FollowImpl) GetFollowingCnt(userId int64) (int64, error) {
	// 查看Redis中是否有关注数。
	if cnt, err := redis.RdbFollowing.SCard(redis.Ctx, strconv.Itoa(int(userId))).Result(); cnt > 0 {
		// 更新过期时间。
		redis.RdbFollowing.Expire(redis.Ctx, strconv.Itoa(int(userId)), expireTime)
		return cnt - 1, err
	}
	// 用SQL查询。
	ids, err := repository.NewFollowRepo().GetFollowingList(userId)

	if nil != err {
		return 0, err
	}
	// 更新Redis中的followers和followPart
	go addFollowingToRedis(int(userId), ids)

	return int64(len(ids)), err
}
func addFollowingToRedis(userId int, ids []int64) {
	redis.RdbFollowing.SAdd(redis.Ctx, strconv.Itoa(userId), -1)
	for i, id := range ids {
		redis.RdbFollowing.SAdd(redis.Ctx, strconv.Itoa(userId), id)
		redis.RdbFollowingPart.SAdd(redis.Ctx, strconv.Itoa(userId), id)
		redis.RdbFollowingPart.SAdd(redis.Ctx, strconv.Itoa(userId), -1)
		// 更新过期时间
		redis.RdbFollowingPart.Expire(redis.Ctx, strconv.Itoa(userId),
			expireTime+time.Duration((i%10)<<8))
	}
	// 更新following的过期时间
	redis.RdbFollowing.Expire(redis.Ctx, strconv.Itoa(userId), expireTime)
}

// AddFollowRelation 给定当前用户和目标对象id，添加他们之间的关注关系。
func (*FollowImpl) AddFollowRelation(userId int64, targetId int64) (bool, error) {
	// 加信息打入消息队列。
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(int(userId)))
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(int(targetId)))
	rabbitmq.RmqFollowAdd.Publish(sb.String())
	// 记录日志
	util.Log.Debug("消息打入成功。")
	// 更新redis信息。
	return updateRedisWithAdd(userId, targetId)
}

// 添加关注时，设置Redis
func updateRedisWithAdd(userId int64, targetId int64) (bool, error) {
	/*
		1-Redis是否存在followers_targetId.
		2-Redis是否存在following_userId.
		3-Redis是否存在following_part_userId.
	*/
	// step1
	targetIdStr := strconv.Itoa(int(targetId))
	if cnt, _ := redis.RdbFollowers.SCard(redis.Ctx, targetIdStr).Result(); 0 != cnt {
		redis.RdbFollowers.SAdd(redis.Ctx, targetIdStr, userId)
		redis.RdbFollowers.Expire(redis.Ctx, targetIdStr, expireTime)
	}
	// step2
	followingUserIdStr := strconv.Itoa(int(userId))
	if cnt, _ := redis.RdbFollowing.SCard(redis.Ctx, followingUserIdStr).Result(); 0 != cnt {
		redis.RdbFollowing.SAdd(redis.Ctx, followingUserIdStr, targetId)
		redis.RdbFollowing.Expire(redis.Ctx, followingUserIdStr, expireTime)
	}
	// step3
	followingPartUserIdStr := followingUserIdStr
	redis.RdbFollowingPart.SAdd(redis.Ctx, followingPartUserIdStr, targetId)
	// 可能是第一次给改用户加followingPart的关注者，需要加上-1防止脏读。
	redis.RdbFollowingPart.SAdd(redis.Ctx, followingPartUserIdStr, -1)
	redis.RdbFollowingPart.Expire(redis.Ctx, followingPartUserIdStr, expireTime)
	return true, nil
}

// DeleteFollowRelation 给定当前用户和目标用户id，删除其关注关系。
func (*FollowImpl) DeleteFollowRelation(userId int64, targetId int64) (bool, error) {
	// 加信息打入消息队列。
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(int(userId)))
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(int(targetId)))
	rabbitmq.RmqFollowDel.Publish(sb.String())
	// 记录日志
	util.Log.Debug("消息打入成功。")
	// 更新redis信息。
	return updateRedisWithDel(userId, targetId)
}

// 当取关时，更新redis里的信息
func updateRedisWithDel(userId int64, targetId int64) (bool, error) {
	/*
		1-Redis是否存在followers_targetId.
		2-Redis是否存在following_userId.
		2-Redis是否存在following_part_userId.
	*/
	// step1
	targetIdStr := strconv.Itoa(int(targetId))
	if cnt, _ := redis.RdbFollowers.SCard(redis.Ctx, targetIdStr).Result(); 0 != cnt {
		redis.RdbFollowers.SRem(redis.Ctx, targetIdStr, userId)
		redis.RdbFollowers.Expire(redis.Ctx, targetIdStr, expireTime)
	}
	// step2
	followingIdStr := strconv.Itoa(int(userId))
	if cnt, _ := redis.RdbFollowing.SCard(redis.Ctx, followingIdStr).Result(); 0 != cnt {
		redis.RdbFollowing.SRem(redis.Ctx, followingIdStr, targetId)
		redis.RdbFollowing.Expire(redis.Ctx, followingIdStr, expireTime)
	}
	// step3
	followingPartUserIdStr := followingIdStr
	if cnt, _ := redis.RdbFollowingPart.Exists(redis.Ctx, followingPartUserIdStr).Result(); 0 != cnt {
		redis.RdbFollowingPart.SRem(redis.Ctx, followingPartUserIdStr, targetId)
		redis.RdbFollowingPart.Expire(redis.Ctx, followingPartUserIdStr, expireTime)
	}
	return true, nil
}

// GetFollowing 根据当前用户id来查询他的关注者列表。
func (f *FollowImpl) GetFollowing(userId int64) ([]User, error) {
	return getFollowing(userId)
}

// 设置Redis关于所有关注的信息。
func setRedisFollowing(userId int64, users []User) {
	/*
		1-设置following_userId的所有关注id。
		2-设置following_part_id关注信息。
	*/
	// 加上-1防止脏读
	followingIdStr := strconv.Itoa(int(userId))
	redis.RdbFollowing.SAdd(redis.Ctx, followingIdStr, -1)
	// 设置过期时间
	redis.RdbFollowing.Expire(redis.Ctx, followingIdStr, expireTime)
	for i, user := range users {
		redis.RdbFollowing.SAdd(redis.Ctx, followingIdStr, user.Id)

		redis.RdbFollowingPart.SAdd(redis.Ctx, followingIdStr, user.Id)
		redis.RdbFollowingPart.SAdd(redis.Ctx, followingIdStr, -1)
		// 随机设置过期时间
		redis.RdbFollowingPart.Expire(redis.Ctx, followingIdStr, expireTime+
			time.Duration((i%10)<<8))
	}
}

// 从数据库查所有关注用户信息。
func getFollowing(userId int64) ([]User, error) {
	followingLists, err := repository.NewFollowRepo().GetFollowingList(userId)
	if err != nil {
		return make([]User, 1), err
	}
	if len(followingLists) == 0 {
		return make([]User, 1), nil
	}
	users := make([]User, 0)
	// 遍历id列表，得到user对象，添加到user切片中
	for _, following := range followingLists {
		UserImpl := new(UserImpl)
		user, err := UserImpl.GetUserById(following)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	// 返回关注对象列表。
	return users, nil
}

// GetFollowers 根据当前用户id来查询他的粉丝列表。

func GetFollowers(userId int64) ([]User, error) {
	// 获取粉丝的id数组。
	ids, err := repository.NewFollowRepo().GetFollowersList(userId)
	// 查询出错
	if nil != err {
		return make([]User, 1), err
	}
	// 没得粉丝
	if len(ids) == 0 {
		return make([]User, 1), nil
	}
	users := make([]User, 0)
	// 根据每个id来查询用户信息。
	// 遍历id列表，得到user对象，添加到user切片中
	for _, follower := range ids {
		userImpl := new(UserImpl)
		user, err := userImpl.GetUserById(follower)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
func (f *FollowImpl) getUserById(ids []int64, userId int64) ([]User, error) {
	len := len(ids)
	if len > 0 {
		len -= 1
	}
	users := make([]User, len)
	var wg sync.WaitGroup
	wg.Add(len)
	i, j := 0, 0
	for ; i < len; j++ {
		// 越过-1
		if ids[j] == -1 {
			continue
		}
		//开启协程来查。
		go func(i int, idx int64) {
			defer wg.Done()
			users[i], _ = f.GetUserByIdWithCurId(idx, userId)
		}(i, ids[i])
		i++
	}
	wg.Wait()
	// 返回粉丝列表。
	return users, nil
}

// GetFollowers 根据当前用户id来查询他的粉丝列表。
func (f *FollowImpl) GetFollowers(userId int64) ([]User, error) {
	return GetFollowers(userId)
}
