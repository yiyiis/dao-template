package config

import (
	"backend/pkg/db"
	"backend/pkg/jwt"
	"backend/pkg/log"
	"fmt"
	"github.com/spf13/viper"
)

type Server struct {
	Port int    `yaml:"Port"`
	IP   string `yaml:"IP"`
}

type Config struct {
	DbConf db.Config  `yaml:"DBConf"`
	Server Server     `yaml:"Server"`
	Auth   jwt.Config `yaml:"Auth"`
	Log    log.Config `yaml:"Log"`
}

func LoadConfig(confPath string) Config {
	viper.SetConfigFile(confPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("读取配置文件失败: %v", err))
	}

	var conf Config
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Sprintf("解析配置文件失败: %v", err))
	}

	return conf
}
