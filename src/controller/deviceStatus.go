package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"src/model"
	"strconv"
	"strings"
	"time"
)

func DeviceStatusView(c *gin.Context) {
	deviceid := c.Query("deviceid")
	channel := c.Query("channel")
	if deviceid != "" && channel != "" {
		//指定设备的全部状态数据
		statusDevLog(c, deviceid, channel)
	} else {
		//设备最新的状态信息
		statusDevNew(c)
	}
}

//获取所有设备最新的状态信息
func statusDevNew(c *gin.Context) {
	var respData RespData
	// 接收数据
	offsetStr := c.Query("offset") //接受offset limit keys对数据进行分页和搜索，获取数据总条数
	limitStr := c.Query("limit")
	ridStr := c.Query("rid")                   //小区ID
	typeidStr := c.Query("typeid")             //设备类型ID
	statustypeidStr := c.Query("statustypeid") //状态类型ID
	startTime := c.Query("start")              //搜索开始时间
	endTime := c.Query("end")                  //搜索结束时间
	//搜索的数据
	var keys model.StatusSearchKeys
	// 基本数据验证
	// 分页
	if offsetStr != "" && limitStr != "" {
		//获取查询条件下的分页数据
		_, errOf := strconv.ParseInt(offsetStr, 10, 64)
		_, errLi := strconv.ParseInt(limitStr, 10, 64)
		if errOf != nil || errLi != nil {
			respData.Error = "参数错误"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
	}
	//小区ID、设备类型ID、//状态类型ID
	var err error
	var statustypeid int64
	if ridStr != "" {
		_, err = strconv.ParseInt(ridStr, 10, 64)
		if err != nil {
			respData.Error = "小区ID必须是数字"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
	}
	if typeidStr != "" {
		_, err = strconv.ParseInt(typeidStr, 10, 64)
		if err != nil {
			respData.Error = "设备类型ID必须是数字"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
	}
	if statustypeidStr != "" {
		statustypeid, err = strconv.ParseInt(statustypeidStr, 10, 64)
		if err != nil {
			respData.Error = "状态类型ID必须是数字"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
	}
	// 时间验证
	keys.StartTime, keys.EndTime, err = isSeachTimeIegal(startTime, endTime)
	if err != nil {
		respData.Error = "提交的搜索时间参数错误"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}
	// 数据库查询验证
	//餐厅ID、设备类型ID、//状态类型ID
	keys.Rid = ridStr

	if typeidStr != "" {
		keys.TypeId = typeidStr
	}
	if statustypeidStr != "" {
		var st model.StatusType
		st.Id = statustypeid
		stdata, err := st.One()
		if err != nil {
			respData.Error = "状态类型ID不存在"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
		if stdata.Parentid == "0" {
			var ids []int64
			st.Parentid = "0x" + fmt.Sprint(stdata.Id)
			ids, err := st.GetDataByParentid()
			if err != nil || len(ids) == 0 {
				respData.Error = "获取状态类型子集失败"
				c.JSON(http.StatusOK, &respData)
				return
			}
			keys.StatusTypeIds = ids
			keys.StatusTypeJudge = 2
		} else {
			keys.StatusTypeId = statustypeidStr
			keys.StatusTypeJudge = 1
		}
	}

	// 获取最新设备状态IDS (MAX(`id`))
	// 1.获取状态类型IDS (status=1)
	var st model.StatusType
	stData, err := st.All()
	if err != nil {
		respData.Error = "获取状态类型失败"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}
	var stIds []int64
	for _, v := range stData {
		stIds = append(stIds, v.Id)
	}
	stIdsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(stIds)), ","), "[]")

	// 2.通过stIds 获取最新设备状态IDS (MAX(`id`))
	var ds model.DeviceStatus
	dsIds, err := ds.GetDataByMaxId(stIdsStr)
	if err != nil {
		respData.Error = "获取设备状态失败"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}
	dsIdsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(dsIds)), ","), "[]")
	// 根据最新设备状态IDS获取设备
	data, err := ds.GetLatestStatusDeviceInfo(dsIdsStr, offsetStr, limitStr, keys)
	if err != nil {
		respData.Error = "获取设备状态失败2"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}
	//处理data数据
	for i, v := range data {
		// 获取状态类型数据
		for _, n := range stData {
			if v.StatusTypeId == n.Id {
				data[i].TypeName = n.StatusName
				data[i].TypeDes = n.Des
				data[i].TypeBrand = n.Brand
				if n.Parentid != "0" {
					str := n.Parentid[len("0x"):len([]byte(n.Parentid))]
					data[i].ParentTypeId, err = strconv.ParseInt(str, 10, 64)
					if err != nil {
						respData.Error = "状态类型数据转换错误"
						c.JSON(http.StatusBadRequest, &respData)
						return
					}
				}
			}
		}
		if data[i].ParentTypeId != 0 {
			for _, n := range stData {
				if data[i].ParentTypeId == n.Id {
					data[i].ParentTypeName = n.StatusName
					data[i].ParentTypeDes = n.Des
				}
			}
		}

		// 时间戳转换日期格式
		timestamp, err := strconv.ParseInt(data[i].Time, 10, 64)
		if err != nil {
			respData.Error = "时间格式转换错误"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
		tm := time.Unix(timestamp, 0)
		data[i].Time = tm.Format("2006-01-02 15:04:05")
	}

	num, err := model.GetDeviceStatusNums(dsIdsStr, keys)
	if err != nil {
		log.Println(err)
		respData.Error = err.Error()
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	//返回数据给前端
	respData.Data = struct {
		Num  int64
		Data []model.DeviceStatus
	}{
		Num:  num,
		Data: data,
	}
	c.JSON(http.StatusOK, &respData)
	return
}

//获取指定设备的全部状态数据
func statusDevLog(c *gin.Context, deviceidStr, channelStr string) {
	var respData RespData
	var err error = nil
	// 接收数据
	offsetStr := c.Query("offset") //接受offset limit keys对数据进行分页和搜索，获取数据总条数
	limitStr := c.Query("limit")
	statustypeidStr := c.Query("statustypeid") //状态类型ID
	startTime := c.Query("start")              //搜索开始时间
	endTime := c.Query("end")                  //搜索结束时间
	//搜索的数据
	var keys model.StatusSearchKeys
	// 基本数据验证
	// deviceid-channel
	deviceid, errDid := strconv.ParseInt(deviceidStr, 10, 64)
	channel, errChl := strconv.ParseInt(channelStr, 10, 64)
	if errDid != nil || errChl != nil {
		respData.Error = "设备ID和通道ID参数不合法"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}
	// 分页
	if offsetStr != "" && limitStr != "" {
		//获取查询条件下的分页数据
		_, errOf := strconv.ParseInt(offsetStr, 10, 64)
		_, errLi := strconv.ParseInt(limitStr, 10, 64)
		if errOf != nil || errLi != nil {
			respData.Error = "分页参数不合法"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
	}
	//状态类型ID
	var statustypeid int64
	if statustypeidStr != "" {
		statustypeid, err = strconv.ParseInt(statustypeidStr, 10, 64)
		if err != nil {
			respData.Error = "状态类型ID必须是数字"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
	}
	// 时间验证
	keys.StartTime, keys.EndTime, err = isSeachTimeIegal(startTime, endTime)
	if err != nil {
		respData.Error = "提交的搜索时间参数错误"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}
	// 数据库查询验证
	// deviceid-channel
	var di model.DeviceInfo
	di.Parentid = deviceid
	di.Channel = channel
	deviceinfo, err := di.GetDataByDeviceId()
	if err != nil {
		respData.Error = "通过deviceID与channel查询，设备不存在"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}
	//状态类型ID
	if statustypeidStr != "" {
		var st model.StatusType
		st.Id = statustypeid
		stdata, err := st.One()
		if err != nil {
			respData.Error = "状态类型ID不存在"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
		if stdata.Parentid == "0" {
			var ids []int64
			st.Parentid = "0x" + fmt.Sprint(stdata.Id)
			ids, err := st.GetDataByParentid()
			if err != nil || len(ids) == 0 {
				respData.Error = "获取状态类型子集失败"
				c.JSON(http.StatusOK, &respData)
				return
			}
			keys.StatusTypeIds = ids
			keys.StatusTypeJudge = 2
		} else {
			keys.StatusTypeId = statustypeidStr
			keys.StatusTypeJudge = 1
		}
	}
	// 获取所有设备状态
	var ds model.DeviceStatus
	ds.DeviceId = deviceid
	ds.Channel = channel
	data, err := ds.DevStuLogs(offsetStr, limitStr, keys)
	// 获取状态类型
	var st model.StatusType
	stData, err := st.All()
	if err != nil {
		respData.Error = "获取状态类型失败"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}
	//处理data数据
	for i, v := range data {
		data[i].DeviceName = deviceinfo.Devicename
		data[i].DeviceIP = deviceinfo.Ip
		for _, n := range stData {
			if v.StatusTypeId == n.Id {
				data[i].TypeName = n.StatusName
				data[i].TypeDes = n.Des
				data[i].TypeBrand = n.Brand
				if n.Parentid != "0" {
					str := n.Parentid[len("0x"):len([]byte(n.Parentid))]
					data[i].ParentTypeId, err = strconv.ParseInt(str, 10, 64)
					if err != nil {
						respData.Error = "状态类型数据转换错误"
						c.JSON(http.StatusBadRequest, &respData)
						return
					}
				}
			}
		}
		if data[i].ParentTypeId != 0 {
			for _, n := range stData {
				if data[i].ParentTypeId == n.Id {
					data[i].ParentTypeName = n.StatusName
					data[i].ParentTypeDes = n.Des
				}
			}
		}

		// 时间戳转换日期格式
		timestamp, err := strconv.ParseInt(data[i].Time, 10, 64)
		if err != nil {
			respData.Error = "时间格式转换错误"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
		tm := time.Unix(timestamp, 0)
		data[i].Time = tm.Format("2006-01-02 15:04:05")
	}

	num, err := ds.GetDeviceStatusNumsByDevId(keys)
	if err != nil {
		log.Println(err)
		respData.Error = err.Error()
		c.JSON(http.StatusBadRequest, &respData)
		return
	}
	//返回数据给前端
	respData.Data = struct {
		Num  int64
		Data []model.DeviceStatus
	}{
		Num:  num,
		Data: data,
	}
	c.JSON(http.StatusOK, &respData)
	return
}

/***********************************************************************************************************/

//获取指定设备的全部状态数据 =========》验证时间函数
func isSeachTimeIegal(startTime, endTime string) (stTime string, edTime string, err error) {
	stTime, edTime, err = "", "", nil     //初始化一下
	if startTime != "" && endTime != "" { //如果提交的时间参数中没有值则表示查询时不带此条件
		loc, _ := time.LoadLocation("Local")
		start_time, errT := time.ParseInLocation("2006-01-02 15:04:05", startTime, loc)
		if errT != nil {
			err = errT
			return
		}
		end_time, errT := time.ParseInLocation("2006-01-02 15:04:05", endTime, loc)
		if errT != nil {
			err = errT
			return
		}
		if start_time.Unix() > end_time.Unix() {
			err = errors.New("查询开始时间必须小于查询结束时间")
			return
		}
		stTime = fmt.Sprint(start_time.Unix())
		edTime = fmt.Sprint(end_time.Unix())
	}
	return
}
