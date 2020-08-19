package router

import (
	"github.com/gin-gonic/gin"

	"src/controller"
)

var Router *gin.Engine

func SetRouter() {
	Router = gin.Default()
	// TODO 静态图片访问
	Router.Static("/img", "../public/img/")
	// 跨域
	Router.Use(CorsHandler())
	// 捕获 panic
	Router.Use(ErrHandler())

	Router.GET("/", controller.Index)

	//Router.GET("/user/list", controller.UserList)
	//Router.POST("/user/update", controller.UserUpdate)
	//Router.POST("/user/password", controller.UserPass)
	//Router.POST("/user/status", controller.UserStatus)

	/*********************************************************************/
	Router.GET("/restaurnt/list", controller.RestaurntList)
	Router.POST("/restaurnt/add", controller.RestaurntAdd)
	Router.POST("/restaurnt/edit", controller.RestaurntEdit)
	Router.POST("/restaurnt/del", controller.RestaurntDel)

	Router.Use(AuthRequired())

}
