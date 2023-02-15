package repository

import (
	"io"
	"log"
	"net/url"
	"tiktok_demo/config"
	"tiktok_demo/middleware/minio"
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

// TableName
func (TableVideo) TableName() string {
	return "videos"
}

// GetVideosByAuthorId
// 根据作者的id来查询对应数据库数据，并TableVideo返回切片
func GetVideosByAuthorId(authorId int64) ([]TableVideo, error) {
	//建立结果集接收
	var data []TableVideo
	//初始化db
	//Init()
	result := DB.Where(&TableVideo{AuthorId: authorId}).Find(&data)
	//如果出现问题，返回对应到空，并且返回error
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

// GetVideoByVideoId
// 依据VideoId来获得视频信息
func GetVideoByVideoId(videoId int64) (TableVideo, error) {
	var tableVideo TableVideo
	tableVideo.Id = videoId
	//Init()
	result := DB.First(&tableVideo)
	if result.Error != nil {
		return tableVideo, result.Error
	}
	return tableVideo, nil

}

// GetVideosByLastTime
// 依据一个时间，来获取这个时间之前的一些视频
func GetVideosByLastTime(lastTime time.Time) ([]TableVideo, error) {
	videos := make([]TableVideo, config.VideoCount)
	result := DB.Where("publish_time>?", lastTime).Order("publish_time desc").Limit(config.VideoCount).Find(&videos)
	if result.Error != nil {
		return videos, result.Error
	}
	return videos, nil
}

// FileMinio
// 上传文件到minio
func FileMinio(bucketName string, objectName string, file io.Reader, contentType string, objectSize int64) error {
	//转到minio相对路线下
	err := minio.UploadFile(bucketName, objectName, file, contentType, objectSize)
	//_, err := minio.UploadLocalFile(bucketName, objectName, filePath, contentType)
	if err != nil {
		log.Println("上传%v类型%v至minio失败！！！", contentType, objectName)
	} else {
		log.Println("上传%v类型%v至minio成功！！！", contentType, objectName)
	}
	log.Println("上传成功！！！！！")
	return nil
}

// 获取视频的URL
func GetfileURL(bucketName string, fileName string) (*url.URL, error) {
	// 过期时间为一天
	var expires time.Duration = 0
	fileURL, err := minio.GetFileUrl(bucketName, fileName, expires)
	if err != nil {
		log.Println("GetURL false!!!")
	}
	return fileURL, err
}

// Save 保存视频记录
func Save(videoURL string, pictureURL string, userId int64, title string) error {
	var video TableVideo
	video.PublishTime = time.Now()
	video.PlayUrl = videoURL
	video.CoverUrl = pictureURL
	video.AuthorId = userId
	video.Title = title
	result := DB.Save(&video)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetVideoIdsByAuthorId
// 通过作者id来查询发布的视频id切片集合
func GetVideoIdsByAuthorId(authorId int64) ([]int64, error) {
	var id []int64
	//通过pluck来获得单独的切片
	result := DB.Model(&TableVideo{}).Where("author_id", authorId).Pluck("id", &id)
	//如果出现问题，返回对应到空，并且返回error
	if result.Error != nil {
		return nil, result.Error
	}
	return id, nil
}
