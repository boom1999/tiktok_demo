package repository

import (
	"errors"
	"tiktok_demo/util"
)

// TableMessage 消息-数据库中的结构体
type TableMessage struct {
	Id         int64  `gorm:"column:id;not null;type:bigint(20) primary key auto_increment"`
	FromUserId int64  `gorm:"column:from_user_id;not null;type:bigint(20)"`
	ToUserId   int64  `gorm:"column:to_user_id;not null;type:bigint(20)"`
	Content    string `gorm:"column:content;not null;type:varchar(255)"`
	CreateTime string `gorm:"column:Create_time;not null;type:varchar(255)"`
}

// TableName 修改映射名
func (tableMessage TableMessage) TableName() string {
	return "messages"
}

func InsertTableMessage(tableMessage TableMessage) error {
	err := DB.Model(TableMessage{}).Create(&tableMessage).Error
	if err != nil {
		// 函数返回提示错误信息
		util.Log.Error("InsertTableMessage: return insert Message failed" + errors.New("insert message failed").Error())
		return errors.New("insert message failed")
	}
	util.Log.Debug("InsertTableMessage: return success")
	return nil
}

func GetMessageList(userId int64, toUserId int64) ([]TableMessage, []TableMessage, error) {
	// userId 发送给 toUserId  && toUserId 发送给 userId
	var tableMessageListFrom []TableMessage
	var tableMessageListTo []TableMessage

	resultFrom := DB.Model(TableMessage{}).Where(map[string]interface{}{"From_user_id": userId, "to_user_id": toUserId}).
		Order("Create_time desc").Find(&tableMessageListFrom)
	resultTo := DB.Model(TableMessage{}).Where(map[string]interface{}{"From_user_id": toUserId, "to_user_id": userId}).
		Order("Create_time desc").Find(&tableMessageListTo)
	if resultFrom.Error != nil && resultTo.Error != nil {
		util.Log.Error("GetMessageList false" + errors.New("get comment list failed").Error())
		return tableMessageListFrom, tableMessageListTo, errors.New("get comment list failed")
	}
	util.Log.Debug("GetMessageList successful")
	return tableMessageListFrom, tableMessageListTo, nil
}
