package controller

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"src/model"
)

func Data(c *gin.Context) {
	//view := c.Query("view")
	//switch view {
	//case "summary":
	//	//pageOne(c)
	//case "area":
	//	//changedPageTwo(c)
	//case "status":
	//	//pageThree(c)
	//case "online":
	//	changedPageOne(c)
	//}
	changedPageOne(c)

}

type AbnormalData struct {
	UnusualInfos []UnusualInfo
	Totals       Total
}
type UnusualInfo struct {
	AreaName       string           //小区名称 表1和表3中都有
	Alarm          int64            //报警数
	Exception      int64            //异常数  表1和表3中都有
	Supervise      int64            //督办数
	Per            string           //表1中的百分比
	DailyAlarm     map[string]int64 //表2 ["0"]=>2,["6"]=>1,["12"]=>0,["18"]=>5,
	DailyException map[string]int64 //表4 ["0"]=>2,["6"]=>1,["12"]=>0,["18"]=>5,
}
type Total struct {
	AlarmTotal          int64
	MonthAlarmTotal     int64
	ExceptionTotal      int64
	MonthExceptionTotal int64
	SuperviseTotal      int64
	MonthSuperviseTotal int64
}

/*func pageOne(c *gin.Context) {
	var respData RespData
	var data AbnormalData
	// 获取所有小区
	var level int64 = 5
	var area model.SysArea
	area.Level = level
	areas, err := area.GetAreaByLevel()
	if err != nil || len(areas) == 0 {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	var areaIds []int64
	for _, v := range areas {
		areaIds = append(areaIds, v.Areaid)
	}

	// 获取所有异常IDS
	var st model.StatusType
	exIds, err := st.GetExceptionIds()
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 获取所有小区异常总数
	var ds model.DeviceStatus
	exceptionAmount, err := ds.GetSumByAreaIdsExceptionIds(areaIds, exIds)
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 获取所有报警IDS
	AlarmIds, err := st.GetAlarmIds()
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}

	operationIds, err := st.GetOperationIds()
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 2.获取当日报警个数
	// 获取当日时间戳，时间格式从 00:00:00 开始
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 00:00:00", time.Local)
	// 00:00:00
	startTime := t.Unix()
	// 06:00:00
	timeAdd6 := t.Unix() + (60 * 60 * 6)
	// 12:00:00
	timeAdd12 := t.Unix() + (60 * 60 * 12)
	// 18:00:00
	timeAdd18 := t.Unix() + (60 * 60 * 18)
	// 24:00:00
	endTime := t.Unix() + (60 * 60 * 24)

	//var unusualinfo []UnusualInfo
	// 1.各个小区的异常总数
	data2, err := ds.GetSumByAreaGroup(exIds)
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 2.获取报警的所有deviceStatus数据
	AlarmDeviceStatusData, err := ds.GetDataByStatusType(AlarmIds)
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 3.异常设备详情数据（报警、异常、督办数对应一个小区为一组）
	// 各个小区的报警总数
	data1, err := ds.GetSumByAreaGroup(AlarmIds)
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 各小区督办总数
	var sp model.Supervise
	supervises, err := sp.GetSuperviseSunByAreaGroup()
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 4.当日异常数量（x轴为时间点）
	ExcepDeviceStatusData, err := ds.GetDataByStatusType(exIds)
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	// --- Total
	// 各个小区的操作总数
	data3, err := ds.GetSumByAreaGroup(operationIds)
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	var alarmtotal int64 = 0
	var exptotal int64 = 0
	var supervtotal int64 = 0
	var operationtotal int64 = 0
	for _, v := range areas {
		var tmp UnusualInfo
		//1.各小区异常数量
		var ExceptionAmount int64 = 0
		for _, q := range data2 {
			if strings.Contains(fmt.Sprint(q.AreaId), fmt.Sprint(v.Areaid)) {
				ExceptionAmount += q.Sum
			}
		}
		//1.各小区异常百分比
		exceptionRatio := float64(ExceptionAmount) / float64(exceptionAmount) * 100
		ratio, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", exceptionRatio), 64)
		exceptionRatioStr := fmt.Sprint(ratio) + "%"
		//exceptionRatioStr := fmt.Sprint(math.Trunc(exceptionRatio*100*1e2+0.5)*1e-2) + "%"

		tmp.AreaName = v.Areaname
		//1.各小区异常数量和百分比
		tmp.Exception = ExceptionAmount
		tmp.Per = exceptionRatioStr
		exptotal += ExceptionAmount

		//2.当日报警数量（x轴为时间点）
		dailyAlarm := make(map[string]int64)
		var zeroToSixAlarmNum int64 = 0
		var sixToTwelveAlarmNum int64 = 0
		var twelveToEighteenAlarmNum int64 = 0
		var eighteenTozeroAlarmNum int64 = 0
		for _, n := range AlarmDeviceStatusData {
			if strings.Contains(fmt.Sprint(n.AreaId), fmt.Sprint(v.Areaid)) {
				time, _ := strconv.ParseInt(n.Time, 10, 64)
				if time > startTime && time <= timeAdd6 {
					zeroToSixAlarmNum++
				}
				if time > timeAdd6 && time <= timeAdd12 {
					sixToTwelveAlarmNum++
				}
				if time > timeAdd12 && time <= timeAdd18 {
					twelveToEighteenAlarmNum++
				}
				if time > timeAdd18 && time <= endTime {
					eighteenTozeroAlarmNum++
				}
			}
		}
		dailyAlarm["6"] = zeroToSixAlarmNum
		dailyAlarm["12"] = sixToTwelveAlarmNum
		dailyAlarm["18"] = twelveToEighteenAlarmNum
		dailyAlarm["0"] = eighteenTozeroAlarmNum
		tmp.DailyAlarm = dailyAlarm

		//3.异常设备详情数据（报警、异常、督办数对应一个小区为一组）
		//各小区报警数量
		var AlarmAmount int64 = 0
		for _, q := range data1 {
			if strings.Contains(fmt.Sprint(q.AreaId), fmt.Sprint(v.Areaid)) {
				AlarmAmount += q.Sum
			}
		}
		tmp.Alarm = AlarmAmount
		alarmtotal += AlarmAmount
		//各小区督办数量
		var SuperviseAmount int64 = 0
		for _, q := range supervises {
			if strings.Contains(fmt.Sprint(q.AreaId), fmt.Sprint(v.Areaid)) {
				SuperviseAmount += q.Sum
			}
		}
		tmp.Supervise = SuperviseAmount
		supervtotal += SuperviseAmount

		//4.当日异常数量（x轴为时间点）
		dailyException := make(map[string]int64)
		var zeroToSixExcepNum int64 = 0
		var sixToTwelveExcepNum int64 = 0
		var twelveToEighteenExcepNum int64 = 0
		var eighteenTozeroExcepNum int64 = 0
		for _, n := range ExcepDeviceStatusData {
			if strings.Contains(fmt.Sprint(n.AreaId), fmt.Sprint(v.Areaid)) {
				time, _ := strconv.ParseInt(n.Time, 10, 64)
				if time > startTime && time <= timeAdd6 {
					zeroToSixExcepNum++
				}
				if time > timeAdd6 && time <= timeAdd12 {
					sixToTwelveExcepNum++
				}
				if time > timeAdd12 && time <= timeAdd18 {
					twelveToEighteenExcepNum++
				}
				if time > timeAdd18 && time <= endTime {
					eighteenTozeroExcepNum++
				}
			}
		}
		dailyException["6"] = zeroToSixExcepNum
		dailyException["12"] = sixToTwelveExcepNum
		dailyException["18"] = twelveToEighteenExcepNum
		dailyException["0"] = eighteenTozeroExcepNum
		tmp.DailyException = dailyException

		data.UnusualInfos = append(data.UnusualInfos, tmp)

		// ---
		var operationAmount int64 = 0
		for _, q := range data3 {
			if strings.Contains(fmt.Sprint(q.AreaId), fmt.Sprint(v.Areaid)) {
				operationAmount += q.Sum
			}
		}
		operationtotal += operationAmount
	}

	var total Total
	// 获取总报警数量
	total.AlarmTotal = alarmtotal + exptotal + operationtotal
	// 获取总异常数量
	total.ExceptionTotal = exptotal
	// 获取总督办数量
	total.SuperviseTotal = supervtotal

	// 获取当月第一天时间和最后一天时间
	firstDayMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1).Format("2006-01-02") + " 00:00:00"
	firstDayMonthTtam, _ := time.ParseInLocation("2006-01-02 15:04:05", firstDayMonth, time.Local)
	monthStartTime := firstDayMonthTtam.Unix()
	lastDayMonth := fmt.Sprint(time.Now().AddDate(0, 0, -time.Now().Day()+1).AddDate(0, 1, 0).Format("2006-01-02") + " 00:00:00")
	lastDayMonthTtam, _ := time.ParseInLocation("2006-01-02 15:04:05", lastDayMonth, time.Local)
	monthEndTime := lastDayMonthTtam.Unix()

	// 获取所有报警数据
	statuses, err := ds.GetDeviceStatusException(areaIds, AlarmIds)
	var monthalarmtotal int64 = 0
	for _, v := range statuses {
		time, _ := strconv.ParseInt(v.Time, 10, 64)
		if time > monthStartTime && time <= monthEndTime {
			monthalarmtotal++
		}
	}

	// 获取所有异常数据
	deviceStatuses, err := ds.GetDeviceStatusException(areaIds, exIds)
	var monthexceptiontotal int64 = 0
	for _, v := range deviceStatuses {
		time, _ := strconv.ParseInt(v.Time, 10, 64)
		if time > monthStartTime && time <= monthEndTime {
			monthexceptiontotal++
		}
	}
	total.MonthExceptionTotal = monthexceptiontotal

	// 获取所有操作数据
	data4, err := ds.GetDeviceStatusException(areaIds, operationIds)
	var monthoperationtotal int64 = 0
	for _, v := range data4 {
		time, _ := strconv.ParseInt(v.Time, 10, 64)
		if time > monthStartTime && time <= monthEndTime {
			monthoperationtotal++
		}
	}

	total.MonthAlarmTotal = monthalarmtotal + monthexceptiontotal + monthoperationtotal

	// 获取所有督办数据
	total.MonthSuperviseTotal, err = sp.GetCountByAreaIds(areaIds, monthStartTime, monthEndTime)
	data.Totals = total

	respData.Data = data
	c.JSON(http.StatusOK, &respData)
	return
}*/

