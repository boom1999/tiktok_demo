package repository

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"tiktok_demo/config"
)

var DB *gorm.DB
var Conf = config.GetConfig()

func Init() *gorm.DB {
	host := Conf.Mysql.Host
	port := Conf.Mysql.Port
	database := Conf.Mysql.Database
	username := Conf.Mysql.Username
	password := Conf.Mysql.Password
	charset := "utf8"

	var err error
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database, err:" + err.Error())
	}
	return DB
}

func GetDataBase() *gorm.DB {
	return DB
}
