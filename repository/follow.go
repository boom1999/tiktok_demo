package repository

import (
	"sync"
	"tiktok_demo/util"
)

type Follow struct {
	Id         int64 `gorm:"column:id;not null;type:bigint(20) primary key auto_increment"`
	UserId     int64 `gorm:"column:user_id;not null;type:bigint(20)"`
	FollowerId int64 `gorm:"column:follower_id;not null;type:bigint(20)"`
	Cancel     int8  `gorm:"column:cancel;not null;default:0;type:tinyint(4)"`
}

func (Follow) TableName() string {
	return "follows"
}

type FollowCURD struct {
}

var (
	followCURD *FollowCURD
	followOnce sync.Once
)

// NewFollowRepo Init and return followRepo object
func NewFollowRepo() *FollowCURD {
	followOnce.Do(func() {
		followCURD = &FollowCURD{}
	})
	return followCURD
}

/*
	CURD
*/

// FindRelation 查看userId是否关注targetId
func (*FollowCURD) FindRelation(userId int64, targetId int64) (*Follow, error) {
	follow := Follow{}
	if err := DB.
		Where("user_id = ? and follower_id = ? and cancel = ?", targetId, userId, 0).
		Take(&follow).Error; err != nil {
		if "record not found" == err.Error() {
			return nil, nil
		}
		util.Log.Error(err.Error())
		return nil, err
	}

	return &follow, nil
}

// GetFollowerCnt 查看userId的粉丝人数
func (*FollowCURD) GetFollowerCnt(userId int64) (int64, error) {
	var cnt int64
	if err := DB.Model(Follow{}).Where("user_id = ? and cancel = ?", userId, 0).Count(&cnt).Error; err != nil {
		util.Log.Error(err.Error())
		return 0, err
	}
	return cnt, nil
}

// GetFollowingCnt 查看userId的关注列表人数
func (*FollowCURD) GetFollowingCnt(userId int64) (int64, error) {
	var cnt int64
	if err := DB.Model(Follow{}).Where("follower_id = ? and cancel = ?", userId, 0).Count(&cnt).Error; err != nil {
		util.Log.Error(err.Error())
		return 0, err
	}
	return cnt, nil
}

// InsertFollowRelation 给userId插入关注关系
func (*FollowCURD) InsertFollowRelation(userId int64, targetId int64) (bool, error) {
	followInsert := Follow{
		UserId:     userId,
		FollowerId: targetId,
		Cancel:     0,
	}
	if err := DB.Select("UserId", "FollowerId", "Cancel").Create(&followInsert).Error; nil != err {
		util.Log.Error(err.Error())
		return false, err
	}
	return true, nil
}

// FindEverFollowing userId是否关注过targetId
func (*FollowCURD) FindEverFollowing(userId int64, targetId int64) (*Follow, error) {
	follow := Follow{}
	if err := DB.
		Where("user_id = ? and follower_id = ?", targetId, userId).
		Where("cancel = ? or cancel = ?", 0, 1).
		Take(&follow).Error; err != nil {
		if "record not found" == err.Error() {
			return nil, nil
		}
		util.Log.Error(err.Error())
		return nil, err
	}

	return &follow, nil
}

// UpdateFollowRelation 更新关注信息
func (*FollowCURD) UpdateFollowRelation(userId int64, targetId int64, cancel int8) (bool, error) {
	if err := DB.Model(Follow{}).
		Where("user_id = ? and follower_id = ?", userId, targetId).
		Update("cancel", cancel).Error; nil != err {
		util.Log.Error(err.Error())
		return false, err
	}
	return true, nil
}

// GetFollowingList 查询userId的关注列表，ids
func (*FollowCURD) GetFollowingList(userId int64) ([]int64, error) {
	var ids []int64
	if err := DB.
		Model(Follow{}).Where("follower_id = ?", userId).
		Pluck("user_id", &ids).Error; nil != err {
		if "record not found" == err.Error() {
			return nil, nil
		}
		util.Log.Error(err.Error())
		return nil, err
	}
	return ids, nil
}

// GetFollowersList  查询userId的粉丝列表，ids
func (*FollowCURD) GetFollowersList(userId int64) ([]int64, error) {
	var ids []int64
	if err := DB.
		Model(Follow{}).
		Where("user_id = ? and cancel = ?", userId, 0).
		Pluck("follower_id", &ids).Error; nil != err {
		if "record not found" == err.Error() {
			return nil, nil
		}
		util.Log.Error(err.Error())
		return nil, err
	}
	return ids, nil
}
