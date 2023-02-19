package repository

import "tiktok_demo/util"

// TableUser <---> User struct in mysql
type TableUser struct {
	Id       int64  `gorm:"column:id;not null;type:bigint(20) primary key auto_increment"`
	Username string `gorm:"column:username;not null;type:varchar(255)"`
	Password string `gorm:"column:password;not null;type:varchar(255)"`
}

// TableName 修改映射名
func (tableUser TableUser) TableName() string {
	return "users"
}

// GetTableUserList 获取所有TableUser对象
func GetTableUserList() ([]TableUser, error) {
	var tableUsers []TableUser
	err := DB.Find(&tableUsers).Error
	if err != nil {
		util.Log.Error(err.Error())
		return tableUsers, err
	}
	return tableUsers, nil
}

// GetTableUserByUserName 根据userName获取TableUser对象
func GetTableUserByUserName(userName string) (TableUser, error) {
	tableUser := TableUser{}
	err := DB.Where("username = ?", userName).First(&tableUser).Error
	if err != nil {
		util.Log.Error(err.Error())
		return tableUser, err
	}
	return tableUser, nil
}

// GetTableUserByUserId 根据userId获取TableUser对象
func GetTableUserByUserId(userId int64) (TableUser, error) {
	tableUser := TableUser{}
	err := DB.Where("id = ?", userId).First(&tableUser).Error
	if err != nil {
		util.Log.Error(err.Error())
		return tableUser, err
	}
	return tableUser, nil
}

// InsertTableUser 插入新用户
func InsertTableUser(newUser *TableUser) bool {
	err := DB.Create(&newUser).Error
	if err != nil {
		util.Log.Error(err.Error())
		return false
	}
	return true
}