type AbnormalStatistics struct {
	DataAreaUnusuals []DataAreaUnusual
	ExcepTotal       Exceptotal
}
type DataAreaUnusual struct {
	AreaName  string
	DayData   Daily
	MonthData Month
}
type Daily struct {
	Alarm     int64
	Exception int64
	Operation int64
}
type Month struct {
	Alarm     int64
	Exception int64
	Operation int64
}
type Exceptotal struct {
	AlarmTotal          int64
	MonthAlarmTotal     int64
	ExceptionTotal      int64
	MonthExceptionTotal int64
	OperationTotal      int64
	MonthOperationTotal int64
}

/*func pageTwo(c *gin.Context) {
	var respData RespData
	var data AbnormalStatistics

	// 获取所有小区
	var level int64 = 5
	var area model.SysArea
	area.Level = level
	areas, err := area.GetAreaByLevel()
	if err != nil || len(areas) == 0 {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 获取所有报警IDS
	var st model.StatusType
	AlarmIds, err := st.GetAlarmIds()
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 获取所有异常IDS
	exIds, err := st.GetExceptionIds()
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 获取所有操作IDS
	operationIds, err := st.GetOperationIds()
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 获取当日时间戳，时间格式从 00:00:00 开始到 24:00:00 结束
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 00:00:00", time.Local)
	dayStartTime := t.Unix()
	dayEndTime := t.Unix() + (60 * 60 * 24)

	// 获取当月第一天时间和最后一天时间
	datetime := c.Query("datetime")
	var monthStartTime int64
	var monthEndTime int64
	if datetime == "" {
		firstDayMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1).Format("2006-01-02") + " 00:00:00"
		firstDayMonthTtam, _ := time.ParseInLocation("2006-01-02 15:04:05", firstDayMonth, time.Local)
		monthStartTime = firstDayMonthTtam.Unix()
		lastDayMonth := fmt.Sprint(time.Now().AddDate(0, 0, -time.Now().Day()+1).AddDate(0, 1, 0).Format("2006-01-02") + " 00:00:00")
		lastDayMonthTtam, _ := time.ParseInLocation("2006-01-02 15:04:05", lastDayMonth, time.Local)
		monthEndTime = lastDayMonthTtam.Unix()
	} else {
		matched, _ := regexp.MatchString("^\\d{4}-(((0([1-9]))|([1-9]))|(1(0|1|2)))$", datetime)
		if matched == false {
			respData.Error = "时间格式不正确"
			c.JSON(http.StatusOK, &respData)
			return
		}
		date := datetime + "-01 00:00:00"
		firstDayMonthTtam, _ := time.ParseInLocation("2006-1-2 15:4:5", date, time.Local)
		monthStartTime = firstDayMonthTtam.Unix()
		lastDayMonthTtam := fmt.Sprint(firstDayMonthTtam.AddDate(0, 0, -firstDayMonthTtam.Day()+1).AddDate(0, 1, 0).Format("2006-01-02") + " 00:00:00")
		lastDayMonthTta, _ := time.ParseInLocation("2006-1-2 15:4:5", lastDayMonthTtam, time.Local)
		monthEndTime = lastDayMonthTta.Unix()
	}

	// 各个小区的报警总数
	var dss model.DeviceStatus
	data1, err := dss.GetSumByAreaGroup(AlarmIds)
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 各个小区的异常总数
	data2, err := dss.GetSumByAreaGroup(exIds)
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 各个小区的操作总数
	data3, err := dss.GetSumByAreaGroup(operationIds)
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	//var data []DataAreaUnusual
	var alarmtotal int64 = 0
	var monthalarmtotal int64 = 0
	var exceptiontotal int64 = 0
	var monthexceptiontotal int64 = 0
	var operationtotal int64 = 0
	var monthoperationtotal int64 = 0
	for _, v := range areas {
		var tmp DataAreaUnusual
		// 区域名称
		tmp.AreaName = v.Areaname
		// 获取设备状态数据
		var dayData Daily
		var monthData Month
		var ds model.DeviceStatus
		ds.AreaId = v.Areaid
		// 获取当日报警
		dayAlarmNum, err := ds.GetAmountByParamAndTime(AlarmIds, dayStartTime, dayEndTime)
		if err != nil {
			respData.Error = "暂无数据"
			c.JSON(http.StatusOK, &respData)
			return
		}
		dayData.Alarm = dayAlarmNum
		// 获取当月报警
		monthAlarmNum, err := ds.GetAmountByParamAndTime(AlarmIds, monthStartTime, monthEndTime)
		if err != nil {
			respData.Error = "暂无数据"
			c.JSON(http.StatusOK, &respData)
			return
		}
		monthData.Alarm = monthAlarmNum
		monthalarmtotal += monthAlarmNum

		// 获取当日异常
		dayExceptionNum, err := ds.GetAmountByParamAndTime(exIds, dayStartTime, dayEndTime)
		if err != nil {
			respData.Error = "暂无数据"
			c.JSON(http.StatusOK, &respData)
			return
		}
		dayData.Exception = dayExceptionNum
		// 获取当月异常
		monthExceptionNum, err := ds.GetAmountByParamAndTime(exIds, monthStartTime, monthEndTime)
		if err != nil {
			respData.Error = "暂无数据"
			c.JSON(http.StatusOK, &respData)
			return
		}
		monthData.Exception = monthExceptionNum
		monthexceptiontotal += monthExceptionNum
		// 获取当日操作
		dayOperationNum, err := ds.GetAmountByParamAndTime(operationIds, dayStartTime, dayEndTime)
		if err != nil {
			respData.Error = "暂无数据"
			c.JSON(http.StatusOK, &respData)
			return
		}
		dayData.Operation = dayOperationNum
		// 获取当月操作
		monthOperationNum, err := ds.GetAmountByParamAndTime(operationIds, monthStartTime, monthEndTime)
		if err != nil {
			respData.Error = "暂无数据"
			c.JSON(http.StatusOK, &respData)
			return
		}
		monthData.Operation = monthOperationNum
		monthoperationtotal += monthOperationNum
		tmp.DayData = dayData
		tmp.MonthData = monthData
		data.DataAreaUnusuals = append(data.DataAreaUnusuals, tmp)

		//各小区报警数量
		var AlarmAmount int64 = 0
		for _, q := range data1 {
			if strings.Contains(fmt.Sprint(q.AreaId), fmt.Sprint(v.Areaid)) {
				AlarmAmount += q.Sum
			}
		}
		alarmtotal += AlarmAmount

		//各小区异常数量
		var ExceptionAmount int64 = 0
		for _, q := range data2 {
			if strings.Contains(fmt.Sprint(q.AreaId), fmt.Sprint(v.Areaid)) {
				ExceptionAmount += q.Sum
			}
		}
		exceptiontotal += ExceptionAmount

		// 各小区操作数量
		var operationAmount int64 = 0
		for _, q := range data3 {
			if strings.Contains(fmt.Sprint(q.AreaId), fmt.Sprint(v.Areaid)) {
				operationAmount += q.Sum
			}
		}
		operationtotal += operationAmount
	}

	var exceptotal Exceptotal
	exceptotal.AlarmTotal = alarmtotal
	exceptotal.MonthAlarmTotal = monthalarmtotal
	exceptotal.ExceptionTotal = exceptiontotal
	exceptotal.MonthExceptionTotal = monthexceptiontotal
	exceptotal.OperationTotal = operationtotal
	exceptotal.MonthOperationTotal = monthoperationtotal
	data.ExcepTotal = exceptotal

	respData.Data = data
	c.JSON(http.StatusOK, &respData)
	return
}*/

