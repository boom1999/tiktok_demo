package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"tiktok_demo/config"
	"tiktok_demo/repository"
	"tiktok_demo/util"
	"time"
)

type VideoServiceImpl struct {
	UserService
	LikeService
	CommentService
}

// Feed
// 通过传入时间戳，当前用户的id，返回对应的视频数组，以及视频数组中最早的发布时间
// 获取视频数组大小是可以控制的，在config中的videoCount变量
func (videoService VideoServiceImpl) Feed(lastTime time.Time, userId int64) ([]Video, time.Time, error) {
	//创建对应返回视频的切片数组，提前将切片的容量设置好，可以减少切片扩容的性能
	videos := make([]Video, 0, config.VideoCount)
	//根据传入的时间，获得传入时间前n个视频，可以通过config.videoCount来控制
	tableVideos, err := repository.GetVideosByLastTime(lastTime)
	if err != nil {
		log.Printf("方法dao.GetVideosByLastTime(lastTime) 失败：%v", err)
		return nil, time.Time{}, err
	}
	log.Printf("方法dao.GetVideosByLastTime(lastTime) 成功：%v", tableVideos)
	//将数据通过copyVideos进行处理，在拷贝的过程中对数据进行组装
	err = videoService.copyVideos(&videos, &tableVideos, userId)
	if err != nil {
		log.Printf("方法videoService.copyVideos(&videos, &tableVideos, userId) 失败：%v", err)
		return nil, time.Time{}, err
	}
	log.Printf("方法videoService.copyVideos(&videos, &tableVideos, userId) 成功")
	//返回数据，同时获得视频中最早的时间返回
	var t time.Time
	return videos, t, nil
}

// GetVideo
// 传入视频id获得对应的视频对象，注意还需要传入当前登录用户id
func (videoService *VideoServiceImpl) GetVideo(videoId int64, userId int64) (Video, error) {
	//初始化video对象
	var video Video

	//从数据库中查询数据，如果查询不到数据，就直接失败返回，后续流程就不需要执行了
	data, err := repository.GetVideoByVideoId(videoId)
	if err != nil {
		log.Printf("方法dao.GetVideoByVideoId(videoId) 失败：%v", err)
		return video, err
	} else {
		log.Printf("方法dao.GetVideoByVideoId(videoId) 成功")
	}

	//插入从数据库中查到的数据
	videoService.creatVideo(&video, &data, userId)
	return video, nil
}

// Publish
// 将传入的视频流保存在minio服务器中，并存储在mysql表中
func (videoService *VideoServiceImpl) Publish(data *multipart.FileHeader, userId int64, title string, c *gin.Context) error {
	//将视频流上传到视频服务器，保存视频链接

	video, err := data.Open()
	if err != nil {
		log.Printf("方法data.Open() 失败%v", err)
		return err
	}
	defer video.Close()
	log.Printf("方法data.Open() 成功")
	//生成一个uuid作为视频的名字
	//videoName := uuid.NewV4().String()
	var videoName strings.Builder
	videoName.WriteString(strconv.FormatInt(userId, 10))
	videoName.WriteString("_")
	videoName.WriteString(strconv.FormatInt(util.GetCurrentTimeMillisecond(), 10))
	videoName.WriteString(".mp4")
	log.Printf("生成视频名称%v", videoName)

	videoPath := config.VideoPath + videoName.String()
	videoBucketName := config.Config.Minio.VideoBuckets
	if err := c.SaveUploadedFile(data, videoPath); err != nil {
		log.Printf("上传到临时地址%v", err)
		return err
	}
	err = repository.FileMinio(videoBucketName, videoName.String(), videoPath, "mp4")
	if err != nil {
		log.Printf("方法repository.VideoMinio(video, videoName.String(), videoSize)%v", err)
		return err
	}
	log.Printf("repository.VideoMinio(video, videoName.String(), videoSize)成功")

	//获取视频第一帧
	picName, err := GetImageFile(videoPath)
	if err != nil {
		log.Printf("获取视频第一帧失败%v", err)
		return err
	}
	log.Printf("获取视频第一帧成功%v", picName)
	picPath := config.PicPath + picName

	// TODO 在服务器上执行ffmpeg 从视频流中获取第一帧截图，并将图片上传到minio上
	pictureBucketName := config.Config.Minio.PicBuckets

	err = repository.FileMinio(pictureBucketName, picName, picPath, "jpg")
	if err != nil {
		log.Printf("方法repository.VideoMinio(image, pictureName.String(), pictureBucketName, videoSize) 失败%v", err)
		return err
	}

	// 向队列中添加消息
	/*ffmpeg.Ffchan <- ffmpeg.Ffmsg{
		videoName,
		imageName,
	}*/
	//组装并持久化
	videoURL, err := repository.GetfileURL(videoBucketName, videoName.String())
	if err != nil {
		log.Printf("方法repository.GetfileURL(videoBucketName, videoName.String()) 失败%v", err)
	}
	pictureURL, err := repository.GetfileURL(pictureBucketName, picName)
	if err != nil {
		log.Printf("方法repository.GetfileURL(pictureBucketName, pictureName.String()) 失败%v", err)
	}
	err = repository.Save(videoURL.String(), pictureURL.String(), userId, title)
	if err != nil {
		log.Printf("方法repository.Save(videoURL.String(), pictureURL.String(), userId, title) 失败%v", err)
		return err
	}
	log.Printf("方法repository.Save(videoURL.String(), pictureURL.String(), userId, title) 成功")
	return nil
}

