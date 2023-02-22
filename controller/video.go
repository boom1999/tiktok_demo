package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"tiktok_demo/service"
	"tiktok_demo/util"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []service.Video `json:"video_list"`
	NextTime  int64           `json:"next_time,omitempty"`
}

type VideoListResponse struct {
	Response
	VideoList []service.Video `json:"video_list"`
}

// Feed /feed/
func Feed(c *gin.Context) {
	inputTime := c.Query("latest_time")
	util.Log.Debug("debug", zap.String("acquired latest_time", inputTime))
	var lastTime time.Time
	if (inputTime != "0") && (inputTime != "") {
		me, _ := strconv.ParseInt(inputTime, 10, 64)
		lastTime = time.Unix(me, 0)
	} else {
		lastTime = time.Now()
	}
	if lastTime.Year() > time.Now().Year() {
		lastTime = time.Now()
	}
	util.Log.Debug("debug", zap.Time("acquired timestamp", lastTime))
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	util.Log.Debug("debug", zap.Int64("acquired userId", userId))
	videoService := GetVideo()
	feed, nextTime, err := videoService.Feed(lastTime, userId)
	if err != nil {
		util.Log.Error("call videoService.Feed(lastTime, userId) failed" + err.Error())
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "获取视频流失败"},
		})
		return
	}
	util.Log.Debug("call videoService.Feed(lastTime, userId) success")
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "获取视频流成功"},
		VideoList: feed,
		NextTime:  nextTime.Unix(),
	})
}

// Publish /publish/action/
func Publish(c *gin.Context) {
	data, err := c.FormFile("data")
	if err != nil {
		util.Log.Error("acquired video streaming failed" + err.Error())
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	title, _ := c.GetPostForm("title")
	if err != nil {
		util.Log.Error("acquired video streaming failed" + err.Error())
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	util.Log.Debug("debug", zap.Int64("acquired userId", userId))
	//title := c.PostForm("title")
	util.Log.Debug("debug", zap.String("acquired video title", title))

	videoService := GetVideo()
	err = videoService.Publish(data, userId, title, c)
	if err != nil {
		util.Log.Error("call videoService.Publish(data, userId, title) failed" + err.Error())
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	util.Log.Debug("call videoService.Publish(data, userId) success")

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "uploaded successfully",
	})
}

// PublishList /publish/list/
func PublishList(c *gin.Context) {
	user_Id, _ := c.GetQuery("user_id")
	userId, _ := strconv.ParseInt(user_Id, 10, 64)
	util.Log.Debug("debug", zap.Int64("acquired userId", userId))
	curId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	util.Log.Debug("debug", zap.Int64("acquired current userId", curId))
	videoService := GetVideo()
	list, err := videoService.List(userId, curId)
	if err != nil {
		util.Log.Error("call videoService.List(userId, curId) failed" + err.Error())
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "获取视频列表失败"},
		})
		return
	}
	util.Log.Debug("call videoService.List(userId, curId) success")
	c.JSON(http.StatusOK, VideoListResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "获取视频列表成功"},
		VideoList: list,
	})
}

// GetVideo 拼装videoService
func GetVideo() service.VideoServiceImpl {
	var userService service.UserImpl
	//var followService service.FollowServiceImp
	var videoService service.VideoServiceImpl
	var likeService service.LikeServiceImpl
	var commentService service.CommentServiceImpl
	//userService.FollowService = &followService
	//userService.LikeService = &likeService
	//followService.UserService = &userService
	//likeService.VideoService = &videoService
	commentService.UserService = &userService
	videoService.CommentService = &commentService
	videoService.LikeService = &likeService
	videoService.UserService = &userService
	return videoService
}
