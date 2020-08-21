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
	Router.POST("/user/login", controller.Login)
	Router.Use(AuthRequired())
	Router.POST("/user/logout", controller.Logout)

	Router.GET("/device/list", controller.DeviceInfoView)
	Router.POST("/device/add", controller.DeviceInfoAddOne)
	Router.POST("/device/edit", controller.DeviceInfoEdit)
	Router.POST("/device/del", controller.DeviceInfoDel)
	Router.GET("/device/parent", controller.GetParentIdData)
	Router.GET("/device/parent/data", controller.GetParentData)
	Router.GET("/device/state", controller.DeviceState)

	Router.GET("/device/status", controller.DeviceStatusView)

	Router.GET("/restaurnt/list", controller.RestaurntList)
	Router.POST("/restaurnt/add", controller.RestaurntAdd)
	Router.POST("/restaurnt/edit", controller.RestaurntEdit)
	Router.POST("/restaurnt/del", controller.RestaurntDel)

	Router.GET("/employee/list", controller.EmployeeList)
	Router.POST("/employee/add", controller.EmployeeAdd)
	Router.POST("/employee/edit", controller.EmployeeEdit)
	Router.POST("/employee/del", controller.EmployeeDel)

	Router.GET("/purchase/list", controller.PurchaseList)
	Router.POST("/purchase/add", controller.PurchaseAdd)
	Router.POST("/purchase/edit", controller.PurchaseEdit)
	Router.POST("/purchase/del", controller.PurchaseDel)

	// data 统计
	Router.GET("/data", controller.Data)
	// 地图获取小区信息
	Router.GET("/map/housedistrict", controller.Housedistrict)
	Router.GET("/map/coordinate", controller.Coordinate)

	//---
	Router.GET("/device/type", controller.DeviceTypeView)
	Router.GET("/status/type/pulldown", controller.StatusTypePullDown)
	// 区块
	Router.GET("/block/info", controller.BlockView)
}
