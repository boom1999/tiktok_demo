package repository

import (
	"errors"
	"log"
	"tiktok_demo/config"
)

// Comment
type Comment struct {
	Id          int64  `gorm:"column:id;not null;type:bigint(20) primary key auto_increment"`
	UserId      int64  `gorm:"column:user_id;not null;type:bigint(20)"`
	VideoId     int64  `gorm:"column:video_id;not null;type:bigint(20)"`
	CommentText string `gorm:"column:comment_text;not null;type:varchar(255)"`
	CreateDate  string `gorm:"column:create_date;not null;type:varchar(255)"`
	Cancel      int32  `gorm:"column:cancel;not null;default:0;type:int(20)"`
}

// TableName
func (comment Comment) TableName() string {
	return "comments"
}

// Count
//  查询Comment数量
func Count(videoId int64) (int64, error) {
	log.Println("CommentDao-Count: running")
	//InitDataBase()
	var count int64
	//数据库中查询评论数量
	err := DB.Model(Comment{}).Where(map[string]interface{}{"video_id": videoId, "cancel": config.ValidComment}).Count(&count).Error
	if err != nil {
		log.Println("CommentDao-Count: return count failed")
		return -1, errors.New("find comments count failed")
	}
	log.Println("CommentDao-Count: return count success")
	return count, nil
}

// CommentIdList
func CommentIdList(videoId int64) ([]string, error) {
	var commentIdList []string
	err := DB.Model(Comment{}).Select("id").Where("video_id = ?", videoId).Find(&commentIdList).Error
	if err != nil {
		log.Println("CommentIdList:", err)
		return nil, err
	}
	return commentIdList, nil
}

// InsertComment
func InsertComment(comment Comment) (Comment, error) {
	log.Println("CommentDao-InsertComment: running")
	//数据库中插入一条评论信息

	err := DB.Model(Comment{}).Create(&comment).Error
	if err != nil {
		log.Println("CommentDao-InsertComment: return create comment failed")
		return Comment{}, errors.New("create comment failed")
	}
	log.Println("CommentDao-InsertComment: return success")
	return comment, nil
}

// DeleteComment
func DeleteComment(id int64) error {
	log.Println("CommentDao-DeleteComment: running")
	var commentInfo Comment
	//先查询是否有此评论
	result := DB.Model(Comment{}).Where(map[string]interface{}{"id": id, "cancel": config.ValidComment}).First(&commentInfo)
	if result.RowsAffected == 0 {
		log.Println("CommentDao-DeleteComment: return del comment is not exist")
		return errors.New("del comment is not exist")
	}
	//数据库中删除评论-更新评论状态为-1
	err := DB.Model(Comment{}).Where("id = ?", id).Update("cancel", config.InvalidComment).Error
	if err != nil {
		log.Println("CommentDao-DeleteComment: return del comment failed")
		return errors.New("del comment failed")
	}
	log.Println("CommentDao-DeleteComment: return success")
	return nil
}

// GetCommentList
func GetCommentList(videoId int64) ([]Comment, error) {
	log.Println("CommentDao-GetCommentList: running")
	//数据库中查询评论信息list
	var commentList []Comment
	result := DB.Model(Comment{}).Where(map[string]interface{}{"video_id": videoId, "cancel": config.ValidComment}).
		Order("create_date desc").Find(&commentList)
	//若此视频没有评论信息，返回空列表，不报错
	if result.RowsAffected == 0 {
		log.Println("CommentDao-GetCommentList: return there are no comments")
		return nil, nil
	}
	//若获取评论列表出错
	if result.Error != nil {
		log.Println(result.Error.Error())
		log.Println("CommentDao-GetCommentList: return get comment list failed")
		return commentList, errors.New("get comment list failed")
	}
	log.Println("CommentDao-GetCommentList: return commentList success")
	return commentList, nil
}
