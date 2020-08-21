package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"src/global"
	"src/model"
)

//新增设备类型
/*func DeviceTypeAdd(c *gin.Context) {

	var respData RespData //返回的结构体

	name := c.PostForm("typename")
	des := c.PostForm("des")

	var dt model.DeviceType
	//验证参数
	err := filter.MSet(
		filter.El(&dt.TypeName, filter.FromString(name, "设备类型名称名称").
			Required().MinLength(4).MaxLength(32).IsLetterOrDigit(),
		),
		filter.El(&dt.Des, filter.FromString(des, "设备类型描述描述").
			Required(),
		),
	)
	if err != nil {
		respData.Error = err.Error()
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	global.FilteredSQLInject(dt.Des)

	err = dt.Add()
	if err != nil {
		log.Println(err)
		respData.Error = "新增设备类型失败"
		c.JSON(http.StatusOK, &respData)
		return
	}

	respData.Data = "新增设备类型成功"
	c.JSON(http.StatusOK, &respData)
	return
}*/

//删除设备类型
/*func DeviceTypeDel(c *gin.Context) {

	var respData RespData //返回的结构体

	typeid, err := strconv.ParseInt(c.PostForm("typeid"), 10, 64)
	if err != nil {
		respData.Error = "设备类型ID必须是数字"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	//修改状态
	var dt model.DeviceType
	dt.TypeId = typeid
	//判断type id是否存在
	_, err = dt.One()
	if err != nil {
		log.Println(err)
		respData.Error = "设备类型ID不存在"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}
	dt.Status = 0 //0代表删除
	err = dt.Del()
	if err != nil {
		log.Println(err)
		respData.Error = "删除设备类型成功"
		c.JSON(http.StatusOK, &respData)
		return
	}

	respData.Data = "删除设备类型成功"
	c.JSON(http.StatusOK, &respData)
	return
}*/

//修改设备类型
/*func DeviceTypeEdit(c *gin.Context) {

	var respData RespData //返回的结构体

	id, err := strconv.ParseInt(c.PostForm("typeid"), 10, 64)
	if err != nil && id != 0 {
		respData.Error = "设备类型ID必须是数字"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}
	name := c.PostForm("typename")
	des := c.PostForm("des")

	var dt model.DeviceType
	//验证参数
	err = filter.MSet(
		filter.El(&dt.TypeName, filter.FromString(name, "设备类型名称").
			Required().MinLength(4).MaxLength(32).IsLetterOrDigit(),
		),
		filter.El(&dt.Des, filter.FromString(des, "设备类型描述").
			Required(),
		),
	)
	if err != nil {
		respData.Error = err.Error()
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	global.FilteredSQLInject(dt.Des)

	dt.TypeId = id
	//判断roleid是否存在
	_, err = dt.One()
	if err != nil {
		log.Println(err)
		respData.Error = "设备类型ID不存在"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	//修改数据
	err = dt.Edit()
	if err != nil {
		log.Println(err)
		respData.Error = "修改设备类型失败"
		c.JSON(http.StatusOK, &respData)
		return
	}

	respData.Data = "修改设备类型成功"
	c.JSON(http.StatusOK, &respData)
	return
}*/

//获取设备信息
func DeviceTypeView(c *gin.Context) {

	//判断是请求角色列表还是请求指定角色信息
	id := c.Query("typeid")
	if id == "" {
		//请求角色list
		deviceTypePage(c)
	} else {
		//请求角色详情
		deviceTypeById(c, id)
	}
}

// 分页
func deviceTypePage(c *gin.Context) {
	var respData RespData
	// 偏移量
	var offset int64
	var err error
	offsetStr := c.Query("offset")
	if offsetStr != "" {
		offset, err = strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			respData.Error = "偏移量必须是数字"
			c.JSON(http.StatusOK, &respData)
			return
		}
	}

	// 每页显示的条数
	var limit int64
	limitStr := c.Query("limit")
	if limitStr != "" {
		limit, err = strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			respData.Error = "每页显示的条数必须是数字"
			c.JSON(http.StatusOK, &respData)
			return
		}
	}

	// 关键字
	keys := c.Query("keys")
	global.FilteredSQLInject(keys)

	// 获取数据
	var deviceType model.DeviceType
	var data []model.DeviceType
	if offsetStr == "" || limitStr == "" {
		data, err = deviceType.All(keys)
	} else {
		data, err = deviceType.Pages(offset, limit, keys)
	}
	if err != nil {
		log.Println(err)
		respData.Error = "获取设备类型列表失败"
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 获取数据总数
	num, _ := deviceType.Count(keys)

	res := struct {
		Num  int64              `json:"num"`
		Data []model.DeviceType `json:"data"`
	}{
		Num:  num,
		Data: data,
	}

	respData.Data = res
	c.JSON(http.StatusOK, &respData)
	return
}

//获取单个设备类型详情
func deviceTypeById(c *gin.Context, id string) {
	var respData RespData //返回的结构体
	var dt model.DeviceType
	var err error = nil
	dt.TypeId, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		respData.Error = "设备类型ID必须是数字"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	data, err := dt.One()
	if err != nil {
		log.Println(err)
		respData.Error = "设备类型不存在"
		c.JSON(http.StatusOK, &respData)
		return
	}

	respData.Data = data
	c.JSON(http.StatusOK, &respData)
	return
}
