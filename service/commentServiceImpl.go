package service

import (
	"sort"
	"strconv"
	"sync"
	"tiktok_demo/config"
	"tiktok_demo/middleware/rabbitmq"
	"tiktok_demo/middleware/redis"
	"tiktok_demo/repository"
	"tiktok_demo/util"
	"time"

	"go.uber.org/zap"
)

type CommentServiceImpl struct {
	UserService
}

// CountFromVideoId
func (c CommentServiceImpl) CountFromVideoId(videoId int64) (int64, error) {
	//先在缓存中查
	cnt, err := redis.RdbVCid.SCard(redis.Ctx, strconv.FormatInt(videoId, 10)).Result()
	if err != nil {
		util.Log.Error("count from redis error:" + err.Error())
	}
	util.Log.Debug("info", zap.Int64("comment count redis :", cnt))

	//1.缓存中查到了数量，则返回数量值-1（去除0值）
	if cnt != 0 {
		return cnt - 1, nil
	}
	//2.缓存中查不到则去数据库查
	cntDao, err1 := repository.Count(videoId)
	util.Log.Debug("info", zap.Int64("comment count dao :", cntDao))
	if err1 != nil {
		util.Log.Error("comment count dao err:" + err1.Error())
		return 0, nil
	}
	//将评论id切片存入redis
	go func() {
		cList, _ := repository.CommentIdList(videoId)
		_, _err := redis.RdbVCid.SAdd(redis.Ctx, strconv.Itoa(int(videoId)), config.DefaultRedisValue).Result()
		if _err != nil {
			util.Log.Error("redis save one vId - cId 0 failed")
			return
		}
		_, err := redis.RdbVCid.Expire(redis.Ctx, strconv.Itoa(int(videoId)),
			time.Duration(config.Config.OneDayOfHours.OneMonth)*time.Second).Result()
		if err != nil {
			util.Log.Error("redis save one vId - cId expire failed")
		}
		for _, commentId := range cList {
			insertRedisVideoCommentId(strconv.Itoa(int(videoId)), commentId)
		}
		util.Log.Debug("count comment save ids in redis")
	}()
	//返回结果
	return cntDao, nil
}

// Send
func (c CommentServiceImpl) Send(comment repository.Comment) (CommentInfo, error) {
	util.Log.Debug("CommentService-Send: running") //函数已运行
	var commentInfo repository.Comment
	commentInfo.VideoId = comment.VideoId         //评论视频id传入
	commentInfo.UserId = comment.UserId           //评论用户id传入
	commentInfo.CommentText = comment.CommentText //评论内容传入
	commentInfo.Cancel = config.ValidComment      //评论状态，0，有效
	commentInfo.CreateDate = comment.CreateDate   //评论时间

	commentRtn, err := repository.InsertComment(commentInfo)
	if err != nil {
		return CommentInfo{}, err
	}

	impl := UserImpl{
		FollowService: &FollowImpl{},
	}
	userData, err2 := impl.GetUserByIdWithCurId(comment.UserId, comment.UserId)
	if err2 != nil {
		return CommentInfo{}, err2
	}

	commentData := CommentInfo{
		Id:         commentRtn.Id,
		UserInfo:   userData,
		Content:    commentRtn.CommentText,
		CreateDate: commentRtn.CreateDate,
	}

	go func() {
		insertRedisVideoCommentId(strconv.Itoa(int(comment.VideoId)), strconv.Itoa(int(commentRtn.Id)))
		util.Log.Debug("send comment save in redis")
	}()
	//返回结果
	return commentData, nil
}

// DelComment
func (c CommentServiceImpl) DelComment(commentId int64) error {
	util.Log.Debug("CommentService-DelComment: running") //函数已运行
	//1.先查询redis，若有则删除，返回客户端-再go协程删除数据库；无则在数据库中删除，返回客户端。
	n, err := redis.RdbCVid.Exists(redis.Ctx, strconv.FormatInt(commentId, 10)).Result()
	if err != nil {
		util.Log.Error(err.Error())
	}
	if n > 0 { //在缓存中有此值，则找出来删除，然后返回
		vid, err1 := redis.RdbCVid.Get(redis.Ctx, strconv.FormatInt(commentId, 10)).Result()
		if err1 != nil { //没找到，返回err
			util.Log.Error("info", zap.String("redis find CV err:", err1.Error()))
		}
		//删除，两个redis都要删除
		del1, err2 := redis.RdbCVid.Del(redis.Ctx, strconv.FormatInt(commentId, 10)).Result()
		if err2 != nil {
			util.Log.Error(err2.Error())
		}
		del2, err3 := redis.RdbVCid.SRem(redis.Ctx, vid, strconv.FormatInt(commentId, 10)).Result()
		if err3 != nil {
			util.Log.Error(err3.Error())
		}
		util.Log.Debug("info", zap.Int64("del comment in Redis success:del1=", del1)) //del1、del2代表删除了几条数据
		util.Log.Debug("info", zap.Int64("del comment in Redis success:del2=", del2))
		//评论id传入消息队列
		rabbitmq.RmqCommentDel.Publish(strconv.FormatInt(commentId, 10))
		return nil
	}
	//不在内存中，则直接走数据库删除
	return repository.DeleteComment(commentId)
}