// List
// 通过userId来查询对应用户发布的视频，并返回对应的视频数组
func (videoService *VideoServiceImpl) List(userId int64, curId int64) ([]Video, error) {
	//依据用户id查询所有的视频，获取视频列表
	data, err := repository.GetVideosByAuthorId(userId)
	if err != nil {
		log.Printf("方法dao.GetVideosByAuthorId(%v)失败:%v", userId, err)
		return nil, err
	}
	log.Printf("方法dao.GetVideosByAuthorId(%v)成功:%v", userId, data)
	//提前定义好切片长度
	result := make([]Video, 0, len(data))
	//调用拷贝方法，将数据进行转换
	err = videoService.copyVideos(&result, &data, curId)
	if err != nil {
		log.Printf("方法videoService.copyVideos(&result, &data, %v)失败:%v", userId, err)
		return nil, err
	}
	//如果数据没有问题，则直接返回
	return result, nil
}

// 该方法可以将数据进行拷贝和转换，并从其他方法获取对应的数据
func (videoService *VideoServiceImpl) copyVideos(result *[]Video, data *[]repository.TableVideo, userId int64) error {
	for _, temp := range *data {
		var video Video
		//将video进行组装，添加想要的信息,插入从数据库中查到的数据
		videoService.creatVideo(&video, &temp, userId)
		*result = append(*result, video)
	}
	return nil
}

// 将video进行组装，添加想要的信息,插入从数据库中查到的数据
func (videoService *VideoServiceImpl) creatVideo(video *Video, data *repository.TableVideo, userId int64) {
	//建立协程组，当这一组的携程全部完成后，才会结束本方法
	var wg sync.WaitGroup
	wg.Add(4)
	var err error
	video.TableVideo = *data
	//插入Author，这里需要将视频的发布者和当前登录的用户传入，才能正确获得isFollow，
	//如果出现错误，不能直接返回失败，将默认值返回，保证稳定
	go func() {
		video.Author, err = videoService.GetUserByIdWithCurId(data.AuthorId, userId)
		if err != nil {
			log.Printf("方法videoService.GetUserByIdWithCurId(data.AuthorId, userId) 失败：%v", err)
		} else {
			log.Printf("方法videoService.GetUserByIdWithCurId(data.AuthorId, userId) 成功")
		}
		wg.Done()
	}()

	//插入点赞数量，同上所示，不将nil直接向上返回，数据没有就算了，给一个默认就行了
	go func() {
		video.FavoriteCount, err = videoService.FavouriteCount(data.Id)
		if err != nil {
			log.Printf("方法videoService.FavouriteCount(data.ID) 失败：%v", err)
		} else {
			log.Printf("方法videoService.FavouriteCount(data.ID) 成功")
		}
		wg.Done()
	}()

	//获取该视屏的评论数字
	go func() {
		video.CommentCount, err = videoService.CountFromVideoId(data.Id)
		if err != nil {
			log.Printf("方法videoService.CountFromVideoId(data.ID) 失败：%v", err)
		} else {
			log.Printf("方法videoService.CountFromVideoId(data.ID) 成功")
		}
		wg.Done()
	}()

	//获取当前用户是否点赞了该视频
	go func() {
		video.IsFavorite, err = videoService.IsFavourite(video.Id, userId)
		if err != nil {
			log.Printf("方法videoService.IsFavourit(video.Id, userId) 失败：%v", err)
		} else {
			log.Printf("方法videoService.IsFavourit(video.Id, userId) 成功")
		}
		wg.Done()
	}()

	wg.Wait()
}

// GetVideoIdList
// 通过一个作者id，返回该用户发布的视频id切片数组
func (videoService *VideoServiceImpl) GetVideoIdList(authorId int64) ([]int64, error) {
	//直接调用dao层方法获取id即可
	id, err := repository.GetVideoIdsByAuthorId(authorId)
	if err != nil {
		log.Printf("方法dao.GetVideoIdsByAuthorId(%v) 失败：%v", authorId, err)
		return nil, err
	} else {
		log.Printf("方法dao.GetVideoIdsByAuthorId(%v) 成功", authorId)
	}
	return id, nil
}

func GetImageFile(videoPath string) (string, error) {
	temp := strings.Split(videoPath, "/")
	videoName := temp[len(temp)-1]
	b := []byte(videoName)
	videoName = string(b[:len(b)-3]) + "jpg"
	picpath := "/tmp/file/picture/"
	picName := filepath.Join(picpath, videoName)
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-ss", "1", "-f", "image2", "-t", "0.01", "-y", picName)
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return videoName, nil
}
