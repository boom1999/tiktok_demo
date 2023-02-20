package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"tiktok_demo/service"
	"tiktok_demo/util"

	"github.com/gin-gonic/gin"
)

type RelationActionResp struct {
	Response
}

type FollowingResp struct {
	Response
	UserList []service.User `json:"user_list,omitempty"`
}

type FollowerResp struct {
	Response
	UserList []service.User `json:"user_list,omitempty"`
}

// RelationAction Handle follow and unfollow requests
func RelationAction(ctx *gin.Context) {
	// Get login info
	userId, err1 := strconv.ParseInt(ctx.GetString("userId"), 10, 64)
	targetId, err2 := strconv.ParseInt(ctx.Query("to_user_id"), 10, 64)
	actionType, err3 := strconv.ParseInt(ctx.Query("action_type"), 10, 64)
	// Check if Info is obtained successfully
	fmt.Println("userId, targetId, actionType: ", userId, targetId, actionType)
	// 1-follow, 2-deleteFollow
	if nil != err1 || nil != err2 || nil != err3 || actionType < 1 || actionType > 2 {
		util.Log.Error("get ctx failed")
		ctx.JSON(http.StatusOK, RelationActionResp{
			Response{
				StatusCode: -1,
				StatusMsg:  "RelationAction format error",
			},
		})
		return
	}
	fsi := service.NewFSIInstance()
	switch {
	// follow
	case 1 == actionType:
		go func() {
			_, err := fsi.AddFollowRelation(userId, targetId)
			if err != nil {
				util.Log.Error("Follow failed")
			} else {
				util.Log.Error("Follow succeed")
			}
		}()
	// delete follow
	case 2 == actionType:
		go func() {
			_, err := fsi.DeleteFollowRelation(userId, targetId)
			if err != nil {
				util.Log.Error("Delete follow failed")

			} else {
				util.Log.Error("Delete follow succeed")
			}
		}()
	}
	ctx.JSON(http.StatusOK, RelationActionResp{
		Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
	})
}

// GetFollowingList 处理获取关注列表请求
func GetFollowingList(ctx *gin.Context) {
	userId, err := strconv.ParseInt(ctx.Query("user_id"), 10, 64)
	if nil != err {
		ctx.JSON(http.StatusOK, FollowingResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "userId format error",
			},
			UserList: nil,
		})
		return
	}
	fsi := service.NewFSIInstance()
	users, err := fsi.GetFollowing(userId)
	// Get followingList failed
	if err != nil {
		ctx.JSON(http.StatusOK, FollowingResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "Get followingList failed",
			},
			UserList: nil,
		})
		return
	}
	// Get followingList succeed
	util.Log.Debug("Get followingList succeed。")
	ctx.JSON(http.StatusOK, FollowingResp{
		UserList: users,
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
	})
}

// GetFollowersList 处理获取粉丝列表请求
func GetFollowersList(ctx *gin.Context) {
	userId, err := strconv.ParseInt(ctx.Query("user_id"), 10, 64)
	if nil != err {
		ctx.JSON(http.StatusOK, FollowingResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "userId format error",
			},
			UserList: nil,
		})
		return
	}
	fsi := service.NewFSIInstance()
	users, err := fsi.GetFollowers(userId)
	// Get followerList failed
	if err != nil {
		ctx.JSON(http.StatusOK, FollowingResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "Get followerList failed",
			},
			UserList: nil,
		})
		return
	}
	// Get followerList succeed
	util.Log.Debug("Get followerList succeed。")
	ctx.JSON(http.StatusOK, FollowingResp{
		UserList: users,
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
	})
}
