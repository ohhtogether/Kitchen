package model

import (
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
	Status int64 `json:"status"`
}

// 获取所有数据
func (self Restaurnt) All() (data []Restaurnt, err error) {
	err = global.Dbs.Table("restaurnt").Where("status = ?", 1).Order("id DESC").Find(&data).Error
	return
}

// 分页获取数据
func (self Restaurnt) Page(offset, limit int64) (data []Restaurnt, err error) {
	if offset == 0 || limit == 0 {
		err = global.Dbs.Table("restaurnt").Where("status = ?", 1).Order("id DESC").Find(&data).Error
	} else {
		err = global.Dbs.Table("restaurnt").Where("status = ?", 1).Order("id DESC").Offset(offset).Limit(limit).Find(&data).Error
	}
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

func (self Restaurnt) One() (user Restaurnt, err error) {
	err = global.Dbs.Table("restaurnt").Where("status = ?", 1).First(&user, self.Id).Error
	return
}