type DataStatusInfo struct {
	SuperViseNum Supervisenum
	DeviceAlarm  []Devicealarm
	Exceptions   []Exception
	DeviceStatus []Devicestatus
	Operations   []Operation
	ExceptionNum []Exceptionnum
}

// 图表一
type Supervisenum struct {
	AreaAndSuper []Areaandsuper
	DailyNum     map[int64]int64 //每天督办次数
}
type Areaandsuper struct {
	AreaName  string //小区名字
	SuperVNum int64  //督办次数
}

// 图表二
/*type Devicealarm struct {
	DeviceName     string //设备名称
	MotdetStart    int64  //移动侦测次数
	HideAlarmStart int64  //遮挡次数
}*/
type Devicealarm struct {
	AreaName      string //小区名字
	IllegalAccess int64  //非法访问 116
	RomoteLogin   int64  //远程登录 217
}

// 图表三
type Exception struct {
	AreaName      string //小区名字
	NetBroken     int64  //网络断开次数 121
	IllegalAccess int64  //非法访问 116
	RecError      int64  //录像出错 122
	ViLost        int64  //视频信号丢失115
	AiLost        int64  //音频信号丢失 147
	Else          int64  //其他
}

// 图表四
type Devicestatus struct {
	AreaName      string //小区名称
	DeviceName    string //设备名称
	ExceptionName string //异常名称
	ExceptionTime string //异常时间
}

// 图表五
type Operation struct {
	AreaName        string //小区名字
	StartDvr        int64  //开机 184
	StopDvr         int64  //关机 185
	StopAbnormal    int64  //异常关机 186
	RomoteLogin     int64  //远程登录 217
	RemoteReboot    int64  //远程重启 228
	RemoteFormatHDD int64  //远程格式化硬盘 235
}

// 图表六
type Exceptionnum struct {
	AreaName        string //小区名字
	ExceptionAmount int64
}

