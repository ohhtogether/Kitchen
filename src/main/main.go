package main

import (
	"fmt"
	"log"
	"src/global"
	"src/router"
)

func main() {
	// 记录日志
	global.Log()

	// 捕获 panic
	defer global.GetPanic()

	// 加载配置文件
	global.LoadConfig()

	// 初始化路由
	router.SetRouter()

	// 链接数据库
	global.SetDatabaseS()
	err := global.SetDatabase()
	if err!=nil{
		log.Println(err)
	}

	// 链接Reids
	global.SetRedis()

	// 初始化设置
	ip := global.Config.Section("app").Key("ip").String()
	port := global.Config.Section("app").Key("port").String()
	err = router.Router.Run(fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		log.Println(err)
		return
	}
}
