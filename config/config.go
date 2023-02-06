package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
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

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type OneDayOfHoursConfig struct {
	OneMinute     int64
	OneDayOfHours int64
	OneMonth      int64
	OneYear       int64
}

type Configs struct {
	JWT           JWTConfig
	Mysql         MysqlConfig
	OneDayOfHours OneDayOfHoursConfig
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
	redis := RedisConfig{
		Host:     viper.GetString("redis.host"),
		Port:     viper.GetString("redis.port"),
		Password: viper.GetString("redis.password"),
	}
	Config = Configs{
		JWT:           jwt,
		Mysql:         mysql,
		OneDayOfHours: OneDayOfHours,
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