/*func pageThree(c *gin.Context) {
	var respData RespData
		//// 获取所有小区
		//var level int64 = 5
		//var area model.SysArea
		//area.Level = level
		//areas, err := area.GetAreaByLevel()
		//if err != nil || len(areas) == 0 {
		//	respData.Error = "暂无数据"
		//	c.JSON(http.StatusOK, &respData)
		//	return
		//}
	// 获取当前用户所属区域
	var ua model.SysUserArea
	ua.Userid = global.UserT.Id
	userAreas, err := ua.GetUserAreaByUserId()
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 获取当前用户所属区域ids
	var userareaids []int64
	for _, v := range userAreas {
		userareaids = append(userareaids, v.Areaid)
	}
	// 获取用户所属小区
	var level int64 = 5
	var area model.SysArea
	area.Level = level
	areas, err := area.GetAreaByLevelAndUserAreaIds(userareaids)
	if err != nil || len(areas) == 0 {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	var areaIds []int64
	for _, v := range areas {
		areaIds = append(areaIds, v.Areaid)
	}

	//1.每个小区的督办次数/当前账户权限内的所有小区加起来 每天的督办次数
	var data DataStatusInfo
	// 获取所有异常IDS
	var st model.StatusType
	exIds, err := st.GetExceptionIds()
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}

	datetime := c.Query("datetime")
	var monthStartTime int64
	var monthEndTime int64
	var dayNum int
	if datetime == "" {
		// 获取当月的的时间戳，从月初 00:00:00 开始到月末 24:00:00 结束
		firstDayMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1).Format("2006-01-02") + " 00:00:00"
		firstDayMonthTtam, _ := time.ParseInLocation("2006-01-02 15:04:05", firstDayMonth, time.Local)
		monthStartTime = firstDayMonthTtam.Unix()
		lastDayMonth := fmt.Sprint(time.Now().AddDate(0, 0, -time.Now().Day()+1).AddDate(0, 1, 0).Format("2006-01-02") + " 00:00:00")
		lastDayMonthTtam, _ := time.ParseInLocation("2006-01-02 15:04:05", lastDayMonth, time.Local)
		monthEndTime = lastDayMonthTtam.Unix()

		// 获取当前月有多少天
		date := fmt.Sprint(time.Now().AddDate(0, 0, -time.Now().Day()+1).AddDate(0, 1, -1).Format("2006-01-02"))
		dayNumStr := date[strings.LastIndex(date, "-")+1:]
		dayNum, _ = strconv.Atoi(dayNumStr)
	} else {
		matched, _ := regexp.MatchString("^\\d{4}-(((0([1-9]))|([1-9]))|(1(0|1|2)))$", datetime)
		if matched == false {
			respData.Error = "时间格式不正确"
			c.JSON(http.StatusOK, &respData)
			return
		}
		date := datetime + "-01 00:00:00"
		firstDayMonthTtam, _ := time.ParseInLocation("2006-1-2 15:4:5", date, time.Local)
		monthStartTime = firstDayMonthTtam.Unix()
		lastDayMonthTtam := fmt.Sprint(firstDayMonthTtam.AddDate(0, 0, -firstDayMonthTtam.Day()+1).AddDate(0, 1, 0).Format("2006-01-02") + " 00:00:00")
		lastDayMonthTta, _ := time.ParseInLocation("2006-1-2 15:4:5", lastDayMonthTtam, time.Local)
		monthEndTime = lastDayMonthTta.Unix()

		dated := fmt.Sprint(firstDayMonthTtam.AddDate(0, 0, -firstDayMonthTtam.Day()+1).AddDate(0, 1, -1).Format("2006-01-02"))
		dayNumStr := dated[strings.LastIndex(dated, "-")+1:]
		dayNum, _ = strconv.Atoi(dayNumStr)
	}

	var supervisenum Supervisenum
	// 定义当前月份的map
	monthNum := make(map[int64]int64)
	for i := 1; i <= dayNum; i++ {
		monthNum[int64(i)] = 0
	}

	// 根据 areaIds、开始时间、结束时间，查询每天的督办次数
	var sp model.Supervise
	daySuperVnum, err := sp.GetDailyCountByAreaIds(areaIds, monthStartTime, monthEndTime)
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	for i, v := range daySuperVnum {
		for x := range monthNum {
			if i == x {
				monthNum[x] = v
			}
		}
	}
	supervisenum.DailyNum = monthNum

	//1.每个小区的督办次数
	var areaandsuper []Areaandsuper
	// 按照区域分组，获取督办总数
	supervises, err := sp.GetSuperviseSunByAreaGroup()
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	//2.各小区非法访问、远程登录
	var devicealarm []Devicealarm
	//3.异常详情次数（网络断开，非法访问，录像出错，视频信号丢失，音频信号丢失）
	var exception []Exception
	// 根据statusTypeId、areaId分组，获取device_status总数。
	var ds model.DeviceStatus
	// 根据statusTypeId、areaId分组，获取device_status总数。
	ipList := ""
	if len(global.IpBlackList) == 0 {
		ipList = ""
	} else {
		for _, v := range global.IpBlackList {
			ipList = ipList + "'" + v + "',"
		}
		ipList = ipList[:len(ipList)-1]
	}
	deviceStatuses, err := ds.GetSumByStatusTypeIdAreaId(ipList)
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	//5.操作详情次数（开/关机，异常关机，远程登录，远程重启，格式化。x轴为小区）
	var operation []Operation
	//6.各个小区的异常总数
	var exceptionnum []Exceptionnum
	data2, err := ds.GetSumByAreaGroup(exIds)
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	for _, v := range areas {
		//1.每个小区的督办次数
		var tmpAreaandsuper Areaandsuper
		var SuperVsum int64 = 0
		for _, y := range supervises {
			if strings.Contains(fmt.Sprint(y.AreaId), fmt.Sprint(v.Areaid)) {
				SuperVsum += y.Sum
			}
		}
		tmpAreaandsuper.AreaName = v.Areaname
		tmpAreaandsuper.SuperVNum = SuperVsum
		areaandsuper = append(areaandsuper, tmpAreaandsuper)

		//2.各小区非法访问、远程登录
		var tmpDevicealarm Devicealarm

		//3.异常详情次数（网络断开，非法访问，录像出错，视频信号丢失，音频信号丢失）
		var tmpException Exception
		// 异常状态ID
		var netBroken int64 = 121
		var illegalAccess int64 = 116
		var recError int64 = 122
		var viLost int64 = 115
		var aiLost int64 = 147
		// 申明变量，接收Sum
		var netBrokenSum int64 = 0
		var illegalAccessSum int64 = 0
		var recErrorSum int64 = 0
		var viLostSum int64 = 0
		var aiLostSum int64 = 0

		//5.操作详情次数（开/关机，异常关机，远程登录，远程重启，格式化。x轴为小区）
		// 操作状态ID
		var tmpOperation Operation
		var startDvr int64 = 184
		var stopDvr int64 = 185
		var stopAbnormal int64 = 186
		var romoteLogin int64 = 217
		var remoteReboot int64 = 228
		var remoteFormatHDD int64 = 235
		// 申明变量，接收Sum
		var startDvrSum int64 = 0
		var stopDvrSum int64 = 0
		var stopAbnormalSum int64 = 0
		var romoteLoginSum int64 = 0
		var remoteRebootSum int64 = 0
		var remoteFormatHDDSum int64 = 0

		for _, n := range deviceStatuses {
			//2.各小区非法访问、远程登录

			//3.异常详情次数（网络断开，非法访问，录像出错，视频信号丢失，音频信号丢失）
			if strings.Contains(fmt.Sprint(n.AreaId), fmt.Sprint(v.Areaid)) && n.StatusTypeId == netBroken {
				netBrokenSum += n.Sum
			}
			if strings.Contains(fmt.Sprint(n.AreaId), fmt.Sprint(v.Areaid)) && n.StatusTypeId == illegalAccess {
				illegalAccessSum += n.Sum
			}
			if strings.Contains(fmt.Sprint(n.AreaId), fmt.Sprint(v.Areaid)) && n.StatusTypeId == recError {
				recErrorSum += n.Sum
			}
			if strings.Contains(fmt.Sprint(n.AreaId), fmt.Sprint(v.Areaid)) && n.StatusTypeId == viLost {
				viLostSum += n.Sum
			}
			if strings.Contains(fmt.Sprint(n.AreaId), fmt.Sprint(v.Areaid)) && n.StatusTypeId == aiLost {
				aiLostSum += n.Sum
			}

			//5.操作详情次数（开/关机，异常关机，远程登录，远程重启，格式化。x轴为小区）
			if strings.Contains(fmt.Sprint(n.AreaId), fmt.Sprint(v.Areaid)) && n.StatusTypeId == startDvr {
				startDvrSum += n.Sum
			}
			if strings.Contains(fmt.Sprint(n.AreaId), fmt.Sprint(v.Areaid)) && n.StatusTypeId == stopDvr {
				stopDvrSum += n.Sum
			}
			if strings.Contains(fmt.Sprint(n.AreaId), fmt.Sprint(v.Areaid)) && n.StatusTypeId == stopAbnormal {
				stopAbnormalSum += n.Sum
			}
			if strings.Contains(fmt.Sprint(n.AreaId), fmt.Sprint(v.Areaid)) && n.StatusTypeId == romoteLogin {
				romoteLoginSum += n.Sum
			}
			if strings.Contains(fmt.Sprint(n.AreaId), fmt.Sprint(v.Areaid)) && n.StatusTypeId == remoteReboot {
				remoteRebootSum += n.Sum
			}
			if strings.Contains(fmt.Sprint(n.AreaId), fmt.Sprint(v.Areaid)) && n.StatusTypeId == remoteFormatHDD {
				remoteFormatHDDSum += n.Sum
			}
		}

		// 6.各个小区的异常总数
		var tmpExceptionnum Exceptionnum
		var ExceptionAmount int64 = 0
		for _, q := range data2 {
			if strings.Contains(fmt.Sprint(q.AreaId), fmt.Sprint(v.Areaid)) {
				ExceptionAmount += q.Sum
			}
		}
		tmpExceptionnum.AreaName = v.Areaname
		tmpExceptionnum.ExceptionAmount = ExceptionAmount
		exceptionnum = append(exceptionnum, tmpExceptionnum)

		//2.各小区非法访问、远程登录
		tmpDevicealarm.AreaName = v.Areaname
		tmpDevicealarm.IllegalAccess = illegalAccessSum
		tmpDevicealarm.RomoteLogin = romoteLoginSum
		devicealarm = append(devicealarm, tmpDevicealarm)

		//3.异常详情次数（网络断开，非法访问，录像出错，视频信号丢失，音频信号丢失）
		tmpException.AreaName = v.Areaname
		tmpException.NetBroken = netBrokenSum
		tmpException.IllegalAccess = illegalAccessSum
		tmpException.RecError = recErrorSum
		tmpException.ViLost = viLostSum
		tmpException.AiLost = aiLostSum
		tmpException.Else = ExceptionAmount - (netBrokenSum + illegalAccessSum + recErrorSum + viLostSum + aiLostSum)
		exception = append(exception, tmpException)

		//5.操作详情次数（开/关机，异常关机，远程登录，远程重启，格式化。x轴为小区）
		tmpOperation.AreaName = v.Areaname
		tmpOperation.StartDvr = startDvrSum
		tmpOperation.StopDvr = stopDvrSum
		tmpOperation.StopAbnormal = stopAbnormalSum
		tmpOperation.RomoteLogin = romoteLoginSum
		tmpOperation.RemoteReboot = remoteRebootSum
		tmpOperation.RemoteFormatHDD = remoteFormatHDDSum
		operation = append(operation, tmpOperation)

	}
	//1.每个小区的督办次数
	supervisenum.AreaAndSuper = areaandsuper
	data.SuperViseNum = supervisenum
	//2.各小区非法访问、远程登录
	data.DeviceAlarm = devicealarm
	//3.异常详情次数（网络断开，非法访问，录像出错，视频信号丢失，音频信号丢失）
	data.Exceptions = exception
	//5.操作详情次数（开/关机，异常关机，远程登录，远程重启，格式化。x轴为小区）
	data.Operations = operation
	// 6.各个小区的异常总数
	data.ExceptionNum = exceptionnum

	////2.每个摄像头的移动侦测次数/每个摄像头的遮挡次数
	//// 获取所有设备
	//var di model.DeviceInfo
	//devInfoAll, err := di.GetAll()
	//if err != nil {
	//	respData.Error = "暂无数据"
	//	c.JSON(http.StatusOK, &respData)
	//	return
	//}
	//
	//// 按照设备、状态类型分组，获取所有设备状态总数
	//statuses, err := ds.DeviceStatusSum()
	//if err != nil {
	//	respData.Error = "暂无数据"
	//	c.JSON(http.StatusOK, &respData)
	//	return
	//}
	//// 获取移动侦测次数
	//var motdet int64 = 8
	//// 获取遮挡次数次数
	//var hideAlarmStart int64 = 10
	//var devicealarm []Devicealarm
	//for _, v := range devInfoAll {
	//	var tmpDevicealarm Devicealarm
	//	tmpDevicealarm.DeviceName = v.Devicename
	//	for _, y := range statuses {
	//		if v.Parentid == y.DeviceId && v.Channel == y.Channel && motdet == y.StatusTypeId {
	//			tmpDevicealarm.MotdetStart = y.Sum
	//		}
	//		if v.Parentid == y.DeviceId && v.Channel == y.Channel && hideAlarmStart == y.StatusTypeId {
	//			tmpDevicealarm.HideAlarmStart = y.Sum
	//		}
	//	}
	//	devicealarm = append(devicealarm, tmpDevicealarm)
	//}
	//data.DeviceAlarm = devicealarm

	//4.最新的设备异常信息
	var devicestatus []Devicestatus
	// 获取区域名称
	areaS, _ := area.All("")
	areaMap := make(map[int64]string)
	for _, v := range areaS {
		areaMap[v.Areaid] = v.Areaname
	}
	// 获取设备名称
	var device model.DeviceInfo
	infos, _ := device.GetAll()
	deviceMap := make(map[string]string)
	for _, v := range infos {
		deviceMap[fmt.Sprint(v.Parentid)+"+"+fmt.Sprint(v.Channel)] = v.Devicename
	}
	// 获取异常名称
	statusTypes, _ := st.All()
	statusTypesMap := make(map[int64]string)
	for _, v := range statusTypes {
		statusTypesMap[v.Id] = v.Des
	}
	//var ds model.DeviceStatus
	newExceptionData, err := ds.NewExceptionData(areaIds, exIds)
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	for _, v := range newExceptionData {
		var tmpDevicestatus Devicestatus
		tmpDevicestatus.AreaName = areaMap[v.AreaId]
		tmpDevicestatus.DeviceName = deviceMap[fmt.Sprint(v.DeviceId)+"+"+fmt.Sprint(v.Channel)]
		tmpDevicestatus.ExceptionName = statusTypesMap[v.StatusTypeId]
		tmpDevicestatus.ExceptionTime = v.Time
		devicestatus = append(devicestatus, tmpDevicestatus)
	}
	data.DeviceStatus = devicestatus

	respData.Data = data
	c.JSON(http.StatusOK, &respData)
	return
}*/

