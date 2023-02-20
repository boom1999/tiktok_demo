package service

import (
	"tiktok_demo/repository"
	"tiktok_demo/util"
)

type MessageServiceImpl struct {
	MessageService
}

type MessageService interface {
	Send(comment repository.Comment) (CommentInfo, error)
	GetList(userId int64, toUserId int64) ([]Message, error)
}

// Message 发送消息结构体
type Message struct {
	Id           int64  `json:"id"`
	From_user_id int64  `json:"from_user_id"`
	To_user_id   int64  `json:"to_user_id"`
	Content      string `json:"content"`
	Create_time  string `json:"create_time,omitempty"`
}

func (m MessageServiceImpl) Send(message Message) error {
	// 发送数据准备
	var tableMesssage repository.TableMessage
	tableMesssage.ToUserId = message.To_user_id
	tableMesssage.FromUserId = message.From_user_id
	tableMesssage.CreateTime = message.Create_time
	tableMesssage.Content = message.Content

	// 发送数据存储到数据库
	err := repository.InsertTableMessage(tableMesssage)
	if err != nil {
		return err
	}
	return nil
}

func (m MessageServiceImpl) GetList(userId int64, toUserId int64) ([]Message, error) {
	// 根据 userId toUserId 去数据库中查询消息
	tableMessageListFrom, tableMessageListTo, err := repository.GetMessageList(userId, toUserId)
	if err != nil {
		util.Log.Error("MessageService-GetList: return err: " + err.Error())
		return make([]Message, 0), err
	}
	// 返回的 messageList
	messageList := make([]Message, len(tableMessageListTo)+len(tableMessageListFrom))
	// 不可以直接用 append 追加元素
	var i = 0
	for _, tableMessage := range tableMessageListFrom {
		message := Message{}
		message.Id = tableMessage.Id
		message.Content = tableMessage.Content
		message.Create_time = tableMessage.CreateTime
		message.To_user_id = tableMessage.ToUserId
		message.From_user_id = tableMessage.FromUserId
		messageList[i] = message
		i++
	}
	for _, tableMessage := range tableMessageListTo {
		var message Message
		message.Id = tableMessage.Id
		message.Content = tableMessage.Content
		message.Create_time = tableMessage.CreateTime
		message.To_user_id = tableMessage.ToUserId
		message.From_user_id = tableMessage.FromUserId
		messageList[i] = message
		i++
	}
	return messageList, nil
}
