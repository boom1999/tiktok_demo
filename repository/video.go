package repository

import (
	"time"
)

type TableVideo struct {
	Id          int64     `gorm:"column:id;not null;type:bigint(20) primary key auto_increment"`
	AuthorId    int64     `gorm:"column:author_id;not null;type:bigint(20)"`
	PlayUrl     string    `gorm:"column:play_url;not null;type:varchar(255)"`
	CoverUrl    string    `gorm:"column:cover_url;not null;type:varchar(255)"`
	PublishTime time.Time `gorm:"column:publish_time;not null;type:datetime"`
	Title       string    `gorm:"column:title;not null;type:varchar(255)"`
}

// TableName 修改映射名
func (tableVideo TableVideo) TableName() string {
	return "videos"
}
