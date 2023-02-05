package service

import (
	"log"
	"tiktok_demo/repository"
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
		log.Println("Err:", err.Error())
		return tableUsers
	}
	return tableUsers
}

func (usi *UserImpl) GetTableUserByUserName(userName string) repository.TableUser {
	tableUser, err := repository.GetTableUserByUserName(userName)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return tableUser
	}
	log.Println("Query User Success")
	return tableUser
}

func (usi *UserImpl) GetTableUserByUserId(userId int64) repository.TableUser {
	tableUser, err := repository.GetTableUserByUserId(userId)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return tableUser
	}
	log.Println("Query User Success")
	return tableUser
}

func (usi *UserImpl) InsertTableUser(newUser *repository.TableUser) bool {
	flag := repository.InsertTableUser(newUser)
	if flag == false {
		log.Println("failed insert")
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
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return user, err
	} else {
		log.Println("Query User Success")
	}
	// .... Else 5 items needed to add
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
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return user, err
	} else {
		log.Println("Query User Success")
	}
	// .... Else 5 items needed to add
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
