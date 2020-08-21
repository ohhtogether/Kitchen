package controller

import (
	"log"
	"net/http"
	"src/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PurchaseList(c *gin.Context) {
	purchaseid := c.Query("id")
	if purchaseid == "" {
		purchasePage(c)
	} else {
		purchaseOne(c, purchaseid)
	}
}

// 分页
func purchasePage(c *gin.Context) {
	var respData RespData
	// 接收数据
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
	startTime := c.Query("start") //搜索开始时间
	endTime := c.Query("end")     //搜索结束时间
	ridStr := c.Query("rid")
	// 验证
	var err error
	//var offset int64
	if offsetStr != "" {
		_, err = strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			respData.Error = "偏移量必须是数字"
			c.JSON(http.StatusOK, &respData)
			return
		}
	}
	//var limit int64
	if limitStr != "" {
		_, err = strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			respData.Error = "每页显示的条数必须是数字"
			c.JSON(http.StatusOK, &respData)
			return
		}
	}

	var start, end int64
	if startTime != "" && endTime != "" {
		start, err = strconv.ParseInt(startTime, 10, 64)
		if err != nil {
			respData.Error = "时间必须是数字"
			c.JSON(http.StatusOK, &respData)
			return
		}
		end, err = strconv.ParseInt(endTime, 10, 64)
		if err != nil {
			respData.Error = "时间必须是数字"
			c.JSON(http.StatusOK, &respData)
			return
		}
		if start >= end {
			respData.Error = "查询开始时间必须小于查询结束时间"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
	}

	//if startTime != "" && endTime != "" { //如果提交的时间参数中没有值则表示查询时不带此条件
	//	loc, _ := time.LoadLocation("Local")
	//	start_time, errStart := time.ParseInLocation("2006-1-2 15:4:5", startTime, loc)
	//	end_time, errEnd := time.ParseInLocation("2006-1-2 15:4:5", endTime, loc)
	//	if errStart != nil || errEnd != nil {
	//		respData.Error = "提交的时间错误"
	//		c.JSON(http.StatusBadRequest, &respData)
	//		return
	//	}
	//
	//	if start_time.Unix() >= end_time.Unix() {
	//		respData.Error = "查询开始时间必须小于查询结束时间"
	//		c.JSON(http.StatusBadRequest, &respData)
	//		return
	//	}
	//
	//	start = start_time.Unix()
	//	end = end_time.Unix()
	//}

	rid, err := strconv.ParseInt(ridStr, 10, 64)
	if err != nil || rid <= 0 {
		respData.Error = "餐厅ID必须是数字，并且大于0"
		c.JSON(http.StatusOK, &respData)
		return
	}
	var rest model.Restaurnt
	rest.Id = rid
	_, err = rest.One()
	if err != nil {
		respData.Error = "餐厅ID不存在"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	// 获取数据
	var purchase model.Purchase
	var data []model.Purchase
	if offsetStr == "" || limitStr == "" {
		data, err = purchase.GetDataByRid(rid, start, end)
	} else {
		data, err = purchase.Page(rid, start, end, offsetStr, limitStr)
	}
	if err != nil {
		log.Println(err)
		respData.Error = "获取数据列表失败"
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 获取数据总数
	num, err := purchase.Total(rid, start, end)
	if err != nil {
		log.Println(err)
		respData.Error = "获取数据总条数失败"
		c.JSON(http.StatusOK, &respData)
		return
	}

	res := struct {
		Num  int64            `json:"num"`
		Data []model.Purchase `json:"data"`
	}{
		Num:  num,
		Data: data,
	}

	respData.Data = res
	c.JSON(http.StatusOK, &respData)
	return

}

//获取一条设备信息数据
func purchaseOne(c *gin.Context, purchaseidStr string) {
	var respData RespData
	// 接收数据
	var purchase model.Purchase
	// 验证数字类型
	purchaseid, err := strconv.ParseInt(purchaseidStr, 10, 64)
	if err != nil {
		respData.Error = "ID必须是数字"
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 通过ID查询信息
	purchase.Id = purchaseid
	data, err := purchase.One()
	if err != nil {
		respData.Error = "数据不存在"
		c.JSON(http.StatusOK, &respData)
		return
	}

	respData.Data = data
	c.JSON(http.StatusOK, &respData)
	return
}

func PurchaseAdd(c *gin.Context) {
	var respData RespData

	// 接收数据
	data := c.Request.FormValue("goods")
	datetime := c.Request.FormValue("time")
	ridStr := c.Request.FormValue("rid")

	//loc, _ := time.LoadLocation("Local")
	//location, err := time.ParseInLocation("2006-1-2 15:4:5", datetime, loc)
	//if err != nil {
	//	respData.Error = "提交的时间错误"
	//	c.JSON(http.StatusBadRequest, &respData)
	//	return
	//}
	//times := location.Unix()

	times, err := strconv.ParseInt(datetime, 10, 64)
	if err != nil {
		respData.Error = "时间格式不正确"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	rid, err := strconv.ParseInt(ridStr, 10, 64)
	if err != nil || rid <= 0 {
		respData.Error = "填写正确的餐厅ID参数！"
		c.JSON(http.StatusOK, &respData)
		return
	}
	var rest model.Restaurnt
	rest.Id = rid
	_, err = rest.One()
	if err != nil {
		respData.Error = "餐厅ID不存在"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	var p model.Purchase
	p.Rid = rid
	p.Time = times
	p.Goods = data
	p.Status = 1
	err = p.Add()
	if err != nil {
		respData.Error = "添加失败"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	respData.Data = "添加成功"
	c.JSON(http.StatusOK, &respData)
	return
}

func PurchaseEdit(c *gin.Context) {
	var respData RespData

	// 接收数据
	data := c.Request.FormValue("goods")
	datetime := c.Request.FormValue("time")
	idStr := c.Request.FormValue("id")

	//loc, _ := time.LoadLocation("Local")
	//location, err := time.ParseInLocation("2006-1-2 15:4:5", datetime, loc)
	//if err != nil {
	//	respData.Error = "提交的时间错误"
	//	c.JSON(http.StatusBadRequest, &respData)
	//	return
	//}
	//times := location.Unix()

	times, err := strconv.ParseInt(datetime, 10, 64)
	if err != nil {
		respData.Error = "时间格式不正确"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		respData.Error = "填写正确的ID参数！"
		c.JSON(http.StatusOK, &respData)
		return
	}
	var p model.Purchase
	p.Id = id
	one, err := p.One()
	if err != nil {
		respData.Error = "ID不存在"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	//var p model.Purchase
	one.Goods = data
	one.Time = times
	one.Status = 1
	err = one.Edit()
	if err != nil {
		respData.Error = "修改失败"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	respData.Data = "修改成功"
	c.JSON(http.StatusOK, &respData)
	return
}
func PurchaseDel(c *gin.Context) {
	var respData RespData

	idStr := c.Request.FormValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		respData.Error = "填写正确的ID参数！"
		c.JSON(http.StatusOK, &respData)
		return
	}
	var p model.Purchase
	p.Id = id
	one, err := p.One()
	if err != nil {
		respData.Error = "ID不存在"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	one.Status = 0
	err = one.Edit()
	if err != nil {
		respData.Error = "删除失败"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	respData.Data = "删除成功"
	c.JSON(http.StatusOK, &respData)
	return
}
