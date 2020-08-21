package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"src/model"
)

/*//获取设备信息
func StatusTypeView(c *gin.Context) {
	//判断是请求角色列表还是请求指定角色信息
	id := c.Query("id")
	if id == "" {
		//请求角色list
		statusTypePage(c)
	} else {
		//请求角色详情
		statusTypeById(c, id)
	}
}

//获取单个设备类型详情
func statusTypeById(c *gin.Context, id string) {
	var respData RespData //返回的结构体
	var st model.StatusType
	var err error = nil
	st.Id, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		respData.Error = "设备类型ID必须是数字"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	data, err := st.One()
	if err != nil {
		log.Println(err)
		respData.Error = "设备类型详情不存在"
		c.JSON(http.StatusOK, &respData)
		return
	}
	respData.Data = data
	c.JSON(http.StatusOK, &respData)
	return
}

// 分页
func statusTypePage(c *gin.Context) {
	var respData RespData
	// 偏移量
	var offset int64
	var err error
	offsetStr := c.Query("offset")
	if offsetStr == "" {
		offset = 0
	} else {
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
	if limitStr == "" {
		limit = 5
	} else {
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
	//fmt.Println(keys)

	// 获取数据
	var statusType model.StatusType
	data, err := statusType.Pages(offset, limit, keys)
	if err != nil {
		respData.Error = "获取设备类型详情列表失败"
		c.JSON(http.StatusOK, &respData)
		return
	}

	respData.Data = data
	c.JSON(http.StatusOK, &respData)
	return
}*/

// 下拉列表 id、des
func StatusTypePullDown(c *gin.Context) {
	var respData RespData
	var st model.StatusType
	sts, err := st.All()
	if err != nil {
		respData.Error = "获取状态类型数据失败"
		c.JSON(http.StatusOK, &respData)
		return
	}

	type Sons struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}
	type Parent struct {
		Sons
		Son []Sons `json:"son"`
	}

	var data []Parent
	for _, v := range sts {
		var tmp Parent
		if v.Parentid == "0" {
			tmp.Id = v.Id
			tmp.Name = v.Des
			data = append(data, tmp)
		}
	}

	for _, v := range sts {
		for n, m := range data {
			var son Sons
			if v.Parentid == "0x"+fmt.Sprint(m.Id) {
				son.Id = v.Id
				son.Name = v.Des
				data[n].Son = append(data[n].Son, son)
			}
		}
	}

	//fmt.Println(data)
	respData.Data = data
	c.JSON(http.StatusOK, &respData)
	return

}
