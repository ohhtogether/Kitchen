package model

import (
	"fmt"
	"log"
	"src/global"
	"strings"
)

type DeviceStatus struct {
	Id             int64  `json:"id"`
	DeviceId       int64  `json:"deviceid"`
	DeviceName     string `json:"deviceName"`
	Time           string `json:"time"`
	StatusTypeId   int64  `json:"statusTypeId"`
	TypeName       string `json:"statusName"`
	TypeDes        string `json:"statusDes"`
	TypeBrand      string `json:"typeBrand"`
	NetUser        string `json:"netUser"`
	RemoteAddr     string `json:"remoteAddr"`
	Channel        int64  `json:"channel"`
	Rid            int64  `json:"rid"`
	DeviceIP       string `json:"deviceIP"`
	ParentTypeId   int64  `json:"parentTypeId"`
	ParentTypeName string `json:"parentTypeName"`
	ParentTypeDes  string `json:"parentTypeDes"`
	Sum            int64  `json:"sum"`
	TimeInt        int64  `json:"timeint"`
	TypeId         int64  `json:"typeId"`
}

//设备状态也搜索的关键词
type StatusSearchKeys struct {
	Rid             string
	TypeId          string
	StatusTypeId    string
	StatusTypeIds   []int64
	StatusTypeJudge int64 // 1 = string ; 2 = []int64
	StartTime       string
	EndTime         string
}

/*
	从device_status中获取每个设备的最新一条记录的id数组
*/
func (self DeviceStatus) LatestDeviceStatus() (data []DeviceStatus, err error) {

	return
}

//通过device-channel 获取指定一个设备的状态历史记录
func (self DeviceStatus) DevStuLogs(offset string, limit string, keys StatusSearchKeys) (data []DeviceStatus, err error) {
	sql := fmt.Sprintf("SELECT `id`,`deviceid`,`time`,`statusTypeId`,`netUser`,`remoteAddr`,`channel`,`rid`,`deviceIP` FROM `device_status` WHERE `deviceid`=%v AND `channel`=%v", self.DeviceId, self.Channel)
	if keys.StatusTypeJudge == 1 {
		sql += fmt.Sprintf(" AND `statusTypeId`=%v", keys.StatusTypeId)
	} else if keys.StatusTypeJudge == 2 {
		statusTypeIdStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(keys.StatusTypeIds)), ","), "[]")
		sql += fmt.Sprintf(" AND `statusTypeId` IN (%v)", statusTypeIdStr)
	}
	if keys.StartTime != "" && keys.EndTime != "" {
		sql += fmt.Sprintf(" AND `time`>=%v AND `time`<%v", keys.StartTime, keys.EndTime)
	}
	sql += " ORDER BY `id` DESC"
	if offset != "" && limit != "" {
		sql = fmt.Sprintf("%v limit %v,%v", sql, offset, limit)
	}
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var tmp DeviceStatus
		err = rows.Scan(&tmp.Id, &tmp.DeviceId, &tmp.Time, &tmp.StatusTypeId, &tmp.NetUser, &tmp.RemoteAddr, &tmp.Channel, &tmp.Rid, &tmp.DeviceIP)
		if err != nil {
			log.Println(err)
			return
		}
		data = append(data, tmp)
	}
	return
}

// 获取单个设备的所有设备状态总条数
func (self DeviceStatus) GetDeviceStatusNumsByDevId(keys StatusSearchKeys) (num int64, err error) {
	sql := fmt.Sprintf("SELECT count(`id`) FROM `device_status` WHERE `deviceid`=%v AND `channel`=%v", self.DeviceId, self.Channel)
	if keys.StatusTypeJudge == 1 {
		sql += fmt.Sprintf(" AND `statusTypeId`=%v", keys.StatusTypeId)
	} else if keys.StatusTypeJudge == 2 {
		statusTypeIdStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(keys.StatusTypeIds)), ","), "[]")
		sql += fmt.Sprintf(" AND `statusTypeId` IN (%v)", statusTypeIdStr)
	}
	if keys.StartTime != "" && keys.EndTime != "" {
		sql += fmt.Sprintf(" AND `time`>=%v AND `time`<%v", keys.StartTime, keys.EndTime)
	}
	err = global.Db.QueryRow(sql).Scan(&num)
	return
}

// 获取最新状态IDS (MAX(`id`))
// 根据状态类型ids获取最新的设备状态ids
func (self DeviceStatus) GetDataByMaxId(statusTypeIds string) (dsIds []int64, err error) {
	sql := fmt.Sprintf("SELECT MAX(`id`) FROM `device_status` WHERE `statusTypeId` IN (%v) GROUP BY `deviceid`,`channel`;", statusTypeIds)
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var tmpId int64 = 0
		err = rows.Scan(&tmpId)
		if err != nil {
			log.Println(err)
			return
		}
		dsIds = append(dsIds, tmpId)
	}
	return
}

