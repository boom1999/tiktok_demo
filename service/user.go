package service

import (
	"tiktok_demo/repository"
	"tiktok_demo/util"
)

type User struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
	TotalFavorite int64  `json:"total_favorite,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
}

type UserImpl struct {
	FollowService
	LikeService
}

type UserService interface {
	GetTableUserList() []repository.TableUser

	GetTableUserByUserName(userName string) repository.TableUser

	GetTableUserByUserId(userId int64) repository.TableUser

	InsertTableUser(newUser *repository.TableUser) bool

	GetUserById(id int64) (User, error)

	GetUserByIdWithCurId(id int64, curId int64) (User, error)
}

func (usi *UserImpl) GetTableUserList() []repository.TableUser {
	tableUsers, err := repository.GetTableUserList()
	if err != nil {
		util.Log.Error("Err:" + err.Error())
		return tableUsers
	}
	return tableUsers
}

func (usi *UserImpl) GetTableUserByUserName(userName string) repository.TableUser {
	tableUser, err := repository.GetTableUserByUserName(userName)
	if err != nil {
		util.Log.Error("Err:" + err.Error())
		util.Log.Error("User Not Found")
		return tableUser
	}
	util.Log.Debug("Query User Success")
	return tableUser
}

func (usi *UserImpl) GetTableUserByUserId(userId int64) repository.TableUser {
	tableUser, err := repository.GetTableUserByUserId(userId)
	if err != nil {
		util.Log.Error("Err:" + err.Error())
		util.Log.Error("User Not Found")
		return tableUser
	}
	util.Log.Debug("Query User Success")
	return tableUser
}

func (usi *UserImpl) InsertTableUser(newUser *repository.TableUser) bool {
	flag := repository.InsertTableUser(newUser)
	if flag == false {
		util.Log.Error("failed insert")
		return false
	}
	return true
}

func (usi *UserImpl) GetUserById(id int64) (User, error) {
	user := User{
		Id:            0,
		Name:          "",
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
		TotalFavorite: 0,
		FavoriteCount: 0,
	}
	tableUser, err := repository.GetTableUserByUserId(id)
	if err != nil {
		util.Log.Error("Err:" + err.Error())
		util.Log.Error("User Not Found")
		return user, err
	} else {
		util.Log.Debug("Query User Success")
	}
	// TODO Else 5 items needed to add
	user = User{
		Id:            id,
		Name:          tableUser.Username,
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
		TotalFavorite: 0,
		FavoriteCount: 0,
	}
	return user, nil
}

func (usi *UserImpl) GetUserByIdWithCurId(id int64, curId int64) (User, error) {
	user := User{
		Id:            0,
		Name:          "",
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
		TotalFavorite: 0,
		FavoriteCount: 0,
	}
	tableUser, err := repository.GetTableUserByUserId(id)
	if err != nil {
		util.Log.Error("Err:" + err.Error())
		util.Log.Error("User Not Found")
		return user, err
	} else {
		util.Log.Debug("Query User Success")
	}
	// TODO Else 5 items needed to add
	user = User{
		Id:            id,
		Name:          tableUser.Username,
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
		TotalFavorite: 0,
		FavoriteCount: 0,
	}
	return user, nil
}
