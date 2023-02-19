package service

import (
	"log"
	"strconv"
	"tiktok_demo/repository"
)

type User struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	FollowCount     int64  `json:"follow_count,omitempty"`
	FollowerCount   int64  `json:"follower_count,omitempty"`
	IsFollow        bool   `json:"is_follow"`
	TotalFavorite   int64  `json:"total_favorited,omitempty"`
	FavoriteCount   int64  `json:"favorite_count,omitempty"`
	AvatarUrl       string `json:"avatar,omitempty"`
	BackgroundImage string `json:"background_image,omitempty"`
	WorkCount       int64  `json:"work_count,omitempty"`
	Signature       string `json:"signature,omitempty"`
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
	followingcnt, err := usi.GetFollowingCnt(id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("GetFollowingCnt failed")
		return User{}, err
	}
	followercnt, err := usi.GetFollowerCnt(id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("followercnt failed")
		return User{}, err
	}
	totalfavouritecount, err := usi.TotalFavourite(id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("Get total favourite count failed")
		return User{}, err
	}
	favouritevideocount, err := usi.FavouriteVideoCount(id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("Get favourite video count failed")
		return User{}, err
	}

	user = User{
		Id:              id,
		Name:            tableUser.Username,
		FollowCount:     followingcnt,
		FollowerCount:   followercnt,
		IsFollow:        false,
		TotalFavorite:   totalfavouritecount,
		FavoriteCount:   favouritevideocount,
		AvatarUrl:       AvatarById(id),
		BackgroundImage: AvatarById(id),
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
	followingcnt, err := usi.GetFollowingCnt(id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("GetFollowingCnt failed")
		return User{}, err
	}
	followercnt, err := usi.GetFollowerCnt(id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("followercnt failed")
		return User{}, err
	}
	isfollowing, err := usi.IsFollowing(curId, id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("GetIsFollowing failed")
		return User{}, err
	}
	totalfavouritecount, err := usi.TotalFavourite(id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("Get total favourite count failed")
		return User{}, err
	}
	favouritevideocount, err := usi.FavouriteVideoCount(id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("Get favourite video count failed")
		return User{}, err
	}

	user = User{
		Id:              id,
		Name:            tableUser.Username,
		FollowCount:     followingcnt,
		FollowerCount:   followercnt,
		IsFollow:        isfollowing,
		TotalFavorite:   totalfavouritecount,
		FavoriteCount:   favouritevideocount,
		AvatarUrl:       AvatarById(id),
		BackgroundImage: AvatarById(id),
	}
	return user, nil
}

func AvatarById(id int64) string {
	return "https://api.multiavatar.com/" + strconv.FormatInt(id, 10) + ".png?apikey=uRiGCxXZwPK9h4"
}
