package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"tiktok_demo/service"
)

type FriendListResponse struct {
	Response
	FriendList []service.FriendUser `json:"user_list,omitempty"`
}

// GetFriendList route: /douyin/relation/friend/list/
func GetFriendList(ctx *gin.Context) {
	userId := ctx.Query("user_id")
	id, err := strconv.ParseInt(userId, 10, 64)
	if nil != err {
		ctx.JSON(http.StatusOK, FriendListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "userId format error",
			},
			FriendList: nil,
		})
		return
	}
	fsi := service.NewFLInstance()
	u, err := fsi.GetFriendListByUserId(id)
	if err != nil {
		ctx.JSON(http.StatusOK, FriendListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "FriendList error"},
		})
	} else {
		ctx.JSON(http.StatusOK, FriendListResponse{
			Response:   Response{StatusCode: 0, StatusMsg: "ok"},
			FriendList: u,
		})
	}
}
