package service

import (
	"log"
	"sync"
)

type FriendListImpl struct {
	FollowImpl
}

var (
	friendListImpl *FriendListImpl
	friendListOnce sync.Once
)

// NewFSIInstance 生成并返回FollowServiceImp结构体单例变量。
func NewFLInstance() *FriendListImpl {
	friendListOnce.Do(
		func() {
			friendListImpl = &FriendListImpl{}
		})
	return friendListImpl
}

// 根据用户id 查询 朋友列表
func (f *FriendListImpl) GetFriendListByUserId(id int64) ([]User, error) {
	var FriendList []User
	following, err := f.GetFollowing(id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("Get following list failed")
		return nil, err
	}
	// 遍历用户的关注列表，主要源于以下两点
	// 1.普通用户的关注和粉丝都少
	// 2.对于大v,关注很少，但粉丝很多，如果遍历粉丝列表，时间和空间复杂度将会上升
	for _, following := range following {
		// 对该用户的关注用户，查询这个关注用户是否关注了该用户
		isFollowing, err := f.IsFollowing(following.Id, id)
		if err != nil {
			log.Println("Err:", err.Error())
			log.Println("GetIsFollowing failed")
			return nil, err
		}
		// 如果两者存在双向关注关系，即符合朋友定义，添加到该用户的 朋友列表
		if isFollowing {
			FriendList = append(FriendList, following)
		}
	}
	return FriendList, nil
}