// GetList
func (c CommentServiceImpl) GetList(videoId int64, userId int64) ([]CommentInfo, error) {
	util.Log.Debug("CommentService-GetList: running")

	//先查评论，再循环查用户信息：
	//先查询评论列表信息
	commentList, err := repository.GetCommentList(videoId)
	if err != nil {
		util.Log.Error("CommentService-GetList: return err: " + err.Error()) //函数返回提示错误信息
		return nil, err
	}
	//当前有0条评论
	if commentList == nil {
		return nil, nil
	}

	//定义切片长度
	commentInfoList := make([]CommentInfo, len(commentList))

	wg := &sync.WaitGroup{}
	wg.Add(len(commentList))
	idx := 0
	for _, comment := range commentList {
		//调用方法组装评论信息，再append
		var commentData CommentInfo
		//将评论信息进行组装，添加想要的信息,插入从数据库中查到的数据
		go func(comment repository.Comment) {
			oneComment(&commentData, &comment, userId)
			//组装list
			//commentInfoList = append(commentInfoList, commentData)
			commentInfoList[idx] = commentData
			idx = idx + 1
			wg.Done()
		}(comment)
	}
	wg.Wait()
	//评论排序-按照主键排序
	sort.Sort(CommentSlice(commentInfoList))

	//协程查询redis中是否有此记录，无则将评论id切片存入redis
	go func() {
		//1.先在缓存中查此视频是否已有评论列表
		cnt, err1 := redis.RdbVCid.SCard(redis.Ctx, strconv.FormatInt(videoId, 10)).Result()
		if err1 != nil { //若查询缓存出错，则打印
			//return 0, err
			util.Log.Error("info", zap.String("count from redis error:", err.Error()))
		}
		//2.缓存中查到了数量大于0，则说明数据正常，不用更新缓存
		if cnt > 0 {
			return
		}
		//3.缓存中数据不正确，更新缓存：
		//先在redis中存储一个-1 值，防止脏读
		_, _err := redis.RdbVCid.SAdd(redis.Ctx, strconv.Itoa(int(videoId)), config.DefaultRedisValue).Result()
		if _err != nil { //若存储redis失败，则直接返回
			util.Log.Error("redis save one vId - cId 0 failed")
			return
		}
		//设置key值过期时间
		_, err2 := redis.RdbVCid.Expire(redis.Ctx, strconv.Itoa(int(videoId)),
			time.Duration(config.Config.OneDayOfHours.OneMonth)*time.Second).Result()
		if err2 != nil {
			util.Log.Error("redis save one vId - cId expire failed")
		}
		//将评论id循环存入redis
		for _, _comment := range commentInfoList {
			insertRedisVideoCommentId(strconv.Itoa(int(videoId)), strconv.Itoa(int(_comment.Id)))
		}
		util.Log.Debug("comment list save ids in redis")
	}()

	util.Log.Debug("CommentService-GetList: return list success") //函数执行成功，返回正确信息
	return commentInfoList, nil
}

// 在redis中存储video_id对应的comment_id 、 comment_id对应的video_id
func insertRedisVideoCommentId(videoId string, commentId string) {
	//在redis-RdbVCid中存储video_id对应的comment_id
	_, err := redis.RdbVCid.SAdd(redis.Ctx, videoId, commentId).Result()
	if err != nil { //若存储redis失败-有err，则直接删除key
		util.Log.Error("redis save send: vId - cId failed, key deleted")
		redis.RdbVCid.Del(redis.Ctx, videoId)
		return
	}
	//在redis-RdbCVid中存储comment_id对应的video_id
	_, err = redis.RdbCVid.Set(redis.Ctx, commentId, videoId, 0).Result()
	if err != nil {
		util.Log.Error("redis save one cId - vId failed")
	}
}

// 此函数用于给评论赋值：
func oneComment(comment *CommentInfo, com *repository.Comment, userId int64) {
	var wg sync.WaitGroup
	wg.Add(1)
	//根据评论用户id和当前用户id，查询评论用户信息
	impl := UserImpl{
		FollowService: &FollowImpl{},
	}
	var err error
	comment.Id = com.Id
	comment.Content = com.CommentText
	comment.CreateDate = com.CreateDate
	comment.UserInfo, err = impl.GetUserByIdWithCurId(com.UserId, userId)
	if err != nil {
		util.Log.Error("CommentService-GetList: GetUserByIdWithCurId return err: " + err.Error()) //函数返回提示错误信息
	}
	wg.Done()
	wg.Wait()
}

// CommentSlice 排序准备工作
type CommentSlice []CommentInfo

func (a CommentSlice) Len() int { //重写Len()方法
	return len(a)
}
func (a CommentSlice) Swap(i, j int) { //重写Swap()方法
	a[i], a[j] = a[j], a[i]
}
func (a CommentSlice) Less(i, j int) bool { //重写Less()方法
	return a[i].Id > a[j].Id
}
