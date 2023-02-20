package repository

import (
	"errors"
	"tiktok_demo/config"
	"tiktok_demo/util"
)

// Like 表的结构
type Like struct {
	Id      int64 `gorm:"column:id;not null;type:bigint(20) primary key auto_increment"` // 自增id
	UserId  int64 `gorm:"column:user_id;not null;type:bigint(20)"`                       // 点赞的用户id
	VideoId int64 `gorm:"column:video_id;not null;type:bigint(20)"`                      // 视频id
	Cancel  int8  `gorm:"column:cancel;not null;default:0;type:tinyint(4)"`              // 是否点赞，0为点赞，1为取消赞
}

// TableName 修改表名映射
func (Like) TableName() string {
	return "likes"
}

// GetLikeUserIdList 根据 videoId 获取点赞用户的 userId
func GetLikeUserIdList(videoId int64) ([]int64, error) {
	var likeUserIdList []int64
	// 查询 likes 表对应视频id 的点赞用户
	err := DB.Model(Like{}).Where(map[string]interface{}{"video_id": videoId, "cancel": config.IsLike}).
		Pluck("user_id", &likeUserIdList).Error
	if err != nil {
		util.Log.Error("get likeUserIdList failed" + err.Error())
		return nil, errors.New("get likeUserIdList failed")
	} else {
		return likeUserIdList, nil
	}
}

// UpdateLike 根据 userId、videoId、actionType 进行点赞或者取消赞
func UpdateLike(userId int64, videoId int64, actionType int32) error {
	// 更新当前用户观看视频的点赞状态
	err := DB.Model(Like{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).
		Update("cancel", actionType).Error
	if err != nil {
		util.Log.Error("update data fail" + err.Error())
		return errors.New("update data fail")
	}
	return nil
}

// InsertLike 插入点赞数据
func InsertLike(likeData Like) error {
	// 创建点赞数据，默认为点赞，cancel为 0
	err := DB.Model(Like{}).Create(&likeData).Error
	if err != nil {
		util.Log.Error("insert data fail" + err.Error())
		return errors.New("insert data fail")
	}
	return nil
}

// GetLikeInfo 根据 userId、videoId 查询点赞信息
func GetLikeInfo(userId int64, videoId int64) (Like, error) {
	var likeInfo Like
	// 根据 userid、videoId 查询是否有该条信息
	err := DB.Model(Like{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).
		First(&likeInfo).Error
	if err != nil {
		if err.Error() == "record not found" {
			util.Log.Error("can't find data" + err.Error())
			return Like{}, nil
		} else {
			util.Log.Error("get likeInfo failed" + err.Error())
			return likeInfo, errors.New("get likeInfo failed")
		}
	}
	return likeInfo, nil
}

// GetLikeVideoIdList 根据 userId 查询所有点赞的 videoId
func GetLikeVideoIdList(userId int64) ([]int64, error) {
	var likeVideoIdList []int64
	err := DB.Model(Like{}).Where(map[string]interface{}{"user_id": userId, "cancel": config.IsLike}).
		Pluck("video_id", &likeVideoIdList).Error
	if err != nil {
		if err.Error() == "record not found" {
			util.Log.Error("there are no likeVideoId" + err.Error())
			return likeVideoIdList, nil
		} else {
			util.Log.Error("get likeVideoIdList failed" + err.Error())
			return likeVideoIdList, errors.New("get likeVideoIdList failed")
		}
	}
	return likeVideoIdList, nil
}
