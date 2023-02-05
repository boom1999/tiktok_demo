package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"tiktok_demo/middleware/jwt"
	"tiktok_demo/repository"
	"tiktok_demo/service"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserInfoResponse struct {
	Response
	User service.User `json:"user"`
}

// Register Post  route: douyin/user/register/
func Register(ctx *gin.Context) {
	userName := ctx.Query("username")
	password := ctx.Query("password")

	usi := service.UserImpl{}
	u := usi.GetTableUserByUserName(userName)

	if userName == u.Username {
		ctx.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		newUser := repository.TableUser{
			Username: userName,
			Password: jwt.PswEnCode(password),
		}
		if usi.InsertTableUser(&newUser) != true {
			log.Println("Insert Data Fail")
		} else {
			log.Println("Insert Data Success")
		}
		token, _ := jwt.GenToken(userName)
		log.Println("registered Id: ", u.Id)
		ctx.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0, StatusMsg: "OK"},
			UserId:   u.Id,
			Token:    token,
		})
	}
}

// Login Post route: douyin/user/login/
func Login(ctx *gin.Context) {
	userName := ctx.Query("username")
	password := ctx.Query("password")
	encodedPassword := jwt.PswEnCode(password)
	log.Println("encodedPassword: ", encodedPassword)

	usi := service.UserImpl{}
	u := usi.GetTableUserByUserName(userName)
	if u.Username == "" {
		ctx.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User Doesn't Exist"},
		})
	}
	if encodedPassword == u.Password {
		token, _ := jwt.GenToken(userName)
		log.Println("login Id: ", u.Id)
		ctx.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0, StatusMsg: "OK"},
			UserId:   u.Id,
			Token:    token,
		})
	} else {
		ctx.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Password Error"},
		})
	}
}

// GetUserInfo route: douyin/user/
func GetUserInfo(ctx *gin.Context) {
	userId := ctx.Query("user_id")
	id, _ := strconv.ParseInt(userId, 10, 64)
	usi := service.UserImpl{}

	u, err := usi.GetUserById(id)
	if err != nil {
		ctx.JSON(http.StatusOK, UserInfoResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User Doesn't Exist"},
		})
	} else {
		ctx.JSON(http.StatusOK, UserInfoResponse{
			Response: Response{StatusCode: 0},
			User:     u,
		})
	}
}
