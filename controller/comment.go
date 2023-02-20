package controller

import (
	"log"
	"net/http"
	"strconv"
	"tiktok_demo/config"
	"tiktok_demo/repository"
	"tiktok_demo/service"
	"tiktok_demo/util"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommentListResponse
type CommentListResponse struct {
	StatusCode  int32                 `json:"status_code"`
	StatusMsg   string                `json:"status_msg,omitempty"`
	CommentList []service.CommentInfo `json:"comment_list,omitempty"`
}

// CommentActionResponse
type CommentActionResponse struct {
	StatusCode int32               `json:"status_code"`
	StatusMsg  string              `json:"status_msg,omitempty"`
	Comment    service.CommentInfo `json:"comment"`
}

// CommentAction comment/action/
func CommentAction(c *gin.Context) {
	util.Log.Debug("CommentController-Comment_Action: running")
	//getuserId
	id, _ := c.Get("userId")
	userid, _ := id.(string)
	userId, err := strconv.ParseInt(userid, 10, 64)

	log.Printf("err:%v", err)
	log.Printf("userId:%v", userId)
	//error
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "comment userId json invalid",
		})
		util.Log.Error("CommentController-Comment_Action: return comment userId json invalid" + err.Error()) //函数返回userId无效
		return
	}

	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	//error
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "comment videoId json invalid",
		})
		util.Log.Error("CommentController-Comment_Action: return comment videoId json invalid" + err.Error()) //函数返回视频id无效
		return
	}
	//actionType
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 32)
	//error
	if err != nil || actionType < 1 || actionType > 2 {
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "comment actionType json invalid",
		})
		util.Log.Error("CommentController-Comment_Action: return actionType json invalid" + err.Error()) //评论类型数据无效

		return
	}

	commentService := new(service.CommentServiceImpl)
	if actionType == 1 { //actionType为1
		content := c.Query("comment_text")
		var sendComment repository.Comment
		sendComment.UserId = userId
		sendComment.VideoId = videoId
		sendComment.CommentText = content
		timeNow := time.Now()
		sendComment.CreateDate = timeNow.Format(config.DateTime)
		//sendCommen
		commentInfo, err := commentService.Send(sendComment)
		//sendCommen failed
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "send comment failed",
			})
			util.Log.Error("CommentController-Comment_Action: return send comment failed" + err.Error()) //发表失败
			return
		}

		//send comment success:
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: 0,
			StatusMsg:  "send comment success",
			Comment:    commentInfo,
		})
		util.Log.Debug("CommentController-Comment_Action: return Send success") //发表评论成功，返回正确信息
		return
	} else { //delete
		commentId, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "delete commentId invalid",
			})
			util.Log.Error("CommentController-Comment_Action: return commentId invalid") //评论id格式错误
			return
		}
		err = commentService.DelComment(commentId)
		if err != nil { //delete comment failed
			str := err.Error()
			c.JSON(http.StatusOK, CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  str,
			})
			util.Log.Error("CommentController-Comment_Action: return delete comment failed") //删除失败
			return
		}
		//删除评论成功
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: 0,
			StatusMsg:  "delete comment success",
		})

		util.Log.Debug("CommentController-Comment_Action: return delete success") //函数执行成功，返回正确信息
		return
	}
}

// CommentList
func CommentList(c *gin.Context) {
	util.Log.Debug("CommentController-Comment_List: running") //函数已运行
	//获取userId
	id, _ := c.Get("userId")
	userid, _ := id.(string)
	userId, err := strconv.ParseInt(userid, 10, 64)

	//获取videoId
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	//错误处理
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: -1,
			StatusMsg:  "comment videoId json invalid",
		})
		util.Log.Error("CommentController-Comment_List: return videoId json invalid") //视频id格式有误
		return
	}
	util.Log.Debug("debug", zap.Int64("videoId", videoId))

	commentService := new(service.CommentServiceImpl)
	commentList, err := commentService.GetList(videoId, userId)
	if err != nil { //return list false
		c.JSON(http.StatusOK, CommentListResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		util.Log.Error("CommentController-Comment_List: return list false") //查询列表失败
		return
	}

	//获取评论列表成功
	c.JSON(http.StatusOK, CommentListResponse{
		StatusCode:  0,
		StatusMsg:   "get comment list success",
		CommentList: commentList,
	})
	util.Log.Debug("CommentController-Comment_List: return success") //成功返回列表
	return
}
