package controller

import (
	"net/http"
	"strconv"
	"tiktok_demo/service"
	"tiktok_demo/util"

	"github.com/gin-gonic/gin"
)

type likeResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type GetFavouriteListResponse struct {
	StatusCode int32           `json:"status_code"`
	StatusMsg  string          `json:"status_msg,omitempty"`
	VideoList  []service.Video `json:"video_list,omitempty"`
}

// FavoriteAction 点赞或者取消赞操作
func FavoriteAction(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	actionType, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)

	likeService := new(service.LikeServiceImpl)
	err := likeService.FavouriteAction(userId, videoId, int32(actionType))
	// 官方文档有矛盾 1-成功 0-失败
	if err == nil {
		util.Log.Debug("favourite action success")
		c.JSON(http.StatusOK, likeResponse{
			StatusCode: 1,
			StatusMsg:  "favourite action success",
		})
	} else {
		util.Log.Error("favourite action fail" + err.Error())
		c.JSON(http.StatusOK, likeResponse{
			StatusCode: 0,
			StatusMsg:  "favourite action fail",
		})
	}
}

// GetFavouriteList 获取点赞列表
func GetFavouriteList(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	curId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)

	likeService := new(service.LikeServiceImpl)
	videos, err := likeService.GetFavouriteList(userId, curId)
	if err == nil {
		util.Log.Debug("get favouriteList success")
		c.JSON(http.StatusOK, GetFavouriteListResponse{
			StatusCode: 0,
			StatusMsg:  "get favouriteList success",
			VideoList:  videos,
		})
	} else {
		util.Log.Error("get favouriteList fail" + err.Error())
		c.JSON(http.StatusOK, GetFavouriteListResponse{
			StatusCode: 1,
			StatusMsg:  "get favouriteList fail",
		})
	}
}
