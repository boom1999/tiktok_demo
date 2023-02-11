package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type JWTConfig struct {
	Secret string
}

type MysqlConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

type OneDayOfHoursConfig struct {
	OneMinute     int64
	OneDayOfHours int64
	OneMonth      int64
	OneYear       int64
}

type MinioConfig struct {
	Host         string
	Port         string
	RootUser     string
	RootPassword string
	VideoBuckets string
	PicBuckets   string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type Configs struct {
	JWT           JWTConfig
	Mysql         MysqlConfig
	OneDayOfHours OneDayOfHoursConfig
	Minio         MinioConfig
	Redis         RedisConfig
}

var Config Configs

const configFile = "/config/config.yaml"

// LoadConfig Configs contain Mysql config and other settings from ./config.yaml (using Viper).
func LoadConfig() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(errors.New("Getwd() error"))
	}
	configPath := currentDir + configFile

	if !fileExist(configPath) {
		panic(errors.New("configFile not exist"))
	}
	viper.SetConfigName("config")
	viper.SetConfigFile(configPath)
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	jwt := JWTConfig{
		Secret: viper.GetString("jwt.secret"),
	}
	mysql := MysqlConfig{
		Host:     viper.GetString("mysql.host"),
		Port:     viper.GetString("mysql.port"),
		Database: viper.GetString("mysql.database"),
		Username: viper.GetString("mysql.username"),
		Password: viper.GetString("mysql.password"),
	}
	OneDayOfHours := OneDayOfHoursConfig{
		OneDayOfHours: viper.GetInt64("OneDayOfHours.OneDayOfHours"),
		OneMinute:     viper.GetInt64("OneDayOfHours.OneMinute"),
		OneMonth:      viper.GetInt64("OneDayOfHours.OneMonth"),
		OneYear:       viper.GetInt64("OneDayOfHours.OneYear"),
	}
	minio := MinioConfig{
		Host:         viper.GetString("minio.host"),
		Port:         viper.GetString("minio.port"),
		RootUser:     viper.GetString("minio.rootUser"),
		RootPassword: viper.GetString("minio.rootPassword"),
		VideoBuckets: viper.GetString("minio.videoBuckets"),
		PicBuckets:   viper.GetString("minio.picBuckets"),
	}
	redis := RedisConfig{
		Host:     viper.GetString("redis.host"),
		Port:     viper.GetString("redis.port"),
		Password: viper.GetString("redis.password"),
	}
	Config = Configs{
		JWT:           jwt,
		Mysql:         mysql,
		OneDayOfHours: OneDayOfHours,
		Minio:         minio,
		Redis:         redis,
	}
}

func GetConfig() Configs {
	return Config
}

func fileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

const ValidComment = 0   //评论状态：有效
const InvalidComment = 1 //评论状态：取消
const DateTime = "2006-01-02 15:04:05"

const DefaultRedisValue = -1 //redis中key对应的预设值，防脏读

const IsLike = 0     //点赞的状态
const Unlike = 1     //取消赞的状态
const LikeAction = 1 //点赞的行为
const Attempts = 3   //操作数据库的最大尝试次数

// VideoCount 每次获取视频流的数量
const VideoCount = 5
