package model

import (
	"fmt"
	"src/global"
)

type Restaurnt struct {
	// Id
	Id int64 `json:"id"`
	// Name 姓名
	Name string `json:"name"`
	// Time 注册时间
	Time int64 `json:"time"`
	// Phone 联系电话
	Phone string `json:"phone"`
	// Longitude 经度
	Longitude string `json:"longitude"`
	// Latitude 纬度
	Latitude string `json:"latitude"`
	// Address 密码
	Address string `json:"address"`
	// Status 状态
	Status string `json:"status"`
}

// 获取所有数据
func (self Restaurnt) All() (data []Restaurnt, err error) {
	err = global.Dbs.Table("restaurnt").Where("status = ?", 1).Order("id DESC").Find(&data).Error
	return
}

// 分页获取数据
func (self Restaurnt) Page(offset, limit int64) (data []Restaurnt, err error) {
	err = global.Dbs.Table("restaurnt").Where("status = ?", 1).Order("id DESC").Offset(offset).Limit(limit).Find(&data).Error
	return
}

// 获取餐厅总数
func (self Restaurnt) Total() (total int64, err error) {
	err = global.Dbs.Table("restaurnt").Where("status = ?", 1).Count(&total).Error
	return
}

// 新增一条餐厅数据
func (self Restaurnt) Add() (err error) {
	err = global.Dbs.Table("restaurnt").Create(&self).Error
	return
}

// 修改数据(包括删除)
func (self Restaurnt) Edit() (err error) {
	err = global.Dbs.Table("restaurnt").Save(&self).Error
	return
}

func (self Restaurnt) Test() (data []Restaurnt, err error) {
	var sql string = "select `name` from `restaurnt`"

	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		fmt.Println(sql)
		fmt.Println(err)
		return
	}
	var tmp Restaurnt
	for rows.Next() {
		err = rows.Scan(&tmp.Name)
		if err != nil {
			fmt.Println(err)
			return
		}
		data = append(data, tmp)
	}

	return
}
