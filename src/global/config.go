package global

import (
	"log"

	"github.com/go-ini/ini"
)

var Config *ini.File

// 加载配置文件
func LoadConfig() {
	config, err := ini.Load("../config.ini")
	if err != nil {
		log.Println(err)
		return
	}
	Config = config
}
