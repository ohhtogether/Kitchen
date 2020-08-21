package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"src/model"
)

// 地图小区设备定位数据
type MapData struct {
	Area   []MapArea       `json:"area"`
	Device []MapDeviceInfo `json:"device"`
}
type MapArea struct {
	Id         int64  `json:"id"`
	Areaid     int64  `json:"areaid"`
	Areaname   string `json:"areaname"`
	Coordinate string `json:"coordinate"`
	Attribute  int64  `json:"attribute"`
	Icon       string `json:"icon"`
}
type MapDeviceInfo struct {
	Deviceid   int64  `json:"deviceid"`
	Devicename string `json:"devicename"`
	Typeid     int64  `json:"typeid"`
	Coordinate string `json:"coordinate"`
	Video      string `json:"video"`
}

// 获取所有数据
func Housedistrict(c *gin.Context) {
	attribute := c.Query("attribute")
	if attribute == "" {
		GetAll(c)
	} else {
		//GetAllByAttribute(c, attribute)
	}
}

// 地图小区设备定位数据(所有数据)
func GetAll(c *gin.Context) {
	var respData RespData
	// 获取区域数据
	var data MapData
	var rest model.Restaurnt
	restall, err := rest.All()
	if err != nil {
		log.Println(err)
		respData.Error = "获取餐厅数据失败"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}
	//fmt.Println("aaa")
	//return
	var maparea []MapArea
	for _, v := range restall {
		var tmp MapArea
		tmp.Areaid = v.Id
		tmp.Areaname = v.Name
		tmp.Coordinate = fmt.Sprintf("%v,%v", v.Longitude, v.Latitude)
		maparea = append(maparea, tmp)
	}
	data.Area = maparea
	// 获取设备数据
	var ds model.DeviceInfo
	devdata, err := ds.Alls()
	if err != nil {
		log.Println(err)
		respData.Error = "获取数据失败"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}
	var devicedata []MapDeviceInfo
	for _, v := range devdata {
		var tmp MapDeviceInfo
		tmp.Deviceid = v.Deviceid
		tmp.Devicename = v.Devicename
		tmp.Typeid = v.Typeid
		tmp.Coordinate = v.Coordinate
		tmp.Video = v.Video
		devicedata = append(devicedata, tmp)
	}
	data.Device = devicedata
	respData.Data = data
	c.JSON(http.StatusOK, &respData)
	return
}

// 通过attribute获取指定的所有数据
func GetAllByAttribute(c *gin.Context, attr string) {
	var respData RespData
	attribute, err := strconv.ParseInt(attr, 10, 64)
	if err != nil || attribute != 0 && attribute != 1 && attribute != 2 {
		log.Println(err)
		respData.Error = "attribute必须是数字0、1、2"
		c.JSON(http.StatusOK, &respData)
		return
	}
	var hd model.Housedistrict
	data, err := hd.All(attr)
	respData.Data = data
	c.JSON(http.StatusOK, &respData)
	return
}

// 地图设备数据
type MapCoordinate struct {
	Deviceid     int64  `json:"deviceid"`
	Devicename   string `json:"devicename"`
	Address      string `json:"address"`
	Typeid       int64  `json:"typeid"`
	State        int64  `json:"state"`
	StatusTypeId int64  `json:"statusTypeId"`
	TypeDes      string `json:"typeDes"`
	Time         string `json:"time"`
	Vedio        string `json:"vedio"`
}

// 地图设备数据
func Coordinate(c *gin.Context) {
	var respData RespData
	deviceidStr := c.Query("deviceid")
	if deviceidStr == "" {
		respData.Error = "设备ID不能为空"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}
	deviceid, err := strconv.ParseInt(deviceidStr, 10, 64)
	if err != nil || deviceid <= 0 {
		respData.Error = "设备ID必须是数字，并且大于0"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	var device model.DeviceInfo
	device.Deviceid = deviceid
	deviceOne, err := device.GetOneDataByDeviceId()
	if err != nil {
		respData.Error = "设备ID不存在"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	// 返回数据
	var data MapCoordinate
	data.Deviceid = deviceOne.Deviceid
	data.Devicename = deviceOne.Devicename
	data.Address = deviceOne.Ip
	data.Typeid = deviceOne.Typeid
	data.Vedio = deviceOne.Video

	// 获取设备在线/离线状态
	var do model.DeviceOnline
	do.DeviceId = deviceOne.Deviceid
	online, err := do.OneByDeviceid()
	if err != nil {
		log.Println(err)
		fmt.Println(err)
		respData.Error = "获取设备在/离线状态失败"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}
	data.State = online.State

	// 获取设备最新状态
	var ds model.DeviceStatus
	ds.DeviceId = deviceOne.Parentid
	ds.Channel = deviceOne.Channel
	maxid, err := ds.GetMaxId()
	if err != nil || maxid == 0 {
		data.StatusTypeId = 0
		data.TypeDes = ""
		data.Time = ""
	} else {
		ds.Id = maxid
		status, err := ds.One()
		if err != nil {
			log.Println(err)
			respData.Error = "获取设备状态出错"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
		data.StatusTypeId = status.StatusTypeId
		timedata, err := strconv.ParseInt(status.Time, 10, 64)
		tm := time.Unix(timedata, 0)
		data.Time = tm.Format("2006-01-02 15:04:05")
		var st model.StatusType
		st.Id = status.StatusTypeId
		statusType, err := st.One()
		if err != nil {
			data.TypeDes = ""
		} else {
			data.TypeDes = statusType.Des
		}
	}

	respData.Data = data
	c.JSON(http.StatusOK, &respData)
	return
}
