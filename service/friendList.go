package service

import (
	"log"
	"sync"
)

type FriendUser struct {
	User
	Messages string `json:"message,omitempty"`
	MsgType  int64  `json:"msgType"`
}

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
func (f *FriendListImpl) GetFriendListByUserId(id int64) ([]FriendUser, error) {
	var FriendList []FriendUser
	followings, err := f.GetFollowing(id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("Get following list failed")
		return nil, err
	}
	// 遍历用户的关注列表，主要源于以下两点
	// 1.普通用户的关注和粉丝都少
	// 2.对于大v,关注很少，但粉丝很多，如果遍历粉丝列表，时间和空间复杂度将会上升
	for _, following := range followings {
		// 对该用户的关注用户，查询这个关注用户是否关注了该用户
		isFollowing, err := f.IsFollowing(following.Id, id)
		if err != nil {
			log.Println("Err:", err.Error())
			log.Println("GetIsFollowing failed")
			return nil, err
		}
		// 如果两者存在双向关注关系，即符合朋友定义，添加到该用户的 朋友列表
		if isFollowing {
			var latestMessage MessageServiceImpl
			list, err := latestMessage.GetList(following.Id, id)
			if err != nil {
				log.Println("Err:", err.Error())
				log.Println("GetList failed")
				return nil, err
			}
			usi := UserImpl{}
			frienduser, err := usi.GetUserById(following.Id)
			if err != nil {
				log.Println("Err:", err.Error())
				log.Println("GetUserById failed")
				return nil, err
			}
			//判断是接收消息还是发送消息
			var msgType int64
			if list[len(list)-1].From_user_id == following.Id {
				msgType = 0 //接收消息
			} else {
				msgType = 1
			}
			friend := FriendUser{
				User:     frienduser,
				Messages: list[len(list)-1].Content,
				MsgType:  msgType,
			}

			FriendList = append(FriendList, friend)
		}
	}
	return FriendList, nil
}
