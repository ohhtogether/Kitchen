package global

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// 配置日志文件
func Log() {
	// 读取文本文件 OpenFile(文件名, 打开方式, 打开模式和权限位)
	// os.O_CREATE 如果文件不存在将创建一个新文件
	// os.O_WRONLY 只写模式打开文件
	// os.O_APPEND 写操作时将数据附加到文件末尾
	file, err := os.OpenFile("../access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("打开日志文件失败：", err) // 非正常运行导致退出程序
		return
	}

	// Gin框架的路由"请求"日志写入file文件中
	// [GIN] 2020/05/08 - 15:49:36 | 200 |            0s |             ::1 | POST      /login
	gin.DefaultWriter = file

	// SetOutput 设置标准logger的输出目的地, log.Println()写入file文件中
	log.SetOutput(file)
	// SetFlags 设置标准logger的输出选项
	// Ldate  日期：2009/01/23
	// Ltime  时间：01:23:23
	// Llongfile  文件全路径名+行号： /a/b/c/d.go:23
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
}

// 配置捕获panic
func GetPanic() {
	if err := recover(); err != nil {
		log.Println(err)
		return
	}
}
