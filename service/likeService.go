package service

// LikeService 接口定义
type LikeService interface {
	// FavouriteAction 当前用户对视频进行点赞或取消点赞操作。1：点赞，2：取消点赞
	FavouriteAction(userId int64, videoId int64, actionType int32) error
	// GetFavouriteList 获取当前用户的所有点赞视频
	GetFavouriteList(userId int64, curId int64) ([]Video, error)

	// IsFavourite 根据 当前视频id 和 用户id 判断是否点赞了该视频
	IsFavourite(videoId int64, userId int64) (bool, error)
	// FavouriteCount 根据当前视频id获取点赞该视频的数量
	FavouriteCount(videoId int64) (int64, error)
	// TotalFavourite 根据 用户id 获取该用户总共被点赞的数量
	TotalFavourite(userId int64) (int64, error)
	// FavouriteVideoCount 根据 用户id 获取该用户点赞视频的数量
	FavouriteVideoCount(userId int64) (int64, error)
}
