package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"src/model"
	"strconv"
)

//餐厅列表
func RestaurntList(c *gin.Context) {
	var respData RespData
	// 接收分页数据
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
	// 验证
	var offset, limit int64 = 0, 0
	if offsetStr != "" && limitStr != "" {
		var errOs, errLi error = nil, nil
		offset, errOs = strconv.ParseInt(offsetStr, 10, 64)
		limit, errLi = strconv.ParseInt(limitStr, 10, 64)
		if errOs != nil || errLi != nil {
			respData.Error = "分页参数错误"
			c.JSON(http.StatusOK, &respData)
			return
		}
	}
	//查询数据
	var rest model.Restaurnt
	data, err := rest.Page(offset, limit)
	if err != nil {
		respData.Error = "数据出错"
		c.JSON(http.StatusOK, &respData)
		return
	}
	fmt.Println(err)
	fmt.Println(data)
	// 获取数据总数
	num, err := rest.Total()
	if err != nil {
		respData.Error = "数据出错"
		c.JSON(http.StatusOK, &respData)
		return
	}

	res := struct {
		Num  int64             `json:"num"`
		Data []model.Restaurnt `json:"data"`
	}{
		Num:  num,
		Data: data,
	}
	respData.Data = res
	c.JSON(http.StatusOK, &respData)
	return
}

//添加一条餐厅信息
func RestaurntAdd(c *gin.Context) {
	var respData RespData

	respData.Data = "添加一条餐厅信息"
	c.JSON(http.StatusOK, &respData)
	return
}

//修改一条餐厅信息
func RestaurntEdit(c *gin.Context) {
	var respData RespData

	respData.Data = "修改一条餐厅信息"
	c.JSON(http.StatusOK, &respData)
	return
}

//删除一条餐厅信息
func RestaurntDel(c *gin.Context) {
	var respData RespData

	respData.Data = "删除一条餐厅信息"
	c.JSON(http.StatusOK, &respData)
	return
}
