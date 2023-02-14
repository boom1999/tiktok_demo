package repository

import (
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
func FileMinio(bucketName string, objectName string, filePath string, contentType string) error {
	//转到minio相对路线下
	//err := minio.UploadFile(bucketName, objectName, file, objectSize)
	_, err := minio.UploadLocalFile(bucketName, objectName, filePath, contentType)
	if err != nil {
		log.Println("转到路径video失败！！！")
	} else {
		log.Println("转到路径video成功！！！")
	}
	log.Println("上传视频成功！！！！！")
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

// ImageFTP
// 将图片传入FTP服务器中，但是这里要注意图片的格式随着名字一起给,同时调用时需要自己结束流
//func ImageFTP(file io.Reader, imageName string) error {
//	//转到video相对路线下
//	err := ftp.MyFTP.Cwd("images")
//	if err != nil {
//		log.Println("转到路径images失败！！！")
//		return err
//	}
//	log.Println("转到路径images成功！！！")
//	if err = ftp.MyFTP.Stor(imageName, file); err != nil {
//		log.Println("上传图片失败！！！！！")
//		return err
//	}
//	log.Println("上传图片成功！！！！！")
//	return nil
//}

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
