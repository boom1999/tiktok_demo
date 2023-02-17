package service

import (
	"tiktok_demo/repository"
	"time"
)

// CommentService 接口定义
type CommentService interface {
	// CountFromVideoId
	CountFromVideoId(id int64) (int64, error)
	// Send
	Send(comment repository.Comment) (CommentInfo, error)
	// DelComment
	DelComment(commentId int64) error
	// GetList
	GetList(videoId int64, userId int64) ([]CommentInfo, error)
}

// CommentInfo 
type CommentInfo struct {
	Id         int64  `json:"id,omitempty"`
	UserInfo   User   `json:"user,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type CommentData struct {
	Id            int64     `json:"id,omitempty"`
	UserId        int64     `json:"user_id,omitempty"`
	Name          string    `json:"name,omitempty"`
	FollowCount   int64     `json:"follow_count"`
	FollowerCount int64     `json:"follower_count"`
	IsFollow      bool      `json:"is_follow"`
	Content       string    `json:"content,omitempty"`
	CreateDate    time.Time `json:"create_date,omitempty"`
}
