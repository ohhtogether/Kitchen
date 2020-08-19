package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"src/controller"
	"src/global"
)

//验证token 验证访问接口权限
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "OPTIONS" {
			var respData controller.RespData
			//判断token
			_, err := global.ParseToken(c.Request)
			if err != nil {

				respData.Error = err.Error()
				c.JSON(http.StatusUnauthorized, &respData)
				c.Abort()
			}
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}
