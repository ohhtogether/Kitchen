package controller

import (
	"errors"
	"log"
	"net/http"
	"src/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
	//接受数据
	name := c.PostForm("name")
	phone := c.PostForm("phone")
	longitude := c.PostForm("longitude")
	latitude := c.PostForm("latitude")
	address := c.PostForm("address")

	//验证提交的参数非空
	if len(name) < 1 || len(name) > 128 || len(phone) < 1 || len(phone) > 16 || len(address) > 128 || longitude == "" || latitude == "" {
		respData.Error = "提交参数长度出错"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	var rest model.Restaurnt
	rest.Name = name
	rest.Address = address
	rest.Latitude = latitude
	rest.Longitude = longitude
	rest.Phone = phone
	rest.Time = time.Now().Unix()
	rest.Status = 1
	//开始添加数据
	err := rest.Add()
	if err != nil {
		respData.Error = "添加数据失败"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}
	respData.Data = "添加一条餐厅信息成功"
	c.JSON(http.StatusOK, &respData)
	return
}

//修改一条餐厅信息
func RestaurntEdit(c *gin.Context) {
	var respData RespData

	//接受数据
	id := c.PostForm("id")
	name := c.PostForm("name")
	phone := c.PostForm("phone")
	longitude := c.PostForm("longitude")
	latitude := c.PostForm("latitude")
	address := c.PostForm("address")

	//验证提交的参数
	rid, err := strconv.ParseInt(id, 10, 64)

	if len(name) < 1 || len(name) > 128 || len(phone) < 1 || len(phone) > 16 || len(address) > 128 || longitude == "" || latitude == "" || err != nil || rid <= 0 {
		respData.Error = "提交参数出错"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	var rest model.Restaurnt
	rest.Id = rid
	rest.Name = name
	rest.Address = address
	rest.Latitude = latitude
	rest.Longitude = longitude
	rest.Phone = phone
	rest.Time = time.Now().Unix()

	//判断是否有需修改的餐厅的记录
	one, err := rest.One()
	if err != nil {
		respData.Error = "没有此餐厅的信息"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}
	rest.Status = one.Status
	//将需要修改的数据修改入数据库
	err = rest.Edit()
	if err != nil {
		respData.Data = "修改餐厅信息失败"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}
	//返回
	respData.Data = "修改餐厅信息成功"
	c.JSON(http.StatusOK, &respData)
	return
}

//删除一条餐厅信息
func RestaurntDel(c *gin.Context) {
	var respData RespData
	//接受数据
	id := c.PostForm("id")

	//验证提交的参数
	rid, err := strconv.ParseInt(id, 10, 64)
	if err != nil || rid <= 0 {
		respData.Error = "提交参数出错"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	var rest model.Restaurnt
	rest.Id = rid
	//判断是否有需修改的餐厅的记录
	one, err := rest.One()
	if err != nil {
		respData.Error = "没有此餐厅的信息"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}
	//将其状态改为删除
	one.Status = 0
	//将需要修改的数据修改入数据库
	err = one.Edit()
	if err != nil {
		respData.Data = "删除餐厅信息失败"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}
	//返回
	respData.Data = "删除餐厅信息成功"
	c.JSON(http.StatusOK, &respData)
	return
}

// id对name
func RestMap() (restMap map[int64]string, err error) {
	var rest model.Restaurnt
	all, err := rest.All()
	if err != nil {
		log.Println(err)
		err = errors.New("获取餐厅数据失败")
		return
	}

	restMap = make(map[int64]string)
	for _, v := range all {
		restMap[v.Id] = v.Name
	}
	return
}
