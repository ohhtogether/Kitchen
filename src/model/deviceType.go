package model

import (
	"fmt"
	"log"

	"src/conf"
	"src/global"
)

type DeviceType struct {
	TypeId   int64  `json:"typeid"`
	TypeName string `json:"typename"`
	Des      string `json:"des"`
	Status   int64  `json:"status"`
}

//新增device type 数据
func (self DeviceType) Add() (err error) {
	var sql string = "insert into `device_type`(`typename`,`des`) value(?,?);"
	_, err = global.Db.Exec(sql, self.TypeName, self.Des)
	if err != nil {
		log.Println(err)
	}
	return
}

//修改device type 数据
func (self DeviceType) Edit() (err error) {
	var sql = "update `device_type` set `typename`=?,`des`=? where `typeid`=?"
	stmt, err := global.Db.Prepare(sql)
	_, err = stmt.Exec(self.TypeName, self.Des, self.TypeId)
	if err != nil {
		log.Println(err)
	}
	return
}

//删除device type 数据
func (self DeviceType) Del() (err error) {
	var sql = "update `device_type` set `status`=? where `typeid`=?"
	stmt, err := global.Db.Prepare(sql)
	_, err = stmt.Exec(self.Status, self.TypeId)
	if err != nil {
		log.Println(err)
	}
	return
}

//获取所有device type 数据
func (self DeviceType) All(keys string) (data []DeviceType, err error) {
	var sql string = fmt.Sprintf("select `typeid`,`typename`,`des`,`status` from `device_type` where `status` = %v  order by `typeid` desc", conf.DEVICE_TYPE_STATUS_NORMAL)
	if keys != "" {
		sql = fmt.Sprintf("select `typeid`,`typename`,`des`,`status` from `device_type` where `status` = %v and (`typename` like '%v' or `des` like '%v') order by `typeid` desc", conf.DEVICE_TYPE_STATUS_NORMAL, "%"+keys+"%", "%"+keys+"%")
	}
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	var tmp DeviceType
	for rows.Next() {
		err = rows.Scan(&tmp.TypeId, &tmp.TypeName, &tmp.Des, &tmp.Status)
		if err != nil {
			log.Println(err)
			return
		}
		data = append(data, tmp)
	}
	return
}

//获取device type 数据 分页
func (self DeviceType) Pages(offset int64, limit int64, keys string) (data []DeviceType, err error) {
	var sql string = fmt.Sprintf("select `typeid`,`typename`,`des`,`status` from `device_type` where `status` = %v  order by `typeid` desc limit %v, %v", conf.DEVICE_TYPE_STATUS_NORMAL, offset, limit)
	if keys != "" {
		sql = fmt.Sprintf("select `typeid`,`typename`,`des`,`status` from `device_type` where `status` = %v and (`typename` like '%v' or `des` like '%v') order by `typeid` desc limit %v, %v", conf.DEVICE_TYPE_STATUS_NORMAL, "%"+keys+"%", "%"+keys+"%", offset, limit)
	}
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	var tmp DeviceType
	for rows.Next() {
		err = rows.Scan(&tmp.TypeId, &tmp.TypeName, &tmp.Des, &tmp.Status)
		if err != nil {
			log.Println(err)
			return
		}
		data = append(data, tmp)
	}
	return
}

// 获取总条数
func (self DeviceType) Count(keys string) (num int64, err error) {
	var sql string = fmt.Sprintf("select count(`typeid`) from `device_type` where `status` = %v", conf.DEVICE_TYPE_STATUS_NORMAL)
	if keys != "" {
		sql = fmt.Sprintf("select count(`typeid`) from `device_type` where `status` = %v and (`typename` like '%v' or `des` like '%v')", conf.DEVICE_TYPE_STATUS_NORMAL, "%"+keys+"%", "%"+keys+"%")
	}
	err = global.Db.QueryRow(sql).Scan(&num)
	if err != nil {
		log.Println(err)
	}
	return
}

//获取一条device type数据
func (self DeviceType) One() (data DeviceType, err error) {
	var sql string = fmt.Sprintf("select `typeid`,`typename`,`des`,`status` from `device_type` where `status`=1 and `typeid`='%v'", self.TypeId)
	err = global.Db.QueryRow(sql).Scan(&data.TypeId, &data.TypeName, &data.Des, &data.Status)
	if err != nil {
		log.Println(err)
	}
	return
}
