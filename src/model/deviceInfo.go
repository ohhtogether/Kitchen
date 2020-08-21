package model

import (
	"fmt"
	"log"
	"src/conf"
	"src/global"
)

type DeviceInfo struct {
	Deviceid         int64  `json:"deviceid"`
	Devicename       string `json:"devicename"`
	Typeid           int64  `json:"typeid"`
	Rid              int64  `json:"rid"`
	Brand            string `json:"brand"`
	Model            string `json:"model"`
	SerialNumber     string `json:"serial_number"`
	ManufactureDate  string `json:"manufacture_date"`
	InstallationDate string `json:"installation_date"`
	AcceptanceDate   string `json:"acceptance_date"`
	InspectionDate   string `json:"inspection_date"`
	MaintenanceDate  string `json:"maintenance_date"`
	Ip               string `json:"ip"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	Port             int64  `json:"port"`
	Longitude        string `json:"longitude"`
	Latitude         string `json:"latitude"`
	Status           int64  `json:"status"`
	Channel          int64  `json:"channel"`
	Parentid         int64  `json:"parentid"`
	State            int64  `json:"state"`
	Coordinate       string `json:"coordinate"`
	Video            string `json:"video"`
	Rname            string `json:"rname"`
}

//添加一条设备信息
func (self DeviceInfo) AddOne() (id int64, err error) {
	var sql string = "insert into `device_info` (`devicename`,`typeid`,`rid`,`brand`,`model`,`serial_number`,`manufacture_date`,`installation_date`,`acceptance_date`,`inspection_date`,`maintenance_date`,`ip`,`username`,`password`,`port`,`longitude`,`latitude`,`status`,`channel`,`parentid`,`video`) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	result, err := global.Db.Exec(sql, self.Devicename, self.Typeid, self.Rid, self.Brand, self.Model, self.SerialNumber, self.ManufactureDate, self.InstallationDate, self.AcceptanceDate, self.InspectionDate, self.MaintenanceDate, self.Ip, self.Username, self.Password, self.Port, self.Longitude, self.Latitude, self.Status, self.Channel, self.Parentid, self.Video)
	if err != nil {
		log.Println(err)
		return
	}
	id, err = result.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	return
}

/*//批量添加设备信息
func AddMultiple(data []DeviceInfo) (err error) {
	// 拼接批量添加数据
	str := ""
	for _, deviceInfo := range data {
		str = str + "('" + deviceInfo.Devicename + "','" + strconv.FormatInt(deviceInfo.Typeid, 10) + "','" + strconv.FormatInt(deviceInfo.Areaid, 10) + "','" + deviceInfo.Brand + "','" + deviceInfo.Model + "','" + deviceInfo.SerialNumber + "','" + deviceInfo.ManufactureDate + "','" + deviceInfo.InstallationDate + "','" + deviceInfo.AcceptanceDate + "','" + deviceInfo.InspectionDate + "','" + deviceInfo.MaintenanceDate + "','" + deviceInfo.Ip + "','" + deviceInfo.Username + "','" + deviceInfo.Password + "','" + strconv.FormatInt(deviceInfo.Port, 10) + "','" + deviceInfo.Longitude + "','" + deviceInfo.Latitude + "','" + strconv.FormatInt(deviceInfo.Channel, 10) + "','" + strconv.FormatInt(deviceInfo.Parentid, 10) + "'),"
	}
	// 截取多余的逗号
	//str = str[0 : len(str)-1]
	str = str[:strings.LastIndex(str, ",")]
	// sql
	var sql string = "insert into `device_info` (`devicename`,`typeid`,`areaid`,`brand`,`model`,`serial_number`,`manufacture_date`,`installation_date`,`acceptance_date`,`inspection_date`,`maintenance_date`,`ip`,`username`,`password`,`port`,`longitude`,`latitude`,`channel`,`parentid`) value " + str
	// 入库
	stmt, err := global.Db.Prepare(sql)
	_, err = stmt.Exec()
	if err != nil {
		log.Println(err)
	}
	return
}*/

//修改一条设备信息
func (self DeviceInfo) Edit() (err error) {
	var sql string = "update `device_info` set `devicename`=?,`typeid`=?,`rid`=?,`brand`=?,`model`=?,`serial_number`=?,`manufacture_date`=?,`installation_date`=?,`acceptance_date`=?,`inspection_date`=?,`maintenance_date`=?,`ip`=?,`username`=?,`password`=?,`port`=?,`longitude`=?,`latitude`=?,`status`=?,`channel`=?,`parentid`=?,`video`=? where `deviceid`=?"
	stmt, err := global.Db.Prepare(sql)
	_, err = stmt.Exec(self.Devicename, self.Typeid, self.Rid, self.Brand, self.Model, self.SerialNumber, self.ManufactureDate, self.InstallationDate, self.AcceptanceDate, self.InspectionDate, self.MaintenanceDate, self.Ip, self.Username, self.Password, self.Port, self.Longitude, self.Latitude, self.Status, self.Channel, self.Parentid, self.Video, self.Deviceid)
	if err != nil {
		log.Println(err)
	}
	return
}

//获取当前用户所有设备信息
/*func (self DeviceInfo) All(keys string) (data []DeviceInfo, err error) {
	// 获取当前用户所属的区域
	var userarea SysUserArea
	userarea.Userid = global.UserT.Id
	areas, err := userarea.GetDataByUserId()
	if err != nil || len(areas) == 0 {
		err = errors.New("当前用户未绑定区域")
		return
	}
	like := "("
	for _, v := range areas {
		like += "`areaid` like '" + fmt.Sprint(v.Areaid) + "%' or "
	}
	like = like[:strings.LastIndex(like, " or")]
	like += ")"

	var sql string = "select `deviceid`,`devicename`,`typeid`,`areaid`,`brand`,`model`,`serial_number`,`manufacture_date`,`installation_date`,`acceptance_date`,`inspection_date`,`maintenance_date`,`ip`,`username`,`password`,`port`,`longitude`,`latitude`,`status`,`channel`,`parentid` from device_info where `status` = " + fmt.Sprint(conf.DEVICE_STATUS_NORMAL) + " and " + like + " order by `deviceid` desc"
	if keys != "" {
		sql = "select `deviceid`,`devicename`,`typeid`,`areaid`,`brand`,`model`,`serial_number`,`manufacture_date`,`installation_date`,`acceptance_date`,`inspection_date`,`maintenance_date`,`ip`,`username`,`password`,`port`,`longitude`,`latitude`,`status`,`channel`,`parentid` from device_info where `status` = " + fmt.Sprint(conf.DEVICE_STATUS_NORMAL) + " and " + like + " and (`devicename` like '%" + fmt.Sprint(keys) + "%' or `brand` like '%" + fmt.Sprint(keys) + "%') order by `deviceid` desc"
	}
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	var tmp DeviceInfo
	for rows.Next() {
		err = rows.Scan(&tmp.Deviceid, &tmp.Devicename, &tmp.Typeid, &tmp.Areaid, &tmp.Brand, &tmp.Model, &tmp.SerialNumber, &tmp.ManufactureDate, &tmp.InstallationDate, &tmp.AcceptanceDate, &tmp.InspectionDate, &tmp.MaintenanceDate, &tmp.Ip, &tmp.Username, &tmp.Password, &tmp.Port, &tmp.Longitude, &tmp.Latitude, &tmp.Status, &tmp.Channel, &tmp.Parentid)
		if err != nil {
			log.Println(err)
			return
		}
		data = append(data, tmp)
	}
	return
}*/
func (self DeviceInfo) GetDataByActiveUser(keys, Sql string) (data []DeviceInfo, err error) {
	var sql string = "select `deviceid`,`devicename`,`typeid`,`rid`,`brand`,`model`,`serial_number`,`manufacture_date`,`installation_date`,`acceptance_date`,`inspection_date`,`maintenance_date`,`ip`,`username`,`password`,`port`,`longitude`,`latitude`,`status`,`channel`,`parentid`,`video` from `device_info` where `status`=1 " + Sql + " order by `deviceid` desc"
	if keys != "" {
		sql = "select `deviceid`,`devicename`,`typeid`,`rid`,`brand`,`model`,`serial_number`,`manufacture_date`,`installation_date`,`acceptance_date`,`inspection_date`,`maintenance_date`,`ip`,`username`,`password`,`port`,`longitude`,`latitude`,`status`,`channel`,`parentid`,`video` from `device_info` where `status`=1 " + Sql + " and (`devicename` like '%" + fmt.Sprint(keys) + "%' or `brand` like '%" + fmt.Sprint(keys) + "%') order by `deviceid` desc"
	}

	err = global.Dbs.Table("device_info").Raw(sql).Scan(&data).Error
	return
}

/*func (self DeviceInfo) Pages(offset int64, limit int64, keys string) (data []DeviceInfo, err error) {
	// 获取当前用户所属的区域
	var userarea SysUserArea
	userarea.Userid = global.UserT.Id
	areas, err := userarea.GetDataByUserId()
	if err != nil || len(areas) == 0 {
		err = errors.New("当前用户未绑定区域")
		return
	}
	like := "("
	for _, v := range areas {
		like += "`areaid` like '" + fmt.Sprint(v.Areaid) + "%' or "
	}
	like = like[:strings.LastIndex(like, " or")]
	like += ")"

	var sql string = "select `deviceid`,`devicename`,`typeid`,`areaid`,`brand`,`model`,`serial_number`,`manufacture_date`,`installation_date`,`acceptance_date`,`inspection_date`,`maintenance_date`,`ip`,`username`,`password`,`port`,`longitude`,`latitude`,`status`,`channel`,`parentid` from device_info where `status` = " + fmt.Sprint(conf.DEVICE_STATUS_NORMAL) + " and " + like + " order by `deviceid` desc limit " + fmt.Sprint(offset) + ", " + fmt.Sprint(limit)
	if keys != "" {
		sql = "select `deviceid`,`devicename`,`typeid`,`areaid`,`brand`,`model`,`serial_number`,`manufacture_date`,`installation_date`,`acceptance_date`,`inspection_date`,`maintenance_date`,`ip`,`username`,`password`,`port`,`longitude`,`latitude`,`status`,`channel`,`parentid` from device_info where `status` = " + fmt.Sprint(conf.DEVICE_STATUS_NORMAL) + " and " + like + " and (`devicename` like '%" + fmt.Sprint(keys) + "%' or `brand` like '%" + fmt.Sprint(keys) + "%') order by `deviceid` desc limit " + fmt.Sprint(offset) + ", " + fmt.Sprint(limit)
	}
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	var tmp DeviceInfo
	for rows.Next() {
		err = rows.Scan(&tmp.Deviceid, &tmp.Devicename, &tmp.Typeid, &tmp.Areaid, &tmp.Brand, &tmp.Model, &tmp.SerialNumber, &tmp.ManufactureDate, &tmp.InstallationDate, &tmp.AcceptanceDate, &tmp.InspectionDate, &tmp.MaintenanceDate, &tmp.Ip, &tmp.Username, &tmp.Password, &tmp.Port, &tmp.Longitude, &tmp.Latitude, &tmp.Status, &tmp.Channel, &tmp.Parentid)
		if err != nil {
			log.Println(err)
			return
		}
		data = append(data, tmp)
	}
	return
}*/
func (self DeviceInfo) PagesByActiveUser(offset int64, limit int64, keys, Sql string) (data []DeviceInfo, err error) {
	var sql string = "select `deviceid`,`devicename`,`typeid`,`rid`,`brand`,`model`,`serial_number`,`manufacture_date`,`installation_date`,`acceptance_date`,`inspection_date`,`maintenance_date`,`ip`,`username`,`password`,`port`,`longitude`,`latitude`,`status`,`channel`,`parentid`,`video` from device_info where `status`=1 " + Sql + " order by `deviceid` desc limit " + fmt.Sprint(offset) + ", " + fmt.Sprint(limit)
	if keys != "" {
		sql = "select `deviceid`,`devicename`,`typeid`,`rid`,`brand`,`model`,`serial_number`,`manufacture_date`,`installation_date`,`acceptance_date`,`inspection_date`,`maintenance_date`,`ip`,`username`,`password`,`port`,`longitude`,`latitude`,`status`,`channel`,`parentid`,`video` from device_info where `status`=1 " + Sql + " and (`devicename` like '%" + fmt.Sprint(keys) + "%' or `brand` like '%" + fmt.Sprint(keys) + "%') order by `deviceid` desc limit " + fmt.Sprint(offset) + ", " + fmt.Sprint(limit)
	}
	err = global.Dbs.Table("device_info").Raw(sql).Scan(&data).Error
	return
}

// 获取总条数
func (self DeviceInfo) Count(keys string) (num int64, err error) {
	var sql string = fmt.Sprintf("select count(`deviceid`) from `device_info` where `status` = %v", conf.DEVICE_STATUS_NORMAL)
	if keys != "" {
		sql = fmt.Sprintf("select count(`deviceid`) from `device_info` where `status` = %v and (`devicename` like '%v' or `brand` like '%v')", conf.DEVICE_STATUS_NORMAL, "%"+keys+"%", "%"+keys+"%")
	}
	err = global.Db.QueryRow(sql).Scan(&num)
	if err != nil {
		log.Println(err)
	}
	return
}
func (self DeviceInfo) CountByActiveUser(keys, Sql string) (num int64, err error) {
	var sql string = "select count(`deviceid`) from `device_info` where `status`=1 " + Sql
	if keys != "" {
		sql = "select count(`deviceid`) from `device_info` where `status`=1 and (`devicename` like '%" + keys + "%' or `brand` like '%" + keys + "%') " + Sql
	}
	err = global.Db.QueryRow(sql).Scan(&num)
	if err != nil {
		log.Println(err)
	}
	return
}

// 通过DeviceId获取一条信息
func (self DeviceInfo) GetDataByDevId() (data DeviceInfo, err error) {
	var sql string = fmt.Sprintf("select `deviceid`,`devicename`,`typeid`,`rid`,`brand`,`model`,`serial_number`,`manufacture_date`,`installation_date`,`acceptance_date`,`inspection_date`,`maintenance_date`,`ip`,`username`,`password`,`port`,`longitude`,`latitude`,`status`,`channel`,`parentid`,`video` from device_info where `deviceid` = %v and `status` = %v", self.Deviceid, conf.DEVICE_STATUS_NORMAL)
	err = global.Db.QueryRow(sql).Scan(&data.Deviceid, &data.Devicename, &data.Typeid, &data.Rid, &data.Brand, &data.Model, &data.SerialNumber, &data.ManufactureDate, &data.InstallationDate, &data.AcceptanceDate, &data.InspectionDate, &data.MaintenanceDate, &data.Ip, &data.Username, &data.Password, &data.Port, &data.Longitude, &data.Latitude, &data.Status, &data.Channel, &data.Parentid, &data.Video)
	if err != nil {
		log.Println(err)
	}
	return
}

// 通过deviceid获取一条信息
func (self DeviceInfo) GetOneDataByDeviceId() (data DeviceInfo, err error) {
	var sql string = fmt.Sprintf("select `deviceid`,`devicename`,`typeid`,`rid`,`brand`,`model`,`serial_number`,`manufacture_date`,`installation_date`,`acceptance_date`,`inspection_date`,`maintenance_date`,`ip`,`username`,`password`,`port`,`longitude`,`latitude`,`status`,`channel`,`parentid`,`video` from device_info where `deviceid` = %v and `status` = %v", self.Deviceid, conf.DEVICE_STATUS_NORMAL)
	err = global.Db.QueryRow(sql).Scan(&data.Deviceid, &data.Devicename, &data.Typeid, &data.Rid, &data.Brand, &data.Model, &data.SerialNumber, &data.ManufactureDate, &data.InstallationDate, &data.AcceptanceDate, &data.InspectionDate, &data.MaintenanceDate, &data.Ip, &data.Username, &data.Password, &data.Port, &data.Longitude, &data.Latitude, &data.Status, &data.Channel, &data.Parentid, &data.Video)
	if err != nil {
		log.Println(err)
	}
	return
}

// 通过parentid和channel获取一条信息
func (self DeviceInfo) GetDataByDeviceId() (data DeviceInfo, err error) {
	var sql string = fmt.Sprintf("select `deviceid`,`devicename`,`typeid`,`rid`,`brand`,`model`,`serial_number`,`manufacture_date`,`installation_date`,`acceptance_date`,`inspection_date`,`maintenance_date`,`ip`,`username`,`password`,`port`,`longitude`,`latitude`,`status`,`channel`,`parentid`,`video` from device_info where `parentid` = %v and `channel` = %v and `status` = %v", self.Parentid, self.Channel, conf.DEVICE_STATUS_NORMAL)
	err = global.Db.QueryRow(sql).Scan(&data.Deviceid, &data.Devicename, &data.Typeid, &data.Rid, &data.Brand, &data.Model, &data.SerialNumber, &data.ManufactureDate, &data.InstallationDate, &data.AcceptanceDate, &data.InspectionDate, &data.MaintenanceDate, &data.Ip, &data.Username, &data.Password, &data.Port, &data.Longitude, &data.Latitude, &data.Status, &data.Channel, &data.Parentid, &data.Video)
	if err != nil {
		log.Println(err)
	}
	return
}

// 通过ParentIdAndChannel获取一条信息
func (self DeviceInfo) GetDataByParentIdAndChannel() (data DeviceInfo, err error) {
	var sql string = fmt.Sprintf("select `deviceid`,`devicename`,`typeid`,`rid`,`brand`,`model`,`serial_number`,`manufacture_date`,`installation_date`,`acceptance_date`,`inspection_date`,`maintenance_date`,`ip`,`username`,`password`,`port`,`longitude`,`latitude`,`status`,`channel`,`parentid`,`video` from `device_info` where `parentid` = %v and `channel`= %v and `status` = %v", self.Parentid, self.Channel, conf.DEVICE_STATUS_NORMAL)
	err = global.Db.QueryRow(sql).Scan(&data.Deviceid, &data.Devicename, &data.Typeid, &data.Rid, &data.Brand, &data.Model, &data.SerialNumber, &data.ManufactureDate, &data.InstallationDate, &data.AcceptanceDate, &data.InspectionDate, &data.MaintenanceDate, &data.Ip, &data.Username, &data.Password, &data.Port, &data.Longitude, &data.Latitude, &data.Status, &data.Channel, &data.Parentid, &data.Video)
	return
}

// 获取所有
func (self DeviceInfo) GetAll() (data []DeviceInfo, err error) {
	var sql string = "select `deviceid`,`devicename`,`typeid`,`rid`,`brand`,`model`,`serial_number`,`manufacture_date`,`installation_date`,`acceptance_date`,`inspection_date`,`maintenance_date`,`ip`,`username`,`password`,`port`,`longitude`,`latitude`,`status`,`channel`,`parentid`,`video` from `device_info` where `status` = " + fmt.Sprint(conf.DEVICE_STATUS_NORMAL)
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	var tmp DeviceInfo
	for rows.Next() {
		err = rows.Scan(&tmp.Deviceid, &tmp.Devicename, &tmp.Typeid, &tmp.Rid, &tmp.Brand, &tmp.Model, &tmp.SerialNumber, &tmp.ManufactureDate, &tmp.InstallationDate, &tmp.AcceptanceDate, &tmp.InspectionDate, &tmp.MaintenanceDate, &tmp.Ip, &tmp.Username, &tmp.Password, &tmp.Port, &tmp.Longitude, &tmp.Latitude, &tmp.Status, &tmp.Channel, &tmp.Parentid, &tmp.Video)
		if err != nil {
			log.Println(err)
			return
		}
		data = append(data, tmp)
	}
	return
}

// 通过typeid获取数据
func (self DeviceInfo) GetDataByTypeid() (data []DeviceInfo, err error) {
	var sql string = "select `deviceid`,`devicename`,`typeid`,`rid`,`brand`,`model`,`serial_number`,`manufacture_date`,`installation_date`,`acceptance_date`,`inspection_date`,`maintenance_date`,`ip`,`username`,`password`,`port`,`longitude`,`latitude`,`status`,`channel`,`parentid`,`video` from `device_info` where `status` = 1 and `typeid` = " + fmt.Sprint(self.Typeid)
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	var tmp DeviceInfo
	for rows.Next() {
		err = rows.Scan(&tmp.Deviceid, &tmp.Devicename, &tmp.Typeid, &tmp.Rid, &tmp.Brand, &tmp.Model, &tmp.SerialNumber, &tmp.ManufactureDate, &tmp.InstallationDate, &tmp.AcceptanceDate, &tmp.InspectionDate, &tmp.MaintenanceDate, &tmp.Ip, &tmp.Username, &tmp.Password, &tmp.Port, &tmp.Longitude, &tmp.Latitude, &tmp.Status, &tmp.Channel, &tmp.Parentid, &tmp.Video)
		if err != nil {
			log.Println(err)
			return
		}
		data = append(data, tmp)
	}
	return
}

// 获取所有父级数据
func (self DeviceInfo) GetDataByChannel() (data []DeviceInfo, err error) {
	var sql string = "select `deviceid`,`devicename`,`typeid`,`rid`,`brand`,`model`,`serial_number`,`manufacture_date`,`installation_date`,`acceptance_date`,`inspection_date`,`maintenance_date`,`ip`,`username`,`password`,`port`,`longitude`,`latitude`,`status`,`channel`,`parentid`,`video` from `device_info` where `status` = 1 and `channel` = " + fmt.Sprint(self.Channel)
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	var tmp DeviceInfo
	for rows.Next() {
		err = rows.Scan(&tmp.Deviceid, &tmp.Devicename, &tmp.Typeid, &tmp.Rid, &tmp.Brand, &tmp.Model, &tmp.SerialNumber, &tmp.ManufactureDate, &tmp.InstallationDate, &tmp.AcceptanceDate, &tmp.InspectionDate, &tmp.MaintenanceDate, &tmp.Ip, &tmp.Username, &tmp.Password, &tmp.Port, &tmp.Longitude, &tmp.Latitude, &tmp.Status, &tmp.Channel, &tmp.Parentid, &tmp.Video)
		if err != nil {
			log.Println(err)
			return
		}
		data = append(data, tmp)
	}
	return
}

// 获取设备在线/离线状态
func (self DeviceInfo) GetDeviceState() (data []DeviceInfo, err error) {
	sql := "SELECT a.`deviceId`, a.`devicename`, a.`rid`, b.`state` FROM `device_info` AS `a` LEFT JOIN `device_online` AS `b` ON a.`deviceid`=b.`deviceid` WHERE a.`status`=1;"
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	var tmp DeviceInfo
	for rows.Next() {
		err = rows.Scan(&tmp.Deviceid, &tmp.Devicename, &tmp.Rid, &tmp.State)
		if err != nil {
			log.Println(err)
			return
		}
		data = append(data, tmp)
	}
	return
}

// 获取地图小区设备定位数据
func (self DeviceInfo) Alls() (data []DeviceInfo, err error) {
	var sql string = "select `deviceid`,`devicename`,`typeid`,`rid`,`brand`,`model`,`serial_number`,`manufacture_date`,`installation_date`,`acceptance_date`,`inspection_date`,`maintenance_date`,`ip`,`username`,`password`,`port`,`longitude`,`latitude`,`status`,`channel`,`parentid`,`coordinate`,`video` from `device_info` where `status`=1 AND `coordinate`!=0"
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	var tmp DeviceInfo
	for rows.Next() {
		err = rows.Scan(&tmp.Deviceid, &tmp.Devicename, &tmp.Typeid, &tmp.Rid, &tmp.Brand, &tmp.Model, &tmp.SerialNumber, &tmp.ManufactureDate, &tmp.InstallationDate, &tmp.AcceptanceDate, &tmp.InspectionDate, &tmp.MaintenanceDate, &tmp.Ip, &tmp.Username, &tmp.Password, &tmp.Port, &tmp.Longitude, &tmp.Latitude, &tmp.Status, &tmp.Channel, &tmp.Parentid, &tmp.Coordinate, &tmp.Video)
		if err != nil {
			log.Println(err)
			return
		}
		data = append(data, tmp)
	}
	return
}
