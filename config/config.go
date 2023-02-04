package config

import "github.com/spf13/viper"

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
	OneDayOfHours int64
	OneMinute     int64
	OneMonth      int64
	OneYear       int64
}

type Configs struct {
	JWT           JWTConfig
	Mysql         MysqlConfig
	OneDayOfHours OneDayOfHoursConfig
}

var Config Configs

// LoadConfig Configs contain Mysql config and other settings from ./config.yaml (using Viper).
func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigFile("./config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	jwt := JWTConfig{
		Secret: viper.GetString("jwt.secretKey"),
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
	Config = Configs{
		JWT:           jwt,
		Mysql:         mysql,
		OneDayOfHours: OneDayOfHours,
	}
}

func GetConfig() Configs {
	return Config
}
