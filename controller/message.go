package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok_demo/config"
	"tiktok_demo/service"
	"tiktok_demo/util"
	"time"
)

// MessageChatResponse
// 消息记录返回参数
type MessageChatResponse struct {
	StatusCode  int32             `json:"status_code"`
	StatusMsg   string            `json:"status_msg,omitempty"`
	MessageList []service.Message `json:"message_list"`
}

// MessageActionResponse
// 发送消息返回参数
type MessageActionResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

// MessageAction
// 发送消息 /message/action/
func MessageAction(c *gin.Context) {
	//获取 From_user_id
	id, _ := c.Get("userId")
	userid, _ := id.(string)
	fromUserId, err := strconv.ParseInt(userid, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, MessageActionResponse{
			StatusCode: -1,
			StatusMsg:  "message fromUserId json invalid",
		})
		util.Log.Error("messageController-MessageAction: fromUserId json invalid" + err.Error())
		return
	}
	// 获取 to_user_id
	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, MessageActionResponse{
			StatusCode: -1,
			StatusMsg:  "message toUserId json invalid",
		})
		util.Log.Error("messageController-MessageAction: toUserId json invalid" + err.Error())
		return
	}
	// 获取 actionType
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, MessageActionResponse{
			StatusCode: -1,
			StatusMsg:  "message actionType json invalid",
		})
		util.Log.Error("messageController-MessageAction: actionType json invalid" + err.Error())
		return
	}
	// 获取 content
	content := c.Query("content")
	if err != nil {
		c.JSON(http.StatusOK, MessageActionResponse{
			StatusCode: -1,
			StatusMsg:  "message content json invalid",
		})
		util.Log.Error("messageController-MessageAction: content json invalid" + err.Error())
		return
	}
	//调用service层发送消息函数
	messageService := new(service.MessageServiceImpl)
	if actionType == 1 {
		// 发送消息数据准备
		var message service.Message
		message.From_user_id = fromUserId
		message.To_user_id = toUserId
		message.Content = content
		timeNow := time.Now()
		message.Create_time = timeNow.Format(config.DateTime)
		// 发送消息
		err := messageService.Send(message)
		if err != nil {
			c.JSON(http.StatusOK, MessageActionResponse{
				StatusCode: -1,
				StatusMsg:  "send message false",
			})
			util.Log.Error("messageController-MessageAction:" + err.Error())
			return
		}
		// 发送消息成功
		c.JSON(http.StatusOK, MessageActionResponse{
			StatusCode: 0,
			StatusMsg:  "send message successful",
		})
		util.Log.Debug("send message successful")
		return

	} else {
		c.JSON(http.StatusOK, MessageActionResponse{
			StatusCode: -1,
			StatusMsg:  "message actionType json invalid",
		})
		util.Log.Debug("messageController-MessageAction: actionType json invalid")
		return
	}
}

// MessageChat
// 获取聊天记录 /message/chat/
func MessageChat(c *gin.Context) {
	//获取userId
	id, _ := c.Get("userId")
	userid, _ := id.(string)
	userId, err := strconv.ParseInt(userid, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, MessageChatResponse{
			StatusCode:  -1,
			StatusMsg:   "message userId json invalid",
			MessageList: make([]service.Message, 0),
		})
		util.Log.Error("messageController-MessageChat: userId json invalid" + err.Error())
		return
	}
	// 获取 to_user_id
	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, MessageChatResponse{
			StatusCode:  -1,
			StatusMsg:   "message to_user_id json invalid",
			MessageList: make([]service.Message, 0),
		})
		util.Log.Error("messageController-MessageChat: to_user_id json invalid" + err.Error())
		return
	}
	// 调用 service 层消息记录函数
	messageService := new(service.MessageServiceImpl)
	message_list, err := messageService.GetList(userId, toUserId)
	if err != nil {
		c.JSON(http.StatusOK, MessageChatResponse{
			StatusCode:  -1,
			StatusMsg:   "message GetList false",
			MessageList: make([]service.Message, 0),
		})
		util.Log.Error("messageController-MessageChat: message GetList false" + err.Error())
		return
	}
	// 获取聊天记录成功
	c.JSON(http.StatusOK, MessageChatResponse{
		StatusCode:  0,
		StatusMsg:   "get messageList success",
		MessageList: message_list,
	})
	util.Log.Debug("messageController-MessageChat: return success") //成功返回列表
	return
}
