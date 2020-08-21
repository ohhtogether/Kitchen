package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dxvgef/filter"
	"github.com/gin-gonic/gin"

	"src/conf"
	"src/global"
	"src/model"
)

//新增一条设备信息
func DeviceInfoAddOne(c *gin.Context) {
	var respData RespData

	// 接收数据
	// 必填
	devicename := c.Request.FormValue("devicename")
	typeidStr := c.Request.FormValue("typeid")
	ridStr := c.Request.FormValue("rid")
	brand := c.Request.FormValue("brand")
	mode := c.Request.FormValue("model")
	serial_number := c.Request.FormValue("serial_number")
	ip := c.Request.FormValue("ip")
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	portStr := c.Request.FormValue("port")
	channelStr := c.Request.FormValue("channel")
	parentidStr := c.Request.FormValue("parentid")
	// 选填
	manufacture_date := c.Request.FormValue("manufacture_date")
	installation_date := c.Request.FormValue("installation_date")
	acceptance_date := c.Request.FormValue("acceptance_date")
	inspection_date := c.Request.FormValue("inspection_date")
	maintenance_date := c.Request.FormValue("maintenance_date")
	longitude := c.Request.FormValue("longitude")
	latitude := c.Request.FormValue("latitude")
	vedio := c.Request.FormValue("vedio")

	// 验证数据
	// 基本数据验证
	var deviceInfo model.DeviceInfo
	err := filter.MSet(
		filter.El(&deviceInfo.Devicename,
			filter.FromString(devicename, "设备名称").
				Required().IsLetterOrDigit("设备名称不能含有特殊字符"),
		),
		filter.El(&typeidStr,
			filter.FromString(typeidStr, "设备类型ID").
				Required().IsDigit().MinInteger(1),
		),
		filter.El(&ridStr,
			filter.FromString(ridStr, "设备区域ID").
				Required().IsDigit(),
		),
		filter.El(&deviceInfo.Brand,
			filter.FromString(brand, "设备品牌").
				Required().IsLetterOrDigit("设备品牌不能含有特殊字符"),
		),
		filter.El(&mode,
			filter.FromString(mode, "设备型号").
				Required(),
		),
		filter.El(&serial_number,
			filter.FromString(serial_number, "设备序列号").
				Required(),
		),
		filter.El(&deviceInfo.Ip,
			filter.FromString(ip, "设备ip").
				Required().IsIP(),
		),
		filter.El(&username,
			filter.FromString(username, "设备登录账号").
				Required().IsLetterOrDigit(),
		),
		filter.El(&password,
			filter.FromString(password, "设备登录密码").
				Required(),
		),
		filter.El(&portStr,
			filter.FromString(portStr, "设备端口").
				Required().MinInteger(1).MaxInteger(65535),
		),
		filter.El(&channelStr,
			filter.FromString(channelStr, "设备通道号").
				Required().IsDigit(),
		),
		filter.El(&parentidStr,
			filter.FromString(parentidStr, "设备父级ID").
				IsDigit(),
		),
		filter.El(&deviceInfo.ManufactureDate,
			filter.FromString(manufacture_date, "设备出厂日期").
				IsDigit(),
		),
		filter.El(&deviceInfo.InstallationDate,
			filter.FromString(installation_date, "设备安装日期").
				IsDigit(),
		),
		filter.El(&deviceInfo.AcceptanceDate,
			filter.FromString(acceptance_date, "设备验收日期").
				IsDigit(),
		),
		filter.El(&deviceInfo.InspectionDate,
			filter.FromString(inspection_date, "设备巡检日期").
				IsDigit(),
		),
		filter.El(&deviceInfo.MaintenanceDate,
			filter.FromString(maintenance_date, "设备维修日期").
				IsDigit(),
		),
	)
	if err != nil {
		respData.Error = err.Error()
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 类型ID
	typeid, err := strconv.ParseInt(typeidStr, 10, 64)
	// 区域ID
	rid, _ := strconv.ParseInt(ridStr, 10, 64)
	// 设备型号
	global.FilteredSQLInject(mode) // 过滤sql注入
	deviceInfo.Model = mode
	// 设备序列号
	global.FilteredSQLInject(serial_number)
	deviceInfo.SerialNumber = serial_number
	// 设备登录账号
	global.FilteredSQLInject(username)
	deviceInfo.Username = username
	// 设备登录密码
	global.FilteredSQLInject(password)
	deviceInfo.Password = password
	// 设备端口
	deviceInfo.Port, _ = strconv.ParseInt(portStr, 10, 64)
	// 通道号
	deviceInfo.Channel, _ = strconv.ParseInt(channelStr, 10, 64)
	// 设备父级ID
	parentid, _ := strconv.ParseInt(parentidStr, 10, 64)
	// 设备经度
	if longitude != "" {
		Bool := global.IsNumeric(longitude)
		if Bool != true {
			respData.Error = "设备经度格式不正确"
			c.JSON(http.StatusOK, &respData)
			return
		}
		deviceInfo.Longitude = longitude
	}
	// 设备纬度
	if latitude != "" {
		Bool := global.IsNumeric(latitude)
		if Bool != true {
			respData.Error = "设备纬度格式不正确"
			c.JSON(http.StatusOK, &respData)
			return
		}
		deviceInfo.Latitude = latitude
	}

	// 数据库验证
	// 验证设备类型ID是否存在
	var deviceType model.DeviceType
	deviceType.TypeId = typeid
	_, err = deviceType.One()
	if err != nil {
		respData.Error = "设备类型ID不存在"
		c.JSON(http.StatusOK, &respData)
		return
	}
	deviceInfo.Typeid = typeid
	// 验证设备区域ID是否存在
	var rest model.Restaurnt
	rest.Id = rid
	_, err = rest.One()
	if err != nil {
		respData.Error = "餐厅ID不存在"
		c.JSON(http.StatusOK, &respData)
		return
	}
	deviceInfo.Rid = rid
	// 验证设备父级ID
	if parentid != 0 {
		var di model.DeviceInfo
		di.Deviceid = parentid
		_, err := di.GetDataByDeviceId()
		if err != nil {
			respData.Error = "设备父级ID不存在"
			c.JSON(http.StatusOK, &respData)
			return
		}
		deviceInfo.Parentid = parentid
	}

	// 状态
	deviceInfo.Status = conf.DEVICE_STATUS_NORMAL
	deviceInfo.Video = vedio

	// 入库
	id, err := deviceInfo.AddOne()
	err = err
	if err != nil {
		respData.Data = "添加失败"
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 新增设备在/离线状态
	var do model.DeviceOnline
	do.DeviceId = id
	do.State = 1
	err = do.Add()
	if err != nil {
		log.Println(err)
		respData.Data = "新增设备状态设置失败"
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 修改新增设备父级ID
	if parentid == 0 {
		deviceInfo.Deviceid = id
		deviceInfo.Parentid = id
		err := deviceInfo.Edit()
		if err != nil {
			log.Println(err)
			respData.Data = "新增设备父级ID设置失败"
			c.JSON(http.StatusOK, &respData)
			return
		}
	}

	respData.Data = "添加成功"
	c.JSON(http.StatusOK, &respData)
	return
}

//修改一条设备信息
func DeviceInfoEdit(c *gin.Context) {
	var respData RespData
	// 接收数据
	var deviceInfo model.DeviceInfo
	// 区域 ID
	deviceid := c.PostForm("deviceid")
	// 验证不为空
	if deviceid == "" {
		respData.Error = "设备ID不能为空"
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 验证数字类型
	deviceId, err := strconv.ParseInt(deviceid, 10, 64)
	if err != nil {
		respData.Error = "设备ID必须是数字"
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 通过ID查询信息
	deviceInfo.Deviceid = deviceId
	data, err := deviceInfo.GetOneDataByDeviceId()
	if err != nil {
		respData.Error = "设备信息不存在"
		c.JSON(http.StatusOK, &respData)
		return
	}
	deviceInfo = data

	// 接收数据
	// 必填
	devicename := c.Request.FormValue("devicename")
	global.FilteredSQLInject(devicename)
	typeidStr := c.Request.FormValue("typeid")
	ridStr := c.Request.FormValue("rid")
	brand := c.Request.FormValue("brand")
	mode := c.Request.FormValue("model")
	serial_number := c.Request.FormValue("serial_number")
	ip := c.Request.FormValue("ip")
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	portStr := c.Request.FormValue("port")
	channelStr := c.Request.FormValue("channel")
	parentidStr := c.Request.FormValue("parentid")
	// 选填
	manufacture_date := c.Request.FormValue("manufacture_date")
	installation_date := c.Request.FormValue("installation_date")
	acceptance_date := c.Request.FormValue("acceptance_date")
	inspection_date := c.Request.FormValue("inspection_date")
	maintenance_date := c.Request.FormValue("maintenance_date")
	longitude := c.Request.FormValue("longitude")
	latitude := c.Request.FormValue("latitude")
	vedio := c.Request.FormValue("vedio")

	// 验证数据
	// 基本数据验证
	err = filter.MSet(
		filter.El(&deviceInfo.Devicename,
			filter.FromString(devicename, "设备名称").
				Required(),
		),
		filter.El(&typeidStr,
			filter.FromString(typeidStr, "设备类型ID").
				Required().IsDigit().MinInteger(1),
		),
		filter.El(&ridStr,
			filter.FromString(ridStr, "设备区域ID").
				Required().IsDigit(),
		),
		filter.El(&deviceInfo.Brand,
			filter.FromString(brand, "设备品牌").
				Required().IsLetterOrDigit("设备品牌不能含有特殊字符"),
		),
		filter.El(&mode,
			filter.FromString(mode, "设备型号").
				Required(),
		),
		filter.El(&serial_number,
			filter.FromString(serial_number, "设备序列号").
				Required(),
		),
		filter.El(&deviceInfo.Ip,
			filter.FromString(ip, "设备ip").
				Required().IsIP(),
		),
		filter.El(&username,
			filter.FromString(username, "设备登录账号").
				Required().IsLetterOrDigit(),
		),
		filter.El(&password,
			filter.FromString(password, "设备登录密码").
				Required(),
		),
		filter.El(&portStr,
			filter.FromString(portStr, "设备端口").
				Required().MinInteger(1).MaxInteger(65535),
		),
		filter.El(&channelStr,
			filter.FromString(channelStr, "设备通道号").
				Required().IsDigit(),
		),
		filter.El(&parentidStr,
			filter.FromString(parentidStr, "设备父级ID").
				Required().IsDigit(),
		),
		filter.El(&deviceInfo.ManufactureDate,
			filter.FromString(manufacture_date, "设备出厂日期").
				IsDigit(),
		),
		filter.El(&deviceInfo.InstallationDate,
			filter.FromString(installation_date, "设备安装日期").
				IsDigit(),
		),
		filter.El(&deviceInfo.AcceptanceDate,
			filter.FromString(acceptance_date, "设备验收日期").
				IsDigit(),
		),
		filter.El(&deviceInfo.InspectionDate,
			filter.FromString(inspection_date, "设备巡检日期").
				IsDigit(),
		),
		filter.El(&deviceInfo.MaintenanceDate,
			filter.FromString(maintenance_date, "设备维修日期").
				IsDigit(),
		),
	)
	if err != nil {
		respData.Error = err.Error()
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 类型ID
	typeid, _ := strconv.ParseInt(typeidStr, 10, 64)
	// 区域ID
	rid, _ := strconv.ParseInt(ridStr, 10, 64)
	// 设备型号
	global.FilteredSQLInject(mode) // 过滤sql注入
	deviceInfo.Model = mode
	// 设备序列号
	global.FilteredSQLInject(serial_number)
	deviceInfo.SerialNumber = serial_number
	// 设备登录账号
	global.FilteredSQLInject(username)
	deviceInfo.Username = username
	// 设备登录密码
	global.FilteredSQLInject(password)
	deviceInfo.Password = password
	// 设备端口
	deviceInfo.Port, _ = strconv.ParseInt(portStr, 10, 64)
	// 通道号
	deviceInfo.Channel, _ = strconv.ParseInt(channelStr, 10, 64)
	// 设备父级ID
	parentid, _ := strconv.ParseInt(parentidStr, 10, 64)
	if parentid == 0 {
		deviceInfo.Parentid = parentid
	}
	// 设备经度
	if longitude != "" {
		Bool := global.IsNumeric(longitude)
		if Bool != true {
			respData.Error = "设备经度格式不正确"
			c.JSON(http.StatusOK, &respData)
			return
		}
		deviceInfo.Longitude = longitude
	}
	// 设备纬度
	if latitude != "" {
		Bool := global.IsNumeric(latitude)
		if Bool != true {
			respData.Error = "设备纬度格式不正确"
			c.JSON(http.StatusOK, &respData)
			return
		}
		deviceInfo.Latitude = latitude
	}

	// 数据库验证
	// 验证设备类型ID是否存在
	var deviceType model.DeviceType
	deviceType.TypeId = typeid
	_, err = deviceType.One()
	if err != nil {
		respData.Error = "设备类型ID不存在"
		c.JSON(http.StatusOK, &respData)
		return
	}
	deviceInfo.Typeid = typeid
	// 验证设备区域ID是否存在
	var area model.Restaurnt
	area.Id = rid
	_, err = area.One()
	if err != nil {
		respData.Error = "餐厅ID不存在"
		c.JSON(http.StatusOK, &respData)
		return
	}
	deviceInfo.Rid = rid
	// 验证设备父级ID
	if parentid != 0 {
		var di model.DeviceInfo
		di.Deviceid = parentid
		_, err := di.GetDataByDevId()
		if err != nil {
			respData.Error = "设备父级ID不存在"
			c.JSON(http.StatusOK, &respData)
			return
		}
	}
	deviceInfo.Video = vedio

	// 修改数据库
	err = deviceInfo.Edit()
	if err != nil {
		respData.Error = "修改失败"
		c.JSON(http.StatusOK, &respData)
		return
	}

	respData.Data = "修改成功"
	c.JSON(http.StatusOK, &respData)
	return

}

//删除一条设备信息
func DeviceInfoDel(c *gin.Context) {

	var respData RespData
	// 接收数据
	var deviceInfo model.DeviceInfo
	// 设备 ID
	deviceid := c.PostForm("deviceid")
	// 验证不为空
	if deviceid == "" {
		respData.Error = "设备ID不能为空"
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 验证数字类型
	deviceId, err := strconv.ParseInt(deviceid, 10, 64)
	if err != nil {
		respData.Error = "设备ID必须是数字"
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 通过ID查询信息
	deviceInfo.Deviceid = deviceId
	data, err := deviceInfo.GetOneDataByDeviceId()
	if err != nil {
		respData.Error = "设备信息不存在"
		c.JSON(http.StatusOK, &respData)
		return
	}
	deviceInfo = data
	deviceInfo.Status = conf.DEVICE_STATUS_DELETE
	err = deviceInfo.Edit()
	if err != nil {
		respData.Error = "删除失败"
		c.JSON(http.StatusOK, &respData)
		return
	}

	respData.Data = "删除成功"
	c.JSON(http.StatusOK, &respData)
	return
}

//获取设备信息
func DeviceInfoView(c *gin.Context) {

	deviceid := c.Query("deviceid")
	if deviceid == "" {
		deviceInfoPage(c)
	} else {
		deviceInfoOne(c, deviceid)
	}
}

// 分页
func deviceInfoPage(c *gin.Context) {
	var respData RespData
	// 接收数据
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
	keys := c.Query("keys")
	ridStr := c.Query("rid")
	typeidStr := c.Query("typeid") //设备类型ID
	// 验证
	var err error
	var offset int64
	if offsetStr != "" {
		offset, err = strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			respData.Error = "偏移量必须是数字"
			c.JSON(http.StatusOK, &respData)
			return
		}
	}
	var limit int64
	if limitStr != "" {
		limit, err = strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			respData.Error = "每页显示的条数必须是数字"
			c.JSON(http.StatusOK, &respData)
			return
		}
	}
	if keys != "" {
		global.FilteredSQLInject(keys)
	}
	var rid int64
	if ridStr != "" {
		rid, err = strconv.ParseInt(ridStr, 10, 64)
		if err != nil || rid <= 0 {
			respData.Error = "餐厅ID必须是数字，并且大于0"
			c.JSON(http.StatusOK, &respData)
			return
		}
		var rest model.Restaurnt
		rest.Id = rid
		_, err := rest.One()
		if err != nil {
			respData.Error = "餐厅ID不存在"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
	}
	var typeid int64
	if typeidStr != "" {
		typeid, err = strconv.ParseInt(typeidStr, 10, 64)
		if err != nil {
			respData.Error = "设备类型ID必须是数字"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
	}

	if typeidStr != "" {
		var devtype model.DeviceType
		devtype.TypeId = typeid
		_, err := devtype.One()
		if err != nil {
			respData.Error = "设备类型ID不存在"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
	}
	var Sql string
	if ridStr != "" {
		Sql = fmt.Sprintf("and `rid`=%v", rid)
	}

	// 获取数据
	var deviceInfo model.DeviceInfo
	var data []model.DeviceInfo
	if offsetStr == "" || limitStr == "" {
		data, err = deviceInfo.GetDataByActiveUser(keys, Sql)
	} else {
		data, err = deviceInfo.PagesByActiveUser(offset, limit, keys, Sql)
	}
	if err != nil {
		log.Println(err)
		respData.Error = "获取设备信息列表失败"
		c.JSON(http.StatusOK, &respData)
		return
	}
	restMap, err := RestMap()
	if err != nil {
		respData.Error = err.Error()
		c.JSON(http.StatusOK, &respData)
		return
	}
	for i, dev := range data {
		tManufactureDate, _ := strconv.ParseInt(dev.ManufactureDate, 10, 64)
		if tManufactureDate > 0 {
			tm := time.Unix(tManufactureDate, 0)
			data[i].ManufactureDate = tm.Format("2006-01-02 03:04:05")
		}

		tInstallationDate, _ := strconv.ParseInt(dev.InstallationDate, 10, 64)
		if tInstallationDate > 0 {
			tm := time.Unix(tInstallationDate, 0)
			data[i].InstallationDate = tm.Format("2006-01-02 03:04:05")
		}

		tAcceptanceDate, _ := strconv.ParseInt(dev.AcceptanceDate, 10, 64)
		if tAcceptanceDate > 0 {
			tm := time.Unix(tAcceptanceDate, 0)
			data[i].AcceptanceDate = tm.Format("2006-01-02 03:04:05")
		}

		tInspectionDate, _ := strconv.ParseInt(dev.InspectionDate, 10, 64)
		if tInspectionDate > 0 {
			tm := time.Unix(tInspectionDate, 0)
			data[i].InspectionDate = tm.Format("2006-01-02 03:04:05")
		}

		tMaintenanceDate, _ := strconv.ParseInt(dev.MaintenanceDate, 10, 64)
		if tMaintenanceDate > 0 {
			tm := time.Unix(tMaintenanceDate, 0)
			data[i].MaintenanceDate = tm.Format("2006-01-02 03:04:05")
		}
		data[i].Rname = restMap[dev.Rid]
	}

	// 获取数据总数
	num, _ := deviceInfo.CountByActiveUser(keys, Sql)

	res := struct {
		Num  int64              `json:"num"`
		Data []model.DeviceInfo `json:"data"`
	}{
		Num:  num,
		Data: data,
	}

	respData.Data = res
	c.JSON(http.StatusOK, &respData)
	return

}

//获取一条设备信息数据
func deviceInfoOne(c *gin.Context, deviceid string) {
	var respData RespData
	// 接收数据
	var deviceInfo model.DeviceInfo
	// 验证数字类型
	deviceId, err := strconv.ParseInt(deviceid, 10, 64)
	if err != nil {
		respData.Error = "设备ID必须是数字"
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 通过ID查询信息
	deviceInfo.Deviceid = deviceId
	data, err := deviceInfo.GetOneDataByDeviceId()
	if err != nil {
		respData.Error = "设备信息不存在"
		c.JSON(http.StatusOK, &respData)
		return
	}

	restMap, err := RestMap()
	if err != nil {
		respData.Error = err.Error()
		c.JSON(http.StatusOK, &respData)
		return
	}

	tManufactureDate, _ := strconv.ParseInt(data.ManufactureDate, 10, 64)
	if tManufactureDate > 0 {
		tm := time.Unix(tManufactureDate, 0)
		data.ManufactureDate = tm.Format("2006-01-02 03:04:05")
	}

	tInstallationDate, _ := strconv.ParseInt(data.InstallationDate, 10, 64)
	if tInstallationDate > 0 {
		tm := time.Unix(tInstallationDate, 0)
		data.InstallationDate = tm.Format("2006-01-02 03:04:05")
	}

	tAcceptanceDate, _ := strconv.ParseInt(data.AcceptanceDate, 10, 64)
	if tAcceptanceDate > 0 {
		tm := time.Unix(tAcceptanceDate, 0)
		data.AcceptanceDate = tm.Format("2006-01-02 03:04:05")
	}

	tInspectionDate, _ := strconv.ParseInt(data.InspectionDate, 10, 64)
	if tInspectionDate > 0 {
		tm := time.Unix(tInspectionDate, 0)
		data.InspectionDate = tm.Format("2006-01-02 03:04:05")
	}

	tMaintenanceDate, _ := strconv.ParseInt(data.MaintenanceDate, 10, 64)
	if tMaintenanceDate > 0 {
		tm := time.Unix(tMaintenanceDate, 0)
		data.MaintenanceDate = tm.Format("2006-01-02 03:04:05")
	}
	data.Rname = restMap[data.Rid]
	respData.Data = data
	c.JSON(http.StatusOK, &respData)
	return
}

//获取typeid为1的数据
func GetParentIdData(c *gin.Context) {
	var respData RespData
	var typeid int64 = 1
	var di model.DeviceInfo
	di.Typeid = typeid
	data, err := di.GetDataByTypeid()
	if err != nil {
		respData.Data = "获取失败"
		c.JSON(http.StatusOK, &respData)
		return
	}
	respData.Data = data
	c.JSON(http.StatusOK, &respData)
	return
}

//获取所有父级数据
func GetParentData(c *gin.Context) {
	var respData RespData
	var channel int64 = 0
	var di model.DeviceInfo
	di.Channel = channel
	data, err := di.GetDataByChannel()
	if err != nil {
		respData.Data = "获取失败"
		c.JSON(http.StatusOK, &respData)
		return
	}
	respData.Data = data
	c.JSON(http.StatusOK, &respData)
	return
}

// 获取设备在线/离线状态
type DeviceStateData struct {
	Rid    int64   `json:"rid"`
	Rname  string  `json:"rname"`
	Device []Devic `json:"device"`
}
type Devic struct {
	Deviceid   int64  `json:"deviceid"`
	Devicename string `json:"devicename"`
	State      int64  `json:"state"`
}

// 获取设备在线/离线状态
func DeviceState(c *gin.Context) {
	var respData RespData
	// 接收数据
	ridStr := c.Query("rid")
	deviceidStr := c.Query("deviceid")
	var rid int64 = 0
	var deviceid int64 = 0
	var err error = nil
	// 验证
	// 餐厅ID
	var rest model.Restaurnt
	if ridStr != "" {
		rid, err = strconv.ParseInt(ridStr, 10, 64)
		if err != nil || rid <= 0 {
			respData.Error = "餐厅ID必须是数字，且大于0"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}

		rest.Id = rid
		//areaData, err := area.GetDataById()
		_, err := rest.One()
		if err != nil {
			respData.Error = "餐厅ID不存在"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
	}
	// 设备ID
	if deviceidStr != "" {
		deviceid, err = strconv.ParseInt(deviceidStr, 10, 64)
		if err != nil || deviceid <= 0 {
			respData.Error = "设备ID必须是数字，且大于0"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
		var device model.DeviceInfo
		device.Deviceid = deviceid
		//deviceData, err := device.GetDataByDevId()
		_, err := device.GetDataByDevId()
		if err != nil {
			respData.Error = "设备ID不存在"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
	}
	// 获取数据
	// 获取所有小区
	rests, err := rest.All()
	if err != nil {
		respData.Error = "获取所有餐厅数据失败"
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 获取所有设备状态
	var device model.DeviceInfo
	devstateData, err := device.GetDeviceState()
	if err != nil {
		respData.Error = "获取设备状态失败"
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 1. 获取所有小区的所有设备状态
	if rid == 0 && deviceid == 0 {
		var data []DeviceStateData
		for _, v := range rests {
			var dsTmp DeviceStateData
			dsTmp.Rid = v.Id
			dsTmp.Rname = v.Name
			var dev []Devic
			for _, n := range devstateData {
				if v.Id == n.Rid {
					var devTmp Devic
					devTmp.Deviceid = n.Deviceid
					devTmp.Devicename = n.Devicename
					devTmp.State = n.State
					dev = append(dev, devTmp)
				}
			}
			dsTmp.Device = dev
			data = append(data, dsTmp)
		}
		respData.Data = data
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 2. 获取指定小区的所有设备状态
	if rid > 0 && deviceid == 0 {
		var exist bool = false
		var data DeviceStateData
		for _, v := range rests {
			if v.Id == rid {
				exist = true
				data.Rid = v.Id
				data.Rname = v.Name
			}
		}
		if exist == false {
			respData.Error = "您没有当前小区设备状态的查看权限"
			c.JSON(http.StatusOK, &respData)
			return
		}
		var dev []Devic
		for _, v := range devstateData {
			if v.Rid == rid {
				var devTmp Devic
				devTmp.Deviceid = v.Deviceid
				devTmp.Devicename = v.Devicename
				devTmp.State = v.State
				dev = append(dev, devTmp)
			}
		}
		data.Device = dev
		respData.Data = data
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 3. 获取指定小区指定设备状态
	if rid > 0 && deviceid > 0 {
		exist := false
		var dev []Devic
		for _, v := range devstateData {
			if deviceid == v.Deviceid && rid == v.Rid {
				exist = true
				var devTmp Devic
				devTmp.Deviceid = v.Deviceid
				devTmp.Devicename = v.Devicename
				devTmp.State = v.State
				dev = append(dev, devTmp)
			}
		}
		if exist == false {
			respData.Error = "当前区域ID不包含此设备ID"
			c.JSON(http.StatusOK, &respData)
			return
		}
		areaExist := false
		var data DeviceStateData
		for _, v := range rests {
			if v.Id == rid {
				areaExist = true
				data.Rid = v.Id
				data.Rname = v.Name
			}
		}
		if areaExist == false {
			respData.Error = "您没有当前小区设备状态的查看权限"
			c.JSON(http.StatusOK, &respData)
			return
		}
		data.Device = dev
		respData.Data = data
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 4. 获取指定设备状态
	if rid == 0 && deviceid > 0 {
		var ridTmp int64 = 0
		var dev []Devic
		for _, v := range devstateData {
			if deviceid == v.Deviceid {
				ridTmp = v.Rid
				var devTmp Devic
				devTmp.Deviceid = v.Deviceid
				devTmp.Devicename = v.Devicename
				devTmp.State = v.State
				dev = append(dev, devTmp)
			}
		}
		areaExist := false
		var data DeviceStateData
		for _, v := range rests {
			if v.Id == ridTmp {
				areaExist = true
				data.Rid = v.Id
				data.Rname = v.Name
			}
		}
		if areaExist == false {
			respData.Error = "您没有设备状态的查看权限"
			c.JSON(http.StatusOK, &respData)
			return
		}
		data.Device = dev
		respData.Data = data
		c.JSON(http.StatusOK, &respData)
		return
	}
}

// 接收数据采集端提供的设备在线状态
type JsonDeviceState struct {
	Deviceid int64 `json:"deviceid"`
	State    int64 `json:"state"`
}

// 就逻辑与上面DeviceState关联
func UpdateDeviceState(c *gin.Context) {
	var respData RespData
	// 接收数据
	var jsondata []JsonDeviceState
	err := c.BindJSON(&jsondata)
	if err != nil {
		log.Println(err)
		respData.Error = "获取数据失败"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}
	// 验证数据
	var online []JsonDeviceState
	var offline []JsonDeviceState
	var state int64
	var deviceids []int64
	var ondeviceids []int64
	var offdeviceids []int64
	for i, v := range jsondata {
		i++
		if v.Deviceid <= 0 {
			respData.Error = "第" + fmt.Sprint(i) + "设备ID必须大于0"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
		if v.State != 0 && v.State != 1 {
			respData.Error = "第" + fmt.Sprint(i) + "状态必须是0和1"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
		var dev model.DeviceInfo
		dev.Deviceid = v.Deviceid
		_, err := dev.GetOneDataByDeviceId()
		if err != nil {
			respData.Error = "第" + fmt.Sprint(i) + "设备ID不存在"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
		var tmp JsonDeviceState
		if v.State == 0 {
			tmp.Deviceid = v.Deviceid
			tmp.State = v.State
			offline = append(offline, tmp)
			offdeviceids = append(offdeviceids, v.Deviceid)
		} else if v.State == 1 {
			tmp.Deviceid = v.Deviceid
			tmp.State = v.State
			online = append(online, tmp)
			ondeviceids = append(ondeviceids, v.Deviceid)
		}
		// 获取状态
		state = v.State
		// 获取所有id
		deviceids = append(deviceids, v.Deviceid)
	}
	var do model.DeviceOnline
	if len(online) == 0 || len(offline) == 0 {
		// sql
		deviceSql := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(deviceids)), ","), "[]")
		err := do.BulkEditing(state, deviceSql)
		if err != nil {
			log.Println(err)
			respData.Error = "修改失败"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
	} else {
		onlineSql := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ondeviceids)), ","), "[]")
		offlineSql := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(offdeviceids)), ","), "[]")
		err := do.BulkEditingTx(onlineSql, offlineSql)
		if err != nil {
			log.Println(err)
			respData.Error = "修改失败"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
	}
	respData.Data = "修改成功"
	c.JSON(http.StatusOK, &respData)
	return
}

// 新逻辑
/*func UpdateDeviceState(c *gin.Context) {
	var respData RespData
	// 接收数据
	var jsondata []JsonDeviceState
	err := c.BindJSON(&jsondata)
	if err != nil {
		log.Println(err)
		respData.Error = "获取数据失败"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	// 获取所有状态最新数据
	var do model.DeviceOnline
	onlineData, err := do.AllNew()
	if err != nil {
		respData.Error = "获取所有状态最新数据失败"
		c.JSON(http.StatusOK, &respData)
		return
	}
	fmt.Println(onlineData)

	// 验证数据
	var addOnline []model.DeviceOnline
	//var offline []JsonDeviceState
	//var state int64
	//var deviceids []int64
	//var ondeviceids []int64
	//var offdeviceids []int64
	for i, v := range jsondata {
		i++
		if v.Deviceid <= 0 {
			respData.Error = "第" + fmt.Sprint(i) + "设备ID必须大于0"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
		if v.State != 0 && v.State != 1 {
			respData.Error = "第" + fmt.Sprint(i) + "状态必须是0和1"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
		var dev model.DeviceInfo
		dev.Deviceid = v.Deviceid
		_, err := dev.GetOneDataByDeviceId()
		if err != nil {
			respData.Error = "第" + fmt.Sprint(i) + "设备ID不存在"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}

		// 组织需要新增的状态数据集合
		for _, m := range onlineData {
			if m.DeviceId == v.Deviceid && m.State != v.State {
				var temp model.DeviceOnline
				temp.DeviceId = v.Deviceid
				temp.State = v.State
				temp.Time = time.Now().Unix()
				addOnline = append(addOnline, temp)
			}
		}

		//var tmp JsonDeviceState
		//if v.State == 0 {
		//	tmp.Deviceid = v.Deviceid
		//	tmp.State = v.State
		//	offline = append(offline, tmp)
		//	offdeviceids = append(offdeviceids, v.Deviceid)
		//} else if v.State == 1 {
		//	tmp.Deviceid = v.Deviceid
		//	tmp.State = v.State
		//	online = append(online, tmp)
		//	ondeviceids = append(ondeviceids, v.Deviceid)
		//}
		//// 获取状态
		//state = v.State
		//// 获取所有id
		//deviceids = append(deviceids, v.Deviceid)
	}

	if len(addOnline) == 0 {
		respData.Error = "状态无需更新"
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 添加状态
	err = do.Adds(addOnline)
	if err != nil {
		log.Println(err)
		respData.Error = "添加失败"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	respData.Data = "添加成功"
	c.JSON(http.StatusOK, &respData)
	return

	//var do model.DeviceOnline
	//if len(online) == 0 || len(offline) == 0 {
	//	// sql
	//	deviceSql := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(deviceids)), ","), "[]")
	//	err := do.BulkEditing(state, deviceSql)
	//	if err != nil {
	//		log.Println(err)
	//		respData.Error = "修改失败"
	//		c.JSON(http.StatusBadRequest, &respData)
	//		return
	//	}
	//} else {
	//	onlineSql := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ondeviceids)), ","), "[]")
	//	offlineSql := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(offdeviceids)), ","), "[]")
	//	err := do.BulkEditingTx(onlineSql, offlineSql)
	//	if err != nil {
	//		log.Println(err)
	//		respData.Error = "修改失败"
	//		c.JSON(http.StatusBadRequest, &respData)
	//		return
	//	}
	//}
	//respData.Data = "修改成功"
	//c.JSON(http.StatusOK, &respData)
	//return
}*/
