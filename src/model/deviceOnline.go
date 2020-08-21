package model

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"src/global"
)

type DeviceOnline struct {
	Id          int64 `json:"id"`
	DeviceId    int64 `json:"deviceid"`
	State       int64 `json:"state"`
	Time        int64 `json:"time"`
	OfflineTime int64 `json:"offline_time"`
}

func (self DeviceOnline) Add() (err error) {
	sql := "INSERT INTO `device_online` (`deviceid`,`state`) VALUE (?,?);"
	_, err = global.Db.Exec(sql, self.DeviceId, self.State)
	if err != nil {
		log.Println(err)
	}
	return
}

// 批量修改设备状态
func (self DeviceOnline) BulkEditing(state int64, deviceSql string) (err error) {
	sql := "UPDATE `device_online` SET `state`=" + fmt.Sprint(state) + " WHERE `deviceid` IN (" + deviceSql + ")"
	stmt, err := global.Db.Prepare(sql)
	_, err = stmt.Exec()
	if err != nil {
		log.Println(err)
	}
	return
}

// 批量修改设备状态 - 事务
func (self DeviceOnline) BulkEditingTx(onlineSql, offlineSql string) (err error) {
	// 开始事务
	tx, err := global.Db.Begin()
	defer tx.Rollback()
	if err != nil {
		log.Println(err)
		return err
	}
	// 将设备批量修改为离线
	offsql := "UPDATE `device_online` SET `state`=0 WHERE `deviceid` IN (" + offlineSql + ")"
	_, err = tx.Exec(offsql)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}
	// 将设备批量修改为在线
	onsql := "UPDATE `device_online` SET `state`=1 WHERE `deviceid` IN (" + onlineSql + ")"
	_, err = tx.Exec(onsql)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}
	// 提交
	err = tx.Commit()
	return
}

// 获取所有最新状态数据
func (self DeviceOnline) AllNew() (data []DeviceOnline, err error) {
	sql := "SELECT `id`, `deviceid`, `state`, `time` FROM `device_online` WHERE `id` IN (SELECT MAX(`id`) FROM `device_online` GROUP BY `deviceid`);"
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var tmp DeviceOnline
		err = rows.Scan(&tmp.Id, &tmp.DeviceId, &tmp.State, &tmp.Time)
		if err != nil {
			log.Println(err)
			return
		}
		data = append(data, tmp)
	}
	return
}

// 添加最新状态
func (self DeviceOnline) Adds(arr []DeviceOnline) (err error) {
	tmp := ""
	for _, v := range arr {
		tmp = tmp + "(" + strconv.Itoa(int(v.DeviceId)) + "," + strconv.Itoa(int(v.State)) + "," + strconv.Itoa(int(v.Time)) + "),"
	}
	if len(tmp) > 0 {
		tmp = tmp[:strings.LastIndex(tmp, ",")]
	}

	sql := "INSERT INTO `device_online` (`deviceid`,`state`,`time`) VALUE " + tmp
	stmt, err := global.Db.Prepare(sql)
	_, err = stmt.Exec()
	if err != nil {
		log.Println(err)
	}
	return
}

/*// 获取所有设备状态
func (self DeviceOnline) All() (data []DeviceOnline, err error) {
	sql := "SELECT `id`,`deviceid`,`state` FROM `device_online`;"
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		return
	}
	var tmp DeviceOnline
	for rows.Next() {
		err = rows.Scan(&tmp.Id, &tmp.DeviceId, &tmp.State)
		if err != nil {
			return
		}
		data = append(data, tmp)
	}
	return
}*/

/*// 获取一条数据
func (self DeviceOnline) OneByDeviceid() (data DeviceOnline, err error) {
	sql := "SELECT `id`,`deviceid`,`state` FROM `device_online` WHERE `deviceid`=" + fmt.Sprint(self.DeviceId)
	err = global.Db.QueryRow(sql).Scan(&data.Id, &data.DeviceId, &data.State)
	return
}*/

// 获取最新一条状态
func (self DeviceOnline) OneByDeviceid() (data DeviceOnline, err error) {
	sql := fmt.Sprintf("SELECT `id`, `deviceid`, `state`, `time` FROM `device_online` WHERE `id`=(SELECT MAX(`id`) FROM `device_online` where `deviceid`=%v);", self.DeviceId)
	err = global.Db.QueryRow(sql).Scan(&data.Id, &data.DeviceId, &data.State, &data.Time)
	return
}

// 获取指定时间内所有离线的设备id
func (self DeviceOnline) OfflineTimeIds(start, end int64) (ids []int64, err error) {
	sql := fmt.Sprintf("SELECT `deviceid` FROM `device_online` WHERE `state`=0 AND `time`>=%v AND `time`<%v GROUP BY `deviceid`;", start, end)
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	var tmp int64
	for rows.Next() {
		err = rows.Scan(&tmp)
		if err != nil {
			log.Println(err)
			return
		}
		ids = append(ids, tmp)
	}
	return
}

func (self DeviceOnline) All() (data []DeviceOnline, err error) {
	sql := "SELECT `id`, `deviceid`, `state`, `time` FROM `device_online`"
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var tmp DeviceOnline
		err = rows.Scan(&tmp.Id, &tmp.DeviceId, &tmp.State, &tmp.Time)
		if err != nil {
			log.Println(err)
			return
		}
		data = append(data, tmp)
	}
	return
}