// 根据最新状态IDS获取设备
func (self DeviceStatus) GetLatestStatusDeviceInfo(Ids, offset, limit string, keys StatusSearchKeys) (data []DeviceStatus, err error) {
	sql := fmt.Sprintf("SELECT a.`id`,a.`time`,a.`statusTypeId`,a.`netUser`,a.`remoteAddr`,a.`channel`,a.`rid`,a.`deviceid`,b.`devicename`,b.`typeid`,b.`ip` FROM `device_status` AS `a` LEFT JOIN `device_info` AS `b` ON a.`deviceid`=b.`parentid` AND a.`channel`=b.`channel` WHERE a.`id` IN (%v) AND b.`status`=1", Ids)
	if keys.Rid != "" {
		sql += fmt.Sprintf(" AND a.`rid`=%v", keys.Rid)
	}
	if keys.TypeId != "" {
		sql += fmt.Sprintf(" AND b.`typeid`=%v", keys.TypeId)
	}
	if keys.StatusTypeJudge == 1 {
		sql += fmt.Sprintf(" AND a.`statusTypeId`=%v", keys.StatusTypeId)
	} else if keys.StatusTypeJudge == 2 {
		statusTypeIdStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(keys.StatusTypeIds)), ","), "[]")
		sql += fmt.Sprintf(" AND a.`statusTypeId` IN (%v)", statusTypeIdStr)
	}
	if keys.StartTime != "" && keys.EndTime != "" {
		sql += fmt.Sprintf(" AND a.`time`>=%v AND a.`time`<%v", keys.StartTime, keys.EndTime)
	}
	if offset != "" && limit != "" {
		sql = fmt.Sprintf("%v limit %v,%v", sql, offset, limit)
	}

	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var temp DeviceStatus
		err = rows.Scan(&temp.Id, &temp.Time, &temp.StatusTypeId, &temp.NetUser, &temp.RemoteAddr, &temp.Channel, &temp.Rid, &temp.DeviceId, &temp.DeviceName, &temp.TypeId, &temp.DeviceIP)
		if err != nil {
			log.Println(err)
			return
		}
		data = append(data, temp)
	}
	return
}

//获取总条数
func GetDeviceStatusNums(Ids string, keys StatusSearchKeys) (num int64, err error) {
	sql := fmt.Sprintf("SELECT count(a.`id`) FROM `device_status` AS `a` LEFT JOIN `device_info` AS `b` ON a.`deviceid`=b.`parentid` AND a.`channel`=b.`channel` WHERE a.`id` IN (%v) AND b.`status`=1", Ids)
	if keys.Rid != "" {
		sql += fmt.Sprintf(" AND a.`rid`=%v", keys.Rid)
	}
	if keys.TypeId != "" {
		sql += fmt.Sprintf(" AND b.`typeid`=%v", keys.TypeId)
	}
	if keys.StatusTypeJudge == 1 {
		sql += fmt.Sprintf(" AND a.`statusTypeId`=%v", keys.StatusTypeId)
	} else if keys.StatusTypeJudge == 2 {
		statusTypeIdStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(keys.StatusTypeIds)), ","), "[]")
		sql += fmt.Sprintf(" AND a.`statusTypeId` IN (%v)", statusTypeIdStr)
	}
	if keys.StartTime != "" && keys.EndTime != "" {
		sql += fmt.Sprintf(" AND a.`time`>=%v AND a.`time`<%v", keys.StartTime, keys.EndTime)
	}
	err = global.Db.QueryRow(sql).Scan(&num)
	return
}

// 通过状态id获取数据
func (self DeviceStatus) GetDataByStatusTypeId() (data []DeviceStatus, err error) {
	sql := "select `id`,`deviceid`,`time`,`statusTypeId`,`netUser`,`remoteAddr`,`channel`,`rid`,`deviceIP` from `device_status` where `statusTypeId`=" + fmt.Sprint(self.StatusTypeId)
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var tmp DeviceStatus
		err = rows.Scan(&tmp.Id, &tmp.DeviceId, &tmp.Time, &tmp.StatusTypeId, &tmp.NetUser, &tmp.RemoteAddr, &tmp.Channel, &tmp.Rid, &tmp.DeviceIP)
		if err != nil {
			log.Println(err)
			return
		}
		data = append(data, tmp)
	}
	return
}

// 通过设备parentid和channel获取最新状态id
func (self DeviceStatus) GetMaxId() (maxid int64, err error) {
	sql := "SELECT MAX(`id`) AS `id` FROM `device_status` WHERE `deviceid`=" + fmt.Sprint(self.DeviceId) + " AND `channel`=" + fmt.Sprint(self.Channel) + ";"
	err = global.Db.QueryRow(sql).Scan(&maxid)
	return
}

func (self DeviceStatus) One() (data DeviceStatus, err error) {
	var sql string = fmt.Sprintf("select `id`,`deviceid`,`time`,`statusTypeId`,`netUser`,`remoteAddr`,`channel`,`rid`,`deviceIP` from `device_status` where `id` = %v", self.Id)
	err = global.Db.QueryRow(sql).Scan(&data.Id, &data.DeviceId, &data.Time, &data.StatusTypeId, &data.NetUser, &data.RemoteAddr, &data.Channel, &data.Rid, &data.DeviceIP)
	if err != nil {
		log.Println(err)
	}
	return
}

// 通过时间范围获取数据
func (self DeviceStatus) GetDataByTime(startTime, endTime string) (data []DeviceStatus, err error) {
	sql := "select `deviceid`,`time`,`statusTypeId`,`channel` from `device_status` where `time` > " + startTime + " and `time` <= " + endTime
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var tmp DeviceStatus
		err = rows.Scan(&tmp.DeviceId, &tmp.Time, &tmp.StatusTypeId, &tmp.Channel)
		if err != nil {
			log.Println(err)
			return
		}
		data = append(data, tmp)
	}
	return
}