// 修改后的第一个页面
type StatusAnalysis struct {
	StatusAnalysisTotal PageOneTotal
	StatusAnalysisData  []PageOneStatus
}
type PageOneTotal struct {
	CommunityTotal int64
	DeviceTotal    int64
	OnlineTotal    int64
	OfflineTotal   int64
}
type PageOneStatus struct {
	AreaName     string
	AreaId       int64
	OfflineNum   int64
	DeviceStatus []DeviceStatu
	WeekOffline  []Weekoffline
}
type DeviceStatu struct {
	DeviceId   int64
	DeviceName string
	State      int64
}
type Weekoffline struct {
	Date       int64
	OfflineNum int64
}

/*func changedPageOne(c *gin.Context) {
	var respData RespData
	// 获取当前用户所属区域
	var ua model.SysUserArea
	ua.Userid = global.UserT.Id
	userAreas, err := ua.GetUserAreaByUserId()
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 获取当前用户所属区域ids
	var userareaids []int64
	for _, v := range userAreas {
		userareaids = append(userareaids, v.Areaid)
	}
	// 获取用户所属小区
	var level int64 = 5
	var area model.SysArea
	area.Level = level
	areas, err := area.GetAreaByLevelAndUserAreaIds(userareaids)
	if err != nil || len(areas) == 0 {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 获取所有设备
	var di model.DeviceInfo
	devInfoAll, err := di.GetAll()
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}

	var pagetotal PageOneTotal
	// 获取小区总数
	areaTotal := int64(len(areas))
	pagetotal.CommunityTotal = areaTotal

	// 设备总数
	var deviceTotal int64
	// 获取当前所有小区的所有设备
	var devices []model.DeviceInfo
	var device model.DeviceInfo
	for _, v := range areas {
		for _, n := range devInfoAll {
			if v.Areaid == n.Areaid {
				deviceTotal++
				device = n
				devices = append(devices, device)
			}
		}
	}
	// 设备总数
	pagetotal.DeviceTotal = deviceTotal

	// 获取离线总数
	// 1. 获取所有设备的最新状态
	var ds model.DeviceStatus
	latestDeviceStatus, err := ds.LatestDeviceStatus()
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 2. 获取当前所有小区的所有设备最新状态
	//var deviceStatuses []model.DeviceStatus
	//var deviceStatus model.DeviceStatus
	//
	type DeviceData struct {
		Devicename   string
		StatusTypeId int64
		AreaId       int64
	}
	var devicedata []DeviceData
	var devicedatatmp DeviceData
	for _, v := range devices {
		for _, n := range latestDeviceStatus {
			if v.Parentid == n.DeviceId && v.Channel == n.Channel {
				//deviceStatus = n
				//deviceStatuses = append(deviceStatuses, deviceStatus)

				devicedatatmp.AreaId = v.Areaid
				devicedatatmp.Devicename = v.Devicename
				devicedatatmp.StatusTypeId = n.StatusTypeId
				devicedata = append(devicedata, devicedatatmp)
			}
		}
	}
	// 离线状态id
	var offline int64 = 123
	var offlineTotal int64 = 0
	for _, v := range devicedata {
		if v.StatusTypeId == offline {
			offlineTotal++
		}
	}
	// 离线总数
	pagetotal.OfflineTotal = offlineTotal
	// 在线总数
	pagetotal.OnlineTotal = deviceTotal - offlineTotal

	// 根据小区划分摄像头
	var pageonestatus []PageOneStatus

	// 图表③近七天离线数据
	// 1. 获取时间戳
	date := time.Now().Format("2006-01-02")
	datetime, _ := time.ParseInLocation("2006-01-02 15:04:05", date+" 00:00:00", time.Local)
	sevenDayStart := datetime.Unix()
	sevenDayEnd := sevenDayStart + 60*60*24
	sixDayStart := sevenDayStart - 60*60*24
	fiveDayStart := sixDayStart - 60*60*24
	fourDayStart := fiveDayStart - 60*60*24
	threeDayStart := fourDayStart - 60*60*24
	twoDayStart := threeDayStart - 60*60*24
	oneDayStart := twoDayStart - 60*60*24
	//fmt.Println(time.Unix(sevenDayStart, 0).Format("2006-01-02 15:04:05"))
	//fmt.Println(time.Unix(sixDayStart, 0).Format("2006-01-02 15:04:05"))
	// 2. 获取所有离线数据
	ds.StatusTypeId = offline
	offlineData, err := ds.GetDataByStatusTypeId()
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	for _, v := range areas {
		var pageonestatustmp PageOneStatus
		var tmp DeviceStatu
		var offlinenum int64 = 0
		for _, n := range devicedata {
			if v.Areaid == n.AreaId {
				if n.StatusTypeId == 123 {
					tmp.DeviceName = n.Devicename
					tmp.Status = "离线"
					// 图表①在线离线列表
					pageonestatustmp.DeviceStatus = append(pageonestatustmp.DeviceStatus, tmp)
					offlinenum++
				} else {
					tmp.DeviceName = n.Devicename
					tmp.Status = "在线"
					// 图表①在线离线列表
					pageonestatustmp.DeviceStatus = append(pageonestatustmp.DeviceStatus, tmp)
				}
			}

		}
		// 图表①在线离线列表
		pageonestatustmp.AreaName = v.Areaname
		pageonestatustmp.AreaId = v.Areaid
		// 图表②离线数量
		pageonestatustmp.OfflineNum = offlinenum
		// 图表③近七天离线数据
		var one int64 = 0
		var two int64 = 0
		var three int64 = 0
		var four int64 = 0
		var five int64 = 0
		var six int64 = 0
		var seven int64 = 0
		for _, n := range offlineData {
			if v.Areaid == n.AreaId {
				times, _ := strconv.ParseInt(n.Time, 10, 64)
				switch {
				case times >= oneDayStart && times < twoDayStart:
					one++
				case times >= twoDayStart && times < threeDayStart:
					two++
				case times >= threeDayStart && times < fourDayStart:
					three++
				case times >= fourDayStart && times < fiveDayStart:
					four++
				case times >= fiveDayStart && times < sixDayStart:
					five++
				case times >= sixDayStart && times < sevenDayStart:
					six++
				case times >= sevenDayStart && times < sevenDayEnd:
					seven++
				}
			}
		}
		//var tmpWeek Weekoffline
		var week Weekoffline
		week.Date = oneDayStart
		week.OfflineNum = one
		pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
		week.Date = twoDayStart
		week.OfflineNum = two
		pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
		week.Date = threeDayStart
		week.OfflineNum = three
		pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
		week.Date = fourDayStart
		week.OfflineNum = four
		pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
		week.Date = fiveDayStart
		week.OfflineNum = five
		pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
		week.Date = sixDayStart
		week.OfflineNum = six
		pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
		week.Date = sevenDayStart
		week.OfflineNum = seven
		pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)

		//week[oneDayStart] = one
		//week[twoDayStart] = two
		//week[threeDayStart] = three
		//week[fourDayStart] = four
		//week[fiveDayStart] = five
		//week[sixDayStart] = six
		//week[sevenDayStart] = seven
		//pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
		pageonestatus = append(pageonestatus, pageonestatustmp)
	}

	var data StatusAnalysis
	data.StatusAnalysisData = pageonestatus
	data.StatusAnalysisTotal = pagetotal
	respData.Data = data
	c.JSON(http.StatusOK, &respData)
	return
}*/
func changedPageOne(c *gin.Context) {
	var respData RespData

	//// 获取当前用户所属区域
	//var ua model.SysUserArea
	//ua.Userid = global.UserT.Id
	//userAreas, err := ua.GetUserAreaByUserId()
	//if err != nil {
	//	respData.Error = "暂无数据"
	//	c.JSON(http.StatusOK, &respData)
	//	return
	//}
	//// 获取当前用户所属区域ids
	//var userareaids []int64
	//for _, v := range userAreas {
	//	userareaids = append(userareaids, v.Areaid)
	//}
	//// 获取用户所属小区
	//var level int64 = 5
	//var area model.SysArea
	//area.Level = level
	//areas, err := area.GetAreaByLevelAndUserAreaIds(userareaids)
	//if err != nil || len(areas) == 0 {
	//	respData.Error = "暂无数据"
	//	c.JSON(http.StatusOK, &respData)
	//	return
	//}

	// 获取所有餐厅
	var rest model.Restaurnt
	restall, err := rest.All()
	if err != nil {
		respData.Error = "获取所有餐厅数据失败"
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 获取所有设备
	var di model.DeviceInfo
	devInfoAll, err := di.GetAll()
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}

	var pagetotal PageOneTotal
	// 获取小区总数
	areaTotal := int64(len(restall))
	pagetotal.CommunityTotal = areaTotal

	// 设备总数
	var deviceTotal int64
	// 获取当前所有小区的所有设备
	var devices []model.DeviceInfo
	var device model.DeviceInfo
	for _, v := range restall {
		for _, n := range devInfoAll {
			if v.Id == n.Rid {
				deviceTotal++
				device = n
				devices = append(devices, device)
			}
		}
	}
	// 设备总数
	pagetotal.DeviceTotal = deviceTotal

	// 获取所有设备的状态
	var do model.DeviceOnline
	devonlinedata, err := do.AllNew()
	if err != nil {
		log.Println(err)
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 添加设备状态
	for i, v := range devices {
		for _, n := range devonlinedata {
			if v.Deviceid == n.DeviceId {
				devices[i].State = n.State
			}
		}
	}

	// 获取在线设备数量
	var offlineTotal int64 = 0
	for _, v := range devices {
		if v.State == 0 {
			offlineTotal++
		}
	}
	// 离线总数
	pagetotal.OfflineTotal = offlineTotal
	// 在线总数
	pagetotal.OnlineTotal = deviceTotal - offlineTotal

	// 图表③近七天离线数据
	// 离线状态id
	var offlineStatus int64 = 123
	// 1. 获取时间戳
	date := time.Now().Format("2006-01-02")
	datetime, _ := time.ParseInLocation("2006-01-02 15:04:05", date+" 00:00:00", time.Local)
	sevenDayStart := datetime.Unix()
	sevenDayEnd := sevenDayStart + 60*60*24
	sixDayStart := sevenDayStart - 60*60*24
	fiveDayStart := sixDayStart - 60*60*24
	fourDayStart := fiveDayStart - 60*60*24
	threeDayStart := fourDayStart - 60*60*24
	twoDayStart := threeDayStart - 60*60*24
	oneDayStart := twoDayStart - 60*60*24
	// 2. 获取所有离线数据
	var ds model.DeviceStatus
	ds.StatusTypeId = offlineStatus
	offlineData, err := ds.GetDataByStatusTypeId()
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	var pageonestatus []PageOneStatus

	for _, v := range restall {
		var pageonestatustmp PageOneStatus
		pageonestatustmp.AreaId = v.Id
		pageonestatustmp.AreaName = v.Name
		var offline int64 = 0
		var tmp DeviceStatu
		for _, n := range devices {
			if v.Id == n.Rid {
				tmp.DeviceId = n.Deviceid
				tmp.DeviceName = n.Devicename
				tmp.State = n.State
				pageonestatustmp.DeviceStatus = append(pageonestatustmp.DeviceStatus, tmp)
				if n.State == 0 {
					offline++
				}
			}
		}
		pageonestatustmp.OfflineNum = offline

		var one int64 = 0
		var two int64 = 0
		var three int64 = 0
		var four int64 = 0
		var five int64 = 0
		var six int64 = 0
		var seven int64 = 0
		for _, n := range offlineData {
			if v.Id == n.Rid {
				times, _ := strconv.ParseInt(n.Time, 10, 64)
				switch {
				case times >= oneDayStart && times < twoDayStart:
					one++
				case times >= twoDayStart && times < threeDayStart:
					two++
				case times >= threeDayStart && times < fourDayStart:
					three++
				case times >= fourDayStart && times < fiveDayStart:
					four++
				case times >= fiveDayStart && times < sixDayStart:
					five++
				case times >= sixDayStart && times < sevenDayStart:
					six++
				case times >= sevenDayStart && times < sevenDayEnd:
					seven++
				}
			}
		}
		//var tmpWeek Weekoffline
		var week Weekoffline
		week.Date = oneDayStart
		week.OfflineNum = one
		pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
		week.Date = twoDayStart
		week.OfflineNum = two
		pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
		week.Date = threeDayStart
		week.OfflineNum = three
		pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
		week.Date = fourDayStart
		week.OfflineNum = four
		pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
		week.Date = fiveDayStart
		week.OfflineNum = five
		pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
		week.Date = sixDayStart
		week.OfflineNum = six
		pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
		week.Date = sevenDayStart
		week.OfflineNum = seven
		pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
		pageonestatus = append(pageonestatus, pageonestatustmp)
	}

	var data StatusAnalysis
	data.StatusAnalysisData = pageonestatus
	data.StatusAnalysisTotal = pagetotal
	respData.Data = data
	c.JSON(http.StatusOK, &respData)
	return

	//// 获取离线总数
	//// 1. 获取所有设备的最新状态
	//var ds model.DeviceStatus
	//latestDeviceStatus, err := ds.LatestDeviceStatus()
	//if err != nil {
	//	respData.Error = "暂无数据"
	//	c.JSON(http.StatusOK, &respData)
	//	return
	//}
	//// 2. 获取当前所有小区的所有设备最新状态
	////var deviceStatuses []model.DeviceStatus
	////var deviceStatus model.DeviceStatus
	////
	//type DeviceData struct {
	//	Devicename   string
	//	StatusTypeId int64
	//	AreaId       int64
	//}
	//var devicedata []DeviceData
	//var devicedatatmp DeviceData
	//for _, v := range devices {
	//	for _, n := range latestDeviceStatus {
	//		if v.Parentid == n.DeviceId && v.Channel == n.Channel {
	//			//deviceStatus = n
	//			//deviceStatuses = append(deviceStatuses, deviceStatus)
	//
	//			devicedatatmp.AreaId = v.Areaid
	//			devicedatatmp.Devicename = v.Devicename
	//			devicedatatmp.StatusTypeId = n.StatusTypeId
	//			devicedata = append(devicedata, devicedatatmp)
	//		}
	//	}
	//}
	// 离线状态id
	//var offline int64 = 123
	//var offlineTotal int64 = 0
	//for _, v := range devicedata {
	//	if v.StatusTypeId == offline {
	//		offlineTotal++
	//	}
	//}
	//// 离线总数
	//pagetotal.OfflineTotal = offlineTotal
	//// 在线总数
	//pagetotal.OnlineTotal = deviceTotal - offlineTotal
	//
	//// 根据小区划分摄像头
	////var pageonestatus []PageOneStatus
	////------------------------------------------------------
	//// 图表③近七天离线数据
	//// 1. 获取时间戳
	//date := time.Now().Format("2006-01-02")
	//datetime, _ := time.ParseInLocation("2006-01-02 15:04:05", date+" 00:00:00", time.Local)
	//sevenDayStart := datetime.Unix()
	//sevenDayEnd := sevenDayStart + 60*60*24
	//sixDayStart := sevenDayStart - 60*60*24
	//fiveDayStart := sixDayStart - 60*60*24
	//fourDayStart := fiveDayStart - 60*60*24
	//threeDayStart := fourDayStart - 60*60*24
	//twoDayStart := threeDayStart - 60*60*24
	//oneDayStart := twoDayStart - 60*60*24
	//////fmt.Println(time.Unix(sevenDayStart, 0).Format("2006-01-02 15:04:05"))
	//////fmt.Println(time.Unix(sixDayStart, 0).Format("2006-01-02 15:04:05"))
	////// 2. 获取所有离线数据
	//var ds model.DeviceStatus
	//ds.StatusTypeId = offline
	//offlineData, err := ds.GetDataByStatusTypeId()
	//if err != nil {
	//	respData.Error = "暂无数据"
	//	c.JSON(http.StatusOK, &respData)
	//	return
	//}
	//for _, v := range areas {
	//	var pageonestatustmp PageOneStatus
	//	var tmp DeviceStatu
	//	var offlinenum int64 = 0
	//	for _, n := range devicedata {
	//		if v.Areaid == n.AreaId {
	//			if n.StatusTypeId == 123 {
	//				tmp.DeviceName = n.Devicename
	//				tmp.Status = "离线"
	//				// 图表①在线离线列表
	//				pageonestatustmp.DeviceStatus = append(pageonestatustmp.DeviceStatus, tmp)
	//				offlinenum++
	//			} else {
	//				tmp.DeviceName = n.Devicename
	//				tmp.Status = "在线"
	//				// 图表①在线离线列表
	//				pageonestatustmp.DeviceStatus = append(pageonestatustmp.DeviceStatus, tmp)
	//			}
	//		}
	//
	//	}
	//	// 图表①在线离线列表
	//	pageonestatustmp.AreaName = v.Areaname
	//	pageonestatustmp.AreaId = v.Areaid
	//	// 图表②离线数量
	//	pageonestatustmp.OfflineNum = offlinenum
	// ---------------------------------------------------------
	//	// 图表③近七天离线数据
	//	var one int64 = 0
	//	var two int64 = 0
	//	var three int64 = 0
	//	var four int64 = 0
	//	var five int64 = 0
	//	var six int64 = 0
	//	var seven int64 = 0
	//	for _, n := range offlineData {
	//		if v.Areaid == n.AreaId {
	//			times, _ := strconv.ParseInt(n.Time, 10, 64)
	//			switch {
	//			case times >= oneDayStart && times < twoDayStart:
	//				one++
	//			case times >= twoDayStart && times < threeDayStart:
	//				two++
	//			case times >= threeDayStart && times < fourDayStart:
	//				three++
	//			case times >= fourDayStart && times < fiveDayStart:
	//				four++
	//			case times >= fiveDayStart && times < sixDayStart:
	//				five++
	//			case times >= sixDayStart && times < sevenDayStart:
	//				six++
	//			case times >= sevenDayStart && times < sevenDayEnd:
	//				seven++
	//			}
	//		}
	//	}
	//	//var tmpWeek Weekoffline
	//	var week Weekoffline
	//	week.Date = oneDayStart
	//	week.OfflineNum = one
	//	pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
	//	week.Date = twoDayStart
	//	week.OfflineNum = two
	//	pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
	//	week.Date = threeDayStart
	//	week.OfflineNum = three
	//	pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
	//	week.Date = fourDayStart
	//	week.OfflineNum = four
	//	pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
	//	week.Date = fiveDayStart
	//	week.OfflineNum = five
	//	pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
	//	week.Date = sixDayStart
	//	week.OfflineNum = six
	//	pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
	//	week.Date = sevenDayStart
	//	week.OfflineNum = seven
	//	pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
	//
	//	//week[oneDayStart] = one
	//	//week[twoDayStart] = two
	//	//week[threeDayStart] = three
	//	//week[fourDayStart] = four
	//	//week[fiveDayStart] = five
	//	//week[sixDayStart] = six
	//	//week[sevenDayStart] = seven
	//	//pageonestatustmp.WeekOffline = append(pageonestatustmp.WeekOffline, week)
	//	pageonestatus = append(pageonestatus, pageonestatustmp)
	//}
	//
	//var data StatusAnalysis
	//data.StatusAnalysisData = pageonestatus
	//data.StatusAnalysisTotal = pagetotal
	//respData.Data = data
	//c.JSON(http.StatusOK, &respData)
	//return
}

// 修改后的第二个页面
type IndeAnlysis struct {
	IndeAnlysisTotal PageTwoTotal
	IndeAnlysisData  []PageTwoData
}
type PageTwoTotal struct {
	AlarmStartTotal           int64
	MonthAlarmStartTotal      int64
	StopAbnormalTotal         int64
	MonthStopAbnormalTotal    int64
	RemoteFormatHDDTotal      int64
	MonthRemoteFormatHDDTotal int64
}
type PageTwoData struct {
	AreaName  string
	DayData   Days
	MonthData Months
}
type Days struct {
	AlarmStart      int64
	StopAbnormal    int64
	RemoteFormatHDD int64
}
type Months struct {
	AlarmStart      int64
	StopAbnormal    int64
	RemoteFormatHDD int64
}

/*func changedPageTwo(c *gin.Context) {
	var respData RespData
	// 遮挡报警开始
	var alarmstart int64 = 10
	// 异常关机
	var stopAbnormal int64 = 186
	// 远程格式化硬盘
	var remoteFormatHDD int64 = 235

		//// 获取所有小区
		//var level int64 = 5
		//var area model.SysArea
		//area.Level = level
		//areas, err := area.GetAreaByLevel()
		//if err != nil || len(areas) == 0 {
		//	respData.Error = "暂无数据"
		//	c.JSON(http.StatusOK, &respData)
		//	return
		//}
	// 获取当前用户所属区域
	var ua model.SysUserArea
	ua.Userid = global.UserT.Id
	userAreas, err := ua.GetUserAreaByUserId()
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 获取当前用户所属区域ids
	var userareaids []int64
	for _, v := range userAreas {
		userareaids = append(userareaids, v.Areaid)
	}
	// 获取用户所属小区
	var level int64 = 5
	var area model.SysArea
	area.Level = level
	areas, err := area.GetAreaByLevelAndUserAreaIds(userareaids)
	if err != nil || len(areas) == 0 {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	var areaIds []int64
	for _, v := range areas {
		areaIds = append(areaIds, v.Areaid)
	}

	var total PageTwoTotal
	var ds model.DeviceStatus
	// 遮挡报警开始 - 总数
	ds.StatusTypeId = alarmstart
	alarmstartTotal, err := ds.GetTotalByAreaIdsAndSTid(areaIds)
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	total.AlarmStartTotal = alarmstartTotal

	// 异常关机 - 总数
	ds.StatusTypeId = stopAbnormal
	stopAbnormalTotal, err := ds.GetTotalByAreaIdsAndSTid(areaIds)
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	total.StopAbnormalTotal = stopAbnormalTotal

	// 远程格式化硬盘 - 总数
	ds.StatusTypeId = remoteFormatHDD
	remoteFormatHDDTotal, err := ds.GetTotalByAreaIdsAndSTid(areaIds)
	if err != nil {
		respData.Error = "暂无数据"
		c.JSON(http.StatusOK, &respData)
		return
	}
	total.RemoteFormatHDDTotal = remoteFormatHDDTotal

	// 获取当日时间戳，时间格式从 00:00:00 开始到 24:00:00 结束
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 00:00:00", time.Local)
	dayStartTime := t.Unix()
	dayEndTime := t.Unix() + (60 * 60 * 24)

	// 获取当月第一天时间和最后一天时间
	datetime := c.Query("datetime")
	var monthStartTime int64
	var monthEndTime int64
	if datetime == "" {
		firstDayMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1).Format("2006-01-02") + " 00:00:00"
		firstDayMonthTtam, _ := time.ParseInLocation("2006-01-02 15:04:05", firstDayMonth, time.Local)
		monthStartTime = firstDayMonthTtam.Unix()
		lastDayMonth := fmt.Sprint(time.Now().AddDate(0, 0, -time.Now().Day()+1).AddDate(0, 1, 0).Format("2006-01-02") + " 00:00:00")
		lastDayMonthTtam, _ := time.ParseInLocation("2006-01-02 15:04:05", lastDayMonth, time.Local)
		monthEndTime = lastDayMonthTtam.Unix()
	} else {
		matched, _ := regexp.MatchString("^\\d{4}-(((0([1-9]))|([1-9]))|(1(0|1|2)))$", datetime)
		if matched == false {
			respData.Error = "时间格式不正确"
			c.JSON(http.StatusOK, &respData)
			return
		}
		date := datetime + "-01 00:00:00"
		firstDayMonthTtam, _ := time.ParseInLocation("2006-1-2 15:4:5", date, time.Local)
		monthStartTime = firstDayMonthTtam.Unix()
		lastDayMonthTtam := fmt.Sprint(firstDayMonthTtam.AddDate(0, 0, -firstDayMonthTtam.Day()+1).AddDate(0, 1, 0).Format("2006-01-02") + " 00:00:00")
		lastDayMonthTta, _ := time.ParseInLocation("2006-1-2 15:4:5", lastDayMonthTtam, time.Local)
		monthEndTime = lastDayMonthTta.Unix()
	}

	var chartdata []PageTwoData
	var alarmstartTotalMonth int64 = 0
	var stopAbnormalTotalMonth int64 = 0
	var remoteFormatHDDTotalMonth int64 = 0
	for _, v := range areas {
		var chartTmp PageTwoData
		var dayTmp Days
		var monthTmp Months
		ds.AreaId = v.Areaid

		// 报警遮挡开始 - 当日数据
		ds.StatusTypeId = alarmstart
		dayalarmstartTotal, err := ds.GetDataByAreaSTidTime(dayStartTime, dayEndTime)
		if err != nil {
			respData.Error = "暂无数据"
			c.JSON(http.StatusOK, &respData)
			return
		}
		dayTmp.AlarmStart = dayalarmstartTotal
		// 报警遮挡开始 - 当月数据
		monthalarmstartTotal, err := ds.GetDataByAreaSTidTime(monthStartTime, monthEndTime)
		if err != nil {
			respData.Error = "暂无数据"
			c.JSON(http.StatusOK, &respData)
			return
		}
		monthTmp.AlarmStart = monthalarmstartTotal

		// 异常关机 - 当日数据
		ds.StatusTypeId = stopAbnormal
		daystopAbnormalTotal, err := ds.GetDataByAreaSTidTime(dayStartTime, dayEndTime)
		if err != nil {
			respData.Error = "暂无数据"
			c.JSON(http.StatusOK, &respData)
			return
		}
		dayTmp.StopAbnormal = daystopAbnormalTotal
		// 异常关机 - 当月数据
		monthstopAbnormalTotal, err := ds.GetDataByAreaSTidTime(monthStartTime, monthEndTime)
		if err != nil {
			respData.Error = "暂无数据"
			c.JSON(http.StatusOK, &respData)
			return
		}
		monthTmp.StopAbnormal = monthstopAbnormalTotal

		// 远程格式化硬盘 - 当日数据
		ds.StatusTypeId = remoteFormatHDD
		dayremoteFormatHDDTotal, err := ds.GetDataByAreaSTidTime(dayStartTime, dayEndTime)
		if err != nil {
			respData.Error = "暂无数据"
			c.JSON(http.StatusOK, &respData)
			return
		}
		dayTmp.RemoteFormatHDD = dayremoteFormatHDDTotal
		// 远程格式化硬盘 - 当月数据
		monthremoteFormatHDDTotal, err := ds.GetDataByAreaSTidTime(monthStartTime, monthEndTime)
		if err != nil {
			respData.Error = "暂无数据"
			c.JSON(http.StatusOK, &respData)
			return
		}
		monthTmp.RemoteFormatHDD = monthremoteFormatHDDTotal

		chartTmp.AreaName = v.Areaname
		chartTmp.DayData = dayTmp
		chartTmp.MonthData = monthTmp
		chartdata = append(chartdata, chartTmp)

		alarmstartTotalMonth += monthalarmstartTotal
		stopAbnormalTotalMonth += monthstopAbnormalTotal
		remoteFormatHDDTotalMonth += monthremoteFormatHDDTotal
	}
	total.MonthAlarmStartTotal = alarmstartTotalMonth
	total.MonthStopAbnormalTotal = stopAbnormalTotalMonth
	total.MonthRemoteFormatHDDTotal = remoteFormatHDDTotalMonth

	var data IndeAnlysis
	data.IndeAnlysisTotal = total
	data.IndeAnlysisData = chartdata
	respData.Data = data
	c.JSON(http.StatusOK, &respData)
	return
}*/
